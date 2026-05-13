package hw03frequencyanalysis

func Top10(text string) []string {
	counts := countWords(text)
	words := sortByFrequency(counts)

	if len(words) > 10 {
		words = words[:10]
	}

	return words
}
