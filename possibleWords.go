package main

import (
	"regexp"
	"strings"
)

func getPossibleWords(wrongLetters string, correct []string, wrongBoxLetters []string) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	validLetters := ""
	for _, l := range alphabet {
		if !strings.Contains(wrongLetters, string(l)) {
			validLetters += string(l)
		}
	}
	rxString := ""
	for i, l := range correct {
		if len(l) == 0 {
			if len(wrongBoxLetters[i]) > 0 {
				validSub := ""
				for _, d := range validLetters {
					if !strings.Contains(wrongBoxLetters[i], string(d)) {
						validSub += string(d)
					}
				}
				rxString += "[" + validSub + "]"
			} else {
				rxString += "[" + validLetters + "]"
			}
		} else {
			rxString += string(l)
		}
	}
	rx := regexp.MustCompile(rxString)
	combined := strings.Join(wrongBoxLetters, "")
	rx2 := regexp.MustCompile(".")
	if len(combined) > 0 {
		rx2 = regexp.MustCompile("[" + combined + "]")
	}
	words := readWordlist()
	remainingWords := ""
	for _, w := range words {
		if rx.MatchString(w) && rx2.MatchString(w) {
			if len(remainingWords) > 0 {
				remainingWords += " "
			}
			remainingWords += w
		}
	}
	return remainingWords
}
