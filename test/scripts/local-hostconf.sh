#!/bin/bash
set -eo pipefail

export XRHID_AS="${XRHID_AS:-system}"
export AUTH_TYPE="${AUTH_TYPE:-cert-auth}"

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

INVENTORY_ID=$"$1"
FQDN="$2"
[ "${INVENTORY_ID}" != "" ] || error "INVENTORY_ID is empty"
[ "${FQDN}" != "" ] || error "FQDN is empty"

export X_RH_IDENTITY="${X_RH_IDENTITY:-$(identity_system)}"
export X_RH_IDM_VERSION
unset X_RH_FAKE_IDENTITY
unset CREDS

exec "${REPOBASEDIR}/scripts/curl.sh" -i -X POST -d '{}' "${BASE_URL}/host-conf/${INVENTORY_ID}/${FQDN}"
