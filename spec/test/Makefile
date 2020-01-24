.PHONY: lint wip test fix catalog sync

all: checkbin sync test

sync:
	pipenv sync
lint:
	gherkin-lint
wip:
	KAPOW_DEBUG_TESTS=1 pipenv run behave --stop --wip
test: lint
	pipenv run behave --no-capture --tags=~@skip
fix: lint
	KAPOW_DEBUG_TESTS=1 pipenv run behave --stop --no-capture --tags=~@skip
catalog:
	pipenv run behave --format steps.usage --dry-run --no-summary -q
clean:
	pipenv --rm
checkbin:
	@which kapow >/dev/null || (echo "ERROR: Your kapow binary is not present in PATH" && exit 1)
testpoc: sync
	pipenv run pip install -r ../../testutils/poc/requirements.txt
	PATH=../../testutils/poc:$$PATH KAPOW_CONTROLAPI_URL=http://localhost:8081 KAPOW_DATAAPI_URL=http://localhost:8081 pipenv run behave --no-capture --tags=~@skip
