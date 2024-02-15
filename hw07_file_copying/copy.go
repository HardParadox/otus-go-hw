package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile         = errors.New("unsupported file")
	ErrOffsetExceedsFileSize   = errors.New("offset exceeds file size")
	ErrUnknownOriginalFileSize = errors.New("original file size unknown")
	ErrorOpenFile              = errors.New("file open failed")
	ErrorCreateFile            = errors.New("file create failed")
	ErrorNegativeOffset        = errors.New("offset can not be negative")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if filepath.Ext(fromPath) == "" || filepath.Ext(fromPath) != filepath.Ext(toPath) {
		return ErrUnsupportedFile
	}
	if offset < 0 {
		return ErrorNegativeOffset
	}

	file, err := os.Open(fromPath)
	if err != nil {
		return ErrorOpenFile
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return ErrorOpenFile
	}
	defer file.Close()

	if fileInfo.Size() <= limit {
		limit = fileInfo.Size()
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	if fileInfo.Size() == 0 {
		return ErrUnknownOriginalFileSize
	}

	targetFile, err := os.Create(toPath)
	if err != nil {
		return ErrorCreateFile
	}

	defer targetFile.Close()

	bar := pb.Full.Start64(limit)

	var buf []byte

	if limit == 0 {
		read, err := io.ReadAll(file)
		buf = read[offset:]

		if err != nil {
			return err
		}
	} else {
		var bufLen int64
		if offset+limit > fileInfo.Size() {
			bufLen = fileInfo.Size() - offset
		} else {
			bufLen = limit
		}

		buf = make([]byte, bufLen)

		file.ReadAt(buf, offset)

		bar.Finish()
	}

	_, err = targetFile.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
