package main

import (
	"bufio"
	"os"
)

func readWordlist() []string {
	if fp, err := os.Open("wordlist.txt"); err != nil {
		panic(err)
	} else {
		defer fp.Close()
		var wordlist []string
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			w := string(scanner.Bytes())
			if w != "" {
				wordlist = append(wordlist, w)
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		return wordlist
	}
}
