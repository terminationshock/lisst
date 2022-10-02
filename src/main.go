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
	pageList *PageList
	pageText *PageText
	pageTextVisible bool
	config *Config
}

type PageList struct {
	flex *tview.Flex
	list *tview.List
	status *tview.TextView
	itemList *ItemList
}

type PageText struct {
	flex *tview.Flex
	text *tview.TextView
	status *tview.TextView
}

var config *Config

func main() {
	config = NewConfig()
	input := readFromPipe()
	itemList := NewItemList(input)

	run(itemList, 0, "", "")
}

func PrintHelp() {
	fmt.Println("Usage: COMMAND_IN | " + os.Args[0] + " [OPTIONS] PATTERN [COMMAND]\n")
	fmt.Println("This program reads the output of COMMAND_IN from the pipe and displays all lines as")
	fmt.Println("an interactive list. Each line is matched against a regular expression PATTERN. The")
	fmt.Println("first match in each line is highlighted. When [Enter] is pressed, the given COMMAND")
	fmt.Println("is executed with the highlighted match of the selected line as additional argument.")
	fmt.Println("The placeholder `{}` can be used in COMMAND to insert the match at a given position.")
	fmt.Println("When the COMMAND returns, the list is displayed again.")
	fmt.Println("\nKey bindings:")
	fmt.Println("\n   [q] or [Esc]        Quit")
	fmt.Println("   [Up] and [Down]     Browse lines")
	fmt.Println("   [n]                 Jump to the next line with a match")
	fmt.Println("   [N]                 Jump to the previous line with a match")
	fmt.Println("   [Enter]             Execute COMMAND with the PATTERN match as argument")
	fmt.Println("\nKeywords to replace PATTERN:")
	fmt.Println("\n   --line              Match the whole line")
	fmt.Println("   --git-commit-hash   Match a Git commit hash")
	fmt.Println("   --filename          Match the name of an existing file")
	fmt.Println("   --dirname           Match the name of an existing directory")
	fmt.Println("\nOther keyword OPTIONS:")
	fmt.Println("\n   --show-output       Show the output of COMMAND")
	fmt.Println("   --help              Display this help")
	fmt.Println("\nExamples:")
	fmt.Println("\n   git log --oneline | " + os.Args[0] + " \"\\b[0-9a-z]{7,40}\\b\" git show")
	fmt.Println("                       or")
	fmt.Println("   git log --oneline | " + os.Args[0] + " --git-commit-hash git show")
	fmt.Println("                       will display all commits, highlight all commit hashes, and")
	fmt.Println("                       show the commit details by executing `git show <commit hash>`")
	fmt.Println("                       when [Enter] is pressed.")
	fmt.Println("\n   ls -1 | " + os.Args[0] + " \".+\" less")
	fmt.Println("                       or")
	fmt.Println("   ls -1 | " + os.Args[0] + " --line less")
	fmt.Println("                       will display all files in the directory and show the selected")
	fmt.Println("                       file with `less <file name>` when [Enter] is pressed.")
	fmt.Println("\n   ls -1 | " + os.Args[0] + " --line cp {} ..")
	fmt.Println("                       will display all files in the directory and copy the selected")
	fmt.Println("                       file to the parent directory by executing `cp <file name> ..`")
	fmt.Println("                       when [Enter] is pressed.")
	fmt.Println("\n   grep -r func | ./lisst \"^(.*?):\" vi")
	fmt.Println("                       or")
	fmt.Println("   grep -r func | ./lisst --filename vi")
	fmt.Println("                       will recursively grep for \"func\" in all files, highlight all")
	fmt.Println("                       file names, and open the text editor `vi <file name>` when")
	fmt.Println("                       [Enter] is pressed.")
	fmt.Println("\n   squeue -u $USER | ./lisst --show-output \"^\\s*([0-9]{1,})\\b\" scontrol show job")
	fmt.Println("                       will query SLURM for all running jobs of the current user,")
	fmt.Println("                       highlight all job IDs, and show details of the selected job by")
	fmt.Println("                       executing `scontrol show job <job ID>` when [Enter] is pressed.")
	fmt.Println("                       Note the additional flag `--show-output` to display the output")
	fmt.Println("                       of `scontrol` instead of printing it to the terminal in the")
	fmt.Println("                       background.")
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

func run(itemList *ItemList, selectedIndex int, programExecuted string, programOutput string) {
	ui := initUi()
	ui.fillList(itemList, selectedIndex)
	ui.pageList.setStatus(programExecuted)

	if programExecuted != "" && config.showProgramOutput {
		ui.setText(programExecuted, programOutput)
	}

	err := ui.app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initUi() *Ui {
	ui := &Ui{
		pageList: &PageList{},
		pageText: &PageText{},
	}

	ui.app = tview.NewApplication()
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Keys for quitting the program
		if event.Rune() == 'q' || event.Key() == tcell.KeyEsc {
			if ui.pageTextVisible {
				ui.app.SetRoot(ui.pageList.flex, true)
				ui.pageTextVisible = false
			} else {
				ui.app.Stop()
				os.Exit(0)
			}
		} else if event.Rune() == 'n' && !ui.pageTextVisible {
			ui.pageList.jumpToMatch(true)
		} else if event.Rune() == 'N' && !ui.pageTextVisible {
			ui.pageList.jumpToMatch(false)
		}
		return event
	})

	// Container for the list and its status bar
	ui.pageList.flex = tview.NewFlex()
	ui.pageList.flex.SetDirection(tview.FlexRow)
	ui.pageTextVisible = false

	// List for the matches
	ui.pageList.list = tview.NewList()
	ui.pageList.list.ShowSecondaryText(false)
	ui.pageList.list.SetWrapAround(false)
	ui.pageList.list.SetHighlightFullLine(true)
	ui.pageList.flex.AddItem(ui.pageList.list, 0, 1, true)

	// Invoked when a line is highlighted
	ui.pageList.list.SetChangedFunc(ui.pageList.lineSelected)
	if config.program != "" {
		// Invoked when enter is pressed on a line
		ui.pageList.list.SetSelectedFunc(ui.lineClicked)
	}

	// Status line at the bottom
	ui.pageList.status = tview.NewTextView()
	ui.pageList.status.SetScrollable(false)
	ui.pageList.status.SetWrap(false)
	ui.pageList.flex.AddItem(ui.pageList.status, 2, 1, false)

	// Container for the text and its status bar
	ui.pageText.flex = tview.NewFlex()
	ui.pageText.flex.SetDirection(tview.FlexRow)

	// Text field for command output
	ui.pageText.text = tview.NewTextView()
	ui.pageText.text.SetScrollable(true)
	ui.pageText.text.SetWrap(false)
	ui.pageText.flex.AddItem(ui.pageText.text, 0, 1, true)

	// Status line at the bottom
	ui.pageText.status = tview.NewTextView()
	ui.pageText.status.SetScrollable(false)
	ui.pageText.status.SetWrap(false)
	ui.pageText.status.SetRegions(true)
	ui.pageText.flex.AddItem(ui.pageText.status, 1, 1, false)

	ui.app.SetRoot(ui.pageList.flex, true)
	return ui
}

