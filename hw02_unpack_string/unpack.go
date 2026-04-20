package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

const backslash rune = '\\'

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

func (c *escapedChar) Set(i int) {
	c.pos = i
	c.set = true
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
	esc := escapedChar{
		set: false,
		pos: 0,
	}
loop:
	for i := 1; i < len(runes); i++ {
		curr := runes[i]
		_, currDigOk := isDigit(curr)

		switch {
		// Обрабатываем числа
		case currDigOk:
			isContinue, err := handleDig(runes, i, &esc, &builder)
			if err != nil {
				return "", err
			}
			if isContinue {
				continue
			}
			break loop
			// если текущий элемент бэкслеш
		case isBackslash(curr):
			err := handleBackslash(runes, i, &esc, &builder)
			if err != nil {
				return "", err
			}
			continue
			// Если элемент не число и не бэкслеш
		default:
			err := handleOther(runes, i, &esc, &builder)
			if err != nil {
				return "", err
			}
			continue
		}
	}
	return builder.String(), nil
}

func handleDig(runes []rune, i int, esc *escapedChar, builder *strings.Builder) (bool, error) {
	curr := runes[i]
	prev := runes[i-1]
	isLast := i == len(runes)-1
	dig, _ := isDigit(curr)
	_, prevDigOk := isDigit(prev)

	// Если предыдущий жлемент число
	if prevDigOk {
		// Если предыдущее число экранировано
		if esc.isEscaped(i - 1) {
			str := strings.Repeat(string(prev), dig)
			builder.WriteString(str)
			return true, nil
		}
		return false, ErrInvalidString
	}
	if isBackslash(prev) {
		// Если предыдущий элемент бэкслэш
		if esc.isEscaped(i - 1) {
			// если бэкслэш экранирован
			str := strings.Repeat(string(prev), dig)
			builder.WriteString(str)
			return true, nil
		}
		if isLast {
			builder.WriteString(string(curr))
			return false, nil
		}

		// Если предыдущий жлемент не экранирован
		// То он экранирует число
		esc.Set(i)
		return true, nil
	}
	// Если предыдущий элемент не число и не бэкслэш
	str := strings.Repeat(string(prev), dig)
	builder.WriteString(str)
	if isLast {
		return false, nil
	}
	return true, nil
}

func handleBackslash(runes []rune, i int, esc *escapedChar, builder *strings.Builder) error {
	prev := runes[i-1]
	isLast := i == len(runes)-1
	_, prevDigOk := isDigit(prev)

	// Если текущий индекс последний
	if isLast {
		// Возвращаем ошибку
		return ErrInvalidString
	}
	// Если предудущий элемент экранирован
	if esc.isEscaped(i - 1) {
		str := string(prev)
		builder.WriteString(str)
		return nil
	}
	// Если предыдущий элемент бэкслэш
	if isBackslash(prev) {
		// То экранируем бекслеш
		esc.Set(i)
		return nil
	}

	if prevDigOk {
		return nil
	}
	// Если предыдущий элемент не число и не бекслэш
	str := string(prev)
	builder.WriteString(str)

	return nil
}

func handleOther(runes []rune, i int, esc *escapedChar, builder *strings.Builder) error {
	curr := runes[i]
	prev := runes[i-1]
	isLast := i == len(runes)-1
	_, prevDigOk := isDigit(prev)

	// Если предудущий элемент экранирован
	if esc.isEscaped(i - 1) {
		str := string(prev)
		builder.WriteString(str)
		if isLast {
			builder.WriteString(string(curr))
		}
		return nil
	}
	// Если предыдущий не экранированный бекслеш
	if isBackslash(prev) {
		return ErrInvalidString
	}
	// Если пердыдущий число
	if prevDigOk {
		// Если последний элемент , то добавляем
		if isLast {
			builder.WriteString(string(curr))
		}
		return nil
	}
	// Если предыдущий не число и не бекслеш, то записываем
	builder.WriteString(string(prev))
	//
	if isLast {
		builder.WriteString(string(curr))
	}
	return nil
}

func isDigit(in rune) (int, bool) {
	if in >= '0' && in <= '9' {
		dig, _ := strconv.Atoi(string(in))
		return dig, true
	}
	return 0, false
}

func isBackslash(in rune) bool {
	return in == backslash
}
