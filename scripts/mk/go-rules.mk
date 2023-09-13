##
# Golang rules to build the binaries, tidy dependencies,
# generate vendor directory, download dependencies and clean
# the generated binaries.
##

ifeq (,$(shell ls -1d vendor 2>/dev/null))
MOD_VENDOR :=
else
MOD_VENDOR ?= -mod vendor
endif

.PHONY: install-go-tools
install-go-tools: $(TOOLS) ## Install Go tools

.PHONY: install-tools
install-tools: install-go-tools install-python-tools ## Install tools used to build, test and lint

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

$(BIN) $(TOOLS_BIN):
	mkdir -p $@

# export CGO_ENABLED
# $(BIN)/%: CGO_ENABLED=0
$(BIN)/%: cmd/%/main.go $(BIN)
	go build $(MOD_VENDOR) -o "$@" "$<"

# oapi-codegen is installed from global go.mod to keep it in sync with backend code
$(OAPI_CODEGEN): go.mod go.sum $(BIN)
	go build -modfile "$<" -o "$@" "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"

# golangci-lint is very picky when it comes to dependencies, install it directly
# 1.53 is the latest version that supports Go 1.19
$(GOLANGCI_LINT): $(BIN)
	GOBIN="$(dir $(CURDIR)/$@)" go install "github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53"

$(TOOLS_BIN)/%: $(TOOLS_DEPS)
	go build -modfile "$<" -o "$@" $(shell grep $(notdir $@) tools/tools.go | awk '{print $$2}')

.PHONY: clean
clean: ## Clean binaries and testbin generated
	@[ ! -e "$(BIN)" ] || for item in cmd/*; do rm -vf "$(BIN)/$${item##cmd/}"; done

.PHONY: cleanall
cleanall: ## Clean and remove all binaries
	rm -rf $(BIN) $(TOOLS_BIN)

.PHONY: run
run: build ## Run the api & kafka consumer locally
	"$(BIN)/service"

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: get-deps
get-deps: ## Download golang dependencies
	go get -d ./...

.PHONY: update-deps
update-deps: ## Update all golang dependencies
	go get -u -t ./...
	@# kin-openapi 0.119.0 requires Golang 1.20
	go get github.com/getkin/kin-openapi@v0.118.0
	$(MAKE) tidy

.PHONY: vet
vet:  ## Run go vet ignoring /vendor directory
	go vet $(shell go list ./... | grep -v /vendor/)

.PHONY: go-fmt
go-fmt:  ## Run go fmt ignoring /vendor directory
	go fmt $(shell go list ./... | grep -v /vendor/)

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
API_LIST := api/public.openapi.yaml api/internal.openapi.yaml api/metrics.openapi.yaml
.PHONY: generate-api
generate-api: $(OAPI_CODEGEN) $(API_LIST) ## Generate server stubs from openapi
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

$(API_LIST):
	git submodule update --init

.PHONY: update-api
update-api:
	git submodule update --init --remote
	$(MAKE) generate-api

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
    internal/api/openapi \
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
	$(GODA) graph "github.com/podengo-project/idmsvc-backend/..." | dot -Tsvg -o docs/service-dependencies.svg

.PHONY: coverage
coverage:  ## Printout coverage
	go tool cover -func ./coverage.out
