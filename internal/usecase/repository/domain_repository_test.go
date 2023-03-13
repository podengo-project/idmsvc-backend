package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/repository"
	"github.com/hmsidm/internal/test"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository.DomainRepository
}

func (s *Suite) SetupSuite() {
	var err error
	s.mock, s.DB, err = test.NewSqlMock(&gorm.Session{SkipHooks: true})
	if err != nil {
		s.Suite.FailNow("Error calling gorm.Open: %s", err.Error())
		return
	}
	s.repository = NewDomainRepository()
}

func (s *Suite) TearDownSuite() {
}

func (s *Suite) TestCreate() {
	orgId := "12345"
	token := uuid.NewString()
	tokenExpiration := &time.Time{}
	*tokenExpiration = time.Now().Add(model.DefaultTokenExpiration())
	testUuid := uuid.New()
	t := s.Suite.T()
	currentTime := time.Now()
	var data model.Domain = model.Domain{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		OrgId:                 "12345",
		DomainUuid:            testUuid,
		DomainName:            pointy.String("domain.example"),
		DomainDescription:     pointy.String("My domain example test."),
		DomainType:            pointy.Uint(model.DomainTypeIpa),
		AutoEnrollmentEnabled: pointy.Bool(true),
		IpaDomain: &model.Ipa{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			RealmName: pointy.String("DOMAIN.EXAMPLE"),
			CaCerts: []model.IpaCert{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					IpaID:          1,
					Nickname:       "MYDOMAIN.EXAMPLE IPA CA",
					Issuer:         "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
					Subject:        "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
					SerialNumber:   "1",
					NotValidBefore: currentTime,
					NotValidAfter:  currentTime,
					Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
				},
			},
			Servers: []model.IpaServer{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					IpaID:               1,
					FQDN:                "server1.mydomain.example",
					RHSMId:              "87353f5c-c05c-11ed-9a9b-482ae3863d30",
					HCCEnrollmentServer: true,
					PKInitServer:        true,
					CaServer:            true,
				},
			},
			RealmDomains:    pq.StringArray{"domain.example"},
			Token:           pointy.String(token),
			TokenExpiration: tokenExpiration,
		},
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "domains" ("created_at","updated_at","deleted_at","org_id","domain_uuid","domain_name","domain_description","domain_type","auto_enrollment_enabled","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
		WithArgs(
			data.CreatedAt,
			data.UpdatedAt,
			nil,

			orgId,
			data.DomainUuid,
			data.DomainName,
			data.DomainDescription,
			data.DomainType,
			data.AutoEnrollmentEnabled,
			data.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(data.ID))

	// https://github.com/DATA-DOG/go-sqlmock#matching-arguments-like-timetime
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_domains","token","token_expiration","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
		WithArgs(
			data.IpaDomain.Model.CreatedAt,
			data.IpaDomain.Model.UpdatedAt,
			nil,

			data.IpaDomain.RealmName,
			data.IpaDomain.RealmDomains,
			data.IpaDomain.Token,
			sqlmock.AnyArg(),
			data.IpaDomain.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(data.IpaDomain.ID))

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_certs" ("created_at","updated_at","deleted_at","ipa_id","issuer","nickname","not_valid_after","not_valid_before","pem","serial_number","subject","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).
		WithArgs(
			data.IpaDomain.CaCerts[0].CreatedAt,
			data.IpaDomain.CaCerts[0].UpdatedAt,
			nil,

			data.IpaDomain.CaCerts[0].IpaID,
			data.IpaDomain.CaCerts[0].Issuer,
			data.IpaDomain.CaCerts[0].Nickname,
			data.IpaDomain.CaCerts[0].NotValidAfter,
			data.IpaDomain.CaCerts[0].NotValidBefore,
			data.IpaDomain.CaCerts[0].Pem,
			data.IpaDomain.CaCerts[0].SerialNumber,
			data.IpaDomain.CaCerts[0].Subject,
			data.IpaDomain.CaCerts[0].ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(data.IpaDomain.CaCerts[0].ID))

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_servers" ("created_at","updated_at","deleted_at","ipa_id","fqdn","rhsm_id","ca_server","hcc_enrollment_server","hcc_update_server","pk_init_server","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(
			data.IpaDomain.Servers[0].CreatedAt,
			data.IpaDomain.Servers[0].UpdatedAt,
			nil,

			data.IpaDomain.Servers[0].IpaID,
			data.IpaDomain.Servers[0].FQDN,
			data.IpaDomain.Servers[0].RHSMId,
			data.IpaDomain.Servers[0].CaServer,
			data.IpaDomain.Servers[0].HCCEnrollmentServer,
			data.IpaDomain.Servers[0].HCCUpdateServer,
			data.IpaDomain.Servers[0].PKInitServer,
			data.IpaDomain.Servers[0].ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(data.IpaDomain.Servers[0].ID))

	err := s.repository.Create(s.DB, orgId, &data)
	require.NoError(t, err)
}

func (s *Suite) TestCreateErrors() {
	orgId := "12345"
	testUuid := uuid.New()
	t := s.Suite.T()
	currentTime := time.Now()
	var (
		data model.Domain = model.Domain{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			OrgId:                 "12345",
			DomainUuid:            testUuid,
			DomainName:            pointy.String("domain.example"),
			DomainDescription:     pointy.String("My domain test description"),
			DomainType:            pointy.Uint(model.DomainTypeIpa),
			AutoEnrollmentEnabled: pointy.Bool(true),
			IpaDomain: &model.Ipa{
				CaCerts: []model.IpaCert{},
				Servers: []model.IpaServer{},
			},
		}
		ipaDomainIsNil model.Domain = model.Domain{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			OrgId:                 "12345",
			DomainUuid:            testUuid,
			DomainName:            pointy.String("domain.example"),
			DomainDescription:     pointy.String("My domain test description"),
			DomainType:            pointy.Uint(model.DomainTypeIpa),
			AutoEnrollmentEnabled: pointy.Bool(true),
			IpaDomain:             nil,
		}
		err error
	)

	err = s.repository.Create(nil, "", nil)
	assert.Error(t, err)

	err = s.repository.Create(s.DB, "", nil)
	assert.Error(t, err)

	err = s.repository.Create(s.DB, orgId, nil)
	assert.Error(t, err)

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "domains" ("created_at","updated_at","deleted_at","org_id","domain_uuid","domain_name","domain_description","domain_type","auto_enrollment_enabled","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
		WithArgs(
			data.CreatedAt,
			data.UpdatedAt,
			nil,
			orgId,
			data.DomainUuid,
			data.DomainName,
			data.DomainDescription,
			data.DomainType,
			data.AutoEnrollmentEnabled,
			data.ID,
		).
		WillReturnError(fmt.Errorf("an error happened"))
	err = s.repository.Create(s.DB, orgId, &data)
	assert.Error(t, err)
	assert.Equal(t, "an error happened", err.Error())

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "domains" ("created_at","updated_at","deleted_at","org_id","domain_uuid","domain_name","domain_description","domain_type","auto_enrollment_enabled","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
		WithArgs(
			data.CreatedAt,
			data.UpdatedAt,
			nil,
			orgId,
			data.DomainUuid,
			data.DomainName,
			data.DomainDescription,
			data.DomainType,
			data.AutoEnrollmentEnabled,
			data.ID,
		).
		WillReturnError(fmt.Errorf("an error happened"))
	err = s.repository.Create(s.DB, orgId, &ipaDomainIsNil)
	assert.Error(t, err)
	assert.Equal(t, "data.IpaDomain is nil", err.Error())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
