#!/bin/bash

function error {
    local err=$?
    printf "%s\n" "$1" >&2
    exit $err
}

# local-domains-token.sh

TOKEN="$1"
[ "${TOKEN}" != "" ] || error "TOKEN is empty"

export X_RH_IDENTITY="$( ./tools/bin/xrhidgen -org-id 12345 system -cn "6f324116-b3d2-11ed-8a37-482ae3863d30" -cert-type system | base64 -w0 )"
unset CREDS
export X_RH_IDM_REGISTRATION_TOKEN="$TOKEN"
export X_RH_IDM_VERSION='{"ipa-hcc": "0.9", "ipa": "4.10.0-8.el9_1", "os-release-id": "rhel", "os-release-version-id": "9.1"}'
BASE_URL="http://localhost:8000/api/idmsvc/v1"
./scripts/curl.sh -i -X POST -d @<( cat "test/data/http/register-rhel-idm-domain.json" | sed -e 's/{{subscription_manager_id}}/6f324116-b3d2-11ed-8a37-482ae3863d30/g' ) "${BASE_URL}/domains"
