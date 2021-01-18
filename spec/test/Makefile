.PHONY: all lint wip test fix catalog

all: checkbin test

lint:
	gherkin-lint
wip:
	KAPOW_DEBUG_TESTS=1 behave --stop --wip
test: lint
	behave --no-capture --tags=~@skip
fix: lint
	KAPOW_DEBUG_TESTS=1 behave --stop --no-capture --tags=~@skip
catalog:
	behave --format steps.usage --dry-run --no-summary -q
checkbin:
	@which kapow >/dev/null || (echo "ERROR: Your kapow binary is not present in PATH" && exit 1)
testpoc:
	PATH=../../testutils/poc:$$PATH behave --no-capture --tags=~@skip
