package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppSecret(t *testing.T) {
	var (
		err error
		sec *AppSecrets
	)
	sec, err = NewAppSecrets(GenerateRandomMainSecret())
	assert.NoError(t, err)
	assert.NotNil(t, sec.DomainRegKey)
	assert.NotNil(t, sec.HostConfEncryptionKey)
	assert.NotEmpty(t, sec.HostconfEncryptionId)

	sec, err = NewAppSecrets("short")
	assert.Nil(t, sec)
	assert.Error(t, err)
}
