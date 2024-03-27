#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/stage.inc"

unset X_RH_IDENTITY
unset X_RH_FAKE_IDENTITY

exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/domains"
