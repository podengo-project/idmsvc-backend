#!/bin/bash

function error {
    local err=$?
    printf "%s\n" "$1" >&2
    exit $err
}

# make db-cli <<< "select domain_uuid from domains order by id desc limit 1;\\q"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

export X_RH_IDENTITY="$( ./bin/xrhidgen -org-id 12345 user -is-active=true -is-org-admin=true -user-id test -username test | base64 -w0 )"
unset CREDS
export X_RH_IDM_REGISTRATION_TOKEN="$TOKEN"
export X_RH_IDM_VERSION="$( base64 -w0 <<< '{"ipa-hcc": "0.7", "ipa": "4.10.0-8.el9_1"}' )"
BASE_URL="http://localhost:8000/api/hmsidm/v1"
./scripts/curl.sh -i -X DELETE "${BASE_URL}/domains/${UUID}"
