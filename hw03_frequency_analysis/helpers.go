package hw03frequencyanalysis

import (
	"slices"
	"strings"
	"unicode"
)

func wordCount(strList []string) map[string]int {
	countOfWords := make(map[string]int)
	var processedWord string
	for _, word := range strList {
		var err error
		processedWord, err = processWord(word)
		if err != nil {
			continue
		}
		countOfWords[processedWord]++
	}
	return countOfWords
}

func createWordList(countOfWords map[string]int) (wordList map[int][]string, frequency []int) {
	wordList = make(map[int][]string, 0)
	for key, value := range countOfWords {
		if _, ok := wordList[value]; !ok {
			wordList[value] = []string{}
		}
		sl := wordList[value]
		sl = append(sl, key)

		wordList[value] = sl
		frequency = append(frequency, value)
	}
	// Если слова имеют одинаковую частоту, то должны быть отсортированы **лексикографически**.
	for _, v := range wordList {
		slices.SortFunc(v, cmp)
	}

	frequency = uniqueInts(frequency)

	return wordList, frequency
}

func getRes(wordList map[int][]string, frequency []int) []string {
	if len(frequency) == 0 {
		return []string{}
	}

	slices.SortFunc(frequency, func(el1, el2 int) int {
		switch {
		case el2 > el1:
			return 1
		case el1 == el2:
			return 0
		default:
			return -1
		}
	})

	res := make([]string, 0)
	// * Если есть более 10 самых частотых слов (например 15 разных слов встречаются ровно 133 раза,
	// остальные < 100), то следует вернуть 10 лексикографически первых слов.
	for _, v := range frequency {
		for ind := range wordList[v] {
			if len(res) == 10 {
				break
			}
			el := wordList[v][ind]
			res = append(res, el)
		}
	}
	return res
}

func processWord(word string) (string, error) {
	if word == "" {
		return "", ErrWordIsEmpty
	}
	if word == "-" {
		return "", ErrSingleDashIsNotWord
	}

	cleaned := strings.TrimFunc(word,
		func(r rune) bool {
			return unicode.IsPunct(r) && r != '-'
		})

	if cleaned == "" {
		return "", ErrEmptyWordAfterTrim
	}

	ok, err := isMultDash(cleaned)
	if err != nil {
		return "", err
	}
	if ok {
		return cleaned, nil
	}

	cleaned = strings.TrimFunc(cleaned, unicode.IsPunct)

	if cleaned == "" {
		return "", ErrEmptyWordAfterTrim
	}

	return strings.ToLower(cleaned), nil
}

func isMultDash(word string) (bool, error) {
	wordInRune := []rune(word)
	if len(wordInRune) == 0 {
		return false, ErrWordIsEmpty
	}
	firstRune := wordInRune[0]
	if firstRune != '-' {
		return false, nil
	}
	if len(wordInRune) < 2 {
		return false, nil
	}

	for _, v := range wordInRune[1:] {
		if v != '-' {
			return false, nil
		}
	}

	return true, nil
}

func cmp(el1, el2 string) int {
	return strings.Compare(el1, el2)
}

func uniqueInts(nums []int) []int {
	seen := make(map[int]struct{}, len(nums))
	result := make([]int, 0, len(nums))

	for _, n := range nums {
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		result = append(result, n)
	}

	return result
}
