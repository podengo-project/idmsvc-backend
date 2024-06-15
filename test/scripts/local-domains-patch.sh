#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

export X_RH_IDENTITY="${X_RH_IDENTITY:-$(identity_generator)}"
unset CREDS
unset X_RH_IDM_REGISTRATION_TOKEN

exec "${REPOBASEDIR}/scripts/curl.sh" -i -X PATCH -d @<(sed -e "s/{{createDomain.response.body.domain_id}}/${UUID}/g" < "${REPOBASEDIR}/test/data/http/patch-rhel-idm-domain.json") "${BASE_URL}/domains/${UUID}"
