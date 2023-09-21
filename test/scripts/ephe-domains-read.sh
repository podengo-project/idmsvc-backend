#!/bin/bash
set -eo pipefail

# shellcheck disable=SC1091
source "$(dirname "${BASH_SOURCE[0]}")/ephe.inc"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

unset X_RH_IDENTITY
export X_RH_FAKE_IDENTITY="${X_RH_FAKE_IDENTITY:-$(identity_user)}"
"${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/domains/${UUID}"
