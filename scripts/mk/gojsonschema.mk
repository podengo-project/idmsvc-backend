##
# This file adds rules to generate code from the
# event message specification
#
# This is based on the repository below:
# https://github.com/RedHatInsights/playbook-dispatcher
#
# https://github.com/atombender/go-jsonschema
##

EVENTS := introspect_request

EVENT_SCHEMA_DIR := $(PROJECT_DIR)/internal/api/event
EVENT_MESSAGE_DIR := $(PROJECT_DIR)/api/event

SCHEMA_YAML_FILES := $(wildcard $(EVENT_MESSAGE_DIR)/*.event.yaml)
SCHEMA_JSON_FILES := $(patsubst $(EVENT_MESSAGE_DIR)/%.event.yaml,$(EVENT_SCHEMA_DIR)/%.event.json,$(wildcard $(EVENT_MESSAGE_DIR)/*.event.yaml))
