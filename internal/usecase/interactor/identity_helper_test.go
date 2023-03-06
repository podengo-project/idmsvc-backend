package interactor

import (
	"testing"

	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeIdentity(t *testing.T) {
	id1 := &identity.Identity{
		AccountNumber:         "11111",
		EmployeeAccountNumber: "22222",
		OrgID:                 "12345",
		Type:                  "User",
		User: identity.User{
			Username:  "jdoe",
			Email:     "jdoe@example.com",
			FirstName: "Jhon",
			LastName:  "Doe",
			Active:    true,
			OrgAdmin:  true,
			Locale:    "en",
			UserID:    "1987348",
			Internal:  true,
		},
	}

	s := EncodeIdentity(id1)
	id2, err := DecodeIdentity(s)
	require.NoError(t, err)
	require.NotNil(t, id2)
	assert.Equal(t, *id1, *id2)
}

func TestEncodeIdentity(t *testing.T) {
	var result string

	result = EncodeIdentity(nil)
	assert.Equal(t, "", result)

	result = EncodeIdentity(&identity.Identity{
		OrgID:                 "12345",
		AccountNumber:         "12345",
		EmployeeAccountNumber: "12345",
		Type:                  "User",
		AuthType:              "basic-auth",
		Internal: identity.Internal{
			OrgID:       "12345",
			CrossAccess: false,
		},
	})
}

func TestDecodeIdentity(t *testing.T) {
	var (
		result *identity.Identity
		err    error
	)

	result, err = DecodeIdentity("")
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, "X-Rh-Identity content cannot be an empty string")

	result, err = DecodeIdentity("abc")
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, "illegal base64 data at input byte 0")

	result, err = DecodeIdentity("ewo=")
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, "unexpected end of JSON input")

	result, err = DecodeIdentity(EncodeIdentity(&identity.Identity{
		Internal: identity.Internal{
			OrgID: "12345",
		},
	}))
	assert.NotNil(t, result)
	require.NoError(t, err)
	assert.Equal(t, "12345", result.Internal.OrgID)
}
