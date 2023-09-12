#!/usr/bin/env python3
import argparse
import base64
import random
import string
from random import randbytes
import subprocess
import sys
import requests
import json


CONTENT_TYPE = "application/json"

HEADER_CONTENT_TYPE = "Content-Type"
HEADER_X_RH_IDENTITY = "X-Rh-Identity"
HEADER_X_RH_INSIGHTS_REQUEST_ID = "X-Rh-Insights-Request-Id"
HEADER_X_RH_IDM_VERSION = "X-Rh-Idm-Version"
# FIXME Diferent headers in request/response
HEADER_X_RH_IDM_REGISTRATION_TOKEN = "X-Rh-Idm-Registration-Token"
HEADER_X_RH_IDM_RHELIDM_REGISTER_TOKEN = "X-Rh-Idm-RhelIdm-Register-Token"

DEFAULT_ORG_ID = "12345"

class xrhidgen:
    """Wrapper to call ./tools/bin/xrhidgen binary and get a x-rh-identity header"""
    def __init__(self,
                 *extra_args,
                 org_id=DEFAULT_ORG_ID,
                 account_number=None,
                 auth_type=None,
                 employe_account_number=None,
                 xrhidgen_type=None):
        self.org_id=org_id
        self.account_number=account_number
        self.auth_type=auth_type
        self.employe_account_number=employe_account_number
        self.xrhidgen_type=xrhidgen_type
        self.extra_args=extra_args

    def __call__(self, *args):
        if self.xrhidgen_type is None:
            sys.exit("'xrhidgen_type' is not set for 'xrhidgen'")
        options = ["./tools/bin/xrhidgen"]
        if self.org_id is not None:
            options.extend(['-org-id', self.org_id])
        if self.account_number is not None:
            options.extend(['-account-number', self.org_id])
        if self.auth_type is not None:
            options.extend(['-auth-type', self.org_id])
        if self.employe_account_number is not None:
            options.extend(['-type', self.type_id])
        if self.xrhidgen_type is None or self.xrhidgen_type == '':
            sys.exit("'xrhidgen_type' is None")
        options.append(self.xrhidgen_type)
        # ./tools/bin/xrhidgen -org-id 12345 system -cn "6f324116-b3d2-11ed-8a37-482ae3863d30" -cert-type system
        options.extend(self.extra_args)
        options.extend(args)
        output = subprocess.check_output(options)
        return json.loads(output)

    def __str__(self):
        return json.dumps(self.__call__())


def get_register_data(domain_id, domain_name, subscription_manager_id):
    """Return an example payload to register a domain"""
    title = f"Domain {domain_name}".replace(".", " ")
    description = f"Description Domain {domain_name}".replace(".", " ")
    data = {
        "domain_id": domain_id,
        "title": title,
        "description": description,
        "domain_type": "rhel-idm",
        "domain_name": domain_name,
        "auto_enrollment_enabled": True,
        "rhel-idm": {
            "realm_name": domain_name.upper(),
            "servers": [
            {
                "fqdn": "ipaserver." + domain_name,
                "subscription_manager_id": subscription_manager_id,
                "location": "boston",
                "ca_server": True,
                "hcc_enrollment_server": True,
                "hcc_update_server": True,
                "pkinit_server": True
            },
            {
                "fqdn": "server2." + domain_name,
                "subscription_manager_id": subscription_manager_id,
                "ca_server": False,
                "hcc_enrollment_server": False,
                "hcc_update_server": False,
                "pkinit_server": False
            }
            ],
            "locations": [
            {
                "name": "boston",
                "description": "Boston data center"
            },
            {
                "name": "europe"
            }
            ],
            "ca_certs": [
            {
                "nickname": domain_name.upper() + " IPA CA",
                "issuer": "CN=Certificate Authority, O=" + domain_name.upper(),
                "subject": "CN=My Domain, O=" + domain_name.upper(),
                "serial_number": "1",
                "not_before": "2023-01-31T13:23:36Z",
                "not_after": "2023-01-31T13:23:36Z",
                "pem": "-----BEGIN CERTIFICATE-----\nMIIElzCCAv+gAwIBAgIBATANBgkqhkiG9w0BAQsFADA6MRgwFgYDVQQKDA9ITVNJ\nRE0tREVWLlRFU1QxHjAcBgNVBAMMFUNlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0y\nMzA2MTIwNjEyMThaFw00MzA2MTIwNjEyMThaMDoxGDAWBgNVBAoMD0hNU0lETS1E\nRVYuVEVTVDEeMBwGA1UEAwwVQ2VydGlmaWNhdGUgQXV0aG9yaXR5MIIBojANBgkq\nhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEA/F+63FGVUElkycJ2I5/rOIQ8331bfqp+\nraVuft2wezXj9O60X4DsEXltjMM+Lb3vPpInI6Fjdr74RWiz7YeWRYT8y4AgiZ7O\nrbe1ivvmutZwdA3S3KVoQhfqLUzYKksL7IpLQFuXsOm85GMQsw2SNz0NIlM3Ixjv\nKFyARcFSLzBAlIUHdZwq2e8PKvIcLGjHRGczfBqSviCBKxTTO3S2vRRHFEw8lsmJ\nyqIb8gLLOSRi4GqZfp6RRnr88z7z/xqZc7ffDko3ngjUn1Cynm715Xqftlj3o297\naVQ/Oxgw/ODiQSZl+HnOgrrH4XbM+hVUfxBXydVgPrN8mTrTcY0X03cLqMWCFO6E\n8XAJFkY+1SLOdruHTfdhbmRcp/vvyZ3rcSP9qk75jFPr3iKU5vnbAtbZfGtzk6te\nsG/Y8tRjdLvcKKM9PBa93VA56nN0+RLtOn24/UfiYjYsYQeq1wJnfJUlcrER9X6t\nbX1umBXcwT9FeofJENCZqP3YfU0EH76nAgMBAAGjgacwgaQwHwYDVR0jBBgwFoAU\ntQw3tdMW/Sz+VLsOZaefg4Vnrm0wDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8E\nBAMCAcYwHQYDVR0OBBYEFLUMN7XTFv0s/lS7DmWnn4OFZ65tMEEGCCsGAQUFBwEB\nBDUwMzAxBggrBgEFBQcwAYYlaHR0cDovL2lwYS1jYS5obXNpZG0tZGV2LnRlc3Qv\nY2Evb2NzcDANBgkqhkiG9w0BAQsFAAOCAYEA6JDiMHd8aWSlyIQ8tg/mEH7mIvSz\niXWfygMcyXP5sGRvrE0yo2lbNfr8y3KnOGkNYMqrKJ28VBXAPjx5zLrooHynLYua\nLEsHw6XzvQWiWvcstSkKhcVOGdDqTMhl2XEGvx+LHZYBWKlb7i+L/0fDl0EUestS\ne4Shh63DLJ+7RaMFqoI/CHO/Jer5R4+dIMR8KSTTBhjEGLwN6rsRNI7D7vsyqDV8\ntZmhMHNEo9jtrPR8+tAzp6BaumioukI75nkAXrKiB0GRXI/jRp94VqEZstWcQPqc\nxzRRyR2Htet4AVbUWnSq2TRWIyeIecgPVmHXgDPpFWrwi/hpysXqT9sN/QOsCa3a\n2IpyGeuieProOeXb5lG4pbwePz5dRRlY3WRvhWdQm+dRGRErJt42KC7JAfiYoSmV\nDfJjQL2S11oYZt048ZQFIsUpiSJTmsCLXURIEuccrKT+WXR7D+WNkYm8aJ/4s8Ub\n+B8Vv5GjCTO5LrjgVWGZtxOttN/uJ1ecgZpW\n-----END CERTIFICATE-----\n"
            }
            ],
            "realm_domains": [
                domain_name,
            ]
        }
    }
    return data


