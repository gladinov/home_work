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
		// // uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `qwe\2e`, expected: `qwe2e`},
		{input: `a4q`, expected: `aaaaq`},
		{input: `aaq`, expected: `aaq`},
		{input: `aa1\1`, expected: `aa1`},
		{input: `\1\2\3\4`, expected: `1234`},
		{input: `\\`, expected: `\`},
		{input: `\\\3`, expected: `\3`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			t.Log(result, err)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{
		"3abc",
		"45",
		"aaa10b",
		"1",
		"日2本3語12",
		`qwe\\\3\`,
		`e\\\3\`,
		`\\\3\`,
		`q3\`,
		`q\q`,
		`\`,
		`a\`,
	}
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

func Test_isBackslash(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in   rune
		want bool
	}{
		{"success", 92, true},
		{"false", 97, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(string(tt.in))
			got := isBackslash(tt.in)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_parser_handleDig(t *testing.T) {
	tests := []struct {
		in          string
		i           int
		escapedPrev bool
		want        parseAction
		wantErr     bool
	}{
		{"a5", 1, false, actionStop, false},
		{"65", 1, false, actionStop, true},
		{`\5`, 1, false, actionStop, false},
		{`\56`, 2, true, actionContinue, false},
		{`\\6`, 2, true, actionContinue, false},
		{`\\6ф`, 2, true, actionContinue, false},
		{`\6ф`, 1, false, actionContinue, false},
		{`f6ф`, 1, false, actionContinue, false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			p := newParser(tt.in)
			p.escapedPrev = tt.escapedPrev
			got, gotErr := p.handleDig(tt.i)
			if tt.wantErr {
				require.Error(t, gotErr)
				require.EqualError(t, gotErr, ErrInvalidString.Error())
				return
			}
			require.NoError(t, gotErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_parser_handleDig_parcer_change(t *testing.T) {
	tests := []struct {
		in                       string
		i                        int
		escapedPrev              bool
		wantEscapedPrevAfterFunc bool
	}{
		{"a5", 1, false, false},
		{"65", 1, false, false},
		{`\5`, 1, false, false},
		{`\56`, 2, true, false},
		{`\\6`, 2, true, false},
		{`\\6ф`, 2, true, false},
		{`\6ф`, 1, false, true},
		{`f6ф`, 1, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			p := newParser(tt.in)
			p.escapedPrev = tt.escapedPrev
			_, _ = p.handleDig(tt.i)

			require.Equal(t, tt.wantEscapedPrevAfterFunc, p.escapedPrev)
		})
	}
}

func Test_parser_handleBackslash(t *testing.T) {
	tests := []struct {
		in          string
		i           int
		escapedPrev bool
		wantErr     bool
	}{
		{`\\a`, 1, false, false},
		{`\\\4`, 2, true, false},
		{`\\\4\`, 4, false, true},
		{`\\\`, 2, true, true},
		{`\\\3`, 2, true, false},
		{`\\3`, 1, false, false},
		{`\\3\`, 1, false, false},
		{`\\\\`, 2, true, false},
		{`\\\\`, 3, false, false},
		{`a4\`, 2, false, true},
		{`a4\1`, 2, false, false},
		{`aa\1`, 2, false, false},
		{`aa\`, 2, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			p := newParser(tt.in)
			p.escapedPrev = tt.escapedPrev
			gotErr := p.handleBackslash(tt.i)
			if tt.wantErr {
				require.Error(t, gotErr)
				require.EqualError(t, gotErr, ErrInvalidString.Error())
				return
			}
			require.NoError(t, gotErr)
		})
	}
}

func Test_parser_handleBackslash_parser_change(t *testing.T) {
	tests := []struct {
		in                       string
		i                        int
		escapedPrev              bool
		wantEscapedPrevAfterFunc bool
		wantBuilder              string
	}{
		{
			in:                       `\\a`,
			i:                        1,
			escapedPrev:              false,
			wantEscapedPrevAfterFunc: true,
			wantBuilder:              ``,
		},
		{
			in:                       `\\\3`,
			i:                        2,
			escapedPrev:              true,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              `\`,
		},
		{
			in:                       `\\\\`,
			i:                        2,
			escapedPrev:              true,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              `\`,
		},
		{
			in:                       `\\\\`,
			i:                        3,
			escapedPrev:              false,
			wantEscapedPrevAfterFunc: true,
			wantBuilder:              `\`,
		},
		{
			in:                       `a4\1`,
			i:                        2,
			escapedPrev:              false,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              ``,
		},
		{
			in:                       `aa\1`,
			i:                        2,
			escapedPrev:              false,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              `a`,
		},
		{
			in:                       `\\`,
			i:                        1,
			escapedPrev:              false,
			wantEscapedPrevAfterFunc: true,
			wantBuilder:              `\`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			p := newParser(tt.in)
			p.escapedPrev = tt.escapedPrev

			_ = p.handleBackslash(tt.i)

			require.Equal(t, tt.wantEscapedPrevAfterFunc, p.escapedPrev)
			require.Equal(t, tt.wantBuilder, p.builder.String())
		})
	}
}

func Test_parser_handleOther(t *testing.T) {
	tests := []struct {
		in          string
		i           int
		escapedPrev bool
		wantErr     bool
	}{
		{`ab`, 1, false, false},  // обычный символ после обычного
		{`a5b`, 2, false, false}, // символ после цифры
		{`a\q`, 2, false, true},  // невалидный escape
		{`\\q`, 2, true, false},  // escapedPrev=true, после него обычный символ
		{`\\a`, 1, false, true},  // после backslash идет не цифра и не backslash
		{`12a`, 2, false, false}, // после цифры обычный символ
		{`🙂a`, 1, false, false},  // unicode
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			p := newParser(tt.in)
			p.escapedPrev = tt.escapedPrev

			gotErr := p.handleOther(tt.i)
			if tt.wantErr {
				require.Error(t, gotErr)
				require.EqualError(t, gotErr, ErrInvalidString.Error())
				return
			}

			require.NoError(t, gotErr)
		})
	}
}

func Test_parser_handleOther_parser_change(t *testing.T) {
	tests := []struct {
		in                       string
		i                        int
		escapedPrev              bool
		wantEscapedPrevAfterFunc bool
		wantBuilder              string
	}{
		{
			in:                       `ab`,
			i:                        1,
			escapedPrev:              false,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              `ab`,
		},
		{
			in:                       `abc`,
			i:                        1,
			escapedPrev:              false,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              `a`,
		},
		{
			in:                       `a5b`,
			i:                        2,
			escapedPrev:              false,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              `b`,
		},
		{
			in:                       `\\q`,
			i:                        2,
			escapedPrev:              true,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              `\q`,
		},
		{
			in:                       `\\qa`,
			i:                        2,
			escapedPrev:              true,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              `\`,
		},
		{
			in:                       `🙂a`,
			i:                        1,
			escapedPrev:              false,
			wantEscapedPrevAfterFunc: false,
			wantBuilder:              `🙂a`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			p := newParser(tt.in)
			p.escapedPrev = tt.escapedPrev

			_ = p.handleOther(tt.i)

			require.Equal(t, tt.wantEscapedPrevAfterFunc, p.escapedPrev)
			require.Equal(t, tt.wantBuilder, p.builder.String())
		})
	}
}
