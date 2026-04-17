package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	builder := strings.Builder{}

	runes := []rune(in)

	if len(runes) == 0 {
		return "", nil
	}

	if _, ok := isChar(runes[0]); ok {
		return "", ErrInvalidString
	}

	for i := 1; i < len(runes); i++ {
		curr := runes[i]
		prev := runes[i-1]
		char, ok := isChar(curr)
		if ok {
			if _, prevOk := isChar(prev); prevOk {
				return "", ErrInvalidString
			}
			str := strings.Repeat(string(prev), char)
			builder.WriteString(str)
			if i == len(runes)-1 {
				break
			}
		} else {
			if _, prevOk := isChar(prev); prevOk {
				if i == len(runes)-1 {
					builder.WriteString(string(curr))
				}
				continue
			}
			builder.WriteString(string(prev))
			if i == len(runes)-1 {
				builder.WriteString(string(curr))
			}
		}
	}
	res := builder.String()
	return res, nil
}

func isChar(in rune) (int, bool) {
	if in >= '0' && in <= '9' {
		char, _ := strconv.Atoi(string(in))
		return char, true
	}
	return 0, false
}
