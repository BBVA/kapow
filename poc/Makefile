.PHONY: sync test

all: test

sync:
	pipenv sync

test: sync
	pipenv run make -C ../spec/test
