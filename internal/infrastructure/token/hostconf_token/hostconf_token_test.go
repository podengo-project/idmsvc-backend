package hostconf_token

import (
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
)

func testToken() (jwt.Token, error) {
	return BuildHostconfToken(
		test.Server1.CertUUID,
		test.OrgId,
		test.Server1.InventoryUUID,
		test.Server1.Fqdn,
		test.DomainUUID,
		time.Hour,
	)
}

func TestBuildToken(t *testing.T) {
	tok, err := testToken()
	assert.NoError(t, err)

	assert.Equal(t, tok.Issuer(), TokenIssuer)
	assert.Equal(t, tok.Subject(), test.Server1.CertCN)
	assert.Equal(t, tok.Audience(), []string{AudJoinHost})
	assert.NotEmpty(t, tok.JwtID())
	iat := tok.IssuedAt()
	assert.Equal(t, iat, tok.NotBefore())
	assert.Equal(t, iat.Add(time.Hour), tok.Expiration())

	ifc, ok := tok.Get(ClaimOrgId)
	assert.True(t, ok)
	assert.Equal(t, ifc.(string), test.OrgId)

	ifc, ok = tok.Get(ClaimInventoryId)
	assert.True(t, ok)
	assert.Equal(t, ifc.(string), test.Server1.InventoryId)

	ifc, ok = tok.Get(ClaimFqdn)
	assert.True(t, ok)
	assert.Equal(t, ifc.(string), test.Server1.Fqdn)

	ifc, ok = tok.Get(ClaimDomainId)
	assert.True(t, ok)
	assert.Equal(t, ifc.(string), test.DomainId)
}

func TestSignToken(t *testing.T) {
	tok, err := testToken()
	assert.NoError(t, err)

	exp := time.Now().Add(time.Hour)
	priv1, err := hostconf_jwk.GeneratePrivateJWK(exp)
	assert.NoError(t, err)
	pub1, err := priv1.PublicKey()
	assert.NoError(t, err)

	priv2, err := hostconf_jwk.GeneratePrivateJWK(exp)
	assert.NoError(t, err)
	pub2, err := priv2.PublicKey()
	assert.NoError(t, err)

	privs := []jwk.Key{priv1, priv2}
	pubs := []jwk.Key{pub1, pub2}

	sig, err := SignToken(tok, privs)
	assert.NoError(t, err)

	set := jwk.NewSet()
	for _, pub := range pubs {
		err = set.AddKey(pub)
		assert.NoError(t, err)
	}
	verified, err := jws.Verify(sig, jws.WithKeySet(set))
	assert.NoError(t, err)

	toks, err := jwt.NewSerializer().Serialize(tok)
	assert.NoError(t, err)
	assert.Equal(t, verified, toks)

	verified, err = jws.Verify(sig, jws.WithKey(priv1.Algorithm(), pub1))
	assert.NoError(t, err)
	assert.Equal(t, verified, toks)

	verified, err = jws.Verify(sig, jws.WithKey(priv2.Algorithm(), pub2))
	assert.NoError(t, err)
	assert.Equal(t, verified, toks)
}
