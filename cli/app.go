package cli

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"math/rand"
	"os/exec"
	"reflect"
	"sort"
	"vocabulary-record/dal"
	"vocabulary-record/model"
)

var app = tview.NewApplication()
var flex = tview.NewFlex()
var table = tview.NewTable()
var review = tview.NewFlex()
var reviewTable = tview.NewTable()
var reviewHint = tview.NewTextView()
var typeIn = tview.NewInputField()
var hint = tview.NewTextView()
var pages = tview.NewPages()

var hintTitle = []string{
	"Show Recent Words",
	"Review",
	"Add New Word",
}

var curWord model.Words

func App() {
	hint.SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)
	for index, title := range hintTitle {
		fmt.Fprintf(hint, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index+1, title)
	}

	table.SetBorder(true)
	//table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	//	if event.Key() == tcell.KeyTAB {
	//		app.SetFocus(menu)
	//	}
	//	return event
	//})

	numField := reflect.ValueOf(model.Words{}).NumField()
	for col := 0; col < numField; col++ {
		fieldName := reflect.TypeOf(model.Words{}).Field(col).Name
		reviewTable.SetCell(0, col, tview.NewTableCell(fieldName))
	}
	reviewTable.SetFixed(1, 0)
	reviewTable.SetSelectable(true, true)
	reviewTable.SetSelectedFunc(func(row, column int) {
		if row <= 0 || row >= reviewTable.GetRowCount() ||
			column != 1 {
			return
		}
		word := reviewTable.GetCell(row, column).Text
		url := fmt.Sprintf("https://www.youdao.com/result?word=%s&lang=en", word)
		exec.Command("xdg-open", url).Start()
	})

	reviewHint.SetText(".) familiar  /) unfamiliar")

	review.SetBorder(true)
	review.SetDirection(tview.FlexRow)
	review.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '.' {
			setFamiliar(curWord, true)
			getNewWord()
		} else if event.Rune() == '/' {
			setFamiliar(curWord, false)
			getNewWord()
		}
		return event
	})
	review.SetFocusFunc(func() {
		getNewWord()
	})
	review.AddItem(reviewTable, 0, 3, true)
	review.AddItem(reviewHint, 1, 1, false)

	typeIn.SetBorder(true)
	typeIn.SetLabel("Add New Word: ").
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return !tview.InputFieldInteger(textToCheck, lastChar)
		}).
		SetDoneFunc(func(key tcell.Key) {
			text := typeIn.GetText()
			if text == "" {
				return
			}
			if err := addNewWord(typeIn.GetText()); err != nil {
				panic(err)
			}
			typeIn.SetText("")
		})

	pages.AddPage(hintTitle[0], table, true, true)
	pages.AddPage(hintTitle[1], review, true, true)
	pages.AddPage(hintTitle[2], typeIn, true, true)

	hint.Highlight("1")
	queryRecent()
	pages.SwitchToPage(hintTitle[0])

	flex.SetDirection(tview.FlexRow).
		AddItem(pages, 0, 3, true).
		AddItem(hint, 1, 1, false)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '1' {
			queryRecent()
			pages.SwitchToPage(hintTitle[0])
			hint.Highlight("1")
		} else if event.Rune() == '2' {
			if curWord.Word == "" {
				getNewWord()
			}
			pages.SwitchToPage(hintTitle[1])
			hint.Highlight("2")
		} else if event.Rune() == '3' {
			pages.SwitchToPage(hintTitle[2])
			hint.Highlight("3")
		} else if event.Key() == tcell.KeyEsc {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func queryRecent() {
	words, err := dal.QueryWords(0, 20, dal.WordsConds{})
	if err != nil {
		panic(err)
	}
	table.Clear()
	numField := reflect.ValueOf(model.Words{}).NumField()
	for col := 0; col < numField; col++ {
		fieldName := reflect.TypeOf(model.Words{}).Field(col).Name
		table.SetCell(0, col, tview.NewTableCell(fieldName))
	}
	table.SetFixed(1, 0)
	row := 1
	for _, word := range words {
		for col := 0; col < numField; col++ {
			fieldVal := fmt.Sprintf("%v", reflect.ValueOf(word).Field(col))
			table.SetCell(row, col, tview.NewTableCell(fieldVal))
		}
		row++
	}

	table.SetSelectable(true, true)
	table.SetSelectedFunc(func(row, column int) {
		if row <= 0 || row >= table.GetRowCount() ||
			column != 1 {
			return
		}
		word := table.GetCell(row, column).Text
		url := fmt.Sprintf("https://www.youdao.com/result?word=%s&lang=en", word)
		exec.Command("xdg-open", url).Start()

	})
}

func addNewWord(newWord string) error {
	err := dal.InsertWords([]model.Words{
		{
			Word: newWord,
		},
	})
	return err
}

func setFamiliar(word model.Words, isFamiliar bool) {
	if word.Word == "" {
		return
	}
	if isFamiliar {
		word.FamiliarCnt++
	}
	word.OccurredCnt++
	if _, err := dal.UpdateWords(word); err != nil {
		panic(err)
	}
	idCell := reviewTable.GetCell(1, 0)
	if isFamiliar {
		idCell.SetBackgroundColor(tcell.ColorBlue)
	} else {
		idCell.SetBackgroundColor(tcell.ColorPink)
	}
}

// has perf issue.
// if words number exceeds 10000, there will have a cut-off
func getNewWord() {
	words, err := dal.QueryWords(0, 10000, dal.WordsConds{})
	if err != nil {
		panic(err)
	}
	// sorted by proficiency (ascending order)
	sort.Slice(words, func(i, j int) bool {
		var rateI, rateJ float64
		if words[i].OccurredCnt != 0 {
			rateI = float64(words[i].FamiliarCnt) / float64(words[i].OccurredCnt)
		}
		if words[j].OccurredCnt != 0 {
			rateJ = float64(words[j].FamiliarCnt) / float64(words[j].OccurredCnt)
		}
		if rateI < rateJ {
			return true
		}
		if rateI == rateJ {
			return words[i].OccurredCnt < words[j].OccurredCnt
		}
		return false
	})
	windowSize := 3
	if len(words)/10 > windowSize {
		windowSize = len(words) / 10
	}
	idx := rand.Intn(windowSize)
	if idx >= len(words) {
		idx = 0
	}

	numField := reflect.ValueOf(model.Words{}).NumField()
	oldRow := reviewTable.GetRowCount()
	for row := oldRow - 1; row > 0; row-- {
		for col := 0; col < numField; col++ {
			reviewTable.SetCell(row+1, col, reviewTable.GetCell(row, col))
		}
	}
	for col := 0; col < numField; col++ {
		fieldVal := fmt.Sprintf("%v", reflect.ValueOf(words[idx]).Field(col))
		reviewTable.SetCell(1, col, tview.NewTableCell(fieldVal))
	}

	curWord = words[idx]
}
