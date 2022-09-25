package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Ui struct {
	app *tview.Application
	flex *tview.Flex
	list *tview.List
	status *tview.TextView
	config *Config
	itemList *ItemList
}

var config *Config

func main() {
	input := readFromPipe()
	config = NewConfig(input)
	itemList := NewItemList(input)

	run(itemList, 0)
}

func readFromPipe() []string {
	scanner := bufio.NewScanner(os.Stdin)

	input := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			input = append(input, line)
		}
	}

	if len(input) == 0 {
		fmt.Fprintln(os.Stderr, "Empty input")
		os.Exit(1)
	}

	err := scanner.Err()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}
	return input
}

func run(itemList *ItemList, selectedIndex int) {
	ui := initUi()
	ui.fillList(itemList, selectedIndex)

	err := ui.app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initUi() *Ui {
	ui := &Ui{}
	ui.app = tview.NewApplication()
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Keys for quitting the program
		if event.Rune() == 'q' || event.Key() == tcell.KeyEsc {
			ui.app.Stop()
			os.Exit(0)
		}
		return event
	})

	// Container for the widgets
	ui.flex = tview.NewFlex()
	ui.flex.SetDirection(tview.FlexRow)

	// List for the matches
	ui.list = tview.NewList()
	ui.list.ShowSecondaryText(false)
	ui.list.SetWrapAround(false)
	ui.list.SetHighlightFullLine(true)
	ui.flex.AddItem(ui.list, 0, 1, true)

	// Invoked when a line is highlighted
	ui.list.SetChangedFunc(ui.lineSelected)
	if config.program != "" {
		// Invoked when enter is pressed on a line
		ui.list.SetSelectedFunc(ui.lineClicked)
	}

	// Status line at the bottom
	ui.status = tview.NewTextView()
	ui.status.SetScrollable(false)
	ui.status.SetWrap(false)
	ui.flex.AddItem(ui.status, 2, 1, false)

	ui.app.SetRoot(ui.flex, true)
	return ui
}

func (ui *Ui) fillList(itemList *ItemList, selectedIndex int) {
	ui.itemList = itemList

	for _, item := range ui.itemList.items {
		// Build the list
		ui.list.AddItem(item.display, "", 0, nil)
	}

	if selectedIndex < ui.list.GetItemCount() {
		// Set the cursor to the previous line if possible
		ui.list.SetCurrentItem(selectedIndex)
	}

	// Draw the status line
	ui.setStatus()

	if config.debug {
		ui.app.Stop()
		if config.program != "" && ui.list.GetItemCount() > 0 {
			ui.itemList.Get(0).LaunchProgram()
		} else {
			itemList.Debug()
		}
		os.Exit(0)
	}
}

func (ui *Ui) setStatus() {
	info := ""
	if config.pattern != nil {
		info = fmt.Sprintf("\n%s     Line %d of %d", config.pattern, ui.list.GetCurrentItem() + 1, ui.list.GetItemCount())
	} else {
		info = fmt.Sprintf("\nLine %d of %d", ui.list.GetCurrentItem() + 1, ui.list.GetItemCount())
	}
	ui.status.SetText(info)
}

// Signature of this function must not be changed
func (ui *Ui) lineSelected(index int, _ string, _ string, _ rune) {
	ui.setStatus()
}

// Signature of this function must not be changed
func (ui *Ui) lineClicked(index int, _ string, _ string, _ rune) {
	item := ui.itemList.Get(index)
	if item.match != "" {
		ui.app.Stop()

		item.LaunchProgram()

		// Restart the list view
		run(ui.itemList, index)
	}
}
