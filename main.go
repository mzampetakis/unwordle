package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var NotExists = 'b'
var Exists = 'y'
var Correct = 'g'

var WordsLength int
var AvailableWords []string
var OpenerWord string
var TotalTries = 6
var ShowInfo bool
var letterOccurrence map[string]int

func init() {
	dictionaryFilePath := flag.String("dictionary", "", "The source file to read the available dictionary.")
	opener := flag.String("opener", "", "The opener word to use. (Optional)")
	tries := flag.Int("tries", 6, "The opener word to use. (Optional)")
	showInfo := flag.Bool("info", false, "Show info about each guess. (Optional)")
	flag.Parse()
	if tries != nil {
		if *tries > 0 {
			TotalTries = *tries
		}
	}
	if showInfo != nil {
		ShowInfo = *showInfo
	}
	if len(*opener) > 0 {
		OpenerWord = *opener
	}
	if len(*dictionaryFilePath) == 0 {
		fmt.Println("Provide the source argument.")
		fmt.Println("Usage: `./unwordle --source=dictionary_source_file_path`")
		os.Exit(1)
	}
	sourceFile, err := os.Open(*dictionaryFilePath)
	if err != nil {
		fmt.Println("Provide an existing source file.")
		fmt.Println("Usage: `./unwordle --source=dictionary_source_file_path`")
		os.Exit(1)
	}
	defer sourceFile.Close()
	scanner := bufio.NewScanner(sourceFile)
	AvailableWords = []string{}
	for scanner.Scan() {
		if WordsLength == 0 {
			WordsLength = utf8.RuneCountInString(scanner.Text())
		}
		if utf8.RuneCountInString(scanner.Text()) == WordsLength {
			AvailableWords = append(AvailableWords, strings.ToUpper(scanner.Text()))
		}
	}
	if len(AvailableWords) == 0 {
		fmt.Println("No valid word found in the file.")
		os.Exit(1)
	}
}

// lettersRules contains a map of all 26 letters and for each letter rules about their occurrence are stored
var lettersRules map[string]string

// validWords holds a map of words of the given dictionary.
// When a word is not possible to be chosen it is being deleted.
var validWords map[string]bool

func main() {
	letterOccurrence = findDictionaryLettersOccurrence(AvailableWords)
	if len(OpenerWord) != WordsLength {
		OpenerWord = findGoodOpener(WordsLength, AvailableWords, letterOccurrence)
	}
	validWords = make(map[string]bool)
	for _, word := range AvailableWords {
		validWords[word] = true
	}
	lettersRules = make(map[string]string)
	for letter := range letterOccurrence {
		lettersRules[letter] = ""
	}

	currentWord := OpenerWord
	currentWordScore := 0
	for try := 0; try < TotalTries; try++ {
		tryInfo := ""
		if ShowInfo {
			tryInfo = fmt.Sprintf("| Possibility: 1/%d | Score: %d", len(validWords), currentWordScore)
		}
		fmt.Printf("Try #%d: \t\t%s %s\n", try+1, currentWord, tryInfo)
		response := ""
		for !isValidResponse(response) {
			fmt.Printf("Response (b|y|g): \t")
			response = readResponse()
		}
		if !strings.Contains(response, string(NotExists)) && !strings.Contains(response, string(Exists)) {
			fmt.Println("Hooray! :-)")
			os.Exit(0)
		}
		index := 0
		for _, currentLetter := range currentWord {
			if !strings.Contains(lettersRules[string(currentLetter)],
				string(response[index])+strconv.Itoa(index)+"|") {
				lettersRules[string(currentLetter)] += string(response[index]) + strconv.Itoa(index) + "|"
			}
			index++
		}
		excludeWords()
		if len(validWords) == 0 {
			fmt.Println("No solution found with the given criteria.")
			fmt.Println("Sorry. :-(")
			os.Exit(0)
		} else if len(validWords) == 1 {
			fmt.Println("Found Solution: \t" + currentWord)
			fmt.Println("Hooray! :-)")
			os.Exit(0)
		}
		currentWord, currentWordScore = findGoodWord()
		delete(validWords, currentWord)
	}
	fmt.Println("No solution found.")
	fmt.Println("Sorry. :-(")
}

func readResponse() string {
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = response[:len(response)-1]
	return response
}

// excludeWords checks all the words within the validWords map and remove impossible words base on the rules stored
// in the lettersRules map
func excludeWords() {
	for wordToCheck := range validWords {
		idx := -1
		for _, wordLetter := range wordToCheck {
			idx++
			// letter has no rule
			if len(lettersRules[string(wordLetter)]) == 0 {
				continue
			}
			// for NotExists(b) response
			for _, rule := range strings.Split(lettersRules[string(wordLetter)], "|") {
				if len(rule) == 0 {
					continue
				}
				if rule[:1] == string(NotExists) && !strings.Contains(lettersRules[string(wordLetter)],
					string(Correct)) && !strings.Contains(lettersRules[string(wordLetter)], string(Exists)) {
					delete(validWords, wordToCheck)
					break
				}
				if rule[:1] == string(NotExists) && rule[1:] == strconv.Itoa(idx) {
					delete(validWords, wordToCheck)
					break
				}
				if rule[:1] == string(Exists) && rule[1:] == strconv.Itoa(idx) {
					delete(validWords, wordToCheck)
					break
				}
				// Exclude words that contain letters that are in wrong place (yellow)
				if rule[:1] == string(NotExists) && rule[1:] == strconv.Itoa(idx) {
					delete(validWords, wordToCheck)
					break
				}
			}
		}
	}
}

func findGoodWord() (string, int) {
	goodFit := ""
	goodFitScore := -1
	currentScore := 0
	for validWord := range validWords {
		currentScore = 0
		pos := 0
		for _, wordLetter := range validWord {
			// Estimate a score per word
			for _, rule := range strings.Split(lettersRules[string(wordLetter)], "|") {
				if len(rule) == 0 {
					continue
				}
				if rule[:1] == string(Exists) {
					currentScore++
				}
				if rule[:1] == string(Correct) && rule[1:] == strconv.Itoa(pos) {
					currentScore += WordsLength
				}
			}
			pos++
		}
		if currentScore > goodFitScore {
			goodFit = validWord
			goodFitScore = currentScore
		}
	}
	return goodFit, goodFitScore
}

// isValidResponse checks id the response string has the correct length
// and contains only the letters "b", "y" and "g".
func isValidResponse(response string) bool {
	if len(response) != WordsLength {
		return false
	}
	for _, val := range response {
		if val != NotExists && val != Exists && val != Correct {
			return false
		}
	}
	return true
}
