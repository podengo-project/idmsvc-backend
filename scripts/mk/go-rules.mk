##
# Golang rules to build the binaries, tidy dependencies,
# generate vendor directory, download dependencies and clean
# the generated binaries.
##

CONFIG_PATH ?= $(PROJECT_DIR)/configs/config.yaml
export CONFIG_PATH

# Directory where the built binaries will be generated
BIN ?= $(PROJECT_DIR)/bin
PATH := $(BIN):$(PATH)
export PATH

ifeq (,$(shell ls -1d vendor 2>/dev/null))
MOD_VENDOR :=
else
MOD_VENDOR ?= -mod vendor
endif

.PHONY: install-tools
install-tools: ## Install tools used to build, test and lint
	$(MAKE) install-oapi-codegen
	$(MAKE) install-mockery
	$(MAKE) install-golangci-lint
	$(MAKE) install-pre-commit
	$(MAKE) install-gojsonschema
	$(MAKE) install-goda
	$(MAKE) install-yq
	$(MAKE) install-planter
	$(MAKE) .venv
	source .venv/bin/activate && pip install -U pip && pip install -r requirements-dev.txt

.PHONY: build-all
build-all: ## Generate code and build binaries
	$(MAKE) generate-api
	$(MAKE) generate-event
	$(MAKE) generate-mock
	$(MAKE) generate-diagrams
	$(MAKE) build

