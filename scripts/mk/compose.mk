##
# Rules for local infrastructure using docker-compose or podman-compose
##

ifeq (podman,$(DOCKER))
DOCKER_COMPOSE ?= podman-compose
endif

ifeq (docker,$(DOCKER))
DOCKER_COMPOSE ?= docker-compose
endif

DOCKER_COMPOSE_FILE ?= $(PROJECT_DIR)/deployments/docker-compose.yaml
DOCKER_COMPOSE_PROJECT ?= hmsidm

DOCKER_COMPOSE_VARS_DATABASE=\
	DATABASE_USER="$(DATABASE_USER)" \
	DATABASE_PASSWORD="$(DATABASE_PASSWORD)" \
	DATABASE_NAME="$(DATABASE_NAME)" \
	DATABASE_EXTERNAL_PORT="$(DATABASE_EXTERNAL_PORT)"

DOCKER_COMPOSE_VARS_KAFKA=\
	KAFKA_CONFIG_DIR=$(KAFKA_CONFIG_DIR) \
	KAFKA_DATA_DIR=$(KAFKA_DATA_DIR) \
	ZOOKEEPER_CLIENT_PORT=$(ZOOKEEPER_CLIENT_PORT) \
	KAFKA_TOPICS=$(KAFKA_TOPICS)

DOCKER_COMPOSE_VARS= \
	DOCKER_IMAGE_BASE="$(DOCKER_IMAGE_BASE)" \
	DOCKER_IMAGE_TAG="$(DOCKER_IMAGE_TAG)" \
	$(DOCKER_COMPOSE_VARS_DATABASE) \
	$(DOCKER_COMPOSE_VARS_KAFKA)


.PHONY: compose-up
compose-up: ## Start local infrastructure
	$(DOCKER_COMPOSE_VARS) \
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) -p $(DOCKER_COMPOSE_PROJECT) up -d
	$(MAKE) db-migrate-up

.PHONY: compose-down
compose-down: ## Stop local infrastructure
	$(DOCKER_COMPOSE_VARS) \
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) -p $(DOCKER_COMPOSE_PROJECT) down --volumes

.PHONY: compose-build
compose-build: ## Build the images at docker-compose.yaml
	$(DOCKER_COMPOSE_VARS) \
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) -p $(DOCKER_COMPOSE_PROJECT) build

.PHONY: compose-logs
compose-logs: ## Print out infrastructure logs
	$(DOCKER_COMPOSE_VARS) \
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) -p $(DOCKER_COMPOSE_PROJECT) logs

.PHONY: compose-clean
compose-clean: compose-down  ## Stop and clean local infrastructure
	$(DOCKER) volume prune --force
