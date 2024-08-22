##
# This makefile provide rules to start / stop a
# local Grafana service.
#
# Variables:
#   GRAFANA_VERSION
#
# See the container tags into the link below:
#   https://hub.docker.com/r/grafana/grafana/tags
#
# See also the grafana documentation at:
#   https://grafana.com/docs/grafana/latest/
##

GRAFANA_VERSION ?= latest
GRAFANA_UI_PORT ?= 3000

export GRAFANA_UI_PORT

export GRAFANA_VERSION

ifneq (,$(shell command -v open 2>/dev/null))
OPEN ?= open
endif
ifneq (,$(shell command -v xdg-open 2>/dev/null))
OPEN ?= xdg-open
endif
ifeq (,$(OPEN))
OPEN ?= false
endif

.PHONY: grafana-up
grafana-up: ## Start grafana service (local access at http://localhost:9090)
	$(CONTAINER_ENGINE) volume inspect grafana &> /dev/null || $(CONTAINER_ENGINE) volume create grafana
	$(CONTAINER_ENGINE) network inspect idmsvc_backend || $(CONTAINER_ENGINE) network create idmsvc_backend
	$(CONTAINER_ENGINE) container inspect grafana &> /dev/null || \
	$(CONTAINER_ENGINE) run -d \
	  --rm \
	  --name grafana \
	  --volume "grafana:/grafana:z" \
	  --network idmsvc_backend \
	  --publish $(GRAFANA_UI_PORT):3000 \
	  docker.io/grafana/grafana:$(GRAFANA_VERSION)
	./scripts/init_grafana.py
	@@echo "Grafana is running at http://localhost:$(GRAFANA_UI_PORT)"
	@@echo "Default user+password: admin/admin"

.PHONY: grafana-down
grafana-down:  ## Stop grafana service
	! $(CONTAINER_ENGINE) container inspect grafana &> /dev/null || $(CONTAINER_ENGINE) container stop grafana

.PHONY: grafana-clean
grafana-clean: grafana-down  ## Clean the grafana instance
	! $(CONTAINER_ENGINE) container inspect grafana &> /dev/null || $(CONTAINER_ENGINE) container rm grafana
	! $(CONTAINER_ENGINE) volume inspect grafana &> /dev/null || $(CONTAINER_ENGINE) volume rm grafana

.PHONY: grafana-logs
grafana-logs: ## Tail grafana logs
	$(CONTAINER_ENGINE) container logs --tail 10 -f grafana

.PHONY: grafana-ui
grafana-ui:  ## Open browser with the grafana ui
	$(OPEN) http://localhost:$(GRAFANA_UI_PORT)
