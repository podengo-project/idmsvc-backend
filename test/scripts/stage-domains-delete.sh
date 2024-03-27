#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/stage.inc"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

unset X_RH_IDENTITY
unset X_RH_FAKE_IDENTITY
unset X_RH_IDM_REGISTRATION_TOKEN
X_RH_IDM_VERSION="$IDM_VERSION"
export X_RH_IDM_VERSION

exec "${REPOBASEDIR}/scripts/curl.sh" -i -X DELETE "${BASE_URL}/domains/${UUID}"
