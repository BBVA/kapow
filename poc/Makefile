.PHONY: sync test

all: test

sync:
	pipenv sync --dev

test: sync
	pipenv run make -C ../spec/test
