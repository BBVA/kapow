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

all: test race build

build: deps
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

test: build
	$(GOTEST) -v -coverprofile=$(TMP_DIR)/c.out ./...

race: build
	$(GOTEST) -race -v ./...

coverage: test race
	mkdir -p $(OUTPUT_DIR)
	$(GOTOOL) cover -html=$(TMP_DIR)/c.out -o $(OUTPUT_DIR)/coverage.html

install: build
	CGO_ENABLED=0 go install ./...

acceptance: install
	make -C ./spec/test

deps:
	@echo "deps here"
