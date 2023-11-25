##
# Rules for local infrastructure using docker-compose or podman-compose
##

ifeq (podman,$(CONTAINER_ENGINE))
CONTAINER_COMPOSE ?= podman-compose
CONTAINER_DATABASE_NAME ?= $(COMPOSE_PROJECT)_database_1
else
CONTAINER_COMPOSE ?= docker-compose
CONTAINER_DATABASE_NAME ?= $(COMPOSE_PROJECT)-database-1
endif

COMPOSE_FILE ?= $(PROJECT_DIR)/deployments/docker-compose.yaml
COMPOSE_PROJECT ?= idmsvc

COMPOSE_VARS_DATABASE=\
	DATABASE_USER="$(DATABASE_USER)" \
	DATABASE_PASSWORD="$(DATABASE_PASSWORD)" \
	DATABASE_NAME="$(DATABASE_NAME)" \
	DATABASE_EXTERNAL_PORT="$(DATABASE_EXTERNAL_PORT)"

COMPOSE_VARS_DB_MIGRATE=\
	DATABASE_USER="$(DATABASE_USER)" \
	DATABASE_PASSWORD="$(DATABASE_PASSWORD)" \
	DATABASE_NAME="$(DATABASE_NAME)" \
	DATABASE_HOST="$(DATABASE_HOST)" \
	DATABASE_PORT="$(DATABASE_EXTERNAL_PORT)"

COMPOSE_VARS_KAFKA=\
	KAFKA_CONFIG_DIR=$(KAFKA_CONFIG_DIR) \
	KAFKA_DATA_DIR=$(KAFKA_DATA_DIR) \
	ZOOKEEPER_CLIENT_PORT=$(ZOOKEEPER_CLIENT_PORT) \
	KAFKA_TOPICS=$(KAFKA_TOPICS)

COMPOSE_VARS_APP=\
    APP_SECRET="$(APP_SECRET)" \
    APP_VALIDATE_API=$(APP_VALIDATE_API) \
    APP_TOKEN_EXPIRATION_SECONDS=$(APP_TOKEN_EXPIRATION_SECONDS)

COMPOSE_VARS= \
	CONTAINER_IMAGE_BASE="$(CONTAINER_IMAGE_BASE)" \
	CONTAINER_IMAGE_TAG="$(CONTAINER_IMAGE_TAG)" \
	$(COMPOSE_VARS_DATABASE) \
	$(COMPOSE_VARS_KAFKA) \
	$(COMPOSE_VARS_APP)


.PHONY: compose-up
compose-up: ## Start local infrastructure
	$(COMPOSE_VARS) \
	    $(CONTAINER_COMPOSE) -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT) up -d
	$(MAKE) .compose-wait-db
	$(MAKE) $(COMPOSE_VARS_DB_MIGRATE) db-migrate-up

.PHONY: .compose-wait-db
.compose-wait-db:
	@printf "Waiting database"; \
	while [ "$$( $(CONTAINER_ENGINE) container inspect --format '{{$(CONTAINER_HEALTH_PATH)}}' "$(CONTAINER_DATABASE_NAME)" )" != "healthy" ]; \
	do sleep 1; printf "."; \
	done; \
	printf "\n"

.PHONY: compose-down
compose-down: ## Stop local infrastructure
	$(COMPOSE_VARS) \
	    $(CONTAINER_COMPOSE) -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT) down --volumes

.PHONY: compose-build
compose-build: ## Build the images at docker-compose.yaml
	$(COMPOSE_VARS) \
	    $(CONTAINER_COMPOSE) -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT) build

.PHONY: compose-pull
compose-pull: ## Pull images
	$(COMPOSE_VARS) \
	    $(CONTAINER_COMPOSE) -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT) pull

.PHONY: compose-logs
compose-logs: ## Print out infrastructure logs
	$(COMPOSE_VARS) \
	    $(CONTAINER_COMPOSE) -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT) logs

.PHONY: compose-clean
compose-clean: compose-down  ## Stop and clean local infrastructure
	$(CONTAINER_ENGINE) volume prune --force
