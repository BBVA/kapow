wip:
	KAPOW_DEBUG_TESTS=1 pipenv run behave --stop --wip
test:
	gherkin-lint
	pipenv run behave --no-capture
fix:
	gherkin-lint
	KAPOW_DEBUG_TESTS=1 pipenv run behave --stop --no-capture
catalog:
	pipenv run behave --steps-catalog
