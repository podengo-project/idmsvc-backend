# https://github.com/deepmap/oapi-codegen
OAPI_CODEGEN ?= $(BIN)/oapi-codegen
OAPI_CODEGEN_VERSION ?= v1.12.4

$(OAPI_CODEGEN):
	@{\
		export GOPATH="$(shell mktemp -d "$(PROJECT_DIR)/tmp.XXXXXXXX" 2>/dev/null)" ; \
		echo "Using GOPATH='$${GOPATH}'" ; \
		[ "$${GOPATH}" != "" ] || { echo "error:GOPATH is empty"; exit 1; } ; \
		export GOBIN="$(dir $(OAPI_CODEGEN))" ; \
		echo "Installing 'oapi-codegen' at '$(OAPI_CODEGEN)'" ; \
		go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) ; \
		find "$${GOPATH}" -type d -exec chmod u+w {} \; ; \
		rm -rf "$${GOPATH}" ; \
	}

.PHONY: install-oapi-codegen
install-oapi-codegen: $(OAPI_CODEGEN)
