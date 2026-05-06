package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_cmp(t *testing.T) {
	tests := []struct {
		name string
		el1  string
		el2  string
		want int
	}{
		{
			el1:  "б",
			el2:  "а",
			want: 1,
		},
		{
			el1:  "а",
			el2:  "б",
			want: -1,
		},
		{
			el1:  "б",
			el2:  "б",
			want: 0,
		},
		{
			el1:  "б",
			el2:  "бa",
			want: -1,
		},
		{
			el1:  "бa",
			el2:  "б ",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmp(tt.el1, tt.el2)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_isMultDash(t *testing.T) {
	tests := []struct {
		name string
		word string
		want bool
	}{
		{"singleDash", "-", false},
		{"firstNotDash", "1-", false},
		{"secondNotDash", "-1", false},
		{"twoDash", "--", true},
		{"threeDash", "---", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isMultDash(tt.word)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_processWord(t *testing.T) {
	tests := []struct {
		word    string
		want    string
		wantErr bool
	}{
		{"-", "", true},
		{"---", "---", false},
		{"Нога", "нога", false},
		{"!Нога", "нога", false},
		{"-Нога-", "нога", false},
		{"-нога-", "нога", false},
		{"какой-то", "какой-то", false},
		{"какойто", "какойто", false},
		{"dog,cat", "dog,cat", false},
		{"dog...cat", "dog...cat", false},
		{"dogcat", "dogcat", false},
	}
	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			got, gotErr := processWord(tt.word)
			if tt.wantErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_getRes(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		wordList map[int][]string
		maxValue int
		want     []string
	}{
		{
			name:     "succes",
			wordList: map[int][]string{7: {"a", "b"}},
			maxValue: 7,
			want:     []string{"a", "b"},
		},
		{
			name:     "zero max value",
			wordList: map[int][]string{},
			maxValue: 0,
			want:     []string{},
		},
		{
			name: "more than 10",
			wordList: map[int][]string{
				7:  {"a", "b"},
				8:  {"c", "d", "e"},
				2:  {"f", "g", "l"},
				11: {"r", "q", "p"},
			},
			maxValue: 0,
			want:     []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRes(tt.wordList, tt.maxValue)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_createWordList(t *testing.T) {
	tests := []struct {
		name         string
		countOfWords map[string]int
		want         map[int][]string
		want2        int
	}{
		{
			name:         "empty",
			countOfWords: map[string]int{},
			want:         map[int][]string{},
			want2:        0,
		},
		{
			name:         "single word",
			countOfWords: map[string]int{"success": 1},
			want:         map[int][]string{1: {"success"}},
			want2:        1,
		},
		{
			name: "groups words by count and sorts equal frequency",
			countOfWords: map[string]int{
				"c": 2,
				"a": 2,
				"b": 3,
				"d": 1,
			},
			want: map[int][]string{
				1: {"d"},
				2: {"a", "c"},
				3: {"b"},
			},
			want2: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := createWordList(tt.countOfWords)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want2, got2)
		})
	}
}

func Test_wordCount(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		strList []string
		want    map[string]int
	}{
		{
			name:    "empty",
			strList: []string{},
			want:    map[string]int{},
		},
		{
			name:    "single word",
			strList: []string{"success"},
			want:    map[string]int{"success": 1},
		},
		{
			name:    "repeated words",
			strList: []string{"a", "b", "a", "c", "b", "a"},
			want:    map[string]int{"a": 3, "b": 2, "c": 1},
		},
		{
			name:    "words with punctuation inside",
			strList: []string{"dog,cat", "dog...cat", "dog,cat", "--"},
			want:    map[string]int{"dog,cat": 2, "dog...cat": 1, "--": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wordCount(tt.strList)
			require.Equal(t, tt.want, got)
		})
	}
}
