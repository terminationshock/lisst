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
	text *tview.TextView
	textVisible bool
	status *tview.TextView
	config *Config
	itemList *ItemList
}

var config *Config

func main() {
	config = NewConfig()
	input := readFromPipe()
	itemList := NewItemList(input)

	run(itemList, 0, false, "")
}

func PrintHelp() {
	fmt.Println("Usage: COMMAND_IN | " + os.Args[0] + " [OPTIONS] PATTERN [COMMAND]\n")
	fmt.Println("This program reads the output of COMMAND_IN from the pipe and displays all lines as")
	fmt.Println("an interactive list. Each line is matched against a regular expression PATTERN. The")
	fmt.Println("first match in each line is highlighted. When [Enter] is pressed, the given COMMAND")
	fmt.Println("is executed with the highlighted match of the selected line as additional argument.")
	fmt.Println("When the COMMAND returns, the list is displayed again.")
	fmt.Println("\nKey bindings:")
	fmt.Println("\n   [q] or [Esc]        Quit")
	fmt.Println("   [Up] and [Down]     Browse lines")
	fmt.Println("   [Enter]             Execute COMMAND with PATTERN match as argument")
	fmt.Println("\nKeywords to replace PATTERN:")
	fmt.Println("\n   --line              Match the whole line")
	fmt.Println("   --git-commit-hash   Match a Git commit hash")
	fmt.Println("   --grep-filename     Match the filename prefix in the output of `grep`")
	fmt.Println("\nOther keyword OPTIONS:")
	fmt.Println("\n   --show-output       Show the output of COMMAND")
	fmt.Println("   --help              Display this help")
	fmt.Println("\nExamples:")
	fmt.Println("\n   git log --oneline | " + os.Args[0] + " \"\\b[0-9a-z]{7,40}\\b\" git show")
	fmt.Println("                       or")
	fmt.Println("   git log --oneline | " + os.Args[0] + " --git-commit-hash git show")
	fmt.Println("                       will display all commits, highlight all commit hashes, and")
	fmt.Println("                       execute `git show <commit hash>` when [Enter] is pressed.")
	fmt.Println("\n   ls -1 | " + os.Args[0] + " \".+\" less")
	fmt.Println("                       or")
	fmt.Println("   ls -1 | " + os.Args[0] + " --line less")
	fmt.Println("                       will display all files in the directory and execute")
	fmt.Println("                       `less <file name>` when [Enter] is pressed.")
	fmt.Println("\n   grep -r func | ./lisst \"^(.*?):\" vi")
	fmt.Println("                       or")
	fmt.Println("   grep -r func | ./lisst --grep-filename vi")
	fmt.Println("                       will recursively grep for \"func\" in all files, highlight all")
	fmt.Println("                       file names, and execute `vi <file name>` when [Enter] is")
	fmt.Println("                       pressed.")
	fmt.Println("\n   squeue -u $USER | ./lisst --show-output \"^\\s*([0-9]{1,})\\b\" scontrol show job")
	fmt.Println("                       will query SLURM for all running jobs of the current user,")
	fmt.Println("                       highlight all job IDs, and execute `scontrol show job <job ID>`")
	fmt.Println("                       when [Enter] is pressed. Note the additional flag")
	fmt.Println("                       `--show-output` to display the output of `scontrol` instead")
	fmt.Println("                       of printing it to the terminal in the background.")
	fmt.Println("\nEnvironment variable:")
	fmt.Println("\n   LISST_COLOR         Set this variable to change the color for highlighting, which")
	fmt.Println("                       defaults to \"red\". Assign \"-\" to disable highlighting.")
}

func readFromPipe() []string {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// There is no pipe
		fmt.Fprintln(os.Stderr, "Missing input")
		PrintHelp()
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

func run(itemList *ItemList, selectedIndex int, programExecuted bool, programOutput string) {
	ui := initUi()
	ui.fillList(itemList, selectedIndex)
	ui.setStatus(programExecuted)

	if programExecuted && config.showProgramOutput {
		ui.setText(programOutput)
	}

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
			if ui.textVisible {
				ui.app.SetRoot(ui.flex, true)
				ui.textVisible = false
			} else {
				ui.app.Stop()
				os.Exit(0)
			}
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

	// Text field for command output
	ui.text = tview.NewTextView()
	ui.text.SetScrollable(true)
	ui.text.SetWrap(false)
	ui.textVisible = false

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

	if config.test {
		// Used for the tests
		ui.app.Stop()
		if config.program != "" && ui.list.GetItemCount() > 0 {
			output := ui.itemList.Get(0).LaunchProgram()
			if output != "" {
				fmt.Println(output)
			}
		} else {
			itemList.Print()
		}
		os.Exit(0)
	}
}

func (ui *Ui) setStatus(programExecuted bool) {
	info := "\n"
	space := "     "

	if config.pattern != nil {
		info += fmt.Sprintf("%s%s", config.pattern, space)
	}

	index := ui.list.GetCurrentItem()
	info += fmt.Sprintf("Line %d of %d", index + 1, ui.list.GetItemCount())
	if config.program != "" && !programExecuted {
		info += space + PrintCommand(ui.itemList.Get(index).match)
	}

	ui.status.SetText(info)
}

func (ui *Ui) setText(programOutput string) {
	// Fill the text view with the output of the program
	ui.text.SetText(programOutput)
	ui.app.SetRoot(ui.text, true)
	ui.textVisible = true
}

// Signature of this function must not be changed
func (ui *Ui) lineSelected(index int, _ string, _ string, _ rune) {
	ui.setStatus(false)
}

// Signature of this function must not be changed
func (ui *Ui) lineClicked(index int, _ string, _ string, _ rune) {
	item := ui.itemList.Get(index)
	if item.match != "" {
		ui.app.Stop()

		// Run the program and fetch the output if it is not writing to stdout
		output := item.LaunchProgram()

		// Restart the list view
		run(ui.itemList, index, true, output)
	}
}
