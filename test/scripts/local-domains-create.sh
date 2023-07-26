#!/bin/bash

export X_RH_IDENTITY="$( ./bin/xrhidgen -org-id 12345 user -is-active=true -is-org-admin=true -user-id test -username test | base64 -w0 )"
unset CREDS
BASE_URL="http://localhost:8000/api/idmsvc/v1"
./scripts/curl.sh -i "${BASE_URL}/domains" -X POST -d @./test/data/http/create-rhel-idm-domain.json

