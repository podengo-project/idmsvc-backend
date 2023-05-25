package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateXRhIdentity(t *testing.T) {
	var err error
	v := xrhiAlwaysTrue{}
	require.NotPanics(t, func() {
		err = v.ValidateXRhIdentity(nil)
	})
	assert.NoError(t, err)
}

func TestNewApiServiceValidator(t *testing.T) {
	v := NewApiServiceValidator()
	assert.NotNil(t, v)
}
