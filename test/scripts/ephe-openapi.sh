#!/bin/bash
set -eo pipefail

# shellcheck disable=SC1091
source "$(dirname "${BASH_SOURCE[0]}")/ephe.inc"

unset X_RH_IDENTITY
unset X_RH_FAKE_IDENTITY
"${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/openapi.json"
