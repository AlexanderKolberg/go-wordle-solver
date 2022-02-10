package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ReadWordlist() []string {
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

func main() {
	createLetterForm := func(color tcell.Color, inputSize int, fields int) *tview.Form {
		form := tview.NewForm().SetHorizontal(true).SetFieldBackgroundColor(color)
		for i := 0; i < fields; i++ {
			form.AddInputField("", "", inputSize, nil, nil)
		}
		form.SetBorderPadding(0, 0, 0, 0)
		return form
	}

	correctLetters := createLetterForm(tcell.ColorGreen, 1, 5)
	getCorrectLetters := func() string {
		out := ""
		for i := 0; i < correctLetters.GetFormItemCount(); i++ {
			item := correctLetters.GetFormItem(i).(*tview.InputField).GetText()
			if len(item) > 0 {
				out += item
			} else {
				out += "."
			}
		}
		return out
	}
	wrongBox := createLetterForm(tcell.ColorYellow, 4, 5)
	getWrongBoxLetters := func() []string {
		var out []string
		for i := 0; i < wrongBox.GetFormItemCount(); i++ {
			item := wrongBox.GetFormItem(i).(*tview.InputField).GetText()
			if len(item) > 0 {
				out = append(out, item)
			} else {
				out = append(out, "")
			}
		}
		return out
	}

	wordsView := tview.NewTextView()
	wordsView.SetBorder(true).SetTitle("Possible Words")
	setWords := func(words string) {
		wordsView.SetText(words)
	}

	wrongLettersInput := tview.NewInputField().SetFieldWidth(15).SetFieldBackgroundColor(tcell.ColorGray)

	submitHandler := func() {
		alphabet := "abcdefghijklmnopqrstuvwxyz"
		wrongLetters := wrongLettersInput.GetText()
		validLetters := ""
		for _, l := range alphabet {
			if !strings.Contains(wrongLetters, string(l)) {
				validLetters += string(l)
			}
		}
		rxString := ""
		correct := getCorrectLetters()
		wrongBoxLetters := getWrongBoxLetters()
		for i, l := range correct {
			if l == '.' {
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
		words := ReadWordlist()
		remainingWords := ""
		for _, w := range words {
			if rx.MatchString(w) && rx2.MatchString(w) {
				if len(remainingWords) > 0 {
					remainingWords += " "
				}
				remainingWords += w
			}
		}
		setWords(remainingWords)
	}

	reset := tview.NewButton("Reset")
	reset.SetBackgroundColor(tcell.ColorRed)
	reset.SetSelectedFunc(func() { setWords("reset") })
	submit := tview.NewButton("Submit")
	submit.SetSelectedFunc(submitHandler)
	form := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Wrong letters:"), 0, 1, false).
		AddItem(wrongLettersInput, 0, 1, false).
		AddItem(tview.NewTextView().SetText("Correct Letters:"), 0, 1, false).
		AddItem(correctLetters, 0, 1, false).
		AddItem(tview.NewTextView().SetText("Correct letters, wrong box:"), 0, 1, false).
		AddItem(wrongBox, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(submit, 0, 1, false).
			AddItem(reset, 0, 1, false), 0, 1, false)

	app := tview.NewApplication().EnableMouse(true)
	flex := tview.NewFlex().
		AddItem(form, 0, 1, false).
		AddItem(wordsView, 0, 2, false)
	flex.SetBorder(true).SetTitle("Go wordle solver")
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}