#!/usr/bin/python3
"""Create new stub domain and print ipa-hcc register command
"""
import argparse
import base64
import json
import subprocess

import requests
import requests.auth


REGISTER_JSON = {
    "domain_type": "rhel-idm",
}

IDM_CI_SECRETS = """\
# idm-ci/secrets
export RHC_ENV="ephemeral"
export RHC_ORG="12345"
export RHC_KEY="not-used"
export RH_API_TOKEN="not-used"
export IDMSVC_BACKEND={idmsvc_backend}
export DEV_USERNAME={username}
export DEV_PASSWORD={password}
"""

XRHID = {
    "identity": {
        "auth_type": "basic-auth",
        "employee_account_number": "07214",
        "internal": {"org_id": "12345"},
        "org_id": "12345",
        "type": "User",
        "user": {
            "email": "test@hmsidm.test",
            "first_name": "HMS",
            "is_active": True,
            "is_internal": False,
            "is_org_admin": True,
            "last_name": "IDM",
            "locale": "en",
            "user_id": "test",
            "username": "test",
        },
    }
}

parser = argparse.ArgumentParser(description=__doc__.strip())
parser.add_argument(
    "--compose",
    help="Use compose instead of ephemeral (defaults to localhost:8000)",
    nargs="?",
    const="localhost:8000",
    default=None,
    dest="backend",
)
parser.add_argument(
    "--secrets-file",
    help="write idm-ci secrets file",
    default=None,
)


def oc(*args) -> str:
    cmd = ["oc"]
    cmd.extend(args)
    return subprocess.check_output(cmd, text=True).strip()


def main() -> None:
    args = parser.parse_args()

    headers = {
        "X-Rh-Insights-Request-Id": "test_12345",
    }

    if args.backend is not None:
        username = "compose"
        password = "compose"
        idmsvc_backend = args.backend
        token_url = f"http://{idmsvc_backend}/api/idmsvc/v1/domains/token"
        auth = None
        headers["X-Rh-Identity"] = base64.urlsafe_b64encode(
            json.dumps(XRHID).encode("utf-8")
        )
    else:
        namespace = oc("project", "-q")
        keycloak = oc("get", f"secrets/env-{namespace}-keycloak", "-o", "json")
        secrets = json.loads(keycloak)
        username = base64.b64decode(secrets["data"]["defaultUsername"]).decode("utf-8")
        password = base64.b64decode(secrets["data"]["defaultPassword"]).decode("utf-8")
        idmsvc_backend = oc(
            "get",
            "routes",
            "-l",
            "app=idmsvc-backend",
            "-o",
            "jsonpath={.items[0].spec.host}",
        )
        token_url = f"https://{idmsvc_backend}/api/idmsvc/v1/domains/token"
        auth = requests.auth.HTTPBasicAuth(username, password)

    if args.secrets_file:
        print(f"Writing idm-ci secrets to file '{args.secrets_file}'.")
        with open(args.secrets_file, "w") as f:
            f.write(IDM_CI_SECRETS.format(**locals()))

    resp = requests.post(
        token_url, auth=auth, headers=headers, json=REGISTER_JSON, timeout=10
    )

    resp.raise_for_status()
    j = resp.json()
    token = j["domain_token"]

    print(f"ipa-hcc register --unattended {token}")


if __name__ == "__main__":
    main()
