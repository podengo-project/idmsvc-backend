#
# Include file with common parts shared for local and ephemeral
#

# Troubleshooting:
# - Run with DEBUG=1 to see some traces from curl.sh wrapper
#   $ DEBUG=1 ./test/scripts/local-domain-token.sh
# - Run with more verbose by:
#   $ DEBUG=1 bash -xv ./test/scripts/local-domain-token.sh

function error {
    local err=$?
    printf "ERROR: %s\n" "$*" >&2
    exit $err
}

command -v jq &>/dev/null || error "'jq' tool is not available"
command -v curl &>/dev/null || error "'curl' tool is not available"
command -v base64 &>/dev/null || error "'base64' tool is not available"

ORG_ID="${ORG_ID:-12345}"
# shellcheck disable=SC2034
IDM_VERSION='{"ipa-hcc": "0.7", "ipa": "4.10.0-8.el9_1"}'

SRCDIR="$(dirname "${BASH_SOURCE[0]}")"
# shellcheck disable=SC2034  # ignore unused variable
BASEDIR="$(dirname "$(dirname "${SRCDIR}")")"
REPOBASEDIR="$(git rev-parse --show-toplevel)"
export REPOBASEDIR
export XRHIDGEN="${REPOBASEDIR}/tools/bin/xrhidgen"
# "jwt-auth" "cert-auth"
export AUTH_TYPE="${AUTH_TYPE:-jwt-auth}"

XRHID_CLIENT_ID="${XRHID_CLIENT_ID:-$(uuidgen)}"

if [[ ! -x "${XRHIDGEN}" ]]; then
    error "${XRHIDGEN} is missing, run 'make install-tools'"
    exit 2
fi

base64nowrap() {
    if [ "$(uname -s)" == "Darwin" ]; then
        base64
    else
        base64 -w0
    fi
}

identity_user() {
    "${XRHIDGEN}" -org-id "${ORG_ID}" -auth-type "${AUTH_TYPE}" user -is-active=true -is-org-admin=true -user-id test -username test | base64nowrap
}

identity_system() {
    "${XRHIDGEN}" -org-id "${ORG_ID}" -auth-type "${AUTH_TYPE}" system -cn "6f324116-b3d2-11ed-8a37-482ae3863d30" -cert-type system | base64nowrap
}

identity_service_account() {
    "${XRHIDGEN}" -org-id "${ORG_ID}" -auth-type "${AUTH_TYPE}" service-account -client-id "${XRHID_CLIENT_ID}" -username test | base64nowrap
}

XRHID_AS="${XRHID_AS:-user}"
case "${XRHID_AS}" in
    "user" )
        identity_generator() {
            identity_user "$@"
        }
        ;;
    "system" )
        identity_generator() {
            identity_system "$@"
        }
        ;;
    "service-account" )
        identity_generator() {
            identity_service_account "$@"
        }
        ;;
    * )
        error "XRHID_AS='${XRHID_AS}' not supported"
        ;;
esac
export -f identity_generator
