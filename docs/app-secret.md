# Backend App Secrets

The backend service requires secrets for some cryptographic operations. At
the moment, the backend needs a secret for the

- key for domain registration token. The token is generated with HMAC-SHA256.
- encryption key to store private JWKs in the database. Key data is encrypted
  with AES-GCM.

The security of the token and key encryption depends on the strength of the
secrets. The secret must be unpredictable and should be created with a
cryptographically secure random number generator (CSRNG).

To make configuration and deployment simple, the backend borrows a trick from
TLS 1.3. It uses [RFC 5869](https://datatracker.ietf.org/doc/html/rfc5869)
HKDF to derive secrets from an initial main secret. HKDF stands for HMAC-based
key derivation function. The algorithm extracts a pseudo-random key from an
input value, then expands the PRK into domain-specific keys. It is considered
good practice to use distinct secrets instead of reusing the same secret value
for different parts of the application

## Implementation

The config value `app.secret` / `APP_SECRET` is an URL-safe base64 encoded
byte string without padding (`base64.RawURLEncoding`). The byte string must
be at least 16 bytes long and should be created from a CSRNG.

Next, a PRK is extracted from the app secret with `HKDF_Extract`. The extract
step uses `SHA-256` as hashing algorithm and static salt `"idmsvc-backend"`.

Finally, the PRK is expanded into domain-specific secrets.

```
ikm = base64_rawurl_decode(app_secret)
prk = HKDF_Extract(sha256, secret, "idmsvc-backend")
key = HDKF_Expand(sha256, prk, "domain registration key", 32)
```

## Keys

- Domain registration token secret (input for HMAC-SHA256)
  HKDF info: "domain registration key"
  Length:    32 bytes

- JWK encryption key (input for AES-GCM AEAD)
  HKDF info: "JWK encryption key"
  Length:    16 bytes
