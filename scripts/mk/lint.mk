.PHONY: install-pre-commit
install-pre-commit: $(PRE_COMMIT)

.PHONY: lint
lint: $(PRE_COMMIT) $(GOLANGCI_LINT)
	PATH=$(TOOLS_BIN):$${PATH} $(PYTHON_VENV)/bin/pre-commit run --all-files --verbose
