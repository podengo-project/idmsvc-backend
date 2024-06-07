package model

import (
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
	"gorm.io/gorm"
)

type HostconfJwk interface {
	Build() *model.HostconfJwk
	WithModel(value *gorm.Model) HostconfJwk
	WithKeyId(value string) HostconfJwk
	WithExpiresAt(value time.Time) HostconfJwk
	WithPublicJwk(value string) HostconfJwk
	WithEncryptionId(value string) HostconfJwk
	WithEncryptedJwk(value []byte) HostconfJwk
}

// gormModel is the specific builder implementation
type hostconfJwk struct {
	model.HostconfJwk
}

// NewModel generate a gorm.Model with random information
// overrided by the customized options.
func NewHostconfJwk() HostconfJwk {
	return &hostconfJwk{
		HostconfJwk: model.HostconfJwk{
			Model:        NewModel().Build(),
			KeyId:        "",
			ExpiresAt:    time.Now().UTC().Add(24 * time.Hour),
			PublicJwk:    "",
			EncryptionId: "",
			EncryptedJwk: []byte{},
		},
	}
}

func (b *hostconfJwk) Build() *model.HostconfJwk {
	return (*model.HostconfJwk)(&b.HostconfJwk)
}

func (b *hostconfJwk) WithModel(value *gorm.Model) HostconfJwk {
	b.HostconfJwk.Model = *value
	return b
}

func (b *hostconfJwk) WithKeyId(value string) HostconfJwk {
	b.HostconfJwk.KeyId = value
	return b
}

func (b *hostconfJwk) WithExpiresAt(value time.Time) HostconfJwk {
	b.HostconfJwk.ExpiresAt = value
	return b
}

func (b *hostconfJwk) WithPublicJwk(value string) HostconfJwk {
	b.HostconfJwk.PublicJwk = value
	return b
}

func (b *hostconfJwk) WithEncryptionId(value string) HostconfJwk {
	b.HostconfJwk.EncryptionId = value
	return b
}

func (b *hostconfJwk) WithEncryptedJwk(value []byte) HostconfJwk {
	b.HostconfJwk.EncryptedJwk = value
	return b
}
