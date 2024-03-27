#
# Include for common parts for stage environment shared between all the scripts
#
# NOTE: Be aware that curl.sh wrapper set options based in the environment
#       variables that has value when it is invoked, and set an environment
#       variable could change the behave on how the request is formed.
#
# See: ./scripts/curl.sh
#
source "$(dirname "${BASH_SOURCE[0]}")/common.inc"

NAMESPACE="idmsvc-stage"
export NAMESPACE

if [ "${https_proxy}" == "" ] && [ "${http_proxy}" == "" ]; then
    error "https_proxy and http_proxy are empty: check your proxy configuration"
fi
if [ "${XRHID_AS}" == "service-account" ]; then
    [ "${CLIENT_ID}" == "" ] && {
        error "CLIENT_ID is required for service-accounts: https://console.stage.redhat.com/iam/service-accounts"
    }
    [ "${CLIENT_SECRET}" == "" ] && {
        error "CLIENT_SECRET is required for service-accounts: https://console.stage.redhat.com/iam/service-accounts"
    }
    CREDS="$(curl "https://sso.stage.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token" \
        -d "grant_type=client_credentials" \
        -d "scope=api.console" \
        -d client_id="${CLIENT_ID}" \
        -d client_secret="${CLIENT_SECRET}" \
            | jq --raw-output .access_token)"
else
    [ "${RH_API_TOKEN}" == "" ] && {
        error "RH_API_TOKEN is empty: Retrieve a token from: https://access.stage.redhat.com/management/api"
    }
    CREDS="$(curl -s -d grant_type=refresh_token \
                    -d client_id=rhsm-api \
                    -d refresh_token="${RH_API_TOKEN}" \
                    https://sso.stage.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token \
            | jq --raw-output .access_token)"
fi
export CREDS

# shellcheck disable=SC2034  # ignore unused variable
BASE_URL="https://console.stage.redhat.com/api/idmsvc/v1"
