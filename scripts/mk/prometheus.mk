##
# This makefile provide rules to start / stop a
# local prometheus service.
#
# Variables:
#   PROMETHEUS_VERSION
#   PROMETHEYS_CONFIG
#
# See the container tags into the link below:
#   https://hub.docker.com/r/prom/prometheus/tags
#
# See also the prometheus documentation at:
#   https://prometheus.io/docs/introduction/overview/
##

PROMETHEUS_VERSION ?= v2.40.2
PROMETHEUS_CONFIG ?= $(PROJECT_DIR)/configs/prometheus.yaml
PROMETHEUS_UI_PORT ?= 9090
export PROMETHEUS_UI_PORT
export PROMETHEUS_CONFIG
export PROMETHEUS_VERSION

ifneq (,$(shell command -v open 2>/dev/null))
OPEN ?= open
endif
ifneq (,$(shell command -v xdg-open 2>/dev/null))
OPEN ?= xdg-open
endif
ifeq (,$(OPEN))
OPEN ?= false
endif

.PHONY: prometheus-up
prometheus-up: ## Start prometheus service (local access at http://localhost:9090)
	$(CONTAINER_ENGINE) volume inspect prometheus &> /dev/null || $(CONTAINER_ENGINE) volume create prometheus
	$(CONTAINER_ENGINE) container inspect prometheus &> /dev/null || \
	$(CONTAINER_ENGINE) run -d \
	  --rm \
	  --name prometheus \
	  --network host \
	  --volume "$(PROMETHEUS_CONFIG):/etc/prometheus/prometheus.yml:ro,z" \
	  --volume "prometheus:/prometheus:z" \
	  docker.io/prom/prometheus:$(PROMETHEUS_VERSION)

.PHONY: prometheus-down
prometheus-down:  ## Stop prometheus service
	! $(CONTAINER_ENGINE) container inspect prometheus &> /dev/null || $(CONTAINER_ENGINE) container stop prometheus

.PHONY: prometheus-clean
prometheus-clean: prometheus-down  ## Clean the prometheus instance
	! $(CONTAINER_ENGINE) container inspect prometheus &> /dev/null || $(CONTAINER_ENGINE) container rm prometheus
	! $(CONTAINER_ENGINE) volume inspect prometheus &> /dev/null || $(CONTAINER_ENGINE) volume rm prometheus

.PHONY: prometheus-logs
prometheus-logs: ## Tail prometheus logs
	$(CONTAINER_ENGINE) container logs --tail 10 -f prometheus

.PHONY: prometheus-ui
prometheus-ui:  ## Open browser with the prometheus ui
	$(OPEN) http://localhost:$(PROMETHEUS_UI_PORT)
