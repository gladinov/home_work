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
	letter := 0
	num := 1

	if len(runes) == 0 {
		return "", nil
	}

	if _, ok := isNum(runes[0]); ok {
		return "", ErrInvalidString
	}

	for num < len(runes) {
		intValue, ok := isNum(runes[num])
		if !ok {
			if num-letter > 1 {
				letter = num
			} else {
				str := string(runes[letter])
				builder.WriteString(str)
				letter++
			}
			num++
			if letter == len(runes)-1 {
				str := string(runes[letter])
				builder.WriteString(str)
			}
		} else {
			if num-letter > 1 {
				return "", ErrInvalidString
			}
			str := string(runes[letter])
			for range intValue {
				builder.WriteString(str)
			}
			num++

		}
	}
	res := builder.String()

	return res, nil
}

func isNum(in rune) (int, bool) {
	str := string(in)

	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}
	return num, true
}
