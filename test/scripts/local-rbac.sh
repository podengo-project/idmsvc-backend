#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

unset X_RH_IDENTITY
unset X_RH_FAKE_IDENTITY
unset CREDS
unset X_RH_IDM_VERSION
BASE_URL="http://localhost:8020/api/rbac/v1"

exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/access?application=idmsvc&offset=0&limit=1000"
