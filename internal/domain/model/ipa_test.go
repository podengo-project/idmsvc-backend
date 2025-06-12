package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

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

	NotBefore := time.Now()
	NotAfter := NotBefore.Add(24 * time.Hour)
	entity = &Ipa{
		Model: gorm.Model{
			ID: 1,
		},
		CaCerts: []IpaCert{
			{
				Nickname:     "MYDOMAIN.EXAMPLE IPA CA",
				Issuer:       "CN=Issuer Authority, O=MYDOMAIN.EXAMPLE",
				Subject:      "CN=Subject, O=MYDOMAIN.EXAMPLE",
				SerialNumber: "1",
				NotBefore:    NotBefore,
				NotAfter:     NotAfter,
				Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
				IpaID:        0,
			},
		},
		Servers: []IpaServer{
			{
				FQDN:                "ipaserver.mydomain.example",
				RHSMId:              pointy.String("547ce70c-9eb5-4783-a619-086aa26f88e5"),
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
