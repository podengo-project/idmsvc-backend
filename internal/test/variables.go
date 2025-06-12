/* Common test constants and variables
 */

package test

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

const (
	DomainId        = "7b160558-8273-5a24-b559-6de3ff053c63"
	OrgId           = "123456"
	DomainName      = "ipa.test"
	RealmName       = "IPA.TEST"
	UserName        = "testuser"
	UserId          = "234"
	UserAccountNr   = "345"
	SystemAccountNr = "456"
)

var (
	DomainType   = public.RhelIdm
	DomainUUID   = uuid.MustParse(DomainId)
	RealmDomains = []string{DomainName, "otherdomain.test"}
	Location1    = public.Location{
		Name:        "sigma",
		Description: pointy.String("Location Sigma"),
	}
	Location2 = public.Location{
		Name: "tau",
		// no description
	}
	IpaCaPublicCert = public.Certificate{
		Nickname:     "IPA.TEST IPA CA",
		Issuer:       "CN=Certificate Authority,O=IPA.TEST",
		Subject:      "CN=Certificate Authority,O=IPA.TEST",
		SerialNumber: "1",
		Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
		NotBefore:    time.Date(2023, 3, 21, 5, 38, 9, 0, time.UTC),
		NotAfter:     time.Date(2043, 3, 21, 5, 38, 9, 0, time.UTC),
	}
	IpaCaModelCert = model.IpaCert{
		Nickname:     IpaCaPublicCert.Nickname,
		Issuer:       IpaCaPublicCert.Issuer,
		Subject:      IpaCaPublicCert.Subject,
		SerialNumber: IpaCaPublicCert.SerialNumber,
		Pem:          IpaCaPublicCert.Pem,
		NotBefore:    IpaCaPublicCert.NotBefore,
		NotAfter:     IpaCaPublicCert.NotAfter,
	}
)

type TestHost struct {
	Fqdn          string
	CertCN        string // subscription manager cert common name
	CertUUID      uuid.UUID
	InventoryId   string
	InventoryUUID uuid.UUID
}

// Create a new test host
// certCN and inventoryId may be empty strings to generate a random UUID
func NewTestHost(fqdn string, certCN string, inventoryId string) TestHost {
	if certCN == "" {
		certCN = uuid.NewString()
	}
	if inventoryId == "" {
		inventoryId = uuid.NewString()
	}
	return TestHost{
		Fqdn:          fqdn,
		CertCN:        certCN,
		CertUUID:      uuid.MustParse(certCN),
		InventoryId:   inventoryId,
		InventoryUUID: uuid.MustParse(inventoryId),
	}
}

var (
	Server1 = NewTestHost(
		"server1.ipa.test",
		"21258fc8-c755-11ed-afc4-482ae3863d30",
		"547ce70c-9eb5-4783-a619-086aa26f88e5",
	)
	Server2 = NewTestHost("server2.ipa.test",
		"5b3ce177-7c02-4ccb-a3d9-037504ded64a",
		"24c82b63-4d8a-4565-b232-0b93913f0c62",
	)
	Client1 = NewTestHost("client1.ipa.test", "", "")
)

// Create XRHID identity for user
func GetUserXRHID(orgId string, userName string, userId string, acountNumber string, admin bool) identity.XRHID {
	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/identities/user.json
	return identity.XRHID{
		Identity: identity.Identity{
			OrgID:         orgId,
			AccountNumber: acountNumber,
			AuthType:      "jwt-auth",
			Type:          "User",
			Internal: identity.Internal{
				OrgID: orgId,
			},
			User: &identity.User{
				Active:    true,
				Internal:  false,
				OrgAdmin:  admin,
				UserID:    userId,
				Username:  userName,
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@ipa.test",
				Locale:    "en_US",
			},
		},
	}
}

// Create XRHID identity for system (cert auth)
func GetSystemXRHID(orgId string, commonName string, acountNumber string) identity.XRHID {
	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/identities/cert.json
	return identity.XRHID{
		Identity: identity.Identity{
			OrgID:         orgId,
			AccountNumber: acountNumber,
			AuthType:      "cert-auth",
			Type:          "System",
			Internal: identity.Internal{
				OrgID: orgId,
			},
			System: &identity.System{
				CommonName: commonName,
				CertType:   "system",
			},
		},
	}

}

func BuildDomainModel(orgID string, domainID uint) *model.Domain {
	currentTime := time.Now()

	return &model.Domain{
		Model: gorm.Model{
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: gorm.DeletedAt{},
		},
		OrgId:                 orgID,
		DomainUuid:            DomainUUID,
		DomainName:            pointy.String(DomainName),
		Title:                 pointy.String("My Domain Example"),
		Description:           pointy.String("Description of My Domain Example"),
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain: &model.Ipa{
			Model: gorm.Model{
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
				DeletedAt: gorm.DeletedAt{},
			},
			RealmName: pointy.String(RealmName),
			CaCerts: []model.IpaCert{
				{
					Model: gorm.Model{
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
						DeletedAt: gorm.DeletedAt{},
					},
					IpaID:        1,
					Issuer:       IpaCaPublicCert.Issuer,
					Nickname:     IpaCaPublicCert.Nickname,
					NotAfter:     currentTime.Add(24 * time.Hour),
					NotBefore:    currentTime,
					SerialNumber: IpaCaPublicCert.SerialNumber,
					Subject:      IpaCaPublicCert.Subject,
					Pem:          IpaCaPublicCert.Pem,
				},
			},
			Servers: []model.IpaServer{
				{
					Model: gorm.Model{
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
						DeletedAt: gorm.DeletedAt{},
					},
					IpaID:               1,
					FQDN:                Server1.Fqdn,
					RHSMId:              pointy.String(Server1.CertCN),
					Location:            pointy.String(Location1.Name),
					CaServer:            true,
					HCCEnrollmentServer: true,
					HCCUpdateServer:     true,
					PKInitServer:        true,
				},
			},
			Locations: []model.IpaLocation{
				{
					Model: gorm.Model{
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
						DeletedAt: gorm.DeletedAt{},
					},
					IpaID:       1,
					Name:        Location1.Name,
					Description: pointy.String(*Location1.Description),
				},
			},
			RealmDomains: pq.StringArray{DomainName},
		},
	}
}

var (
	SystemXRHID  = GetSystemXRHID(OrgId, Server1.CertCN, SystemAccountNr)
	Client1XRHID = GetSystemXRHID(OrgId, Client1.CertCN, SystemAccountNr)
	UserXRHID    = GetUserXRHID(OrgId, UserName, UserId, UserAccountNr, false)
)
