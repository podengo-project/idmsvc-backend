#
# Include for common parts for ephemeral environment shared between all the scripts
#
# NOTE: Be aware that curl.sh wrapper set options based in the environment
#       variables that has value when it is invoked, and set an environment
#       variable could change the behave on how the request is formed.
#
# See: ./scripts/curl.sh
# See: https://consoledot.pages.redhat.com/docs/dev/getting-started/uienvs/uienvironments.html
# See: https://access.stage.redhat.com/management/api
# See: https://access.stage.redhat.com/articles/3626371
#
source "$(dirname "${BASH_SOURCE[0]}")/common.inc"

NAMESPACE="idmsvc-stage"
export NAMESPACE

[ "${RH_API_TOKEN}" == "" ] && {
    error "RH_API_TOKEN is empty: Retrieve a token from: https://access.stage.redhat.com/management/api"
}
CREDS="$(curl -s -d grant_type=refresh_token \
                 -d client_id=rhsm-api \
                 -d refresh_token=${RH_API_TOKEN} \
                 https://sso.stage.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token \
         | jq --raw-output .access_token)"
export CREDS

# shellcheck disable=SC2034  # ignore unused variable
BASE_URL="https://console.stage.redhat.com/api/idmsvc/v1"
