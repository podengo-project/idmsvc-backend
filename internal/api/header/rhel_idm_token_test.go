package header

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeRhelIdmToken(t *testing.T) {
	var (
		token *RhelIdmToken
		err   error
	)
	token, err = DecodeRhelIdmToken("")
	assert.EqualError(t, err, "'data' is empty")
	assert.Nil(t, token)

	token, err = DecodeRhelIdmToken("{}")
	assert.EqualError(t, err, "illegal base64 data at input byte 0")
	assert.Nil(t, token)

	token, err = DecodeRhelIdmToken(base64.StdEncoding.EncodeToString([]byte("{")))
	assert.EqualError(t, err, "unexpected end of JSON input")
	assert.Nil(t, token)

	token, err = DecodeRhelIdmToken(base64.StdEncoding.EncodeToString([]byte("{}")))
	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestEncodeRhelIdmToken(t *testing.T) {
	var (
		tokenStr string
		err      error
	)

	tokenStr, err = EncodeRhelIdmToken(nil)
	assert.Equal(t, "", tokenStr)
	assert.EqualError(t, err, "'data' is nil")

	tokenStr, err = EncodeRhelIdmToken(&RhelIdmToken{
		Secret:     nil,
		Expiration: nil,
	})
	assert.Equal(t, "e30=", tokenStr)
	assert.Nil(t, err)
}
