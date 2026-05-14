package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func countWords(text string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.Fields(text) {
		word, err := processWord(word)
		if err != nil {
			continue
		}
		counts[word]++
	}
	return counts
}

func sortByFrequency(counts map[string]int) []string {
	words := make([]string, 0, len(counts))
	for word := range counts {
		words = append(words, word)
	}

	sort.Slice(words, frequencyLess(words, counts))

	return words
}

func frequencyLess(words []string, counts map[string]int) func(i, j int) bool {
	return func(i, j int) bool {
		if counts[words[i]] == counts[words[j]] {
			return words[i] < words[j]
		}
		return counts[words[i]] > counts[words[j]]
	}
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
