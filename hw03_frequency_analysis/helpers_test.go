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
		{"!-!", "--", ErrEmptyWordAfterTrim},
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

func Test_getRes(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		wordList  map[int][]string
		frequency []int
		want      []string
	}{
		{
			name:      "succes",
			wordList:  map[int][]string{7: {"a", "b"}},
			frequency: []int{7},
			want:      []string{"a", "b"},
		},
		{
			name:      "zero max value",
			wordList:  map[int][]string{},
			frequency: []int{},
			want:      []string{},
		},
		{
			name: "more than 10",
			wordList: map[int][]string{
				7:  {"a", "b"},
				8:  {"c", "d", "e"},
				2:  {"f", "g", "l"},
				11: {"r", "q", "p"},
			},
			frequency: []int{11, 8, 7, 2},
			want:      []string{"r", "q", "p", "c", "d", "e", "a", "b", "f", "g"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRes(tt.wordList, tt.frequency)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_createWordList(t *testing.T) {
	tests := []struct {
		name         string
		countOfWords map[string]int
		want         map[int][]string
		want2        []int
	}{
		{
			name:         "empty",
			countOfWords: map[string]int{},
			want:         map[int][]string{},
			want2:        nil,
		},
		{
			name:         "single word",
			countOfWords: map[string]int{"success": 1},
			want:         map[int][]string{1: {"success"}},
			want2:        []int{1},
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
			want2: []int{2, 3, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := createWordList(tt.countOfWords)
			require.Equal(t, tt.want, got)
			require.ElementsMatch(t, tt.want2, got2)
		})
	}
}

func Test_wordCount(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		strList []string
		want    map[string]int
		err     error
	}{
		{
			name:    "empty",
			strList: []string{},
			want:    nil,
			err:     ErrStrListIsEmpty,
		},
		{
			name:    "single word",
			strList: []string{"success"},
			want:    map[string]int{"success": 1},
			err:     nil,
		},
		{
			name:    "repeated words",
			strList: []string{"a", "b", "a", "c", "b", "a"},
			want:    map[string]int{"a": 3, "b": 2, "c": 1},
			err:     nil,
		},
		{
			name:    "words with punctuation inside",
			strList: []string{"dog,cat", "dog...cat", "dog,cat", "--"},
			want:    map[string]int{"dog,cat": 2, "dog...cat": 1, "--": 1},
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := wordCount(tt.strList)
			if tt.err != nil {
				require.ErrorIs(t, gotErr, tt.err)
				return
			}
			require.NoError(t, gotErr)
			require.Equal(t, tt.want, got)
		})
	}
}
