package header

import (
	b64 "encoding/base64"
	"testing"

	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeIdentity(t *testing.T) {
	id1 := &identity.XRHID{
		Identity: identity.Identity{
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
		},
	}

	s := EncodeXRHID(id1)
	id2, err := DecodeXRHID(s)
	require.NoError(t, err)
	require.NotNil(t, id2)
	assert.Equal(t, *id1, *id2)
}

func TestEncodeIdentity(t *testing.T) {
	var result string

	result = EncodeXRHID(nil)
	assert.Equal(t, "", result)

	result = EncodeXRHID(&identity.XRHID{
		Identity: identity.Identity{
			OrgID:                 "12345",
			AccountNumber:         "12345",
			EmployeeAccountNumber: "12345",
			Type:                  "User",
			AuthType:              "basic-auth",
			Internal: identity.Internal{
				OrgID:       "12345",
				CrossAccess: false,
			},
		},
	})
}

func TestDecodeIdentity(t *testing.T) {
	var (
		result *identity.XRHID
		err    error
	)

	result, err = DecodeXRHID("")
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, "'"+HeaderXRHID+"' is an empty string")

	result, err = DecodeXRHID("abc")
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, "illegal base64 data at input byte 0")

	result, err = DecodeXRHID("ewo=")
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, "unexpected end of JSON input")

	result, err = DecodeXRHID(EncodeXRHID(&identity.XRHID{
		Identity: identity.Identity{
			Internal: identity.Internal{
				OrgID: "12345",
			},
		},
	}))
	assert.NotNil(t, result)
	require.NoError(t, err)
	assert.Equal(t, "12345", result.Identity.Internal.OrgID)
}

func TestSystemIdentity(t *testing.T) {

	payload := `{
		"identity": {
		  "account_number": "11111",
		  "auth_type": "cert-auth",
		  "internal": {
			"auth_time": 900,
			"cross_access": false,
			"org_id": "12345"
		  },
		  "org_id": "12345",
		  "system": {
			"cert_type": "system",
			"cn": "6f324116-b3d2-11ed-8a37-482ae3863d30"
		  },
		  "type": "System"
		}
	  }`

	b64Identity := b64.StdEncoding.EncodeToString([]byte(payload))
	result, err := DecodeXRHID(b64Identity)
	assert.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "12345", result.Identity.OrgID)
	assert.Equal(t, "System", result.Identity.Type)
	assert.Equal(t, "11111", result.Identity.AccountNumber)
	assert.Equal(t, "cert-auth", result.Identity.AuthType)

	assert.Equal(t, float32(900), result.Identity.Internal.AuthTime)
	assert.Equal(t, false, result.Identity.Internal.CrossAccess)
	assert.Equal(t, "12345", result.Identity.Internal.OrgID)

	assert.Equal(t, "system", result.Identity.System.CertType)
}
