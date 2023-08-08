package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNamespaceIDMSVC(t *testing.T) {
	exp := uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://console.redhat.com/api/idmsvc"))
	assert.Equal(t, exp, NamespaceIDMSVC)
}

func TestTokenDomainId(t *testing.T) {
	tok := DomainRegistrationToken("F3kVxQP4sIs.cjbtH-GB8JuszfqrQnnudLoLzJH3zkw5jnhmTgKP_HU")
	exp := uuid.MustParse("681abfd7-18ce-51b3-a9cc-10d386c8dc35")
	domain_id := TokenDomainId(tok)
	assert.Equal(t, exp, domain_id)
}

func TestNewDomainRegistrationToken(t *testing.T) {
	var (
		err        error
		expiration uint64                  = 1691407070973702283
		knownToken DomainRegistrationToken = "F3kVxQP4sIs.cjbtH-GB8JuszfqrQnnudLoLzJH3zkw5jnhmTgKP_HU"
		token      DomainRegistrationToken
		orgId      string = "123456"
		key        []byte = []byte("secretkey")
	)

	token, err = newDomainRegistrationTokenAt(key, orgId, expiration)
	assert.NoError(t, err)
	assert.Equal(t, token, knownToken)

	_, err = NewDomainRegistrationToken(key, orgId, time.Hour)
	assert.NoError(t, err)
}

func TestVerifyDomainRegistrationToken(t *testing.T) {
	var (
		err        error
		expiration uint64                  = 1691407070973702283
		knownToken DomainRegistrationToken = "F3kVxQP4sIs.cjbtH-GB8JuszfqrQnnudLoLzJH3zkw5jnhmTgKP_HU"
		token      DomainRegistrationToken
		orgId      string = "123456"
		otherOrgId string = "789789"
		key        []byte = []byte("secretkey")
		domainId   uuid.UUID
	)
	exp, err := parseDomainRegistrationToken(key, orgId, knownToken)
	assert.NoError(t, err)
	assert.Equal(t, expiration, exp)
	// known token has expired
	domainId, err = VerifyDomainRegistrationToken(key, orgId, token)
	assert.Error(t, err)
	assert.Equal(t, domainId, uuid.Nil)

	token, err = NewDomainRegistrationToken(key, orgId, time.Hour)
	assert.NoError(t, err)
	exp, err = parseDomainRegistrationToken(key, orgId, token)
	assert.NoError(t, err)
	domainId, err = VerifyDomainRegistrationToken(key, orgId, token)
	assert.NoError(t, err)
	assert.NotEqual(t, domainId, uuid.Nil)

	// wrong orgId == invalid signature
	domainId, err = VerifyDomainRegistrationToken(key, otherOrgId, token)
	assert.Error(t, err)
	assert.Equal(t, domainId, uuid.Nil)
}
