#!/bin/bash
set -eo pipefail

# shellcheck disable=SC1091
source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

INVENTORY_ID=$"$1"
FQDN="$2"
[ "${INVENTORY_ID}" != "" ] || error "INVENTORY_ID is empty"
[ "${FQDN}" != "" ] || error "FQDN is empty"

export X_RH_IDENTITY="${X_RH_IDENTITY:-$(identity_user)}"
X_RH_IDM_VERSION="$(idm_version)"
export X_RH_IDM_VERSION
unset X_RH_FAKE_IDENTITY
unset CREDS
"${REPOBASEDIR}/scripts/curl.sh" -i -X POST -d '{}' "${BASE_URL}/host-conf/${INVENTORY_ID}/${FQDN}"
