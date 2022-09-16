package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Options struct {
	pattern *regexp.Regexp
	program string
	programArgs []string
	color string
	debug bool
}

type Ui struct {
	app *tview.Application
	flex *tview.Flex
	list *tview.List
	status *tview.TextView
	itemList *ItemList
}

var options *Options

func main() {
	input := readFromPipe()
	parseOptions(input)
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

func parseOptions(input []string) {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println("Usage: COMMAND_IN | " + os.Args[0] + " PATTERN [COMMAND]")
		fmt.Println("\nThis program displays the output of COMMAND_IN as interactive list")
		fmt.Println("and opens COMMAND after a certain match has been selected.")
		fmt.Println("   PATTERN    - Accepts regular expressions")
		fmt.Println("              - Case insensitive (can be toggled with key [c])")
		fmt.Println("   DIRECTORY  - Recursive search within this directory")
		fmt.Println("              - Defaults to the current working directory")
		fmt.Println("\nKey bindings:")
		fmt.Println("   [q] or [Esc]      Quit")
		fmt.Println("   [Up] and [Down]   Select a match")
		fmt.Println("   [Enter]           Open the selected file with the editor (see variable EDITOR)")
		fmt.Println("\nEnvironment variables:")
		fmt.Println("   LISST_COLOR        - Highlight color for matches")
		fmt.Println("                       - Assign \"-\" to disable highlighting")
		fmt.Println("                       - Default: \"red\"")
		os.Exit(0)
	}

	options = &Options {
		pattern: nil,
		program: "",
		programArgs: []string{},
		color: "#ff0000",
		debug: false,
	}

	if len(os.Args) > 1 {
		// Regex pattern
		pattern, err := regexp.Compile(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid regular expression")
			os.Exit(1)
		}
		options.pattern = pattern
	}

	if len(os.Args) > 2 {
		// Program name
		options.program = os.Args[2]
	}

	if len(os.Args) > 3 {
		// Program arguments
		options.programArgs = os.Args[3:]
	}

	color := strings.TrimSpace(os.Getenv("LISST_COLOR"))
	if color != "" {
		// Highlighting color for matches
		options.color = color
	}

	if os.Getenv("LISST_DEBUG") != "" {
		options.debug = true
	}
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
	if options.program != "" {
		// Invoked when enter is pressed on a line
		ui.list.SetSelectedFunc(ui.lineClicked)
	}

	// Status line at the bottom
	ui.status = tview.NewTextView()
	ui.status.SetScrollable(false)
	ui.status.SetDynamicColors(true)
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

	if options.debug {
		ui.app.Stop()
		if options.program != "" && len(ui.itemList.items) > 0 {
			ui.itemList.Get(0).LaunchProgram()
		} else {
			itemList.Debug()
		}
		os.Exit(0)
	}
}

func (ui *Ui) setStatus() {
	info := fmt.Sprintf("\nLine %d of %d", ui.list.GetCurrentItem() + 1, ui.list.GetItemCount())
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
