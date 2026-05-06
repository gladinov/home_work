package hw03frequencyanalysis

import (
	"errors"
	"strings"
)

type parseAction int

var ErrSingleDashIsNotWord = errors.New("single dash is not a valid word")

const (
	actionContinue parseAction = iota
	actionStop
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

func Top10(text string) []string {
	// * Словом считается набор символов, разделенных пробельными символами.
	// Пробельные симоволы это unicode.IsSpace
	strList := strings.Fields(text)
	countOfWords := wordCount(strList)
	if len(countOfWords) == 0 {
		return []string{}
	}

	wordList, maxValue := createWordList(countOfWords)

	res := getRes(wordList, maxValue)

	return res
}
