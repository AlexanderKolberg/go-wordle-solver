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
	correctLetters := createLetterForm(tcell.ColorGreen, 1, 5)
	wrongBox := createLetterForm(tcell.ColorYellow, 4, 5)
	wordsView := tview.NewTextView()
	wordsView.SetBorder(true).SetTitle("Possible Words")
	reset := tview.NewButton("Reset")
	reset.SetBackgroundColor(tcell.ColorRed)
	submit := tview.NewButton("Submit")

	setWords := func(words string) {
		wordsView.SetText(words)
	}

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
		correct := getFormInputs(*correctLetters)
		wrongBoxLetters := getFormInputs(*wrongBox)
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
		setWords(remainingWords)
	}

	reset.SetSelectedFunc(func() { setWords("reset") })
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
