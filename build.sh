#!/bin/bash
set -e

go get github.com/rivo/tview
go get github.com/gdamore/tcell/v2
go build -o lisst src/config.go src/cmd.go src/items.go src/main.go
