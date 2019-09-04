.PHONY: test install acceptance deps

all: acceptance

test: deps
	go test -race -coverprofile=/tmp/c.out github.com/BBVA/kapow/pkg/...
	go tool cover -html=/tmp/c.out -o coverage.html

install: test
	go install github.com/BBVA/kapow/...

acceptance: install
	pipenv run make -C spec/test

deps:
	go install github.com/spf13/cobra
