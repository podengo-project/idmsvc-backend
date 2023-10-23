#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/ephe.inc"

INVENTORY_ID=$"$1"
[ "${INVENTORY_ID}" != "" ] || error "INVENTORY_ID is empty"
FQDN="$2"
[ "${FQDN}" != "" ] || error "FQDN is empty"

unset X_RH_IDENTITY
export X_RH_FAKE_IDENTITY="${X_RH_FAKE_IDENTITY:-$(identity_system)}"
X_RH_IDM_VERSION="$IDM_VERSION"
export X_RH_IDM_VERSION

exec "${REPOBASEDIR}/scripts/curl.sh" -i -X POST -d '{}' "${BASE_URL}/host-conf/${INVENTORY_ID}/${FQDN}"
