#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

export X_RH_IDENTITY="${X_RH_IDENTITY:-$(identity_generator)}"
unset X_RH_FAKE_IDENTITY
unset CREDS

exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/domains/${UUID}"
