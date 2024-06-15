#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

export X_RH_IDENTITY="${X_RH_IDENTITY:-$(identity_generator)}"
unset X_RH_FAKE_IDENTITY
unset CREDS
unset X_RH_IDM_VERSION

exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/domains"
