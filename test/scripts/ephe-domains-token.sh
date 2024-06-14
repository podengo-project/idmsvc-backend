#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/ephe.inc"

unset X_RH_IDENTITY
export X_RH_FAKE_IDENTITY="${X_RH_FAKE_IDENTITY:-$(identity_generator)}"

exec "${REPOBASEDIR}/scripts/curl.sh" -i -X POST -d '{"domain_type": "rhel-idm"}' "${BASE_URL}/domains/token"
