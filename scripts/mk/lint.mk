GOLANGCI_LINT ?= $(BIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.46.2

.PHONY: install-pre-commit
install-pre-commit: $(PRE_COMMIT)

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT)

# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN) $(GOLANGCI_LINT_VERSION)
$(GOLANGCI_LINT):
	@{\
		export GOPATH="$(shell mktemp -d "$(PROJECT_DIR)/tmp.XXXXXXXX" 2>/dev/null)" ; \
		echo "Using GOPATH='$${GOPATH}'" ; \
		[ "$${GOPATH}" != "" ] || { echo "error:GOPATH is empty"; exit 1; } ; \
		export GOBIN="$(dir $(GOLANGCI_LINT))" ; \
		echo "Installing 'golangci-lint' at '$(GOLANGCI_LINT)'" ; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) ; \
		find "$${GOPATH}" -type d -exec chmod u+w {} \; ; \
		rm -rf "$${GOPATH}" ; \
	}

.PHONY: lint
lint: $(PRE_COMMIT) $(GOLANGCI_LINT)
	$(PYTHON_VENV)/bin/pre-commit run --all-files --verbose
