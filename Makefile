.PHONY: lint build test jaillover race coverage install acceptance deps docker clean

GOCMD=go
GOBUILD=$(GOCMD) build -trimpath
GOINSTALL=$(GOCMD) install -trimpath
GOGET=$(GOCMD) get
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOLANGLINT=golangci-lint
PROJECTREPO=github.com/BBVA/kapow

BUILD_DIR=./build
OUTPUT_DIR=./output
TMP_DIR=/tmp
DOCS_DIR=./doc
DOCKER_DIR=./docker

BINARY_NAME=kapow

all: lint test race build

lint:
	$(GOLANGLINT) run

build: deps
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

test: build jaillover
	$(GOTEST) -v -coverprofile=$(TMP_DIR)/c.out ./...

jaillover:
	$(GOINSTALL) -v $(PROJECTREPO)/testutils/$@

race: build
	$(GOTEST) -race -v ./...

coverage: test race
	mkdir -p $(OUTPUT_DIR)
	$(GOTOOL) cover -html=$(TMP_DIR)/c.out -o $(OUTPUT_DIR)/coverage.html

install: build
	CGO_ENABLED=0 $(GOINSTALL) ./...

acceptance: build
	cd ./spec/test && PATH=$(PWD)/build:$$PATH  nix-shell --command make

deps:
	@echo "deps here"

clean:
	rm -rf	 $(BUILD_DIR) $(OUTPUT_DIR) $(DOCKER_DIR)/*
