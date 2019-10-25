.PHONY: lint wip test fix catalog sync

all: sync test

sync:
	pipenv sync
lint:
	gherkin-lint
wip:
	KAPOW_DEBUG_TESTS=1 pipenv run behave --stop --wip
test: lint
	pipenv run behave --no-capture --tags=~@skip
fix: lint
	KAPOW_DEBUG_TESTS=1 pipenv run behave --stop --no-capture
catalog:
	pipenv run behave --format steps.usage --dry-run --no-summary -q
