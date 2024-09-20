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

# used by ipa-hcc backend test
.PHONY: install-xrhidgen
install-xrhidgen: $(XRHIDGEN) ## Install xrhidgen tool

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
# Build by path, not by referring to main.go. It's required to bake VCS
# information into binary, see golang/go#51279.
$(BIN)/%: cmd/%/main.go $(BIN)
	go build -C $(dir $<) $(MOD_VENDOR) -buildvcs=true -o "$(CURDIR)/$@" .

$(TOOLS_BIN)/%: $(TOOLS_DEPS)
	cd tools && GOBIN="$(PROJECT_DIR)/tools/bin" go install $(shell grep $(notdir $@) tools/tools.go | awk '{print $$2}')

.PHONY: clean
clean: ## Clean binaries and testbin generated
	@[ ! -e "$(BIN)" ] || for item in cmd/*; do rm -vf "$(BIN)/$${item##cmd/}"; done

.PHONY: cleanall
cleanall: ## Clean and remove all binaries
	rm -rf $(BIN) $(TOOLS_BIN)

.PHONY: run
run: $(BIN)/service .compose-wait-db ## Run the api & kafka consumer locally
	$(MAKE) mock-rbac-up
	"$(BIN)/service"

# See: https://go.dev/doc/modules/managing-dependencies#synchronizing
.PHONY: tidy
tidy:  ## Synchronize your code's dependencies
	go mod tidy -go=$(GO_VERSION)
	cd tools && go mod tidy -go=$(GO_VERSION)

.PHONY: get-deps
get-deps: ## Download golang dependencies
	go mod download

.PHONY: update-deps
update-deps: ## Update all golang dependencies
	go get -u -t ./...
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
  -e /internal/test \
  -e /internal/interface/ \
  -e /internal/api/metrics \
  -e /internal/api/private \
  -e /internal/api/public

.PHONY: test
test: ## Run unit tests, smoke tests and integration tests
	$(MAKE) test-unit test-smoke

.PHONY: test-unit
test-unit: ## Run unit tests
	CLIENTS_RBAC_BASE_URL="http://localhost:8021/api/rbac/v1" \
	go test -parallel 4 -coverprofile="coverage.out" -covermode count $(MOD_VENDOR) $(shell go list ./... | grep $(TEST_GREP_FILTER) )

.PHONY: test-ci
test-ci: ## Run tests for ci
	CLIENTS_RBAC_BASE_URL="http://localhost:8021/api/rbac/v1" \
	go test $(MOD_VENDOR) ./...

.PHONY: test-smoke
test-smoke:  ## Run smoke tests
	CLIENTS_RBAC_BASE_URL="http://localhost:8021/api/rbac/v1" \
	go test -parallel 1 ./internal/test/smoke/... -test.failfast -test.v

.PHONY: test-perf
test-perf:  ## Run smoke tests
	CLIENTS_RBAC_BASE_URL="http://localhost:8021/api/rbac/v1" \
	go test -parallel 1 ./internal/test/perf/... -test.failfast -test.v -timeout 30m

# Add dependencies from binaries to all the the sources
# so any change is detected for the build rule
$(patsubst cmd/%,$(BIN)/%,$(wildcard cmd/*)): $(shell find $(PROJECT_DIR)/cmd -type f -name '*.go') $(shell find $(PROJECT_DIR)/pkg -type f -name '*.go' 2>/dev/null) $(shell find $(PROJECT_DIR)/internal -type f -name '*.go' 2>/dev/null)

# # Regenerate code when message schema changes
# $(shell find "$(EVENT_MESSAGE_DIR)" -type f -name '*.go'): $(SCHEMA_YAML_FILES)
# 	$(MAKE) gen-event-messages

.PHONY: pr-check
pr-check: ## Run common checks before submitting a PR
	$(MAKE) go-fmt
	$(MAKE) tidy vet
	$(MAKE) generate-api generate-mock install-go-tools
	$(MAKE) build test-unit

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

# Thanks to RHEnvision
# See: https://github.com/RHEnVision/provisioning-backend/blob/main/mk/clients.mk
# Generate HTTP client stubs
CLIENT_LIST := internal/usecase/client/rbac/client.gen.go
# CLIENT_LIST += pkg/public/client.gen.go

.PHONY: generate-client
generate-client: $(OAPI_CODEGEN) $(CLIENT_LIST)  ## Generate client stubs from openapi

configs/rbac_client_api.json:
	curl -s -o ./configs/rbac_client_api.json -z ./configs/rbac_client_api.json https://raw.githubusercontent.com/RedHatInsights/insights-rbac/master/docs/source/specs/openapi.json || true

internal/usecase/client/rbac/client.gen.go: configs/rbac_client_gen_config.yaml configs/rbac_client_api.json
	$(OAPI_CODEGEN) -config configs/rbac_client_gen_config.yaml configs/rbac_client_api.json

# pkg/public/client.gen.go: configs/idmsvc_client_gen_config.yaml api/public.openapi.json
# 	$(OAPI_CODEGEN) -config configs/idmsvc_client_gen_config.yaml api/public.openapi.json

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
	internal/api/metrics \
	internal/interface/repository \
	internal/interface/interactor \
	internal/interface/presenter \
	internal/interface/event \
	internal/interface/client/inventory \
	internal/interface/client/rbac \
	internal/interface/client/pendo \
	internal/handler \
	internal/infrastructure/service \
	internal/infrastructure/event \
	internal/infrastructure/event/handler \
	internal/infrastructure/logger \
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
