.PHONY: all build go-get unittest test clean

exe=lisst

all: build unittest test

build: go-get
	go build -o $(exe) src/*

go-get:
	go get github.com/rivo/tview
	go get github.com/gdamore/tcell/v2

unittest: build
	go test ./...

test: build
	./test.sh

clean:
	rm -f $(exe)
