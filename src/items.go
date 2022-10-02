package main

import (
	"fmt"
	"strings"
)

const tabSize = 4

type ItemList struct {
	items []Item
}

type Item struct {
	original string
	display string
	match string
}

func NewItemList(input []string) *ItemList {
	list := &ItemList {
		items: make([]Item, len(input)),
	}

	for i, line := range input {
		list.items[i].process(line)
	}

	return list
}

func (item *Item) process(line string) {
	item.original = line
	item.display = line

	if config.pattern != nil {
		tokens := config.pattern.FindAllStringSubmatch(line, -1)
		if tokens != nil {
			for _, matches := range tokens {
				// Highlight the first match only
				if item.highlightFirstMatch(matches) {
					break
				}
			}
		}
	}

	// Replace tab characters
	item.display = strings.ReplaceAll(item.display, "\t", strings.Repeat(" ", tabSize))
}

func (item *Item) highlightFirstMatch(matches []string) bool {
	// If there is any submatch, highlight the first submatch, otherwise highlight the entire match
	index := 0
	if len(matches) > 1 {
		index = 1
	}
	match := matches[index]

	// Check the match using the pattern function and highlight it if the result is true
	if config.patternFunc == nil || config.patternFunc(match) {
		item.match = match
		highlighted := strings.Replace(matches[0], item.match, "[" + config.color + "]" + item.match + "[-]", 1)
		item.display = strings.Replace(item.display, matches[0], highlighted, 1)
		return true
	}
	return false
}

func (item *Item) HasMatch() bool {
	return item.match != ""
}

func (item *Item) PrintCommand() string {
	return PrintCommand(item.match)
}

func (item *Item) RunCommand() (string, string) {
	return RunCommand(item.match)
}

func (list *ItemList) Get(index int) *Item {
    return &list.items[index]
}

func (list *ItemList) Print() {
	for _, item := range list.items {
		fmt.Println(item.display)
	}
}
