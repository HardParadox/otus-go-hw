package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(value string) (string, error) {
	if len(value) == 0 {
		return "", nil
	}

	var resultString strings.Builder
	valueInRuneArray := []rune(value)

	for index, char := range valueInRuneArray {
		if unicode.IsDigit(char) {
			if index == 0 {
				return "", ErrInvalidString
			}

			if len(valueInRuneArray) > index+1 && unicode.IsDigit(valueInRuneArray[index+1]) {
				return "", ErrInvalidString
			}

			continue
		}

		if len(valueInRuneArray) > index+1 {
			res, err := strconv.Atoi(string(valueInRuneArray[index+1]))

			if err == nil {
				resultString.WriteString(strings.Repeat(string(char), res))
				continue
			}
		}

		resultString.WriteRune(char)
	}

	return resultString.String(), nil
}
