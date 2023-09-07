# https://github.com/deepmap/oapi-codegen
OAPI_CODEGEN ?= $(BIN)/oapi-codegen
OAPI_CODEGEN_VERSION ?= v1.14.0

$(OAPI_CODEGEN): $(BIN)
	GOBIN="$(dir $(CURDIR)/$@)" go install "github.com/deepmap/oapi-codegen/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)"

.PHONY: install-oapi-codegen
install-oapi-codegen: $(OAPI_CODEGEN)

.PHONY: openapi-sort
openapi-sort: $(PYTHON_VENV)  ## sort and lint OpenAPI YAML files
	$(PYTHON_VENV)/bin/python $(PWD)/scripts/yamlsort.py $(PWD)/api/*.yaml
