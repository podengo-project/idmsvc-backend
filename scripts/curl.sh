#!/bin/bash

##
# curl helper to wrap and add headers automatically based
# into the environment variables defined; the idea is
# define once, use many times, so the curl command can be
# simplified.
#
# NOTE: Do not forget to to 'export MYVAR' for the next
#       one to get a better user experience, else you need
#       to set before the command in the same command,
#       reducing the user experience then.
#
# X_RH_INSIGHTS_REQUEST_ID if it is empty, a random value
#   is used to fill up the field. This will add the header
#   'X-Rh-Insights-Request-Id'.
# X_RH_IDM_VERSION if it is not empty, then the header
#   'X-Rh-Idm-Version' is added to the curl options.
# X_RH_IDENTITY if it is not empty, then the header
#   'X-Rh-Identity' is added to the curl options.
#   NOTE this value should not be specified when launching
#   the request against ephemeral environment, as the
#   values are filled by the api gateway.
# X_RH_FAKE_IDENTITY if it is not empty, then the header
#   'X-Rh-Fake-Identity' is added to the curl options.
#   this is used for development propose only, and it
#   requires APP_IS_FAKE_ENABLED=true when it was deployed.
# X_RH_IDM_REGISTRATION_TOKEN if it is not empty, then
#   the header is filled.
# CREDS if it is not empty, then `-u "${CREDS}"` options
#   are added to the curl command; this is used to reach out
#   the ephemeral environment.
# RH_API_TOKEN if it is not empty, then `-H "Authorization: Bearer ${CREDS}"`
#   options are added to the curl command; this is used to
#   reach out the stage environment.
##

# Uncomment to print verbose traces into stderr
# DEBUG=1
function verbose {
    [ "${DEBUG}" != 1 ] && return 0
    echo "$@" >&2
}

# Initialize the array of options
opts=()

# Generate a X_RH_INSIGHTS_REQUEST_ID if it is not set
if [ "${X_RH_INSIGHTS_REQUEST_ID}" == "" ]; then
    if [ "$(uname -s)" == "Darwin" ]; then
        X_RH_INSIGHTS_REQUEST_ID="test_$(uuidgen | sed 's/[-]//g' | head -c 20)"
    else
        X_RH_INSIGHTS_REQUEST_ID="test_$(sed 's/[-]//g' < "/proc/sys/kernel/random/uuid" | head -c 20)"
    fi
fi
opts+=("-H" "X-Rh-Insights-Request-Id: ${X_RH_INSIGHTS_REQUEST_ID}")
verbose "-H X-Rh-Insights-Request-Id: ${X_RH_INSIGHTS_REQUEST_ID}"

# Optionally add X-Rh-Idm-Version
if [ "${X_RH_IDM_VERSION}" != "" ]; then
    opts+=("-H" "X-Rh-Idm-Version: ${X_RH_IDM_VERSION}")
    verbose "-H X-Rh-Idm-Version: ${X_RH_IDM_VERSION}"
fi

# Optionally add X-Rh-Identity (used for testing in local workstation)
if [ "${X_RH_IDENTITY}" != "" ]; then
    opts+=("-H" "X-Rh-Identity: ${X_RH_IDENTITY}")
    verbose "-H X-Rh-Identity: ${X_RH_IDENTITY}"
fi

# Optionally add X-Rh-Fake-Identity (used for testing in ephemeral)
if [ "${X_RH_FAKE_IDENTITY}" != "" ]; then
    opts+=("-H" "X-Rh-Fake-Identity: ${X_RH_FAKE_IDENTITY}")
    verbose "-H X-Rh-Fake-Identity: ${X_RH_FAKE_IDENTITY}"
fi

# Optionally add X-Rh-Idm-Registration-Token
if [ "${X_RH_IDM_REGISTRATION_TOKEN}" != "" ]; then
    opts+=("-H" "X-Rh-Idm-Registration-Token: ${X_RH_IDM_REGISTRATION_TOKEN}")
    verbose "-H X-Rh-Idm-Registration-Token: ${X_RH_IDM_REGISTRATION_TOKEN}"
fi

# Add Content-Type
opts+=("-H" "Content-Type: application/json")
verbose "-H Content-Type: application/json"

# Optionally add CREDS if it was set (used for testing in ephemeral)
# See: make ephemeral-namespace-describe
if [ "${CREDS}" != "" ]; then
    if [ "${RH_API_TOKEN}" == "" ] && [ "${https_proxy}" == "" ]; then
        opts+=("-u" "${CREDS}")
        # shellcheck disable=SC2016
        verbose '-u "${CREDS}"'
    else
        verbose "https_proxy=${https_proxy:-${http_proxy}}"
        opts+=("-H" "Authorization: Bearer ${CREDS}")
        verbose '-H Authorization: Bearer ${CREDS}'
    fi
fi

# Add the rest of values
opts+=("$@")

verbose /usr/bin/curl "${opts[@]}"
/usr/bin/curl "${opts[@]}"
