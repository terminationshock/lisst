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
		tokens := config.pattern.FindStringSubmatch(line)
		if tokens != nil && len(tokens) > 0 {
			// Highlight first occurrence of regexp in line
			index := 0
			if len(tokens) > 1 {
				index = 1
			}
			item.match = tokens[index]
			highlighted := strings.Replace(tokens[0], item.match, "[" + config.color + "]" + item.match + "[-]", 1)
			item.display = strings.Replace(item.display, tokens[0], highlighted, 1)
		}
	}

	// Replace tab characters
	item.display = strings.ReplaceAll(item.display, "\t", strings.Repeat(" ", tabSize))
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
