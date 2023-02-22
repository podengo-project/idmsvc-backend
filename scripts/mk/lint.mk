ADD_PYTHON_ENV := source .venv/bin/activate &&
GOLANGCI_LINT ?= $(BIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.46.2

.venv:
	python3 -m venv .venv && $(ADD_PYTHON_ENV) pip3 install -U pip

# FIXME fails for the pre-commit hook installed, but launching pre-commit by 'make lint' works
.PHONY: install-pre-commit
install-pre-commit: install-golangci-lint .venv
	$(ADD_PYTHON_ENV) pip3 install pre-commit
	# $(ADD_PYTHON_ENV) pre-commit install --install-hooks --allow-missing-config

.PHONY: install-golangci-lint
install-golangci-lint: $(BIN)/golangci-lint
# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN) $(GOLANGCI_LINT_VERSION)

$(BIN)/golangci-lint:
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
lint: .venv ## Run linters
	$(ADD_PYTHON_ENV) pre-commit run --all-files --verbose
