BINARY_NAME=sentinel

all: build

build:
	go build -o $(BINARY_NAME) ./cmd/hackura/main.go

install:
	go install ./cmd/hackura/main.go

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)

.PHONY: all build install test clean
