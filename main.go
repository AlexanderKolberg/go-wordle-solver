package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	wrongLettersInput := tview.NewInputField().SetFieldWidth(15).SetFieldBackgroundColor(tcell.ColorGray)
	correctLettersForm := createLetterForm(tcell.ColorGreen, 1, 5)
	wrongBoxLettersForm := createLetterForm(tcell.ColorYellow, 4, 5)
	possibleWordsView := tview.NewTextView()
	possibleWordsView.SetBorder(true).SetTitle("Possible Words")
	resetButton := tview.NewButton("Reset")
	resetButton.SetBackgroundColor(tcell.ColorRed)
	submitButton := tview.NewButton("Submit")

	submitHandler := func() {
		wrongLetters := wrongLettersInput.GetText()
		correct := getFormInputs(*correctLettersForm)
		wrongBoxLetters := getFormInputs(*wrongBoxLettersForm)
		words := getWords(wrongLetters, correct, wrongBoxLetters)
		possibleWordsView.SetText(words)
	}

	resetButton.SetSelectedFunc(func() { possibleWordsView.SetText("reset") })
	submitButton.SetSelectedFunc(submitHandler)
	form := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Wrong letters:"), 0, 1, false).
		AddItem(wrongLettersInput, 0, 1, false).
		AddItem(tview.NewTextView().SetText("Correct Letters:"), 0, 1, false).
		AddItem(correctLettersForm, 0, 1, false).
		AddItem(tview.NewTextView().SetText("Correct letters, wrong box:"), 0, 1, false).
		AddItem(wrongBoxLettersForm, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(submitButton, 0, 1, false).
			AddItem(resetButton, 0, 1, false), 0, 1, false)

	app := tview.NewApplication().EnableMouse(true)
	flex := tview.NewFlex().
		AddItem(form, 0, 1, false).
		AddItem(possibleWordsView, 0, 2, false)
	flex.SetBorder(true).SetTitle("Go wordle solver")
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
func getWords(wrongLetters string, correct []string, wrongBoxLetters []string) string {
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
func createLetterForm(color tcell.Color, inputSize int, fields int) *tview.Form {
	form := tview.NewForm().SetHorizontal(true).SetFieldBackgroundColor(color)
	for i := 0; i < fields; i++ {
		form.AddInputField("", "", inputSize, nil, nil)
	}
	form.SetBorderPadding(0, 0, 0, 0)
	return form
}

func getFormInputs(form tview.Form) []string {
	var out []string
	for i := 0; i < form.GetFormItemCount(); i++ {
		item := form.GetFormItem(i).(*tview.InputField).GetText()
		if len(item) > 0 {
			out = append(out, item)
		} else {
			out = append(out, "")
		}
	}
	return out
}