# Meta rule to add dependency on the binaries generated
.PHONY: build
build: $(patsubst cmd/%,$(BIN)/%,$(wildcard cmd/*)) ## Build binaries

# export CGO_ENABLED
# $(BIN)/%: CGO_ENABLED=0
$(BIN)/%: cmd/%/main.go
	@[ -e "$(BIN)" ] || mkdir -p "$(BIN)"
	go build $(MOD_VENDOR) -o "$@" "$<"

.PHONY: clean
clean: ## Clean binaries and testbin generated
	@[ ! -e "$(BIN)" ] || for item in cmd/*; do rm -vf "$(BIN)/$${item##cmd/}"; done
#	@[ ! -e testbin ] || rm -rf testbin

.PHONY: run
run: build ## Run the api & kafka consumer locally
	"$(BIN)/service"

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: get-deps
get-deps: ## Download golang dependencies
	go get -d ./...

.PHONY: vendor
vendor: ## Generate vendor/ directory populated with the dependencies
	go mod vendor

# Exclude /internal/test/mock /internal/api directories because the content is
# generated.
# Exclude /vendor in case it exists
# Exclude /internal/interface directories because only contain interfaces
TEST_GREP_FILTER := -v \
  -e /vendor/ \
  -e /internal/test/mock \
  -e /internal/interface/ \
  -e /internal/api/metrics \
  -e /internal/api/private \
  -e /internal/api/public

.PHONY: test
test: ## Run tests
	go test -coverprofile="coverage.out" -covermode count $(MOD_VENDOR) $(shell go list ./... | grep $(TEST_GREP_FILTER) )

.PHONY: test-ci
test-ci: ## Run tests for ci
	go test $(MOD_VENDOR) ./...

# Add dependencies from binaries to all the the sources
# so any change is detected for the build rule
$(patsubst cmd/%,$(BIN)/%,$(wildcard cmd/*)): $(shell find $(PROJECT_DIR)/cmd -type f -name '*.go') $(shell find $(PROJECT_DIR)/pkg -type f -name '*.go' 2>/dev/null) $(shell find $(PROJECT_DIR)/internal -type f -name '*.go' 2>/dev/null)

# # Regenerate code when message schema changes
# $(shell find "$(EVENT_MESSAGE_DIR)" -type f -name '*.go'): $(SCHEMA_YAML_FILES)
# 	$(MAKE) gen-event-messages


############### TOOLS

# https://github.com/RedHatInsights/playbook-dispatcher/blob/master/Makefile
.PHONY: generate-api
generate-api: $(OAPI_CODEGEN)  ## Generate server stubs from openapi
	# Public API
	$(OAPI_CODEGEN) -generate spec -package public -o internal/api/public/spec.gen.go api/public.openapi.yaml
	$(OAPI_CODEGEN) -generate server -package public -o internal/api/public/server.gen.go api/public.openapi.yaml
	$(OAPI_CODEGEN) -generate types -package public -o internal/api/public/types.gen.go -alias-types api/public.openapi.yaml
	# Internal API # FIXME Update -import-mapping options
	$(OAPI_CODEGEN) -generate spec -package private -o internal/api/private/spec.gen.go api/internal.openapi.yaml
	$(OAPI_CODEGEN) -generate server -package private -o internal/api/private/server.gen.go api/internal.openapi.yaml
	$(OAPI_CODEGEN) -generate types -package private -o internal/api/private/types.gen.go api/internal.openapi.yaml
	# metrics API
	$(OAPI_CODEGEN) -generate spec -package metrics -o internal/api/metrics/spec.gen.go api/metrics.openapi.yaml
	$(OAPI_CODEGEN) -generate server -package metrics -o internal/api/metrics/server.gen.go api/metrics.openapi.yaml
	$(OAPI_CODEGEN) -generate types -package metrics -o internal/api/metrics/types.gen.go api/metrics.openapi.yaml

EVENTS := todo_created
# Generate event types
.PHONY: generate-event
generate-event: $(GOJSONSCHEMA) $(SCHEMA_JSON_FILES)  ## Generate event messages from schemas
	@[ -e "$(EVENT_MESSAGE_DIR)" ] || mkdir -p "$(EVENT_MESSAGE_DIR)"
	for event in $(EVENTS); do \
		$(GOJSONSCHEMA) -p event "$(EVENT_SCHEMA_DIR)/$${event}.event.json" -o "$(PROJECT_DIR)/internal/api/event/$${event}.event.types.gen.go"; \
	done

.PHONY: generate-event-debug
generate-event-debug:
	@echo SCHEMA_JSON_FILES=$(SCHEMA_JSON_FILES)
	@echo GOJSONSCHEMA=$(GOJSONSCHEMA)
	@echo EVENT_MESSAGE_DIR=$(EVENT_MESSAGE_DIR)
	@echo EVENT_SCHEMA_DIR=$(EVENT_SCHEMA_DIR)

# Generic rule to generate the JSON files
$(EVENT_SCHEMA_DIR)/%.event.json: $(EVENT_MESSAGE_DIR)/%.event.yaml
	@[ -e "$(EVENT_MESSAGE_DIR)" ] || mkdir -p "$(EVENT_MESSAGE_DIR)"
	yaml2json "$<" "$@"

# Mockery support
MOCK_DIRS := internal/api/private \
    internal/api/public \
	internal/interface/repository \
	internal/interface/interactor \
	internal/interface/presenter \
	internal/interface/event \
	internal/interface/client \
	internal/handler \
	internal/infrastructure/service \
	internal/infrastructure/event \
	internal/infrastructure/event/handler \
	internal/infrastructure/middleware \

.PHONY: generate-mock
generate-mock: $(MOCKERY)  ## Generate mock by using mockery tool
	for item in $(MOCK_DIRS); do \
	  PKG="$${item##*/}"; \
	  DEST_DIR="internal/test/mock/$${item#*/}"; \
	  [ -e "$${DEST_DIR}" ] || mkdir -p "$${DEST_DIR}"; \
	  $(MOCKERY) \
	    --all \
	    --outpkg "$${PKG}" \
	    --dir "$${item}" \
		--output "$${DEST_DIR}" \
		--case underscore; \
	done

.PHONY: generate-deps
generate-deps: $(GODA)
	$(GODA) graph "github.com/hmsidm/..." | dot -Tsvg -o docs/service-dependencies.svg
