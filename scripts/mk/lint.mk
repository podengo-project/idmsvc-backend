GOLANGCI_LINT ?= $(BIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.46.2

.PHONY: install-pre-commit
install-pre-commit: $(PRE_COMMIT)

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT)

# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN) $(GOLANGCI_LINT_VERSION)
$(GOLANGCI_LINT): $(BIN)
	GOBIN="$(dir $(CURDIR)/$@)" go install "github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)"

.PHONY: lint
lint: $(PRE_COMMIT) $(GOLANGCI_LINT)
	$(PYTHON_VENV)/bin/pre-commit run --all-files --verbose
