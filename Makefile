.PHONY: go-get check unit-tests integration-tests clean

exe=lisst

$(exe): go-get
	go build -o $(exe) src/*

go-get:
	go get -u github.com/rivo/tview
	go get -u github.com/gdamore/tcell/v2

check: unit-tests integration-tests

unit-tests: $(exe)
	go test ./...

integration-tests: $(exe)
	./test.sh

clean:
	rm -f $(exe)
