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
