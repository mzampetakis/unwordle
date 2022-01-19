package main

import (
	"sort"
)

// findGoodOpener calculates the occurrence of each letter within the AvailableWords array
// and using the most occurred letters finds a word that contains most of those letters.
func findGoodOpener() string {
	letterOccurrence := make(map[rune]int)
	for _, availableWord := range AvailableWords {
		for _, letter := range availableWord {
			if _, ok := letterOccurrence[letter]; ok {
				letterOccurrence[letter]++
			} else {
				letterOccurrence[letter] = 1
			}
		}
	}
	letterOccurrenceSorted := make(PairList, len(letterOccurrence))

	i := 0
	for k, v := range letterOccurrence {
		letterOccurrenceSorted[i] = Pair{k, v}
		i++
	}
	sort.Sort(letterOccurrenceSorted)
	mostOccurredLetters := make(map[rune]bool)
	for i := 0; i < WordsLength; i++ {
		mostOccurredLetters[letterOccurrenceSorted[len(letterOccurrenceSorted)-i-1].Key] = true
	}
	promisingWord := ""
	promisingWordScore := 0
	for _, availableWord := range AvailableWords {
		letterExistence := 0
		for k := range mostOccurredLetters {
			mostOccurredLetters[k] = true
		}
		for _, letter := range availableWord {
			if val, ok := mostOccurredLetters[letter]; ok && val {
				mostOccurredLetters[letter] = false
				letterExistence++
			}
		}
		if letterExistence == WordsLength {
			return availableWord
		}
		if letterExistence > promisingWordScore {
			promisingWordScore = letterExistence
			promisingWord = availableWord
		}
	}
	return promisingWord
}

type Pair struct {
	Key   rune
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
