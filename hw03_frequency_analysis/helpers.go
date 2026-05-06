package hw03frequencyanalysis

import (
	"slices"
	"strings"
	"unicode"
)

func wordCount(strList []string) map[string]int {
	countOfWords := make(map[string]int)
	var proccessWord string
	var err error
	for _, word := range strList {
		if taskWithAsteriskIsCompleted {
			proccessWord, err = processWord(word)
			if err != nil {
				continue
			}
			countOfWords[proccessWord]++
		} else {
			countOfWords[word]++
		}
	}
	return countOfWords
}

func createWordList(countOfWords map[string]int) (map[int][]string, int) {
	var maxValue int
	// TODO: Исправить ошибку с длиной массива
	wordList := make(map[int][]string, 0)
	for key, value := range countOfWords {
		if _, ok := wordList[value]; !ok {
			wordList[value] = []string{}
		}
		sl := wordList[value]
		sl = append(sl, key)
		// Если слова имеют одинаковую частоту, то должны быть отсортированы **лексикографически**.
		slices.SortFunc(sl, cmp)
		wordList[value] = sl
		maxValue = max(maxValue, value)
	}

	return wordList, maxValue
}

func getRes(wordList map[int][]string, maxValue int) []string {
	if maxValue == 0 {
		return []string{}
	}
	res := make([]string, 0)
	// * Если есть более 10 самых частотых слов (например 15 разных слов встречаются ровно 133 раза,
	// остальные < 100), то следует вернуть 10 лексикографически первых слов.
	for i := maxValue; condition(res, i); i-- {
		for ind := range wordList[i] {
			if len(res) == 10 {
				break
			}
			el := wordList[i][ind]
			res = append(res, el)
		}
	}
	return res
}

func processWord(word string) (string, error) {
	// * "-" словом не является
	if word == "-" {
		return "", ErrSingleDashIsNotWord
	}
	// * "-------" это слово
	if isMultDash(word) {
		return word, nil
	}

	//  "нога!", "нога", "нога," и " 'нога' " - это одинаковые слова;
	wordWithoutMarks := strings.TrimFunc(word, trimFunc)
	// * "Нога" и "нога" - это одинаковые слова,
	lowerCaseWord := strings.ToLower(wordWithoutMarks)

	return lowerCaseWord, nil
}

func isMultDash(word string) bool {
	wordInRune := []rune(word)
	firstRune := wordInRune[0]
	if firstRune != '-' {
		return false
	}
	if len(wordInRune) < 2 {
		return false
	}

	for _, v := range wordInRune[1:] {
		if v != '-' {
			return false
		}
	}

	return true
}

func trimFunc(in rune) bool {
	if unicode.IsLetter(in) {
		return false
	}
	return true
}

func cmp(el1, el2 string) int {
	return strings.Compare(el1, el2)
}

func condition(res []string, i int) bool {
	if len(res) < 10 && i > 0 {
		return true
	}
	return false
}
