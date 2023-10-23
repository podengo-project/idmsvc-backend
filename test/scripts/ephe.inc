#
# Include for common parts for ephemeral environment shared between all the scripts
#
# NOTE: Be aware that curl.sh wrapper set options based in the environment
#       variables that has value when it is invoked, and set an environment
#       variable could change the behave on how the request is formed.
#
# See: ./scripts/curl.sh
#
source "$(dirname "${BASH_SOURCE[0]}")/common.inc"

NAMESPACE="$(oc project -q)"
export NAMESPACE

username="$( oc get secrets/env-"${NAMESPACE}"-keycloak -o jsonpath='{.data.defaultUsername}' | base64 -d )"
password="$( oc get secrets/env-"${NAMESPACE}"-keycloak -o jsonpath='{.data.defaultPassword}' | base64 -d )"
CREDS="${username}:${password}"
export CREDS

# shellcheck disable=SC2034  # ignore unused variable
BASE_URL="https://$( oc get routes -l app=idmsvc-backend -o jsonpath='{.items[0].spec.host}' )/api/idmsvc/v1"
