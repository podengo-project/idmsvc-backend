#!/bin/bash

export NAMESPACE="$(oc project -q)"
CREDS="$( oc get secrets/env-${NAMESPACE}-keycloak -o jsonpath='{.data.defaultUsername}' | base64 -d )"
CREDS="${CREDS}:$( oc get secrets/env-${NAMESPACE}-keycloak -o jsonpath='{.data.defaultPassword}' | base64 -d )"
export CREDS

unset X_RH_IDENTITY
unset X_RH_FAKE_IDENTITY
BASE_URL="https://$( oc get routes -l app=hmsidm-backend -o jsonpath='{.items[0].spec.host}' )/api/hmsidm/v1"
./scripts/curl.sh -i "${BASE_URL}/domains" -X POST -d @./test/data/http/create-rhel-idm-domain.json

