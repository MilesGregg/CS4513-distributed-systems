package cs4513_go_impl

import (
	"fmt"
	"sort"
	"os"
	"bufio"
	"strings"
	"regexp"
)

/*
Find the top K most common words in a text document.
What is a word?
	A word here only consists of alphanumeric characters, e.g., catch21
	All punctuations and other characters should be removed, e.g. "don't" becomes "dont" or "end." becomes "end"; done before the charThreshold
	A word has to satifies the charThreshold, e.g., if charThreshold = 5  "apple" is a word, but neither "new" or "york" are words
Matching condition
	Matching is case insensitive
Parameters:
- path: file path
- numWords: number of words to return (i.e. k)
- charThreshold: threshold for whether a token qualifies as a word
You should use `checkError` to handle potential errors.
*/
func TopWords(path string, numWords int, charThreshold int) []WordCount {
	// TODO: implement me
	// HINT: You may find the `strings.Fields` and `strings.ToLower` functions helpful
	// HINT: the regex "[^0-9a-zA-Z]+" can be used to spot any non-alphanumeric characters
	// make output word count
	var output []WordCount

	//fmt.Println(path)
	file, errFile := os.Open(path)
	if errFile != nil { // check for error
		fmt.Println(errFile) // printout error
		return nil
	}
	//fmt.Println(string(file))
	defer file.Close()
	// make bufio in order to iterate through the file
	it := bufio.NewScanner(file) // make the scanner
	it.Split(bufio.ScanWords) // split on all of the words

	for it.Scan() { // iterate through all of the words in the bufio stream
		word_lower := strings.ToLower(it.Text()) // make the current word lowercase
		word := regexp.MustCompile(`[^0-9a-zA-Z]+`).ReplaceAllString(word_lower, "") // apply the regex and 

		if len(word) >= charThreshold { // if the word is above the char threshold then append and do the normal routine
			if len(output) == 0 { // if we appending the first word onto the 
				o := append(output, WordCount{Word: word, Count: 1}) // make 
				output = o // append the o to the output array of WordCount
				//fmt.Println(o)
				continue
			} else {
				for i, val := range output { // iterate through the whole output WordCount array
					if word == val.Word { // if the word exists then update it
						new := WordCount{Word: val.Word, Count: val.Count + 1} // update the word count
						output[i] = new // reput this back into the WordCount array
						break
					}
		
					if i == len(output)-1 { // if the word isn't in the output WordCount array then put it in
						o := append(output, WordCount{Word: word, Count: 1}) // insert into the output WordCount array
						output = o // set it as the output WordCount array
					}
				}
			}
		}
	}

	sortWordCounts(output) // sort the output Word Counts array
	var truncated []WordCount // time to truncate the output Word Count
	for i, val := range output { // iterate through the 
		if i == numWords { // if the i equals to the numWords than stop the turncation 
			break
		}
		truncated = append(truncated, val) // append the val onto the truncatation 
		//fmt.Println(truncated[i])
	}

	//fmt.Println(truncated)

	return truncated // output the truncated 
}

/*
Do NOT modify this struct!
A struct that represents how many times a word is observed in a document
*/
type WordCount struct {
	Word  string
	Count int
}

/*
Do NOT modify this function!
*/
func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

/*
Do NOT modify this function!
Helper function to sort a list of word counts in place.
This sorts by the count in decreasing order, breaking ties using the word.
*/
func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}
