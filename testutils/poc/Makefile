.PHONY: sync test

all: test

sync:
	pipenv sync --dev

test: sync
	KAPOW_DATAAPI_URL=http://localhost:8081 pipenv run make -C ../../spec/test

fix:
	KAPOW_DATAAPI_URL=http://localhost:8081 pipenv run make -C ../../spec/test fix

clean:
	pipenv --rm
	make -C ../../spec/test clean
