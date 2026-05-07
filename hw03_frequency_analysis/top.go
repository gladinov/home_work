package hw03frequencyanalysis

import (
	"strings"
)

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
