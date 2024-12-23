package hw03frequencyanalysis

import (
	"slices"
	"strings"
)

func Top10(str string) []string {
	if str == "" {
		return []string{}
	}

	words := make(map[string]int)

	for _, word := range strings.Fields(str) {
		words[word]++
	}

	outputLen := 10
	if len(words) < outputLen {
		outputLen = len(words)
	}
	output := make([]string, outputLen)

	for i := 0; i < outputLen; i++ {
		for word := range words {
			if slices.Contains(output, word) {
				continue
			}

			if output[i] == "" {
				output[i] = word

				continue
			}

			if words[word] < words[output[i]] {
				continue
			}

			if words[word] == words[output[i]] && strings.Compare(output[i], word) <= 0 {
				continue
			}

			output[i] = word
		}
	}

	return output
}
