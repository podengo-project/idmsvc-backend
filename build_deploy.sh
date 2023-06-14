#!/bin/bash
set -e

SERVICE="hmsidm"
BASE_IMAGE="${SERVICE}-backend"
CONTAINER_ENGINE="${CONTAINER_ENGINE:-podman}"

IMAGE="${IMAGE:-quay.io/cloudservices/${BASE_IMAGE}}"
IMAGE_TAG="$(git rev-parse --short=7 HEAD)"
SMOKE_TEST_TAG="latest"

if [[ -z "${QUAY_USER}" || -z "${QUAY_TOKEN}" ]]; then
    echo "QUAY_USER and QUAY_TOKEN must be set"
    exit 1
fi

if [[ -z "${RH_REGISTRY_USER}" || -z "${RH_REGISTRY_TOKEN}" ]]; then
    echo "RH_REGISTRY_USER and RH_REGISTRY_TOKEN  must be set"
    exit 1
fi

# Login to quay
make registry-login \
    CONTAINER_REGISTRY_USER="${QUAY_USER}" \
    CONTAINER_REGISTRY_TOKEN="${QUAY_TOKEN}" \
    CONTAINER_REGISTRY="quay.io"

# Login to registry.redhat
make registry-login \
    CONTAINER_REGISTRY_USER="${RH_REGISTRY_USER}" \
    CONTAINER_REGISTRY_TOKEN="${RH_REGISTRY_TOKEN}" \
    CONTAINER_REGISTRY="registry.redhat.io"

# Build and push
make container-build container-push \
    CONTAINER_BUILD_OPTS=--no-cache \
    CONTAINER_IMAGE_BASE="${IMAGE}" \
    CONTAINER_IMAGE_TAG="${IMAGE_TAG}"

# Push to logged in registries and tag for SHA
"${CONTAINER_ENGINE}" tag "${IMAGE}:${IMAGE_TAG}" "${IMAGE}:${SMOKE_TEST_TAG}"
"${CONTAINER_ENGINE}" push "${IMAGE}:${SMOKE_TEST_TAG}"
