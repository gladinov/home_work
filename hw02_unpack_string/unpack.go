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

	if _, ok := isDigit(runes[0]); ok {
		return "", ErrInvalidString
	}

	for i := 1; i < len(runes); i++ {
		curr := runes[i]
		prev := runes[i-1]
		dig, ok := isDigit(curr)
		isLast := i == len(runes)-1
		_, prevOk := isDigit(prev)

		if ok {
			if prevOk {
				return "", ErrInvalidString
			}
			str := strings.Repeat(string(prev), dig)
			builder.WriteString(str)
			if isLast {
				break
			}
			continue
		}
		if prevOk {
			if isLast {
				builder.WriteString(string(curr))
			}
			continue
		}
		builder.WriteString(string(prev))
		if isLast {
			builder.WriteString(string(curr))
		}
	}
	return builder.String(), nil
}

func isDigit(in rune) (int, bool) {
	if in >= '0' && in <= '9' {
		dig, _ := strconv.Atoi(string(in))
		return dig, true
	}
	return 0, false
}
