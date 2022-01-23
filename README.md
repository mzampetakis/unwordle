# Unrordle

Unrordle is a toy app that tries to solve the [wordle](https://www.powerlanguage.co.uk/wordle/) puzzle and its
variations. Wordle is a game where use has to guess the WORDLE (a word) in 6 tries. Each guess must be a valid 5-letter
word. After each guess, the color of the tiles will change to show how close your guess was to the word.

Wordle is a [Constraint satisfaction problem](https://en.wikipedia.org/wiki/Constraint_satisfaction_problem) and the way
it is solved is using [Constraint propagation](https://en.wikipedia.org/wiki/Constraint_propagation) techniques.

## Prerequisites

In order to build the app an installed version of Go greater than v1.11 with GOMODULES enabled is required.

## Build the app

Clone this repo

```console
git clone https://github.com/mzampetakis/unwordle.git
cd unwordle
```

Build the app

```console
go build .
```

## Run the app

After a successful build, unwordle can be run by executing:

```console
./unwordle --dictionary=dictionary_file_path     
```

Where the `dictionary_file_path` is a file that points to the dictionary to use. The app comes with two dictionaries:

* en_5_letters : a provided dictionary with english words consisted of 5 letters
* gr_5_letters : a provided dictionary with greek words consisted of 5 letters

Dictionary is the source where unwordle retrieves words to use and also expects answers to be one of the contained
words. Dictionaries contain words with the required length that our puzzle requires as a simple list of words containing
a single word per line.

### Parameters

Apart from the `dictionary` parameter, unwordle can be executed with two optional parameters:

* --tries (integer): the amount of tries that a user is allowed to guess the worlde
* --opener (string): specify the first word to try
* --info (boolean): show info about each proposed word (possibility based on remaining words and score of the proposed
  word)

In order to use all of the available parameters we can use the following command:

```console
./unwordle --dictionary=en_5_letters --info=true --opener=arose
```

### Solving a puzzle

In order to solve a wordle puzzle as soon as unworlde starts a proposal will be given for the first word to use. The
first word (opener) is estimated based on the given dictionary. More on this process can be found at the  
`Unwordle Internals` chapter. The first proposal is like this:

```
Try #1: 		arose
```

As first try, the word `arose` is proposed to be used. The `unworlde` awaits for user input based on the results of the
puzzle. A string with the same length as the wordle's length indicating the presence of each letter should be given by
user. Letters accepted are `b`, `y` and `g`. These letters mean:

* b (black): the letter does not exist in the wordle
* y (yellow): the letter exists in the wordle but not at the given place
* g (green): the letter exists in the wordle at this specific place

```
Try #1: 	    	arose
Response (b|y|g): 	gbybb
```

The process continues until all tries are exhausted or a reply all of `g`s is given or only one possible word matches
our criteria.

```
Found Solution: 	thick
Hooray! :-)
```

or

```
No solution found.
Sorry. :-(
```

If we choose to display info for each try we get this output:

```
Try #3: 	dodge | Possibility: 1/125 | Score: 20
```

`1/125` is the possibility to success as unwordle has excluded all other words from the given dictionary. The way that
score is calculated is analyzed at the `Eliminating candidates` chapter.

# Unwordle Internals

## Estimating a good opener

Unworlde is able to estimate a good opener word (the first proposed word) instead os using a random one from the given
dictionary. The process of estimating a good opener goes as follows:

* the occurrence of each letter is calculated based on the given dictionary
* using the most frequent letters we search for a word that contains each one of them

This way the opener word will give more value even from the first try!

## Eliminating candidates

After each proposed word, unwordle requests for the user response based on the wordle result. This input contains
valuable information for each one of the letters of the proposed word. So

* for each letter that does not exist in the wordle (black), unwordle eliminates all dictionary's words that contain
  this letter
* for each letter that exists on the wordle but is placed on wrong place (yellow), unwordle eliminates all dictionary's
  words that contain this letter at ths specific place

So, the maximum score a wolde can achieve is 5 * wordsLength. Doing this process for each input result given, unworlde
manages to exclude as many words as possible from the given dictionary.

## Proposing a good solution

For each try, unwordle tries to propose the most promising word chosen from the words that have not been eliminated from
the input dictionary. For each word within the dictionary, a score is calculated before proposing a new word. The score
is calculated using only the given responses.

* score is incremented by 1 for each letter that exists both in the word and the wordle
* score is incremented by 5 for each letter that exists both in the word and the wordle in the exact place

## Supporting Unicode

In order to support unicode letters within a dictionary unwordle has to iterate over words with caution. Words' letters
cannot be accessed using indexes, so iterations over each word are used in order to get each single letter. Unicode is
necessary to support languages/dictionaries that use letters outside unicode characters.

# Future Improvements

Provided dictionaries have some limitations. Not all words are available to some `wordle` variations. This can cause
unwordle not to be able to find the solution. Apart from this, unwordle might propose to use some word that is not
acceptable by the puzzle game leaving the user unable to move forward.

Rules used to eliminate words from dictionary and also calculate the score can be both enriched. Sophisticated rules can
be applied such as when we have information about all the letters of a word (green and yellows) to exlude all other
words.
