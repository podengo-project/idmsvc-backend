package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

type SystemXRHID interface {
	Build() identity.XRHID
	WithAccountNumber(value string) SystemXRHID
	WithEmployeeAccountNumber(value string) SystemXRHID
	WithOrgID(value string) SystemXRHID
	WithAuthType(value string) SystemXRHID

	WithCommonName(value string) SystemXRHID
	WithCertType(value string) SystemXRHID
}

type systemXRHID identity.XRHID

// NewSystemXRHID create a new builder to get a XRHID for a user
func NewSystemXRHID() SystemXRHID {
	var orgID = strconv.Itoa(int(builder_helper.GenRandNum(1, 100000)))
	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/schema.json
	return &systemXRHID{
		Identity: identity.Identity{
			AccountNumber:         strconv.Itoa(int(builder_helper.GenRandNum(1, 100000))),
			EmployeeAccountNumber: strconv.Itoa(int(builder_helper.GenRandNum(1, 100000))),
			OrgID:                 orgID,
			Type:                  "System",
			AuthType:              "cert-auth",
			Internal: identity.Internal{
				OrgID:    orgID,
				AuthTime: float32(time.Now().Sub(time.Time{})),
			},
			System: &identity.System{
				CommonName: uuid.NewString(),
				CertType:   "system",
			},
		},
	}
}

func (b *systemXRHID) Build() identity.XRHID {
	return identity.XRHID(*b)
}

func (b *systemXRHID) WithOrgID(value string) SystemXRHID {
	b.Identity.OrgID = value
	b.Identity.Internal.OrgID = value
	return b
}

func (b *systemXRHID) WithAuthType(value string) SystemXRHID {
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

func (b *systemXRHID) WithAccountNumber(value string) SystemXRHID {
	b.Identity.AccountNumber = value
	return b
}

func (b *systemXRHID) WithEmployeeAccountNumber(value string) SystemXRHID {
	b.Identity.EmployeeAccountNumber = value
	return b
}

// --- Start specific System data

func (b *systemXRHID) WithCommonName(value string) SystemXRHID {
	b.Identity.System.CommonName = value
	return b
}

func (b *systemXRHID) WithCertType(value string) SystemXRHID {
	b.Identity.System.CertType = value
	return b
}
