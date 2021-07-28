package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var (
		newstr strings.Builder
		letter string
		append string
	)
	for _, sym := range str {
		if unicode.IsDigit(sym) {
			if letter == "" {
				return "", ErrInvalidString
			}

			count, _ := strconv.Atoi(string(sym))

			append = strings.Repeat(letter, count)
			letter = ""
		} else {
			append = letter
			letter = string(sym)
		}
		newstr.WriteString(append)
	}
	newstr.WriteString(letter)
	return newstr.String(), nil
}
