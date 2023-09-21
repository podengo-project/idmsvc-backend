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

ORG_ID="${ORG_ID:-12345}"

export IDENTITY_USER=""   # Use $(identity_user)
export IDENTITY_SYSTEM="" # Use $(identity_system)
export IDM_VERSION=""     # Use $(idm_version)

SRCDIR="$(dirname "${BASH_SOURCE[0]}")"
# shellcheck disable=SC2034  # ignore unused variable
BASEDIR="$(dirname "$(dirname "${SRCDIR}")")"
REPOBASEDIR="$(git rev-parse --show-toplevel)"
export REPOBASEDIR
export XRHIDGEN="${REPOBASEDIR}/tools/bin/xrhidgen"

if [[ ! -x "${XRHIDGEN}" ]]; then
    error "${XRHIDGEN} is missing, run 'make install-tools'"
    exit 2
fi

function identity_user() {
    "${XRHIDGEN}" -org-id "${ORG_ID}" user -is-active=true -is-org-admin=true -user-id test -username test | base64 -w0
}
export -f identity_user     # Needed for making it available in sub-shells

function identity_system() {
    "${XRHIDGEN}" -org-id "${ORG_ID}" system -cn "6f324116-b3d2-11ed-8a37-482ae3863d30" -cert-type system | base64 -w0
}
export -f identity_system   # Needed for making it available in sub-shells

function idm_version() {
    IDM_VERSION='{"ipa-hcc": "0.7", "ipa": "4.10.0-8.el9_1"}'
    printf "%s" "${IDM_VERSION}"
}
export -f idm_version
