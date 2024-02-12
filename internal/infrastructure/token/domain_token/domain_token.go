/* Domain Registration Token
 *
 * see docs/domain-token.md for more information
 */
package domain_token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DomainRegistrationToken string

var (
	// uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://console.redhat.com/api/idmsvc"))
	NamespaceIDMSVC           = uuid.MustParse("2978cc95-31c8-503d-ba8f-581911b6bea0")
	RegisterDomainPersonality = []byte("register domain")
)

// Derive domain id from token string
func TokenDomainId(token DomainRegistrationToken) uuid.UUID {
	return uuid.NewSHA1(NamespaceIDMSVC, []byte(token))
}

// Create a new domain registration token
// The token is signed by *key*, bound to *orgID*, and validate until
// now + validity duration.
func NewDomainRegistrationToken(
	key []byte, domainType string, orgID string, validity time.Duration,
) (token DomainRegistrationToken, expireNS uint64, err error) {
	expireNS = uint64(time.Now().UnixNano() + validity.Nanoseconds())
	tok, err := newDomainRegistrationTokenAt(key, domainType, orgID, uint64(expireNS))
	return tok, expireNS, err
}

// Create a domain registration token that expires at *expireNS* nanoseconds
// after Unix epoch.
func newDomainRegistrationTokenAt(
	key []byte, domainType string, orgID string, expireNS uint64,
) (token DomainRegistrationToken, err error) {
	payload_bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(payload_bytes, expireNS)
	payload_b64 := base64.RawURLEncoding.EncodeToString(payload_bytes)

	sig := mac_digest(key, domainType, orgID, payload_bytes)
	sig_b64 := base64.RawURLEncoding.EncodeToString(sig)

	return DomainRegistrationToken(fmt.Sprintf("%s.%s", payload_b64, sig_b64)), nil
}

// Verify signature, *orgID* binding, and expiration time stamp of a token.
// Returns the domain UUID on success.
func VerifyDomainRegistrationToken(
	key []byte, domainType string, orgID string, token DomainRegistrationToken,
) (domainId uuid.UUID, err error) {
	var expireNS uint64
	if expireNS, err = parseDomainRegistrationToken(key, domainType, orgID, token); err != nil {
		return uuid.Nil, err
	}
	var now uint64 = uint64(time.Now().UnixNano())
	if now > expireNS {
		return uuid.Nil, fmt.Errorf("Token has expired: %d > %d", now, expireNS)
	}
	return TokenDomainId(token), nil
}

// Parse and check signature of token
func parseDomainRegistrationToken(
	key []byte, domainType string, orgID string, token DomainRegistrationToken,
) (expireNS uint64, err error) {
	var (
		payload_bytes []byte
		sig           []byte
	)
	if len(token) > 100 {
		return 0, fmt.Errorf("Token length exceeded")
	}
	parts := strings.Split(string(token), ".")
	if len(parts) != 2 {
		return 0, fmt.Errorf("Invalid token")
	}
	if payload_bytes, err = base64.RawURLEncoding.DecodeString(parts[0]); err != nil {
		return 0, err
	}
	if sig, err = base64.RawURLEncoding.DecodeString(parts[1]); err != nil {
		return 0, err
	}
	expected_sig := mac_digest(key, domainType, orgID, payload_bytes)
	if !hmac.Equal(sig, expected_sig) {
		return 0, fmt.Errorf("Signature mismatch")
	}
	return binary.BigEndian.Uint64(payload_bytes), nil
}

// Calculate keyed MAC digest from orgID and payload
func mac_digest(key []byte, domainType string, orgID string, payload []byte) []byte {
	mac := hmac.New(sha256.New, key)
	// Hash.Write() never returns an error
	mac.Write(RegisterDomainPersonality)
	mac.Write([]byte(domainType))
	mac.Write([]byte(orgID))
	mac.Write(payload)
	return mac.Sum(nil)
}
