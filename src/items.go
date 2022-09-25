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

	if options.pattern != nil {
		tokens := options.pattern.FindAllStringSubmatch(line, 1)
		if len(tokens) > 0 && len(tokens[0]) > 0 {
			// Highlight regexp in line
			last := len(tokens[0]) - 1
			item.match = tokens[0][last]
			item.display = strings.Replace(item.display, item.match, "[" + options.color + "]" + item.match + "[-]", 1)
		}
	}

	// Replace tab characters
	item.display = strings.ReplaceAll(item.display, "\t", strings.Repeat(" ", tabSize))
}

func (item *Item) LaunchProgram() {
	LaunchProgram(item.match)
}

func (list *ItemList) Get(index int) *Item {
    return &list.items[index]
}

func (list *ItemList) Debug() {
	for _, item := range list.items {
		fmt.Println(item.display)
	}
}
