package main

import (
	"sort"
)

// findDictionaryLettersOccurrence calculates a map of letters within the dictionary where each letter is accompanied
// by its occurrence times
func findDictionaryLettersOccurrence(dictionary []string) map[string]int {
	dictionaryLetterOccurrence := make(map[string]int)
	for _, availableWord := range dictionary {
		for _, letter := range availableWord {
			if _, ok := dictionaryLetterOccurrence[string(letter)]; ok {
				dictionaryLetterOccurrence[string(letter)]++
			} else {
				dictionaryLetterOccurrence[string(letter)] = 1
			}
		}
	}
	return dictionaryLetterOccurrence
}

// findGoodOpener finds a word within the dictionary that contains the most occurred letters given in the
// dictionaryLetterOccurrence
func findGoodOpener(wordsLength int, dictionary []string, dictionaryLetterOccurrence map[string]int) string {
	letterOccurrenceSorted := make(PairList, len(dictionaryLetterOccurrence))
	i := 0
	for k, v := range dictionaryLetterOccurrence {
		letterOccurrenceSorted[i] = Pair{string(k), v}
		i++
	}
	sort.Sort(letterOccurrenceSorted)
	mostOccurredLetters := make(map[string]bool)
	for i := 0; i < wordsLength; i++ {
		mostOccurredLetters[letterOccurrenceSorted[len(letterOccurrenceSorted)-i-1].Key] = true
	}
	promisingWord := ""
	promisingWordScore := 0
	for _, availableWord := range dictionary {
		frequentLetterExistence := 0
		for k := range mostOccurredLetters {
			mostOccurredLetters[k] = true
		}
		for _, letter := range availableWord {
			if val, ok := mostOccurredLetters[string(letter)]; ok && val {
				mostOccurredLetters[string(letter)] = false
				frequentLetterExistence++
			}
		}
		// if availableWord contains all the most frequent letters
		if frequentLetterExistence == wordsLength {
			return availableWord
		}
		if frequentLetterExistence > promisingWordScore {
			promisingWordScore = frequentLetterExistence
			promisingWord = availableWord
		}
	}
	return promisingWord
}

type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}
func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PairList) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}