def get_create_stub_data(domain_name, domain_type):
    """Return an example payload to create a domain stub"""
    data = {
        "title": "Domain " + domain_name,
        "description": "Description Domain " + domain_name,
        "domain_type": domain_type,
        "auto_enrollment_enabled": True,
    }
    return data

def get_token(tokenb64):
    data = json.loads(base64.b64decode(tokenb64))
    return data["secret"]

def generate_request_id():
    return ''.join(random.choices(string.ascii_lowercase + string.digits, k=32))

if __name__ == "__main__":
    app_name = "idmsvc"
    service_version = "1"
    base_url = "http://localhost:8000/api/" + app_name + "/v" + service_version

    parser = argparse.ArgumentParser()
    parser.add_argument("-oid", "--org-id", dest="org_id" , help = f"Organization ID (default={DEFAULT_ORG_ID})")
    args = parser.parse_args()

    org_id = DEFAULT_ORG_ID
    if args.org_id:
        org_id = args.org_id

    # x-rh-identity headers
    xrhidgen_user = xrhidgen(
        "-is-active=true",
        "-is-org-admin=true",
        "-user-id", "test",
        "-username", "test",
        org_id=org_id, xrhidgen_type='user')
    xrhidgen_system = xrhidgen("-cn", "6f324116-b3d2-11ed-8a37-482Completed", "--cert-type", "system", org_id=org_id, xrhidgen_type='system')
    b64_identity_user = base64.b64encode(
        json.dumps(xrhidgen_user(),
                sort_keys=True).
        encode()
    ).decode()
    b64_identity_system = base64.b64encode(json.dumps(xrhidgen_system()).encode()).decode()

    # body payload
    domain_list = (
        "myorg.test",
        "department1.myorg.test",
        "department2.myorg.test",
        "department3.myorg.test",
    )

    for domain_name in domain_list:
        data = get_create_stub_data(domain_name, "rhel-idm")
        response = requests.post(
            base_url + "/domains",
            headers=({
                HEADER_X_RH_INSIGHTS_REQUEST_ID: generate_request_id(),
                HEADER_X_RH_IDENTITY: b64_identity_user,
                HEADER_CONTENT_TYPE: CONTENT_TYPE,
            }),
            data=json.dumps(data))
        if response.status_code >= 400:
            sys.exit("stub:%d:%s" % (response.status_code, response.content))
        token = get_token(response.headers["X-Rh-Idm-Rhelidm-Register-Token"])
        ipa_hcc_version = json.dumps({
            "ipa-hcc":"0.7",
            "ipa":"4.10.0-8.el9_1",
            "os-release-id":"rhel",
            "os-release-version-id":"9.1"
        })
        response = json.loads(response.content)

        # Register the domain information
        subscription_manager_id = "b0d5b30c-3765-11ee-804e-482ae3863d30"
        data = get_register_data(response["domain_id"], domain_name, subscription_manager_id)
        response = requests.put(
            "%s/domains/%s/register" % (base_url, response["domain_id"]),
            headers=({
                HEADER_X_RH_INSIGHTS_REQUEST_ID: generate_request_id(),
                HEADER_X_RH_IDENTITY: b64_identity_system,
                HEADER_X_RH_IDM_VERSION: ipa_hcc_version,
                HEADER_X_RH_IDM_REGISTRATION_TOKEN: token,
                HEADER_CONTENT_TYPE: CONTENT_TYPE,
            }),
            data=json.dumps(data))
        if response.status_code >= 400:
            sys.exit("register:%d:%s" % (response.status_code, response.content))
