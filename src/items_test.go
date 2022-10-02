package main

import (
	"regexp"
	"testing"
)

func TestProcessNoRegexp(t *testing.T) {
	lines := []string{"line\twith\ttab", "line with match and mctch"}

	config = &Config{}
	config.pattern = nil
	items := NewItemList(lines)

	if len(items.items) != 2 {
		t.Error("Incorrect number of lines")
	}
	if items.items[0].original != lines[0] || items.items[1].original != lines[1] {
		t.Error("Incorrect original lines")
	}
	if items.items[0].display != "line    with    tab" {
		t.Error("Incorrect processed line with tabs")
	}
	if items.items[0].match != "" || items.items[1].match != "" {
		t.Error("Incorrect matches")
	}
}

func TestProcessRegexp(t *testing.T) {
	lines := []string{"the match", "no match, but the match here and not the match again", "the match, the mbtch, the mctch"}

	config = &Config{}
	config.pattern = regexp.MustCompile("the (m[a-c]tch)")
	items := NewItemList(lines)

	if items.items[0].display != "the [::r]match[::-]" {
		t.Error("Incorrect processed line with single submatch")
	}
	if items.items[1].display != "no match, but the [::r]match[::-] here and not the match again" {
		t.Error("Incorrect processed line with multiple submatches highlighting the second")
	}
	if items.items[2].display != "the [::r]match[::-], the mbtch, the mctch" {
		t.Error("Incorrect processed line with multiple submatches highlighting the first")
	}

	config.pattern = regexp.MustCompile("m[a-c]tch")
	items = NewItemList(lines)

	if items.items[0].display != "the [::r]match[::-]" {
		t.Error("Incorrect processed line with single match")
	}
	if items.items[1].display != "no [::r]match[::-], but the match here and not the match again" {
		t.Error("Incorrect processed line with multiple matches highlighting the first")
	}

	config.pattern = regexp.MustCompile("m[b-c]t(c)h")
	items = NewItemList(lines)

	if items.items[0].display != "the match" {
		t.Error("Incorrect processed line without submatch")
	}
	if items.items[1].display != "no match, but the match here and not the match again" {
		t.Error("Incorrect processed line without any submatch")
	}
	if items.items[2].display != "the match, the mbt[::r]c[::-]h, the mctch" {
		t.Error("Incorrect processed line with multiple submatches highlighting the first")
	}

	config.pattern = regexp.MustCompile("m(at)(c(h))")
	items = NewItemList(lines)

	if items.items[0].display != "the m[::r]at[::-]ch" {
		t.Error("Incorrect processed line with many submatches")
	}
}

func TestProcessRegexpWithFunc(t *testing.T) {
	lines := []string{"the match MATCH"}

	config = &Config{}
	config.pattern = regexp.MustCompile("\\S*")

	config.patternFunc = func(p string) bool {
		return len(p) == 5
	}
	items := NewItemList(lines)

	if items.items[0].display != "the [::r]match[::-] MATCH" {
		t.Error("Incorrect processed line with given function")
	}

	config.patternFunc = func(_ string) bool {
		return false
	}
	items = NewItemList(lines)

	if items.items[0].display != "the match MATCH" {
		t.Error("Incorrect processed line with given function")
	}
}
