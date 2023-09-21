#!/bin/bash
set -eo pipefail

# shellcheck disable=SC1091
source "$(dirname "${BASH_SOURCE[0]}")/ephe.inc"

TOKEN="$1"
[ "${TOKEN}" != "" ] || error "TOKEN is empty"

unset X_RH_IDENTITY
export X_RH_FAKE_IDENTITY="${X_RH_FAKE_IDENTITY:-$(identity_system)}"
export X_RH_IDM_REGISTRATION_TOKEN="${TOKEN}"
X_RH_IDM_VERSION="$(idm_version)"
export X_RH_IDM_VERSION
"${REPOBASEDIR}/scripts/curl.sh" -i -X POST -d @<(sed -e 's/{{subscription_manager_id}}/6f324116-b3d2-11ed-8a37-482ae3863d30/g' < "${REPOBASEDIR}/test/data/http/register-rhel-idm-domain.json") "${BASE_URL}/domains"
