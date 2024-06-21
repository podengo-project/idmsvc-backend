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

type ServiceAccountXRHID interface {
	Build() identity.XRHID
	WithAccountNumber(value string) ServiceAccountXRHID
	WithEmployeeAccountNumber(value string) ServiceAccountXRHID
	WithOrgID(value string) ServiceAccountXRHID
	WithAuthType(value string) ServiceAccountXRHID

	WithClientID(value string) ServiceAccountXRHID
	WithUsername(value string) ServiceAccountXRHID
}

type serviceAccountXRHID identity.XRHID

// NewServiceAccountXRHID create a new builder to get a XRHID for a user
func NewServiceAccountXRHID() ServiceAccountXRHID {
	var orgID = strconv.Itoa(int(builder_helper.GenRandNum(1, 100000)))
	// See: https://github.com/RedHatInsights/identity-schemas/blob/main/3scale/identities/service-account.json
	return &serviceAccountXRHID{
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

func (b *serviceAccountXRHID) Build() identity.XRHID {
	return identity.XRHID(*b)
}

func (b *serviceAccountXRHID) WithOrgID(value string) ServiceAccountXRHID {
	b.Identity.OrgID = value
	b.Identity.Internal.OrgID = value
	return b
}

func (b *serviceAccountXRHID) WithAuthType(value string) ServiceAccountXRHID {
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

func (b *serviceAccountXRHID) WithAccountNumber(value string) ServiceAccountXRHID {
	b.Identity.AccountNumber = value
	return b
}

func (b *serviceAccountXRHID) WithEmployeeAccountNumber(value string) ServiceAccountXRHID {
	b.Identity.EmployeeAccountNumber = value
	return b
}

// --- Start specific ServiceAccount data

func (b *serviceAccountXRHID) WithClientID(value string) ServiceAccountXRHID {
	b.Identity.ServiceAccount.ClientId = value
	return b
}

func (b *serviceAccountXRHID) WithUsername(value string) ServiceAccountXRHID {
	b.Identity.ServiceAccount.Username = value
	return b

}
