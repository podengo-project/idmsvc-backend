"""Reference implementation of domain registration token"""

import base64
import hashlib
import hmac
import time
import uuid


def b64encode(payload: bytes) -> str:
    """URL-safe base64 encoding without padding"""
    return base64.urlsafe_b64encode(payload).rstrip(b"=").decode("ascii")


def b64decode(payload: bytes) -> str:
    """URL-safe base64 decoding without padding"""
    rem = len(payload) % 4
    if rem == 2:
        payload += "=="
    elif rem == 3:
        payload += "="
    return base64.urlsafe_b64decode(payload.encode("utf-8"))


PERSONALITY = b"register domain"

# namespace UUID for IDMSVC
NS_IDMSVC = uuid.uuid5(uuid.NAMESPACE_URL, "https://console.redhat.com/api/idmsvc")


def get_domain_id(token: str) -> uuid.UUID:
    """Get domain UUID from a token string"""
    return uuid.uuid5(NS_IDMSVC, token)


def mac_digest(key: bytes, org_id: str, payload_bytes: bytes) -> bytes:
    """Create MAC from org_id and payload"""
    mac = hmac.HMAC(key, digestmod=hashlib.sha256)
    mac.update(PERSONALITY)
    mac.update(org_id.encode("utf-8"))
    mac.update(payload_bytes)
    return mac.digest()


def generate_token(key: bytes, org_id: str, expires: int = 2 * 60 * 60) -> str:
    """Generate a domain registration token"""
    expiration = time.time_ns() + (expires * 1_000_000_000)
    payload_bytes = expiration.to_bytes(8, "big")
    payload_b64 = b64encode(payload_bytes)

    sig = mac_digest(key, org_id, payload_bytes)
    sig_b64 = b64encode(sig)
    token = f"{payload_b64}.{sig_b64}"
    return token


def validate_token(key: bytes, org_id: str, token: str) -> uuid.UUID:
    """Validate a domain registration token"""
    payload_b64, sig_b64 = token.split(".")
    payload_bytes = b64decode(payload_b64)
    sig = b64decode(sig_b64)
    digest = mac_digest(key, org_id, payload_bytes)
    if not hmac.compare_digest(sig, digest):
        raise ValueError("Invalid signature")
    expiration = int.from_bytes(payload_bytes, "big")
    if time.time_ns() > expiration:
        raise ValueError("token expired")
    return get_domain_id(token)


def test() -> None:
    KEY = b"secretkey"
    ORG_ID = "123456"
    token = generate_token(KEY, ORG_ID)
    print(token)
    domain_id = validate_token(KEY, ORG_ID, token)
    print(domain_id)
    # validate_token(KEY, "123123", token)


if __name__ == "__main__":
    test()
