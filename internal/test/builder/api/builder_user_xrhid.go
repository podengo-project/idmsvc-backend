package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

type UserXRHID interface {
	Build() identity.XRHID
	WithAccountNumber(value string) UserXRHID
	WithEmployeeAccountNumber(value string) UserXRHID
	WithOrgID(value string) UserXRHID
	WithAuthType(value string) UserXRHID

	WithUserID(value string) UserXRHID
	WithUsername(value string) UserXRHID
	WithEmail(value string) UserXRHID
	WithActive(value bool) UserXRHID
	WithInternal(value bool) UserXRHID
	WithOrgAdmin(value bool) UserXRHID
	WithLocale(value string) UserXRHID
	// TODO Add more methods as they are needed
}

type userXRHID identity.XRHID

// NewUserXRHID create a new builder to get a XRHID for a user
func NewUserXRHID() UserXRHID {
	orgID := strconv.Itoa(int(builder_helper.GenRandNum(1, 100000)))
	firstName := builder_helper.GenRandFirstName()
	lastName := builder_helper.GenRandLastName()
	username := fmt.Sprintf("%s.%s",
		strings.ToLower(firstName),
		strings.ToLower(lastName),
	)
	userID := builder_helper.GenRandUserID()
	email := fmt.Sprintf("%s@acme.test", username)

	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/schema.json
	return &userXRHID{
		Identity: identity.Identity{
			AccountNumber:         strconv.Itoa(int(builder_helper.GenRandNum(1, 100000))),
			EmployeeAccountNumber: strconv.Itoa(int(builder_helper.GenRandNum(1, 100000))),
			OrgID:                 orgID,
			Type:                  "User",
			AuthType:              authType[int(builder_helper.GenRandNum(0, 1))],
			Internal: identity.Internal{
				OrgID:    orgID,
				AuthTime: float32(time.Now().Sub(time.Time{})),
			},
			User: &identity.User{
				UserID:    userID,
				Username:  username,
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Active:    builder_helper.GenRandBool(),
			},
		},
	}
}

func (b *userXRHID) Build() identity.XRHID {
	return identity.XRHID(*b)
}

func (b *userXRHID) WithOrgID(value string) UserXRHID {
	b.Identity.OrgID = value
	b.Identity.Internal.OrgID = value
	return b
}

func (b *userXRHID) WithAuthType(value string) UserXRHID {
	switch value {
	case authType[0]:
		b.Identity.AuthType = value
		return b
	case authType[1]:
		b.Identity.AuthType = value
		return b
	default:
		panic(fmt.Sprintf("value='%s' is not valid", value))
	}
}

func (b *userXRHID) WithAccountNumber(value string) UserXRHID {
	b.Identity.AccountNumber = value
	return b
}

func (b *userXRHID) WithEmployeeAccountNumber(value string) UserXRHID {
	b.Identity.EmployeeAccountNumber = value
	return b
}

// --- Start specific User data

func (b *userXRHID) WithUserID(value string) UserXRHID {
	b.Identity.User.UserID = value
	return b
}

func (b *userXRHID) WithUsername(value string) UserXRHID {
	b.Identity.User.Username = value
	return b
}

func (b *userXRHID) WithEmail(value string) UserXRHID {
	b.Identity.User.Email = value
	return b
}

func (b *userXRHID) WithActive(value bool) UserXRHID {
	b.Identity.User.Active = value
	return b
}

func (b *userXRHID) WithInternal(value bool) UserXRHID {
	b.Identity.User.Internal = value
	return b
}

func (b *userXRHID) WithOrgAdmin(value bool) UserXRHID {
	b.Identity.User.OrgAdmin = value
	return b
}

func (b *userXRHID) WithLocale(value string) UserXRHID {
	b.Identity.User.Locale = value
	return b
}
