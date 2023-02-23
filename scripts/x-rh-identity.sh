#!/bin/bash

ORG_ID="$1"
USER_NAME="$2"



function print_out_usage {
	cat <<EOF
Usage: ./scripts/header.sh <org_id> [user_name]
EOF
}

function error {
	local err=$?
	print_out_usage >&2
	printf "error: %s\n" "$*" >&2
	exit $err
}

[ "${ORG_ID}" != "" ] || error "ORG_ID is required and cannot be empty"

case "$( uname -s )" in
"Darwin" )
  BASE64ENC="base64 -b 0"
;;

"Linux" | *)
  BASE64ENC="base64 -w0"
;;
esac

printf "X-Rh-Identity: "
cat <<EOF | jq -c -M | ${BASE64ENC}
{
  "account_number": "11111",
  "employee_account_number": "22222",
  "org_id": "${ORG_ID}",
  "internal": {
    "org_id": "${ORG_ID}"
  },
  "user": {
    "username": "${USER_NAME}",
    "email": "${USER_NAME}@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "is_org_admin": true,
    "is_internal": true,
    "locale": "en",
    "user_id": "1987348"
  },
  "system": {},
  "associate": {
    "Role": null,
    "email": "",
    "givenName": "",
    "rhatUUID": "",
    "surname": ""
  },
  "x509": {
    "subject_dn": "",
    "issuer_dn": ""
  },
  "type": "User",
  "auth_type": ""
}
EOF
