package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var NotExists = 'b'
var Exists = 'y'
var Correct = 'g'

var WordsLength int
var AvailableWords []string
var OpenerWord string
var TotalTries = 6
var ShowInfo bool

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
			WordsLength = len(scanner.Text())
		}
		if len(scanner.Text()) == WordsLength {
			AvailableWords = append(AvailableWords, scanner.Text())
		}
	}
	if len(AvailableWords) == 0 {
		fmt.Println("No valid word found in the file.")
		os.Exit(1)
	}
}

var lettersRules map[string]string
var validWords map[string]bool

func main() {
	if len(OpenerWord) != WordsLength {
		OpenerWord = findGoodOpener()
	}
	validWords = make(map[string]bool)
	for _, word := range AvailableWords {
		validWords[word] = true
	}
	lettersRules = make(map[string]string)
	var alphabet = "a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z"
	for _, letter := range strings.Split(alphabet, ",") {
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
		totalCorrectLetters := 0
		for k, responseVal := range response {
			switch responseVal {
			case NotExists:
				if len(lettersRules[string(currentWord[k])]) == 0 {
					lettersRules[string(currentWord[k])] = string(NotExists)
				}
			case Exists:
				if len(lettersRules[string(currentWord[k])]) == 0 {
					lettersRules[string(currentWord[k])] = string(Exists) + strconv.Itoa(k)
				}
			case Correct:
				totalCorrectLetters++
				if len(lettersRules[string(currentWord[k])]) == 0 || lettersRules[string(currentWord[k])] == string(
					Exists) {
					lettersRules[string(currentWord[k])] = strconv.Itoa(k)
				} else {
					lettersRules[string(currentWord[k])] = lettersRules[string(currentWord[k])] + "|" + strconv.Itoa(k)
				}
			}
		}
		if totalCorrectLetters == WordsLength {
			fmt.Println("Hooray! :-)")
			os.Exit(0)
		}
		currentWord, currentWordScore = findGoodWord()
		if len(validWords) == 1 {
			fmt.Println("Found Solution: \t" + currentWord)
			fmt.Println("Hooray! :-)")
			os.Exit(0)
		}
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

func findGoodWord() (string, int) {
	goodFit := ""
	goodFitScore := -1
	//fmt.Printf("%+v\n", lettersRules)
	for validWord := range validWords {
		currentScore := 0
		for pos, wordLetter := range validWord {
			// Exclude words that contain letters that don't exist (black)
			if lettersRules[string(wordLetter)] == string(NotExists) {
				delete(validWords, validWord)
				break
			}
			// Exclude words that contain letters that are in wrong place (yellow)
			if strings.Contains(lettersRules[string(wordLetter)], string(Exists)) {
				position := lettersRules[string(wordLetter)][1:]
				if position == strconv.Itoa(pos) {
					delete(validWords, validWord)
					break
				}
			}
			// Estimate a score per word
			if len(lettersRules[string(wordLetter)]) > 0 {
				if lettersRules[string(wordLetter)] == strconv.Itoa(pos) {
					currentScore += WordsLength
				} else if strings.Contains(lettersRules[string(wordLetter)], "|") {
					positions := strings.Split(lettersRules[string(wordLetter)], "|")
					for _, position := range positions {
						if position == strconv.Itoa(pos) {
							currentScore += WordsLength
						}
					}
				} else {
					currentScore++
				}
			}
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
			fmt.Println(val)
			return false
		}
	}
	return true
}
