package main

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	var file *os.File
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	environment := make(Environment, len(files))

	for _, f := range files {
		file, err = os.Open(path.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}

		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}

		if fi.Size() == 0 {
			environment[filepath.Base(f.Name())] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		reader := bufio.NewReader(file)

		line, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		cleanStr := bytes.ReplaceAll(line, []byte{0}, []byte("\n"))

		environment[filepath.Base(f.Name())] = EnvValue{
			Value:      strings.TrimRight(string(cleanStr), " \t"),
			NeedRemove: false,
		}
	}

	defer file.Close()
	return environment, nil
}