func (ui *Ui) fillList(itemList *ItemList, selectedIndex int) {
	ui.pageList.itemList = itemList

	for _, item := range ui.pageList.itemList.items {
		// Build the list
		ui.pageList.list.AddItem(item.display, "", 0, nil)
	}

	if selectedIndex < ui.pageList.list.GetItemCount() {
		// Set the cursor to the previous line if possible
		ui.pageList.list.SetCurrentItem(selectedIndex)
	}

	if config.test {
		// Used for the tests
		ui.app.Stop()
		if config.program != "" && ui.pageList.list.GetItemCount() > 0 {
			_, output := ui.pageList.itemList.Get(0).RunCommand()
			if output != "" {
				fmt.Println(output)
			}
		} else {
			itemList.Print()
		}
		os.Exit(0)
	}
}

func (pageList *PageList) setStatus(programExecuted string) {
	info := "\n"
	space := "     "

	if config.pattern != nil {
		info += fmt.Sprintf("%s%s", config.pattern, space)
	}

	index := pageList.list.GetCurrentItem()
	info += fmt.Sprintf("Line %d of %d", index + 1, pageList.list.GetItemCount())
	if config.program != "" && programExecuted == "" {
		info += space + pageList.itemList.Get(index).PrintCommand()
	}

	pageList.status.SetText(info)
}

func (pageList *PageList) jumpToMatch(forward bool) {
	count := pageList.list.GetItemCount()
	if count < 2 {
		return
	}

	index := pageList.list.GetCurrentItem()

	if forward && index < count - 1 {
		for i := index + 1; i < count; i++ {
			if pageList.itemList.Get(i).HasMatch() {
				index = i
				break
			}
		}
	} else if !forward && index > 0 {
		for i := index - 1; i >= 0; i-- {
			if pageList.itemList.Get(i).HasMatch() {
				index = i
				break
			}
		}
	}

	pageList.list.SetCurrentItem(index)
	pageList.setStatus("")
}

func (ui *Ui) setText(programExecuted string, programOutput string) {
	// Fill the text view with the output of the program
	ui.pageText.text.SetText(programOutput)

	// Set the status bar text
	status := fmt.Sprintf("[\"0\"]%s[\"\"]", programExecuted)
	ui.pageText.status.SetText(status)
	ui.pageText.status.Highlight("0")

	// Display the text together with the status bar
	ui.app.SetRoot(ui.pageText.flex, true)
	ui.pageTextVisible = true
}

// Signature of this function must not be changed
func (pageList *PageList) lineSelected(index int, _ string, _ string, _ rune) {
	pageList.setStatus("")
}

// Signature of this function must not be changed
func (ui *Ui) lineClicked(index int, _ string, _ string, _ rune) {
	item := ui.pageList.itemList.Get(index)
	if item.HasMatch() {
		ui.app.Stop()

		// Run the program and fetch the output if it is not writing to stdout
		program, output := item.RunCommand()

		// Restart the list view
		run(ui.pageList.itemList, index, program, output)
	}
}
