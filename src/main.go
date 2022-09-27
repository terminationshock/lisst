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
	config = NewConfig()
	input := readFromPipe()
	itemList := NewItemList(input)

	run(itemList, 0)
}

func printHelp() {
	fmt.Println("Usage: COMMAND_IN | " + os.Args[0] + " PATTERN [COMMAND]\n")
	fmt.Println("This program reads the output of COMMAND_IN from the pipe and displays all lines as")
	fmt.Println("an interactive list. Each line is matched against a regular expression PATTERN. The")
	fmt.Println("first match in each line is highlighted. When [Enter] is pressed, the given COMMAND")
	fmt.Println("is executed with the highlighted match of the selected line as additional argument.")
	fmt.Println("When the COMMAND returns, the list is displayed again.")
	fmt.Println("\nExamples:")
	fmt.Println("\n   git log | " + os.Args[0] + " \"\\b[0-9a-z]{40}\\b\" git show")
	fmt.Println("   ...will display all commits, highlight all commit hashes, and execute")
	fmt.Println("      `git show <commit hash>` when [Enter] is pressed.")
	fmt.Println("\n   ls -1 | " + os.Args[0] + " \".+\" less")
	fmt.Println("   ...will display all files in the directory and execute `less <file name>` when")
	fmt.Println("      [Enter] is pressed.")
	fmt.Println("\n   grep -r func | ./lisst \"^(.*):\" vi")
	fmt.Println("   ...will recursively grep for \"func\" in all files, highlight all file names,")
	fmt.Println("      and execute `vi <file name>` when [Enter] is pressed.")
	fmt.Println("\n   squeue -u $USER | ./lisst \"^\\s*([0-9]{1,})\\b\" scancel")
	fmt.Println("   ...will query SLURM for all running jobs of the current user, highlight all job IDs,")
	fmt.Println("      and execute `scancel <job ID>` when [Enter] is pressed.")
	fmt.Println("\nKey bindings:")
	fmt.Println("\n   [q] or [Esc]      Quit")
	fmt.Println("   [Up] and [Down]   Browse lines")
	fmt.Println("   [Enter]           Execute COMMAND with PATTERN match as argument")
	fmt.Println("\nEnvironment variable:")
	fmt.Println("\n   LISST_COLOR       Set this variable to change the color for highlighting, which")
	fmt.Println("                     defaults to \"red\". Assign \"-\" to disable highlighting.")
}

func readFromPipe() []string {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// There is no pipe
		fmt.Fprintln(os.Stderr, "Missing input")
		printHelp()
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)

	input := []string{}
	for scanner.Scan() {
		// Read input line by line
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

	if config.test {
		// Used for the tests
		ui.app.Stop()
		if config.program != "" && ui.list.GetItemCount() > 0 {
			ui.itemList.Get(0).LaunchProgram()
		} else {
			itemList.Print()
		}
		os.Exit(0)
	}
}

func (ui *Ui) setStatus() {
	info := "\n"
	space := "     "

	if config.pattern != nil {
		info += fmt.Sprintf("%s%s", config.pattern, space)
	}
	info += fmt.Sprintf("Line %d of %d", ui.list.GetCurrentItem() + 1, ui.list.GetItemCount())
	if config.executed != "" {
		info += space + "Most recently executed: " + config.executed
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
