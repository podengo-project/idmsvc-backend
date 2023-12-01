##
# General rules for interacting with container
# manager (podman or docker).
##

ifneq (,$(shell command podman -v 2>/dev/null))
CONTAINER_ENGINE ?= podman
else
ifneq (,$(shell command docker -v 2>/dev/null))
CONTAINER_ENGINE ?= docker
else
CONTAINER_ENGINE ?= false
endif
endif
export CONTAINER_ENGINE

CONTAINER_HEALTH_PATH ?= .State.Health.Status

ifneq (,$shell(selinuxenabled 2>/dev/null))
CONTAINER_VOL_SUFFIX ?= :Z
else
CONTAINER_VOL_SUFFIX ?=
endif

CONTAINER_REGISTRY_USER ?= $(USER)
CONTAINER_REGISTRY ?= quay.io
CONTAINER_CONTEXT_DIR ?= .
CONTAINERFILE ?= build/package/Dockerfile
CONTAINER_IMAGE_BASE ?= $(CONTAINER_REGISTRY)/$(CONTAINER_REGISTRY_USER)/myapp
CONTAINER_IMAGE_TAG ?= $(shell git rev-parse --short HEAD)
CONTAINER_IMAGE ?= $(CONTAINER_IMAGE_BASE):$(CONTAINER_IMAGE_TAG)
# CONTAINER_BUILD_OPTS
# CONTAINER_ENGINE_OPTS
# CONTAINER_RUN_ARGS

# if go is available, mount user's Go module and build cache to speed up dev builds.
ifneq (,$(shell command go 2>&1 >/dev/null))
USE_GO_CACHE = true
CONTAINER_BUILD_OPTS += -v "$(shell go env GOCACHE):/opt/app-root/src/.cache/go-build$(CONTAINER_VOL_SUFFIX)"
CONTAINER_BUILD_OPTS += -v "$(shell go env GOMODCACHE):/opt/app-root/src/go/pkg/mod$(CONTAINER_VOL_SUFFIX)"
else
USE_GO_CACHE = false
endif

.PHONY: registry-login
registry-login:
	$(CONTAINER_ENGINE) login -u "$(CONTAINER_REGISTRY_USER)" -p "$(CONTAINER_REGISTRY_TOKEN)" $(CONTAINER_REGISTRY)

.PHONY: container-build
container-build: QUAY_EXPIRATION ?= never
container-build:  ## Build image CONTAINER_IMAGE from CONTAINERFILE using the CONTAINER_CONTEXT_DIR
	$(USE_GO_CACHE) && mkdir -p $(shell go env GOCACHE) $(shell go env GOMODCACHE) || true
	$(CONTAINER_ENGINE) build \
	  --label "quay.expires-after=$(QUAY_EXPIRATION)" \
	  $(CONTAINER_BUILD_OPTS) \
	  -t "$(CONTAINER_IMAGE)" \
	  $(CONTAINER_CONTEXT_DIR) \
	  -f "$(CONTAINERFILE)"
	@# prune builder container
	$(CONTAINER_ENGINE) image prune --filter label=idmsvc-backend=builder --force

.PHONY: container-push
container-push:  ## Push image to remote registry
	$(CONTAINER_ENGINE) push "$(CONTAINER_IMAGE)"

.PHONY: container-clean
container-clean:  ## Remove all local images with label=idmsvc-backend
	$(CONTAINER_ENGINE) image prune --filter label=idmsvc-backend --force
	$(CONTAINER_ENGINE) image list --filter label=idmsvc-backend --format "{{.ID}}" | xargs -r $(CONTAINER_ENGINE) image rm

# TODO Indicate in the options the IP assigned to the postgres container
# .PHONY: container-run
# container-run: CONTAINER_ENGINE_OPTS += --env-file .env
# container-run:  ## Run with CONTAINER_ENGINE_OPTS the CONTAINER_IMAGE using CONTAINER_RUN_ARGS as arguments (eg. make container-run CONTAINER_ENGINE_OPTS="-p 9000:9000")
# 	$(CONTAINER_ENGINE) run $(CONTAINER_ENGINE_OPTS) $(CONTAINER_IMAGE) $(CONTAINER_RUN_ARGS)
