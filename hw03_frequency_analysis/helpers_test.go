package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isMultDash(t *testing.T) {
	tests := []struct {
		name string
		word string
		want bool
		err  error
	}{
		{"singleDash", "-", false, nil},
		{"firstNotDash", "1-", false, nil},
		{"secondNotDash", "-1", false, nil},
		{"twoDash", "--", true, nil},
		{"threeDash", "---", true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := isMultDash(tt.word)
			if tt.err != nil {
				require.ErrorIs(t, gotErr, tt.err)
				return
			}
			require.NoError(t, gotErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_processWord(t *testing.T) {
	tests := []struct {
		word string
		want string
		err  error
	}{
		{"-", "", ErrSingleDashIsNotWord},
		{"---", "---", nil},
		{"Нога", "нога", nil},
		{"!Нога", "нога", nil},
		{"-Нога-", "нога", nil},
		{"-нога-", "нога", nil},
		{"какой-то", "какой-то", nil},
		{"какойто", "какойто", nil},
		{"dog,cat", "dog,cat", nil},
		{"dog...cat", "dog...cat", nil},
		{"dogcat", "dogcat", nil},
		{"!!!", "", ErrEmptyWordAfterTrim},
		{"", "", ErrWordIsEmpty},
		{"-,-", "", ErrEmptyWordAfterTrim},
		{"-олег-", "олег", nil},
		{"!--!", "--", nil},
		{"-------,", "-------", nil},
		{"'-------'", "-------", nil},
		{`"-------"`, "-------", nil},
		{"---!", "---", nil},
		{"!-!", "", ErrEmptyWordAfterTrim},
	}
	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			got, gotErr := processWord(tt.word)
			if tt.err != nil {
				require.ErrorIs(t, gotErr, tt.err)
				return
			}
			require.NoError(t, gotErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_sortByFrequency(t *testing.T) {
	tests := []struct {
		name   string
		counts map[string]int
		want   []string
	}{
		{
			name:   "empty",
			counts: map[string]int{},
			want:   []string{},
		},
		{
			name:   "single frequency",
			counts: map[string]int{"b": 7, "a": 7},
			want:   []string{"a", "b"},
		},
		{
			name: "sorts by frequency descending and word ascending",
			counts: map[string]int{
				"a": 7,
				"b": 7,
				"c": 8,
				"d": 8,
				"e": 8,
				"f": 2,
				"g": 2,
				"l": 2,
				"p": 11,
				"q": 11,
				"r": 11,
			},
			want: []string{"p", "q", "r", "c", "d", "e", "a", "b", "f", "g", "l"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortByFrequency(tt.counts)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_countWords(t *testing.T) {
	tests := []struct {
		name string
		text string
		want map[string]int
	}{
		{
			name: "empty",
			text: "",
			want: map[string]int{},
		},
		{
			name: "single word",
			text: "success",
			want: map[string]int{"success": 1},
		},
		{
			name: "repeated words",
			text: "a b a c b a",
			want: map[string]int{"a": 3, "b": 2, "c": 1},
		},
		{
			name: "normalizes words",
			text: "Нога нога, - dog,cat dog,cat --",
			want: map[string]int{"нога": 2, "dog,cat": 2, "--": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countWords(tt.text)
			require.Equal(t, tt.want, got)
		})
	}
}
