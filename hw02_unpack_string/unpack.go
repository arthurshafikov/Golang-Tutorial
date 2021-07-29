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
		newstr    strings.Builder
		letter    string
		append    string
		convert   bool
		echoSlash bool
	)
	for _, sym := range str {
		if unicode.IsDigit(sym) && !convert {
			if letter == "" {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(sym))
			append = strings.Repeat(letter, count)
			letter = ""
		} else {
			append = letter
			if sym == '\\' {
				if letter == "\\" && !echoSlash {
					append = ""
					echoSlash = true
				} else {
					echoSlash = false
				}
			} else {
				if letter == "\\" {
					append = ""
				}
			}
			letter = string(sym)

			convert = false
			if letter == "\\" && !echoSlash {
				convert = true
			}
		}
		newstr.WriteString(append)
	}
	newstr.WriteString(letter)
	return newstr.String(), nil
}
