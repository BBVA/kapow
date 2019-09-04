.PHONY: build test coverage install acceptance deps

GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool

BUILD_DIR=./build
OUTPUT_DIR=./output
TMP_DIR=/tmp

BINARY_NAME=kapow

all: test build

build: deps
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

test: build
	$(GOTEST) -race -coverprofile=$(TMP_DIR)/c.out ./...

coverage: test
	mkdir -p $(OUTPUT_DIR)
	$(GOTOOL) cover -html=$(TMP_DIR)/c.out -o $(OUTPUT_DIR)/coverage.html

install: build
	go install ./...

acceptance: install
	make -C ./spec/test

deps:
	@echo "deps here"
