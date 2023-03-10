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
			RealmName:  pointy.String("DOMAIN.EXAMPLE"),
			CaCerts:    []model.IpaCert{},
			Servers:    []model.IpaServer{},
			RealmNames: "domain.example",
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

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_names","id") VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT ("id") DO UPDATE SET "id"="excluded"."id" RETURNING "id"`)).
		WithArgs(
			data.IpaDomain.Model.CreatedAt,
			data.IpaDomain.Model.UpdatedAt,
			nil,

			data.IpaDomain.RealmName,
			data.IpaDomain.RealmNames,
			data.IpaDomain.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(data.IpaDomain.ID))

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
