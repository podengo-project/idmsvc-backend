#!/bin/bash
set -eo pipefail

export XRHID_AS="${XRHID_AS:-system}"
export AUTH_TYPE="${AUTH_TYPE:-cert-auth}"

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

TOKEN="$1"
[ "${TOKEN}" != "" ] || error "TOKEN is empty"

export X_RH_IDENTITY="${X_RH_IDENTITY:-$(identity_generator)}"
unset CREDS
export X_RH_IDM_REGISTRATION_TOKEN="$TOKEN"
X_RH_IDM_VERSION="$IDM_VERSION"
export X_RH_IDM_VERSION

exec "${REPOBASEDIR}/scripts/curl.sh" -i -X POST -d @<(sed -e 's/{{subscription_manager_id}}/6f324116-b3d2-11ed-8a37-482ae3863d30/g' < "${REPOBASEDIR}/test/data/http/register-rhel-idm-domain.json") "${BASE_URL}/domains"
