#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/ephe.inc"

unset X_RH_IDENTITY
unset X_RH_FAKE_IDENTITY

BASE_URL="${BASE_URL%/*}"  # Remove v1
BASE_URL="${BASE_URL%/*}"  # Remove service name
BASE_URL="${BASE_URL}/rbac/v1"

# Be aware the trailing slash is required
exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/access/?application=idmsvc&offset=0&limit=1000"
