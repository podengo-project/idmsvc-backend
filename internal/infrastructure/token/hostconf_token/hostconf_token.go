package hostconf_token

import (
	"crypto/rand"
	"encoding/base64"

	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
)

// BuildHostconfToken creates a token instance with all claims for
// host configuration token.
func BuildHostconfToken(
	rhsmId public.SubscriptionManagerId,
	orgId string,
	inventoryId public.HostId,
	fqdn public.Fqdn,
	domainId public.DomainId,
	validity time.Duration,
) (tok jwt.Token, err error) {
	// random JTI
	r := make([]byte, 6)
	if _, err = rand.Read(r); err != nil {
		return nil, err
	}
	jti := base64.RawURLEncoding.EncodeToString(r)

	now := time.Now()
	return jwt.NewBuilder().
		Issuer(TokenIssuer).
		Subject(rhsmId.String()).
		Audience([]string{AudJoinHost}).
		JwtID(jti).
		IssuedAt(now).
		NotBefore(now).
		Expiration(now.Add(validity)).
		Claim(ClaimOrgId, orgId).
		Claim(ClaimInventoryId, inventoryId.String()).
		Claim(ClaimFqdn, string(fqdn)).
		Claim(ClaimDomainId, domainId.String()).
		Build()
}

// signToken serializes the token and signs it with all given keys. The return
// value is a JWS in JSON format.
func SignToken(tok jwt.Token, keys []jwk.Key) ([]byte, error) {
	serialized, err := jwt.NewSerializer().Serialize(tok)
	if err != nil {
		return nil, err
	}

	opts := []jws.SignOption{}
	// sign with all keys
	for _, key := range keys {
		opts = append(opts, jws.WithKey(key.Algorithm(), key))
	}
	// always return JSON format (non-compact serialization)
	opts = append(opts, jws.WithJSON())
	return jws.Sign(serialized, opts...)
}

func init() {
	var s string
	jwt.RegisterCustomField(ClaimDomainId, s)
	jwt.RegisterCustomField(ClaimFqdn, s)
	jwt.RegisterCustomField(ClaimInventoryId, s)
	jwt.RegisterCustomField(ClaimOrgId, s)
}
