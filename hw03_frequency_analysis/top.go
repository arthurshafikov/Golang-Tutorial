package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordList map[int][]string

func Top10(text string) (result []string) {

	wordsFrequencies := map[string]int{}
	for _, word := range strings.Fields(text) {
		wordsFrequencies[word]++
	}

	if len(wordsFrequencies) < 1 {
		return nil
	}

	wordList := wordList{}
	for word, freq := range wordsFrequencies {
		wordList[freq] = append(wordList[freq], word)
	}

	frequencies := []int{}
	for freq := range wordList {
		frequencies = append(frequencies, freq)
	}
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i] > frequencies[j]
	})

	for _, freq := range frequencies {
		sort.Strings(wordList[freq])
		result = append(result, wordList[freq]...)
	}

	if len(result) > 10 {
		result = result[:10]
	}

	return result
}
