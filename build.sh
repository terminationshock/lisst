#!/bin/bash
set -e

go get github.com/rivo/tview
go get github.com/gdamore/tcell/v2

go build -o lisst src/*

go test ./...
