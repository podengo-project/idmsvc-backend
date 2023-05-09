PYTHON_VENV := .venv
GOLANGCI_LINT ?= $(BIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.46.2

$(PYTHON_VENV):
	python3 -m venv $(PYTHON_VENV)
	$(PYTHON_VENV)/bin/pip install -U pip setuptools

# FIXME fails for the pre-commit hook installed, but launching pre-commit by 'make lint' works
$(PYTHON_VENV)/bin/pre-commit: $(GOLANGCI_LINT) $(PYTHON_VENV)
	$(PYTHON_VENV)/bin/pip3 install pre-commit
	# $(PYTHON_VENV)/bin/pre-commit install --install-hooks --allow-missing-config

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
lint: $(PYTHON_VENV)/bin/pre-commit
	$(PYTHON_VENV)/bin/pre-commit run --all-files --verbose
