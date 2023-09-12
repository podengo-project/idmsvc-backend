GOLANGCI_LINT ?= $(BIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.46.2

.PHONY: install-pre-commit
install-pre-commit: $(PRE_COMMIT)

.PHONY: lint
lint: $(PRE_COMMIT) $(GOLANGCI_LINT)
	$(PYTHON_VENV)/bin/pre-commit run --all-files --verbose
