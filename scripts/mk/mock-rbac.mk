##
# Rules to automate rbac mock tasks
##

# The tag is set manually as it is not expected to generate
# a new image with every change on the repository.
# Once it is generated, the same image will be used over
# and over again.
MOCK_RBAC_CONTAINER ?= quay.io/podengo/mock-rbac:1.0.0
export MOCK_RBAC_CONTAINER

.PHONY: mock-rbac-build
mock-rbac-build: ## Build rbac mock container
	$(MAKE) container-build CONTAINER_IMAGE="$(MOCK_RBAC_CONTAINER)" CONTAINER_CONTEXT_DIR="$(PROJECT_DIR)" CONTAINER_FILE="$(PROJECT_DIR)/build/mock-rbac/Dockerfile"

.PHONY: mock-rbac-up
mock-rbac-up: ## Start rbac mock using local infra
	$(CONTAINER_COMPOSE) -p idmsvc -f "$(COMPOSE_FILE)" up -d mock-rbac

.PHONY: mock-rbac-down
mock-rbac-down: ## Stop rbac mock using local infra
	$(CONTAINER_COMPOSE) -p idmsvc -f "$(COMPOSE_FILE)" down mock-rbac
