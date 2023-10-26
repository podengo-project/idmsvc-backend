#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

export X_RH_IDENTITY="${X_RH_IDENTITY:-$(identity_user)}"
unset X_RH_FAKE_IDENTITY
unset CREDS

exec "${REPOBASEDIR}/scripts/curl.sh" -i -X POST -d '{"domain_type": "rhel-idm"}' "${BASE_URL}/domains/token"
