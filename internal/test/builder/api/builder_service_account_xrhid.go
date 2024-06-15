package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

type ServiceAcccountXRHID interface {
	Build() identity.XRHID
	WithAccountNumber(value string) ServiceAcccountXRHID
	WithEmployeeAccountNumber(value string) ServiceAcccountXRHID
	WithOrgID(value string) ServiceAcccountXRHID
	WithAuthType(value string) ServiceAcccountXRHID

	WithClientID(value string) ServiceAcccountXRHID
	WithUsername(value string) ServiceAcccountXRHID
}

type serviceAcccountXRHID identity.XRHID

// NewServiceAccountXRHID create a new builder to get a XRHID for a user
func NewServiceAccountXRHID() ServiceAcccountXRHID {
	var orgID = strconv.Itoa(int(builder_helper.GenRandNum(1, 100000)))
	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/schema.json
	return &serviceAcccountXRHID{
		Identity: identity.Identity{
			AccountNumber:         strconv.Itoa(int(builder_helper.GenRandNum(1, 100000))),
			EmployeeAccountNumber: strconv.Itoa(int(builder_helper.GenRandNum(1, 100000))),
			OrgID:                 orgID,
			Type:                  "ServiceAccount",
			AuthType:              "jwt-auth",
			Internal: identity.Internal{
				OrgID:    orgID,
				AuthTime: float32(time.Now().Sub(time.Time{})),
			},
			ServiceAccount: &identity.ServiceAccount{
				ClientId: uuid.NewString(),
				Username: helper.GenRandUsername(),
			},
		},
	}
}

func (b *serviceAcccountXRHID) Build() identity.XRHID {
	return identity.XRHID(*b)
}

func (b *serviceAcccountXRHID) WithOrgID(value string) ServiceAcccountXRHID {
	b.Identity.OrgID = value
	b.Identity.Internal.OrgID = value
	return b
}

func (b *serviceAcccountXRHID) WithAuthType(value string) ServiceAcccountXRHID {
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

func (b *serviceAcccountXRHID) WithAccountNumber(value string) ServiceAcccountXRHID {
	b.Identity.AccountNumber = value
	return b
}

func (b *serviceAcccountXRHID) WithEmployeeAccountNumber(value string) ServiceAcccountXRHID {
	b.Identity.EmployeeAccountNumber = value
	return b
}

// --- Start specific System data

func (b *serviceAcccountXRHID) WithClientID(value string) ServiceAcccountXRHID {
	b.Identity.ServiceAccount.ClientId = value
	return b
}

func (b *serviceAcccountXRHID) WithUsername(value string) ServiceAcccountXRHID {
	b.Identity.ServiceAccount.Username = value
	return b

}
