#!/bin/bash
set -eo pipefail

export XRHID_AS="${XRHID_AS:-system}"
export AUTH_TYPE="${AUTH_TYPE:-cert-auth}"

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

export X_RH_IDENTITY="${X_RH_IDENTITY:-$(identity_generator)}"
unset X_RH_FAKE_IDENTITY
unset CREDS

exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/signing_keys"
