BINARY_NAME=eth2-keystore-converter

all: build

install:
	go mod download

build: install
	go build -v -o bin/$(BINARY_NAME) main.go

clean:
	rm -rf bin

fclean: clean

.PHONY: all build generate test clean
