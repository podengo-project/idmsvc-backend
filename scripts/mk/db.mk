##
# Rules related with the database operations
##

.PHONY: db-migrate-up
db-migrate-up: $(BIN)/db-tool  ## Migrate the database upto the current state
	$(BIN)/db-tool migrate up 0

.PHONY: db-cli
db-cli:  ## Open a cli shell inside the databse container
	$(DOCKER_COMPOSE) \
	  -f "$(DOCKER_COMPOSE_FILE)" \
	  -p "$(DOCKER_COMPOSE_PROJECT)" \
	  exec database psql \
	  "sslmode=disable dbname=$(DATABASE_NAME) user=$(DATABASE_USER) host=$(DATABASE_HOST) port=5432 password=$(DATABASE_PASSWORD)"
