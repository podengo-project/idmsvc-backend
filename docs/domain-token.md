# Domain Registration Token

Domain registration tokens are used to register a new identity domain with
*Directory & Domain Services* on the Hybrid Cloud Console. The token enables
a privileged user to create a domain entry without entering any additional
information in the Console. Instead an administrators runs a single command
on a RHEL IdM server. The `ipa-hcc register` command uses the host's
subscription manager certificate/key pair to authenticate to Console with
mTLS (mutual TLS client cert authentication). The domain registration token
is a short-lived, one-time-use token that authorizes the host to register a
single domain.

Example:

```sh
ipa-hcc register F3kVxQP4sIs.cjbtH-GB8JuszfqrQnnudLoLzJH3zkw5jnhmTgKP_HU
```

The token is structured someward similar like a JWT with a payload and
signature in base64 encoding. However it uses raw binary encoding in order to
archive a short, compact notation.


## Security considerations

* Token must be signed with a secret key. The consumer must verify the
  signature.
* Tokens must expire after a short time span (hours). Therefore the token
  payload must contain either a creation time stamp or an expiration time
  stamp.
* Expiration time stamp must be Y2038-safe.
* Tokens must be unique (with very high probability). Multiple token requests
  in fast succession should result in different tokens. This can be
  achieved by either including a random value or by using an expiration time
  stamp with high change rate (nanoseconds since epoch). The extremely
  unlikely case of a token collision is treated like token reuse.
* Each token must be bound to the organization of the user that requested
  the token. Only hosts from the same organization can use the token to
  register a domain.
* Token reuse must be prevented. A token cannot register more than one
  domain.

The token generator and token consumer are the same system. Therefore no
complex asymmetric signature scheme is required. A symmetric signature with
HMAC-SHA256 is an easy and efficient solution.


## Token definition

### Payload

The payload is just the *expiration* time stamp as a serialized `uin64_t`
number in *big endian* notation. The value encodes nanoseconds since
beginning of the Unix epoch. The value won't overflow until year 2554.
Nanoseconds resolution may look excessive at first glance. However the clock
source makes it extremely unlikely to run into collisions in parallel
requests.

A typical token is valid for a time span of 10 minutes to 2 hours.

### Signature

The token is generated and verified in `idmsvc-backend`. Therefore the
complexity and key management of an asymmetric signature scheme is not
required. A symmetric signature with HMAC-SHA256 is an easy and efficient
solution.

The MAC digest is calculated over a personality string (UTF-8 bytes), the
*domain type* (UTF-8 bytes), the *organization id* (UTF-8 bytes), and the
expiration time stamp (`uint64_t` as 8 bytes in *big endian* notation).
The *personality* string `"register domain"` and *domain type* string bind
the MAC to a specific purpose and domain type. The *organization
id* is included, because the value is not part of the payload. Instead the
organization id is transmitted out-of-band in the `X-Rh-Identity` header.


### HMAC secret key generation and rotation

The secret key for HMAC should be a random key of size *32* bytes
(the output length of *SHA-256* hash function). Keys must be created with
a cryptographically strong pseudo-number generator. See
[RFC 2104](https://datatracker.ietf.org/doc/html/rfc2104#section-3) for
more details.

Tokens are short-lived in the range of minutes to an hour or two. Key
rotation may not be necessary. In case key rotation is preferred, the
application can simply try all known keys. HMAC is cheap to calculate.


### Token representation

The token is represented similar to a JSON Web Token (JWT). A domain
registration token is represented by two URL-safe base64 encoded strings that
are separated by a period `.`. Paddding (`=`) is omitted. The first part is
the payload and the second part is the signature.

Example: `F3kVxQP4sIs.cjbtH-GB8JuszfqrQnnudLoLzJH3zkw5jnhmTgKP_HU`


### Token to domain id

[RFC 4122](https://datatracker.ietf.org/doc/html/rfc4122.html#section-4.3)
defines multiple classes of UUIDs. Version 5 UUIDs are generated from the
SHA-1 digest of a namespace UUID and an arbitrary name. This allows us to
dervice a stable UUID from a token. Since tokens are pseudo-random values,
the derived UUIDs are sufficiently randomized, too. The *domain id* is
derived from the entire token, not just the *payload*, to bind the token
to *organization id*, too.

First we need our own namespace UUID. The RFC defines namespace id
`6ba7b811-9dad-11d1-80b4-00c04fd430c8` for names, which are URL string. This
lets us create a new namespace UUID from a well-known URL.

```python
>>> import uuid
>>> url = "https://console.redhat.com/api/idmsvc"
>>> NS_IDMSVC = uuid.uuid5(uuid.NAMESPACE_URL, url)
>>> NS_IDMSVC
UUID('2978cc95-31c8-503d-ba8f-581911b6bea0')
```

```python
>>> token = "F3kVxQP4sIs.cjbtH-GB8JuszfqrQnnudLoLzJH3zkw5jnhmTgKP_HU"
>>> uuid.uuid5(NS_IDMSVC, token)
UUID('681abfd7-18ce-51b3-a9cc-10d386c8dc35')
```

### Logging

The domain registration token does not contain any information, which user
account has created the token. For security and compliance reason, it
useful to track the user. Therefore the application should log the domain's
UUID, current user acount, and other metadata, whenever a user requests a
domain registration token.


## Attack scenarios

1. A malicious user transmits an ill-formed or overly long token. The
   validation has to limit input size and gracefully fail when input
   cannot be parsed correctly.
2. User sends an expired token. Validation code has to check validity by
   comparing the expiration time stamp with current time stamp.
3. User attempts to register host for a different organization. MAC
   validation fails, because the host's org id does not match the org id
   that was used to create the signature of the token.
4. User attempts to register a second domain with a token. The first
   registration call added a domain to the database. The token of the
   second attempt is the same, so is the derived UUID of the domain in the
   database. A unique constraint on `domain_id` prevents the second insert.
5. User attempts to forge a token. HMAC prevents forgery unless the user
   is able to get hold of the secret key.


## Example values

```python
key = b"secretkey"
domain_type = "rhel-idm"
org_id = "123456"
expiration = 1691662998988903762
token = "F3n-iOZn1VI.wbzIH7v-kRrdvfIvia4nBKAvEpIKGdv6MSIFXeUtqVY"
domain_id = "7b160558-8273-5a24-b559-6de3ff053c63"
```
