package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

type parseAction int

const (
	actionContinue parseAction = iota
	actionStop
)

const backslash rune = '\\'

var ErrInvalidString = errors.New("invalid string")

type parser struct {
	runes       []rune
	builder     strings.Builder
	escapedPrev bool
}

func newParser(in string) *parser {
	return &parser{
		runes:   []rune(in),
		builder: strings.Builder{},
	}
}

func Unpack(in string) (string, error) {
	parser := newParser(in)

	if len(parser.runes) == 0 {
		return "", nil
	}

	if _, ok := isDigit(parser.runes[0]); ok {
		return "", ErrInvalidString
	}

	if len(parser.runes) == 1 && isBackslash(parser.runes[0]) {
		return "", ErrInvalidString
	}

	if len(parser.runes) == 1 {
		parser.builder.WriteString(string(parser.runes[0]))
		return parser.builder.String(), nil
	}

	for i := 1; i < len(parser.runes); i++ {
		curr := parser.runes[i]
		_, currDigOk := isDigit(curr)

		switch {
		case currDigOk:
			action, err := parser.handleDig(i)
			if err != nil {
				return "", err
			}
			if action == actionContinue {
				continue
			}
			return parser.builder.String(), nil

		case isBackslash(curr):
			if err := parser.handleBackslash(i); err != nil {
				return "", err
			}
		default:
			if err := parser.handleOther(i); err != nil {
				return "", err
			}
		}
	}
	return parser.builder.String(), nil
}

func (p *parser) handleDig(i int) (parseAction, error) {
	curr := p.runes[i]
	prev := p.runes[i-1]
	isLast := i == len(p.runes)-1
	dig, _ := isDigit(curr)
	_, prevDigOk := isDigit(prev)

	switch {
	case prevDigOk:
		if p.escapedPrev {
			str := strings.Repeat(string(prev), dig)
			p.builder.WriteString(str)
			p.escapedPrev = false
			return actionContinue, nil
		}
		return actionStop, ErrInvalidString
	case isBackslash(prev):
		if p.escapedPrev {
			str := strings.Repeat(string(prev), dig)
			p.builder.WriteString(str)
			p.escapedPrev = false
			return actionContinue, nil
		}
		if isLast {
			p.builder.WriteString(string(curr))
			return actionStop, nil
		}

		p.escapedPrev = true
		return actionContinue, nil
	default:
		str := strings.Repeat(string(prev), dig)
		p.builder.WriteString(str)
		if isLast {
			return actionStop, nil
		}
		return actionContinue, nil
	}
}

func (p *parser) handleBackslash(i int) error {
	curr := p.runes[i]
	prev := p.runes[i-1]
	isLast := i == len(p.runes)-1
	_, prevDigOk := isDigit(prev)

	if p.escapedPrev {
		str := string(prev)
		p.builder.WriteString(str)
		p.escapedPrev = false
		if isLast {
			return ErrInvalidString
		}
		return nil
	}

	switch {
	case isBackslash(prev):
		p.escapedPrev = true
		if isLast {
			str := string(curr)
			p.builder.WriteString(str)
		}
		return nil
	case prevDigOk:
		if isLast {
			return ErrInvalidString
		}
		return nil
	default:
		if isLast {
			return ErrInvalidString
		}
		str := string(prev)
		p.builder.WriteString(str)

		return nil
	}
}

func (p *parser) handleOther(i int) error {
	curr := p.runes[i]
	prev := p.runes[i-1]
	isLast := i == len(p.runes)-1
	_, prevDigOk := isDigit(prev)

	if p.escapedPrev {
		str := string(prev)
		p.builder.WriteString(str)
		if isLast {
			p.builder.WriteString(string(curr))
		}
		p.escapedPrev = false
		return nil
	}

	switch {
	case isBackslash(prev):
		return ErrInvalidString
	case prevDigOk:
		if isLast {
			p.builder.WriteString(string(curr))
		}
		return nil
	default:
		p.builder.WriteString(string(prev))

		if isLast {
			p.builder.WriteString(string(curr))
		}
		return nil
	}
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
