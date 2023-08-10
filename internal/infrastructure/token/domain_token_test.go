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
	tok := DomainRegistrationToken("F3n-iOZn1VI.wbzIH7v-kRrdvfIvia4nBKAvEpIKGdv6MSIFXeUtqVY")
	exp := uuid.MustParse("7b160558-8273-5a24-b559-6de3ff053c63")
	domain_id := TokenDomainId(tok)
	assert.Equal(t, exp, domain_id)
}

func TestNewDomainRegistrationToken(t *testing.T) {
	var (
		err        error
		expiration uint64                  = 1691662998988903762
		knownToken DomainRegistrationToken = "F3n-iOZn1VI.wbzIH7v-kRrdvfIvia4nBKAvEpIKGdv6MSIFXeUtqVY"
		token      DomainRegistrationToken
		domainType string = "rhel-idm"
		orgId      string = "123456"
		key        []byte = []byte("secretkey")
	)

	token, err = newDomainRegistrationTokenAt(key, domainType, orgId, expiration)
	assert.NoError(t, err)
	assert.Equal(t, token, knownToken)

	_, err = NewDomainRegistrationToken(key, domainType, orgId, time.Hour)
	assert.NoError(t, err)
}

func TestVerifyDomainRegistrationToken(t *testing.T) {
	var (
		err        error
		expiration uint64                  = 1691662998988903762
		knownToken DomainRegistrationToken = "F3n-iOZn1VI.wbzIH7v-kRrdvfIvia4nBKAvEpIKGdv6MSIFXeUtqVY"
		token      DomainRegistrationToken
		domainType string = "rhel-idm"
		orgId      string = "123456"
		otherOrgId string = "789789"
		key        []byte = []byte("secretkey")
		domainId   uuid.UUID
	)
	exp, err := parseDomainRegistrationToken(key, domainType, orgId, knownToken)
	assert.NoError(t, err)
	assert.Equal(t, expiration, exp)
	// known token has expired
	domainId, err = VerifyDomainRegistrationToken(key, domainType, orgId, token)
	assert.Error(t, err)
	assert.Equal(t, domainId, uuid.Nil)

	token, err = NewDomainRegistrationToken(key, domainType, orgId, time.Hour)
	assert.NoError(t, err)
	exp, err = parseDomainRegistrationToken(key, domainType, orgId, token)
	assert.NoError(t, err)
	domainId, err = VerifyDomainRegistrationToken(key, domainType, orgId, token)
	assert.NoError(t, err)
	assert.NotEqual(t, domainId, uuid.Nil)

	// wrong orgId == invalid signature
	domainId, err = VerifyDomainRegistrationToken(key, domainType, otherOrgId, token)
	assert.Error(t, err)
	assert.Equal(t, domainId, uuid.Nil)
}
