#!/bin/bash
#
# Generate a selg-sifned certificate for checking TlsCAPath
#
# Run from the base of the repository by: ./test/scripts/generate-cert.sh
set -e
function error {
    local ret=$?
    printf "ERROR: %s\n" "$*" >&2
    [ "$ret" == "0" ] && ret=1
    exit $ret
}

[ -e .git ] || error "Not running from base repository directory"
[ -e ./cdapp/certs ] || mkdir -p ./cdapp/certs
openssl genrsa -out ./cdapp/key.pem 2048
openssl req -new -key ./cdapp/key.pem -out ./cdapp/csr.pem -subj "/O=Organization/CN=example.com"
openssl x509 -req -days 365 -in ./cdapp/csr.pem -signkey ./cdapp/key.pem -out ./cdapp/certs/cert.pem
openssl x509 -in ./cdapp/certs/cert.pem -text
