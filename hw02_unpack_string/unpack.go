package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

const backslash = 92

var ErrInvalidString = errors.New("invalid string")

type escapedChar struct {
	set bool
	pos int
}

func (c *escapedChar) isEscaped(prev int) bool {
	if c.set && prev == c.pos {
		return true
	}
	return false
}

func Unpack(in string) (string, error) {
	builder := strings.Builder{}

	runes := []rune(in)

	if len(runes) == 0 {
		return "", nil
	}

	if _, ok := isDigit(runes[0]); ok {
		return "", ErrInvalidString
	}

	if len(runes) == 1 {
		builder.WriteString(string(runes[0]))
		return builder.String(), nil
	}
	escapedChar := escapedChar{
		set: false,
		pos: 0,
	}

	for i := 1; i < len(runes); i++ {
		curr := runes[i]
		prev := runes[i-1]
		isLast := i == len(runes)-1
		dig, currDigOk := isDigit(curr)
		_, prevDigOk := isDigit(prev)

		if currDigOk {
			if prevDigOk {
				if escapedChar.isEscaped(i - 1) {
					str := strings.Repeat(string(prev), dig)
					builder.WriteString(str)
					continue
				}
				return "", ErrInvalidString
			}
			if isBackslash(prev) {
				if escapedChar.isEscaped(i - 1) {
					str := strings.Repeat(string(prev), dig)
					builder.WriteString(str)
					continue
				}
			}
			str := strings.Repeat(string(prev), dig)
			builder.WriteString(str)
			if isLast {
				break
			}
			continue
		}

		if isBackslash(curr) {
			// Если предыдущий элемент был
			if isBackslash(prev) {
				builder.WriteString(string(curr))
			}
			continue
		}
		if !currDigOk {
			if isBackslash(prev) {
				if escapedChar.isEscaped(i - 1) {
					builder.WriteString(string(prev))
				} else {
					return "", ErrInvalidString
				}

				continue
			}
			if prevDigOk {
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
	}
	return builder.String(), nil
}

func hanldeDig() {}
func handleBackslash() {
}
func handleOther() {}

func isDigit(in rune) (int, bool) {
	if in >= '0' && in <= '9' {
		dig, _ := strconv.Atoi(string(in))
		return dig, true
	}
	return 0, false
}

func isBackslash(in rune) bool {
	if in == backslash {
		return true
	}
	return false
}
