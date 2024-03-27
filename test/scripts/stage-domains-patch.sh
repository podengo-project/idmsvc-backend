#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/stage.inc"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

unset X_RH_IDENTITY
unset X_RH_FAKE_IDENTITY
unset X_RH_IDM_REGISTRATION_TOKEN
X_RH_IDM_VERSION="$IDM_VERSION"
export X_RH_IDM_VERSION

exec "${REPOBASEDIR}/scripts/curl.sh" -i -X PATCH -d @<(sed -e "s/{{createDomain.response.body.domain_id}}/${UUID}/g" -e 's/{{subscription_manager_id}}/6f324116-b3d2-11ed-8a37-482ae3863d30/g' < "${REPOBASEDIR}/test/data/http/patch-rhel-idm-domain.json") "${BASE_URL}/domains/${UUID}"
