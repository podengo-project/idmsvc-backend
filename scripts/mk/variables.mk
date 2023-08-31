##
# Set default variable values for the project
##
APP ?= idmsvc
export APP

BIN ?= bin
PATH := $(CURDIR)/$(BIN):$(PATH)
export PATH

CONFIG_PATH ?= $(PROJECT_DIR)/configs
export CONFIG_PATH
CONFIG_YAML := $(CONFIG_PATH)/config.yaml

COMPOSE_FILE ?= $(PROJECT_DIR)/deployments/docker-compose.yaml

CONTAINER_IMAGE_BASE ?= quay.io/$(firstword $(subst +, ,$(QUAY_USER)))/$(APP_NAME)-$(APP_COMPONENT)

LOAD_DB_CFG_WITH_YQ := n
ifneq (,$(shell "$(BIN)/yq" --version 2>/dev/null))
ifneq (,$(shell ls -1 "$(CONFIG_YAML)" 2>/dev/null))
LOAD_DB_CFG_WITH_YQ := y
endif
endif

#
# Kafka configuration variables
#

DATABASE_CONTAINER_NAME="database"
ifeq (y,$(LOAD_DB_CFG_WITH_YQ))
$(info info:Trying to load DATABASE configuration from '$(CONFIG_YAML)')
DATABASE_HOST ?= $(shell "$(BIN)/yq" -r -M '.database.host' "$(CONFIG_YAML)")
DATABASE_EXTERNAL_PORT ?= $(shell "$(BIN)/yq" -M '.database.port' "$(CONFIG_YAML)")
DATABASE_NAME ?= $(shell "$(BIN)/yq" -r -M '.database.name' "$(CONFIG_YAML)")
DATABASE_USER ?= $(shell "$(BIN)/yq" -r -M '.database.user' "$(CONFIG_YAML)")
DATABASE_PASSWORD ?= $(shell "$(BIN)/yq" -r -M '.database.password' "$(CONFIG_YAML)")
else
$(info info:Using DATABASE_* defaults)
DATABASE_HOST ?= localhost
DATABASE_EXTERNAL_PORT ?= 5432
DATABASE_NAME ?= idmsvc-db
DATABASE_USER ?= idmsvc-user
DATABASE_PASSWORD ?= idmsvc-secret
endif


#
# Kafka configuration variables
#

# The directory where the kafka data will be stored
KAFKA_DATA_DIR ?= $(PROJECT_DIR)/build/kafka/data

# The directory where the kafka configuration will be
# bound to the containers
KAFKA_CONFIG_DIR ?= $(PROJECT_DIR)/build/kafka/config

# The topics used by the repository
# Updated to follow the pattern used at playbook-dispatcher
KAFKA_TOPICS ?= platform.idmsvc.todo-created

# The group id for the consumers; every consumer subscribed to
# a topic with different group-id will receive a copy of the
# message. In our scenario, any replica of the consumer wants
# only one message to be processed, so we only use a unique
# group id at the moment.
KAFKA_GROUP_ID ?= idmsvc

# Application specific parameters
APP_EXPIRATION_TIME ?= 15
export APP_EXPIRATION_TIME
APP_PAGINATION_DEFAULT_LIMIT ?= 10
export APP_PAGINATION_DEFAULT_LIMIT
APP_PAGINATION_MAX_LIMIT ?= 100
export APP_PAGINATION_MAX_LIMIT
# Enable IS_FAKE_ENABLED for the ephemeral deployment
APP_ACCEPT_X_RH_FAKE_IDENTITY ?= true
export APP_ACCEPT_X_RH_FAKE_IDENTITY
APP_VALIDATE_API ?= true
export APP_VALIDATE_API

# Set the default token expiration in seconds (2 hours)
APP_TOKEN_EXPIRATION_SECONDS ?= 7200
export APP_TOKEN_EXPIRATION_SECONDS

# Use random MAC key for domain registration token
APP_DOMAIN_REG_KEY ?= random
export APP_DOMAIN_REG_KEY
