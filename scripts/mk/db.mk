##
# Rules related with the database operations
##

.PHONY: install-db-tool
install-db-tool: $(BIN)/db-tool

.PHONY: db-migrate-up
db-migrate-up: $(BIN)/db-tool  ## Migrate the database upto the current state and refresh/create JWKs
	$(BIN)/db-tool migrate up 0
	$(BIN)/db-tool jwk refresh

db-migrate-down: $(BIN)/db-tool  ## Migrate the database down by one step
	$(BIN)/db-tool migrate down 1

.PHONY: db-cli
db-cli:  ## Open a cli shell inside the databse container
	$(CONTAINER_COMPOSE) \
	  -f "$(COMPOSE_FILE)" \
	  -p "$(COMPOSE_PROJECT)" \
	  exec database psql \
	  "sslmode=disable dbname=$(DATABASE_NAME) user=$(DATABASE_USER) host=$(DATABASE_HOST) port=5432 password=$(DATABASE_PASSWORD)"
