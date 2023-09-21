#!/bin/bash

function error {
    local err=$?
    printf "%s\n" "$*" >&2
    exit $err
}

INVENTORY_ID=$"$1"
FQDN="$2"
[ "${INVENTORY_ID}" != "" ] || error "INVENTORY_ID is empty"
[ "${FQDN}" != "" ] || error "FQDN is empty"

export NAMESPACE="$(oc project -q)"
CREDS="$( oc get secrets/env-${NAMESPACE}-keycloak -o jsonpath='{.data.defaultUsername}' | base64 -d )"
CREDS="${CREDS}:$( oc get secrets/env-${NAMESPACE}-keycloak -o jsonpath='{.data.defaultPassword}' | base64 -d )"
export CREDS

unset X_RH_IDENTITY
export X_RH_FAKE_IDENTITY="$( ./tools/bin/xrhidgen -org-id ${ORG_ID:-12345} system -cn "3f35fc7f-079c-4940-92ed-9fdc8694a0f3" -cert-type system | base64 -w0 )"
export X_RH_IDM_VERSION='{"ipa-hcc": "0.9", "ipa": "4.10.0-8.el9_1", "os-release-id": "rhel", "os-release-version-id": "9.1"}'
BASE_URL="https://$( oc get routes -l app=idmsvc-backend -o jsonpath='{.items[0].spec.host}' )/api/idmsvc/v1"
./scripts/curl.sh -i -X POST -d '{}' "${BASE_URL}/host-conf/${INVENTORY_ID}/${FQDN}"
