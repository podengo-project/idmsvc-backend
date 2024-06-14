#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/ephe.inc"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

unset X_RH_IDENTITY
export X_RH_FAKE_IDENTITY="${X_RH_FAKE_IDENTITY:-$(identity_generator)}"

exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/domains/${UUID}"
