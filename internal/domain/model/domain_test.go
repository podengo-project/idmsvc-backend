package model

import (
	"testing"

	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDomainTypeString(t *testing.T) {
	assert.Equal(t, DomainTypeUndefinedString, DomainTypeString(DomainTypeUndefined))
	assert.Equal(t, DomainTypeUndefinedString, DomainTypeString(1000))
	assert.Equal(t, DomainTypeIpaString, DomainTypeString(DomainTypeIpa))
}

func TestDomainTypeUint(t *testing.T) {
	assert.Equal(t, DomainTypeUndefined, DomainTypeUint(""))
	assert.Equal(t, DomainTypeUndefined, DomainTypeUint("anything"))
	assert.Equal(t, DomainTypeIpa, DomainTypeUint(DomainTypeIpaString))
}

func TestDomainBeforeCreate(t *testing.T) {
	var item *Domain
	item = &Domain{}
	item.BeforeCreate(nil)

	assert.NotEmpty(t, item.DomainUuid)
	assert.NotEmpty(t, item.CreatedAt)
	assert.NotEmpty(t, item.UpdatedAt)
}

func TestDomainAfterCreate(t *testing.T) {
	var item *Domain

	item = &Domain{}
	assert.EqualError(t, item.AfterCreate(nil), "'Type' is nil")

	item = &Domain{
		Model: gorm.Model{
			ID: 1,
		},
		Type: pointy.Uint(DomainTypeIpa),
		IpaDomain: &Ipa{
			Model: gorm.Model{
				ID: 0,
			},
		},
	}
	require.NoError(t, item.AfterCreate(nil))
	assert.Equal(t, item.ID, item.IpaDomain.ID)
}
