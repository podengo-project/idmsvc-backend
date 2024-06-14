#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

export X_RH_IDENTITY="${X_RH_IDENTITY:-$(identity_generator)}"
unset CREDS
export X_RH_IDM_REGISTRATION_TOKEN="$TOKEN"
unset X_RH_IDM_VERSION

exec "${REPOBASEDIR}/scripts/curl.sh" -i -X DELETE "${BASE_URL}/domains/${UUID}"
