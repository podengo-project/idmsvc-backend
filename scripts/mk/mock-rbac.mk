##
# Rules to automate rbac mock tasks
##

.PHONY: mock-rbac-up
mock-rbac-up: ## Start rbac mock using local infra
	$(CONTAINER_COMPOSE) -p idmsvc -f "$(COMPOSE_FILE)" up -d mock-rbac

.PHONY: mock-rbac-down
mock-rbac-down: ## Stop rbac mock using local infra
	$(CONTAINER_COMPOSE) -p idmsvc -f "$(COMPOSE_FILE)" down mock-rbac
