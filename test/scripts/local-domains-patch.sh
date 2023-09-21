#!/bin/bash

function error {
    local err=$?
    printf "%s\n" "$1" >&2
    exit $err
}

# make db-cli <<< "select domain_uuid from domains order by id desc limit 1;\\q"
# make db-cli <<< "select token from ipas order by id desc limit 1;\\q"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

export X_RH_IDENTITY="$( ./bin/xrhidgen -org-id ${ORG_ID:-12345} user -is-active=true -is-org-admin=true -user-id test -username test | base64 -w0 )"
unset CREDS
unset X_RH_IDM_REGISTRATION_TOKEN
BASE_URL="http://localhost:8000/api/idmsvc/v1"
./scripts/curl.sh -i -X PATCH -d @<( cat "test/data/http/patch-rhel-idm-domain.json" | sed -e "s/{{createDomain.response.body.domain_id}}/${UUID}/g" ) "${BASE_URL}/domains/${UUID}"
