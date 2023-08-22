#!/bin/bash

export NAMESPACE="$(oc project -q)"
CREDS="$( oc get secrets/env-${NAMESPACE}-keycloak -o jsonpath='{.data.defaultUsername}' | base64 -d )"
CREDS="${CREDS}:$( oc get secrets/env-${NAMESPACE}-keycloak -o jsonpath='{.data.defaultPassword}' | base64 -d )"
export CREDS

unset X_RH_IDENTITY
export X_RH_FAKE_IDENTITY="$( ./bin/xrhidgen -org-id 12345 system -cn "6f324116-b3d2-11ed-8a37-482ae3863d30" -cert-type system | base64 -w0 )"
BASE_URL="https://$( oc get routes -l app=idmsvc-backend -o jsonpath='{.items[0].spec.host}' )/api/idmsvc/v1"
./scripts/curl.sh -i -X POST -d '{"domain_type": "rhel-idm"}' "${BASE_URL}/domains/token"
