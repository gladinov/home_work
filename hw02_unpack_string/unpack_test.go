package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "🙃2ф4,1", expected: "🙃🙃фффф,"},
		{input: "!!!!!", expected: "!!!!!"},
		{input: "!5", expected: "!!!!!"},
		{input: "a", expected: "a"},
		{input: "日本語", expected: "日本語"},
		{input: "日2本3語1", expected: "日日本本本語"},
		{input: "日0本0語0", expected: ""},
		{input: "aaф0b", expected: "aab"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "1", "日2本3語12"}
	// tc :=
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func Test_isDig(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in    rune
		want  int
		want2 bool
	}{
		{"success", '5', 5, true},
		{"false", 'e', 0, false},
		{"false", 'a', 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := isDigit(tt.in)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want2, got2)
		})
	}
}
