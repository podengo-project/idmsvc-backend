#!/bin/bash

function error {
    local err=$?
    printf "%s\n" "$*" >&2
    exit $err
}

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

export X_RH_IDENTITY="$( ./tools/bin/xrhidgen -org-id 12345 user -is-active=true -is-org-admin=true -user-id test -username test | base64 -w0 )"
unset X_RH_FAKE_IDENTITY
unset CREDS
BASE_URL="http://localhost:8000/api/idmsvc/v1"
./scripts/curl.sh -i "${BASE_URL}/domains/${UUID}"
