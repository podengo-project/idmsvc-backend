#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

unset X_RH_IDENTITY
unset X_RH_FAKE_IDENTITY
unset CREDS

exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/signing_keys"
