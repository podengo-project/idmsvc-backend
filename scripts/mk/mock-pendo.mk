##
# Rules to automate pendo mock tasks
##

# The tag is set manually as it is not expected to generate
# a new image with every change on the repository.
# Once it is generated, the same image will be used over
# and over again.
MOCK_PENDO_CONTAINER ?= quay.io/podengo/mock-pendo:1.0.0
export MOCK_PENDO_CONTAINER

.PHONY: mock-pendo-build
mock-pendo-build: ## Build pendo mock container
	$(MAKE) container-build CONTAINER_IMAGE="$(MOCK_PENDO_CONTAINER)" CONTAINER_CONTEXT_DIR="$(PROJECT_DIR)" CONTAINER_FILE="$(PROJECT_DIR)/build/mock-pendo/Dockerfile"

.PHONY: mock-pendo-up
mock-pendo-up: ## Start pendo mock using local infra
	@[ -e "$(PROJECT_DIR)/configs/config.yaml" ] || { echo "ERROR:Missed configs/config.yaml check README.md file"; exit 1 ; }
	$(CONTAINER_COMPOSE) -p idmsvc -f "$(COMPOSE_FILE)" up -d mock-pendo

.PHONY: mock-pendo-down
mock-pendo-down: ## Stop pendo mock using local infra
	$(CONTAINER_COMPOSE) -p idmsvc -f "$(COMPOSE_FILE)" down mock-pendo
