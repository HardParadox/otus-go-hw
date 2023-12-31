package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordInArray struct {
	value   string
	counter int32
}

func Top10(value string) []string {
	if len(value) == 0 {
		return []string{}
	}

	splitedString := strings.Fields(value)

	dictionaryWithCountOfWord := countWords(splitedString)

	arrayWords := mapToArray(dictionaryWithCountOfWord)

	sort.Slice(arrayWords, func(i, j int) bool {
		if arrayWords[i].counter == arrayWords[j].counter {
			return arrayWords[i].value < arrayWords[j].value
		}

		return arrayWords[i].counter > arrayWords[j].counter
	})

	var shortArrayWords []wordInArray

	if len(arrayWords) > 10 {
		shortArrayWords = arrayWords[:10]
	} else {
		shortArrayWords = arrayWords
	}

	arrayOfTop10Words := make([]string, len(shortArrayWords))

	for index, word := range shortArrayWords {
		arrayOfTop10Words[index] = word.value
	}

	return arrayOfTop10Words
}

func countWords(arrayOfString []string) map[string]int32 {
	dictionaryWithCountOfWord := make(map[string]int32)

	for _, word := range arrayOfString {
		dictionaryWithCountOfWord[word]++
	}

	return dictionaryWithCountOfWord
}

func mapToArray(initMap map[string]int32) []wordInArray {
	array := []wordInArray{}

	for word, counter := range initMap {
		array = append(array, wordInArray{
			value:   word,
			counter: counter,
		})
	}

	return array
}
