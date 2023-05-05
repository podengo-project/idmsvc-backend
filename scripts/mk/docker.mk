##
# General rules for interacting with container
# manager (podman or docker).
##

ifneq (,$(shell command podman -v 2>/dev/null))
DOCKER ?= podman
DOCKER_HEALTH_PATH ?= .State.Healthcheck.Status
DOCKER_VOL_SUFFIX ?= :Z
else
ifneq (,$(shell command docker -v 2>/dev/null))
DOCKER ?= docker
DOCKER_HEALTH_PATH ?= .State.Health.Status
DOCKER_VOL_SUFFIX ?=
else
DOCKER ?= false
endif
endif
export DOCKER

DOCKER_LOGIN_USER ?= $(USER)
DOCKER_REGISTRY ?= quay.io
DOCKER_CONTEXT_DIR ?= .
DOCKER_DOCKERFILE ?= build/package/Dockerfile
DOCKER_IMAGE_BASE ?= $(DOCKER_REGISTRY)/$(DOCKER_LOGIN_USER)/myapp
DOCKER_IMAGE_TAG ?= $(shell git rev-parse --short HEAD)
DOCKER_IMAGE ?= $(DOCKER_IMAGE_BASE):$(DOCKER_IMAGE_TAG)
# DOCKER_BUILD_OPTS
# DOCKER_OPTS
# DOCKER_RUN_ARGS

# if go is available, mount user's Go module and build cache to speed up dev builds.
ifneq (,$(shell command go 2>&1 >/dev/null))
USE_GO_CACHE = true
DOCKER_BUILD_OPTS += -v "$(shell go env GOCACHE):/opt/app-root/src/.cache/go-build$(DOCKER_VOL_SUFFIX)"
DOCKER_BUILD_OPTS += -v "$(shell go env GOMODCACHE):/opt/app-root/src/go/pkg/mod$(DOCKER_VOL_SUFFIX)"
else
USE_GO_CACHE = false
endif

.PHONY: docker-login
docker-login:
	$(DOCKER) login -u "$(DOCKER_LOGIN_USER)" -p "$(DOCKER_LOGIN_TOKEN)" $(DOCKER_REGISTRY)

.PHONY: docker-build
docker-build: QUAY_EXPIRATION ?= 1d
docker-build:  ## Build image DOCKER_IMAGE from DOCKER_DOCKERFILE using the DOCKER_CONTEXT_DIR
	$(USE_GO_CACHE) && mkdir -p $(shell go env GOCACHE) $(shell go env GOMODCACHE) || true
	$(DOCKER) build \
	  --label "quay.expires-after=$(QUAY_EXPIRATION)" \
	  $(DOCKER_BUILD_OPTS) \
	  -t "$(DOCKER_IMAGE)" \
	  $(DOCKER_CONTEXT_DIR) \
	  -f "$(DOCKER_DOCKERFILE)"
.PHONY: docker-push
docker-push:  ## Push image to remote registry
	$(DOCKER) push "$(DOCKER_IMAGE)"

# TODO Indicate in the options the IP assigned to the postgres container
# .PHONY: docker-run
# docker-run: DOCKER_OPTS += --env-file .env
# docker-run:  ## Run with DOCKER_OPTS the DOCKER_IMAGE using DOCKER_RUN_ARGS as arguments (eg. make docker-run DOCKER_OPTS="-p 9000:9000")
# 	$(DOCKER) run $(DOCKER_OPTS) $(DOCKER_IMAGE) $(DOCKER_RUN_ARGS)
