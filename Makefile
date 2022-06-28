.PHONY: all build link tidy

GO=$(shell which go)
CMD_DIR=./cmd
BIN_DIR=./bin

all: tidy build link

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/bip39gen $(CMD_DIR)

link: build
	@ln -f -s $(BIN_DIR)/bip39gen ~/bin/bip39gen

tidy:
	@go mod tidy