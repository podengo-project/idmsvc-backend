package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDefaultTokenExpiration(t *testing.T) {
	valueOld := DefaultTokenExpiration()
	valueNew := valueOld + time.Hour*24
	assert.Equal(t, valueOld, DefaultTokenExpiration())
	SetDefaultTokenExpiration(valueNew)
	assert.Equal(t, valueNew, DefaultTokenExpiration())
	SetDefaultTokenExpiration(valueOld)
	assert.Equal(t, valueOld, DefaultTokenExpiration())
}

func TestBeforeCreate(t *testing.T) {
	var (
		entity *Ipa
		err    error
	)
	entity = nil
	err = entity.BeforeCreate(nil)
	assert.EqualError(t, err, "'BeforeCreate' cannot be invoked on nil")

	entity = &Ipa{}
	err = entity.BeforeCreate(nil)
	require.NoError(t, err)
	assert.NotNil(t, entity.Token)
	assert.NotEqual(t, "", *entity.Token)
	assert.NotNil(t, entity.TokenExpiration)
	assert.NotEqual(t, time.Time{}, entity.TokenExpiration)
}

func TestAfterCreate(t *testing.T) {
	var (
		entity *Ipa
		err    error
	)
	entity = nil
	err = entity.AfterCreate(nil)
	assert.EqualError(t, err, "'AfterCreate' cannot be invoked on nil")

	entity = &Ipa{}
	err = entity.AfterCreate(nil)
	require.NoError(t, err)

	notValidBefore := time.Now()
	notValidAfter := notValidBefore.Add(24 * time.Hour)
	entity = &Ipa{
		Model: gorm.Model{
			ID: 1,
		},
		CaCerts: []IpaCert{
			{
				Nickname:       "MYDOMAIN.EXAMPLE IPA CA",
				Issuer:         "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
				Subject:        "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
				SerialNumber:   "1",
				NotValidBefore: notValidBefore,
				NotValidAfter:  notValidAfter,
				Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
				IpaID:          0,
			},
		},
		Servers: []IpaServer{
			{
				FQDN:                "ipaserver.mydomain.example",
				RHSMId:              "547ce70c-9eb5-4783-a619-086aa26f88e5",
				CaServer:            true,
				HCCEnrollmentServer: true,
				PKInitServer:        true,
				IpaID:               0,
			},
		},
	}
	err = entity.AfterCreate(nil)
	require.NoError(t, err)
	assert.Equal(t, uint(1), entity.CaCerts[0].IpaID)
	assert.Equal(t, uint(1), entity.Servers[0].IpaID)
}
