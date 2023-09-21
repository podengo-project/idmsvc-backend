#!/bin/bash
set -eo pipefail

# shellcheck disable=SC1091
source "$(dirname "${BASH_SOURCE[0]}")/ephe.inc"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

unset X_RH_IDENTITY
export X_RH_FAKE_IDENTITY="${X_RH_FAKE_IDENTITY:-$(identity_user)}"
export X_RH_IDM_REGISTRATION_TOKEN="${TOKEN}"
X_RH_IDM_VERSION="$(idm_version)"
export X_RH_IDM_VERSION
"${REPOBASEDIR}/scripts/curl.sh" -i -X DELETE "${BASE_URL}/domains/${UUID}"
