#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/ephe.inc"

unset X_RH_IDENTITY
export X_RH_FAKE_IDENTITY="${X_RH_FAKE_IDENTITY:-$(identity_user)}"

exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/domains"
