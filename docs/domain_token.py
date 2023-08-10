"""Reference implementation of domain registration token"""

import base64
import hashlib
import hmac
import sys
import time
import typing
import uuid

PERSONALITY = b"register domain"

# namespace UUID for IDMSVC
# uuid.uuid5(uuid.NAMESPACE_URL, "https://console.redhat.com/api/idmsvc")
NS_IDMSVC = uuid.UUID("2978cc95-31c8-503d-ba8f-581911b6bea0")

# Python 3.6 does not have time.time_ns()
if sys.version_info < (3, 7):

    def _time_ns():
        return int(time.time() * 1_000_000_000)

else:
    _time_ns = time.time_ns


def token_domain_id(token: str) -> uuid.UUID:
    """Get domain UUID from a token string"""
    return uuid.uuid5(NS_IDMSVC, token)


def generate_token(
    key: bytes, domain_type: str, org_id: str, *, validity: int = 2 * 60 * 60
) -> typing.Tuple[str, int]:
    """Generate a domain registration token"""
    expires = int(_time_ns() + (validity * 1_000_000_000))
    token = _generate_token_ns(key, domain_type, org_id, expires)
    return token, expires


def _generate_token_ns(key: bytes, domain_type: str, org_id: str, expires: int) -> str:
    payload_bytes = expires.to_bytes(8, "big")
    payload_b64 = _b64encode(payload_bytes)
    sig = _mac_digest(key, domain_type, org_id, payload_bytes)
    sig_b64 = _b64encode(sig)
    return f"{payload_b64}.{sig_b64}"


def validate_token(key: bytes, domain_type: str, org_id: str, token: str) -> uuid.UUID:
    """Validate a domain registration token"""
    expires = _validate_token_sig(key, domain_type, org_id, token)
    if _time_ns() > expires:
        raise ValueError("token expired")
    return token_domain_id(token)


def _validate_token_sig(key: bytes, domain_type: str, org_id: str, token: str) -> int:
    payload_b64, sig_b64 = token.split(".")
    payload_bytes = _b64decode(payload_b64)
    sig = _b64decode(sig_b64)
    digest = _mac_digest(key, domain_type, org_id, payload_bytes)
    if not hmac.compare_digest(sig, digest):
        raise ValueError("Invalid signature")
    return int.from_bytes(payload_bytes, "big")


def _b64encode(payload: bytes) -> str:
    """URL-safe base64 encoding without padding"""
    return base64.urlsafe_b64encode(payload).rstrip(b"=").decode("ascii")


def _b64decode(payload: str) -> bytes:
    """URL-safe base64 decoding without padding"""
    rem = len(payload) % 4
    if rem == 2:
        payload += "=="
    elif rem == 3:
        payload += "="
    return base64.urlsafe_b64decode(payload.encode("utf-8"))


def _mac_digest(
    key: bytes, domain_type: str, org_id: str, payload_bytes: bytes
) -> bytes:
    """Create MAC from org_id and payload"""
    mac = hmac.HMAC(key, digestmod=hashlib.sha256)
    mac.update(PERSONALITY)
    mac.update(domain_type.encode("utf-8"))
    mac.update(org_id.encode("utf-8"))
    mac.update(payload_bytes)
    return mac.digest()


def test() -> None:
    KEY = b"secretkey"
    ORG_ID = "123456"
    DOMAIN_TYPE = "rhel-idm"
    token, expires = generate_token(KEY, DOMAIN_TYPE, ORG_ID)
    print(token, expires)
    domain_id = validate_token(KEY, DOMAIN_TYPE, ORG_ID, token)
    print(domain_id)
    # validate_token(KEY, "123123", token)


if __name__ == "__main__":
    test()
