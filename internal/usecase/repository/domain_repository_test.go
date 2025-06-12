package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/domain_token"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	builder_model "github.com/podengo-project/idmsvc-backend/internal/test/builder/model"
	test_sql "github.com/podengo-project/idmsvc-backend/internal/test/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

type DomainRepositorySuite struct {
	SuiteBase
	repository *domainRepository
}

// https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
func (s *DomainRepositorySuite) SetupTest() {
	s.SuiteBase.SetupTest()
	s.repository = &domainRepository{}
}

func (s *DomainRepositorySuite) TestNewDomainRepository() {
	t := s.Suite.T()
	assert.NotPanics(t, func() {
		_ = NewDomainRepository()
	})
}

func (s *DomainRepositorySuite) TestCreateIpaDomain() {
	t := s.Suite.T()
	domainID := uint(1)
	data := test.BuildDomainModel(test.OrgId, domainID)
	var (
		err         error
		expectedErr error
	)

	// Check nil
	expectedErr = fmt.Errorf("code=500, message='data' cannot be nil")
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, nil)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Error on INSERT INTO "ipas"
	expectedErr = fmt.Errorf(`error at INSERT INTO "ipas"`)
	test_sql.CreateIpaDomain(1, s.mock, expectedErr, domainID, data.IpaDomain)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Error on INSERT INTO "ipa_certs"
	expectedErr = fmt.Errorf(`error at INSERT INTO "ipa_certs"`)
	test_sql.CreateIpaDomain(2, s.mock, expectedErr, domainID, data.IpaDomain)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Error on INSERT INTO "ipa_servers"
	expectedErr = fmt.Errorf(`error at INSERT INTO "ipa_servers"`)
	test_sql.CreateIpaDomain(3, s.mock, expectedErr, domainID, data.IpaDomain)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Error on INSERT INTO "ipa_locations"
	expectedErr = fmt.Errorf(`error at INSERT INTO "ipa_locations"`)
	test_sql.CreateIpaDomain(4, s.mock, expectedErr, domainID, data.IpaDomain)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Success scenario
	expectedErr = nil
	test_sql.CreateIpaDomain(4, s.mock, nil, domainID, data.IpaDomain)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data.IpaDomain)
	assert.NoError(t, err)
	require.NoError(t, s.mock.ExpectationsWereMet())
}

func (s *DomainRepositorySuite) TestUpdateAgent() {
	t := s.Suite.T()
	orgID := test.OrgId
	domainID := uint(1)
	var (
		data *model.Domain = test.BuildDomainModel(test.OrgId, domainID)
		err  error
	)

	assert.PanicsWithValue(t, "'ctx' is nil", func() {
		_ = s.repository.UpdateAgent(nil, "", nil)
	})
	require.NoError(t, s.mock.ExpectationsWereMet())

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_ = s.repository.UpdateAgent(context.Background(), "", nil)
	})
	require.NoError(t, s.mock.ExpectationsWereMet())

	err = s.repository.UpdateAgent(s.Ctx, "", nil)
	assert.EqualError(t, err, "'orgID' is empty")
	require.NoError(t, s.mock.ExpectationsWereMet())

	err = s.repository.UpdateAgent(s.Ctx, orgID, nil)
	assert.EqualError(t, err, "code=500, message='data' cannot be nil")
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr := fmt.Errorf("Domain.Model.ID cannot be 0")
	s.mock.MatchExpectationsInOrder(true)
	data.Model.ID = 0
	data.IpaDomain.Model.ID = 0
	test_sql.UpdateAgent(0, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = fmt.Errorf("error at record not found")
	s.mock.MatchExpectationsInOrder(true)
	data.Model.ID = domainID
	data.IpaDomain.Model.ID = domainID
	test_sql.UpdateAgent(1, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = fmt.Errorf("error at UPDATE INTO domains")
	s.mock.MatchExpectationsInOrder(true)
	test_sql.UpdateAgent(2, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = fmt.Errorf("error at DELETE FROM ipas")
	s.mock.MatchExpectationsInOrder(true)
	test_sql.UpdateAgent(3, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = fmt.Errorf("error at INSERT INTO ipas")
	s.mock.MatchExpectationsInOrder(true)
	test_sql.UpdateAgent(4, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Success scenario
	expectedErr = nil
	s.mock.MatchExpectationsInOrder(true)
	test_sql.UpdateAgent(4, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.NoError(t, err)
	require.NoError(t, s.mock.ExpectationsWereMet())
}

func (s *DomainRepositorySuite) TestUpdateIpaDomain() {
	var (
		err         error
		expectedErr error
	)
	t := s.Suite.T()
	orgID := "11111"
	domainID := uint(1)
	data := test.BuildDomainModel(orgID, domainID)

	// Wrong arguments: log is nil
	expectedErr = internal_errors.NilArgError("log")
	test_sql.UpdateIpaDomain(0, s.mock, expectedErr, domainID, data)
	err = s.repository.updateIpaDomain(nil, nil, nil)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Wrong arguments: db is nil
	expectedErr = internal_errors.NilArgError("db")
	test_sql.UpdateIpaDomain(0, s.mock, expectedErr, domainID, data)
	err = s.repository.updateIpaDomain(s.Log, nil, nil)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Wrong arguments: dataIPA is nil
	expectedErr = internal_errors.NilArgError("dataIPA")
	test_sql.UpdateIpaDomain(0, s.mock, expectedErr, domainID, data)
	err = s.repository.updateIpaDomain(s.Log, s.DB, nil)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Wrong arguments: Ipa.ID is 0 when trying to update ipa information
	expectedErr = fmt.Errorf("dataIPA.Model.ID cannot be 0")
	test_sql.UpdateIpaDomain(0, s.mock, expectedErr, domainID, data)
	data.IpaDomain.Model.ID = 0
	err = s.repository.updateIpaDomain(s.Log, s.DB, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// database error at DELETE FROM 'ipas'
	expectedErr = fmt.Errorf("database error at DELETE FROM 'ipas'")
	test_sql.UpdateIpaDomain(1, s.mock, expectedErr, domainID, data)
	data.IpaDomain.Model.ID = domainID
	err = s.repository.updateIpaDomain(s.Log, s.DB, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// database error at INSERT INTO 'ipas'
	expectedErr = fmt.Errorf("database error at INSERT INTO 'ipas'")
	test_sql.UpdateIpaDomain(2, s.mock, expectedErr, domainID, data)
	data.IpaDomain.Model.ID = domainID
	err = s.repository.updateIpaDomain(s.Log, s.DB, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// database error at INSERT INTO 'ipa_certs'
	expectedErr = fmt.Errorf("database error at INSERT INTO 'ipa_certs'")
	test_sql.UpdateIpaDomain(3, s.mock, expectedErr, domainID, data)
	data.IpaDomain.Model.ID = domainID
	err = s.repository.updateIpaDomain(s.Log, s.DB, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// database error at INSERT INTO 'ipa_servers'
	expectedErr = fmt.Errorf("database error at INSERT INTO 'ipa_servers'")
	test_sql.UpdateIpaDomain(4, s.mock, expectedErr, domainID, data)
	data.IpaDomain.Model.ID = domainID
	err = s.repository.updateIpaDomain(s.Log, s.DB, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// database error at INSERT INTO 'ipa_locations'
	expectedErr = fmt.Errorf("database error at INSERT INTO 'ipa_locations'")
	test_sql.UpdateIpaDomain(5, s.mock, expectedErr, domainID, data)
	data.IpaDomain.Model.ID = domainID
	err = s.repository.updateIpaDomain(s.Log, s.DB, data.IpaDomain)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Success scenario
	expectedErr = nil
	test_sql.UpdateIpaDomain(5, s.mock, expectedErr, domainID, data)
	err = s.repository.updateIpaDomain(s.Log, s.DB, data.IpaDomain)
	require.NoError(t, err)
	require.NoError(t, s.mock.ExpectationsWereMet())

}

func (s *DomainRepositorySuite) TestList() {
	t := s.T()
	r := &domainRepository{}
	currentTime := time.Now()
	orgID := "11111"
	domainId := uuid.MustParse("3bccb88e-dd25-11ed-99e0-482ae3863d30")
	subscriptionManagerID := "fe106208-dd32-11ed-aa87-482ae3863d30"
	data := model.Domain{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: gorm.DeletedAt{},
		},
		OrgId:                 orgID,
		DomainUuid:            domainId,
		DomainName:            pointy.String("mydomain.example"),
		Title:                 pointy.String("My Domain Example"),
		Description:           pointy.String("Description of My Domain Example"),
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain: &model.Ipa{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
				DeletedAt: gorm.DeletedAt{},
			},
			RealmName: pointy.String("MYDOMAIN.EXAMPLE"),
			CaCerts: []model.IpaCert{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
						DeletedAt: gorm.DeletedAt{},
					},
					IpaID:        1,
					Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
					Nickname:     "MYDOMAIN.EXAMPLE IPA CA",
					NotAfter:     currentTime.Add(24 * time.Hour),
					NotBefore:    currentTime,
					SerialNumber: "1",
					Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
					Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----",
				},
			},
			Servers: []model.IpaServer{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
						DeletedAt: gorm.DeletedAt{},
					},
					IpaID:               1,
					FQDN:                "server1.mydomain.example",
					RHSMId:              pointy.String(subscriptionManagerID),
					Location:            pointy.String("europe"),
					CaServer:            true,
					HCCEnrollmentServer: true,
					HCCUpdateServer:     true,
					PKInitServer:        true,
				},
			},
			RealmDomains: pq.StringArray{"mydomain.example"},
		},
	}

	// Fail on checks
	ctx := context.TODO()
	ctx = app_context.CtxWithLog(ctx, slog.Default())
	ctx = app_context.CtxWithDB(ctx, s.DB)
	output, count, err := r.List(ctx, "", -1, -1)
	assert.EqualError(t, err, "'orgID' is empty")
	assert.Equal(t, int64(0), count)
	assert.Nil(t, output)

	// Return error
	offset := 0
	limit := 5
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "domains" WHERE org_id = $1`)).
		WithArgs(orgID).
		WillReturnError(fmt.Errorf("an error happened"))
	output, count, err = r.List(ctx, orgID, offset, limit)
	assert.EqualError(t, err, "an error happened")
	assert.Equal(t, int64(0), count)
	assert.Nil(t, output)

	// Success case
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "domains" WHERE org_id = $1`)).
		WithArgs(orgID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE org_id = $1 AND "domains"."deleted_at" IS NULL LIMIT $2`)).
		WithArgs(orgID, 5).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"org_id", "domain_uuid", "domain_name",
			"title", "description", "type",
			"auto_enrollment_enabled",
		}).AddRow(
			data.Model.ID,
			data.Model.CreatedAt,
			data.Model.UpdatedAt,
			data.Model.DeletedAt,

			data.OrgId,
			data.DomainUuid,
			data.DomainName,
			data.Title,
			data.Description,
			data.Type,
			data.AutoEnrollmentEnabled,
		))
	output, count, err = r.List(ctx, orgID, offset, limit)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, []model.Domain{
		{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: data.CreatedAt,
				UpdatedAt: data.UpdatedAt,
			},
			AutoEnrollmentEnabled: pointy.Bool(true),
			OrgId:                 data.OrgId,
			DomainUuid:            data.DomainUuid,
			DomainName:            data.DomainName,
			Title:                 data.Title,
			Description:           data.Description,
			Type:                  pointy.Uint(model.DomainTypeIpa),
		},
	}, output)
}

func (s *DomainRepositorySuite) TestFindByID() {
	t := s.T()
	r := &domainRepository{}
	s.mock.MatchExpectationsInOrder(true)

	currentTime := time.Now()
	notBefore := currentTime
	notAfter := currentTime.Add(356 * 24 * time.Hour)
	domainUUID := uuid.MustParse("c5d2c9d0-2b2f-11ee-8ec5-482ae3863d30")
	domainID := uint(1)

	// TODO Use the builder_model.NewDomain(...)
	data := &model.Domain{
		Model: gorm.Model{
			ID:        domainID,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		OrgId:                 "12345",
		DomainUuid:            domainUUID,
		DomainName:            pointy.String("mydomain.example"),
		Title:                 pointy.String("My Example Domain"),
		Description:           pointy.String("My long description for my example domain"),
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain: &model.Ipa{
			Model: gorm.Model{
				ID:        domainID,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			RealmName:    pointy.String("MYDOMAIN.EXAMPLE"),
			RealmDomains: pq.StringArray{"mydomain.example"},
			CaCerts: []model.IpaCert{
				{
					Model: gorm.Model{
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					IpaID:        domainID,
					Issuer:       "issuer",
					Subject:      "subject",
					Nickname:     "nickname",
					NotBefore:    notBefore,
					NotAfter:     notAfter,
					SerialNumber: "18384021308",
					Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
				},
			},
			Servers: []model.IpaServer{
				{
					Model: gorm.Model{
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					IpaID:               domainID,
					FQDN:                "server1.mydomain.example",
					RHSMId:              pointy.String("17a151b6-2b2e-11ee-97d9-482ae3863d30"),
					Location:            pointy.String("boston"),
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
					},
					IpaID:       domainID,
					Name:        "boston",
					Description: pointy.String("Boston data center"),
				},
			},
		},
	}
	dataTypeNil := &model.Domain{
		Model:                 data.Model,
		OrgId:                 data.OrgId,
		DomainUuid:            data.DomainUuid,
		DomainName:            data.DomainName,
		Title:                 data.Title,
		Description:           data.Description,
		Type:                  nil,
		AutoEnrollmentEnabled: data.AutoEnrollmentEnabled,
		IpaDomain:             data.IpaDomain,
	}
	var expectedErr error
	var err error
	var domain *model.Domain

	s.mock.MatchExpectationsInOrder(true)

	// Check one wrong argument
	assert.Panics(t, func() {
		_, _ = r.FindByID(nil, "", uuid.Nil)
	})

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_, _ = r.FindByID(context.Background(), "", uuid.Nil)
	})

	// Check path when an error hapens into the sql statement
	expectedErr = fmt.Errorf(`error at SELECT * FROM "domains"`)
	test_sql.FindByID(1, s.mock, expectedErr, domainID, data)
	domain, err = r.FindByID(s.Ctx, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, domain)

	// Check path when a domain type is NULL
	expectedErr = internal_errors.NilArgError("Type")
	test_sql.FindByID(1, s.mock, nil, domainID, dataTypeNil)
	domain, err = r.FindByID(s.Ctx, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, domain)

	// Check for 'ipas' record not found
	expectedErr = gorm.ErrRecordNotFound
	test_sql.FindByID(1, s.mock, nil, domainID, data)
	test_sql.FindIpaByID(1, s.mock, expectedErr, domainID, data)
	domain, err = r.FindByID(s.Ctx, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, domain)

	// Successful scenario
	expectedErr = nil
	test_sql.FindByID(1, s.mock, nil, domainID, data)
	test_sql.FindIpaByID(4, s.mock, expectedErr, domainID, data)
	domain, err = r.FindByID(s.Ctx, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.NoError(t, err)
	require.NotNil(t, domain)

	assert.Equal(t, data.ID, domain.ID)
	assert.Equal(t, data.CreatedAt, domain.CreatedAt)
	assert.Equal(t, data.DeletedAt, domain.DeletedAt)
	assert.Equal(t, data.Title, domain.Title)
	assert.Equal(t, data.OrgId, domain.OrgId)
	assert.Equal(t, data.AutoEnrollmentEnabled, domain.AutoEnrollmentEnabled)
	assert.Equal(t, data.Description, domain.Description)
	assert.Equal(t, data.DomainName, domain.DomainName)
	assert.Equal(t, data.DomainUuid, domain.DomainUuid)
	assert.Equal(t, data.Type, domain.Type)
}

func (s *DomainRepositorySuite) TestUpdateUser() {
	var (
		err         error
		expectedErr error
	)
	t := s.Suite.T()
	orgID := test.OrgId
	domainID := uint(1)
	data := test.BuildDomainModel(orgID, domainID)
	data.Model.ID = domainID
	data.IpaDomain.Model.ID = domainID
	c := app_context.CtxWithLog(app_context.CtxWithDB(context.Background(), s.DB), slog.Default())

	// ctx is nil
	expectedErr = fmt.Errorf("code=500, message='ctx' cannot be nil")
	err = s.repository.UpdateUser(nil, "", nil)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// orgID is empty
	expectedErr = fmt.Errorf("'orgID' is empty")
	err = s.repository.UpdateUser(c, "", nil)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// data is nil
	expectedErr = fmt.Errorf("code=500, message='data' cannot be nil")
	err = s.repository.UpdateUser(c, orgID, nil)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// error at FindByID
	expectedErr = fmt.Errorf("error at FindByID")
	data.Model.ID = domainID
	test_sql.UpdateUser(1, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateUser(c, orgID, data)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// error at UPDATE INTO 'domains'
	expectedErr = fmt.Errorf("error at UPDATE INTO 'domains'")
	test_sql.UpdateUser(2, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateUser(c, orgID, data)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// successful scenario
	expectedErr = nil
	test_sql.UpdateUser(2, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateUser(c, orgID, data)
	require.NoError(t, err)
	require.NoError(t, s.mock.ExpectationsWereMet())
}

// ---------------- Test for private methods ---------------------

func (s *DomainRepositorySuite) TestCheckCommon() {
	t := s.T()
	r := &domainRepository{}

	err := r.checkCommon(nil, "")
	assert.EqualError(t, err, "code=500, message='db' cannot be nil")

	err = r.checkCommon(s.DB, "")
	assert.EqualError(t, err, "'orgID' is empty")

	err = r.checkCommon(s.DB, "12345")
	assert.NoError(t, err)
}

func (s *DomainRepositorySuite) TestCheckCommonAndUUID() {
	t := s.T()
	r := &domainRepository{}

	err := r.checkCommonAndUUID(nil, "", uuid.Nil)
	assert.EqualError(t, err, "code=500, message='db' cannot be nil")

	err = r.checkCommonAndUUID(s.DB, "", uuid.Nil)
	assert.EqualError(t, err, "'orgID' is empty")

	err = r.checkCommonAndUUID(s.DB, "12345", uuid.Nil)
	assert.EqualError(t, err, "'uuid' is invalid")

	err = r.checkCommonAndUUID(s.DB, "12345", uuid.MustParse("42f7adee-e932-11ed-8d73-482ae3863d30"))
	assert.NoError(t, err)
}

func (s *DomainRepositorySuite) TestCheckCommonAndData() {
	t := s.T()
	r := &domainRepository{}

	err := r.checkCommonAndData(nil, "", nil)
	assert.EqualError(t, err, "code=500, message='db' cannot be nil")

	err = r.checkCommonAndData(s.DB, "", nil)
	assert.EqualError(t, err, "'orgID' is empty")

	err = r.checkCommonAndData(s.DB, "12345", nil)
	assert.EqualError(t, err, "code=500, message='data' cannot be nil")

	err = r.checkCommonAndData(s.DB, "12345", &model.Domain{})
	assert.NoError(t, err)
}

func (s *DomainRepositorySuite) TestCheckCommonAndDataAndType() {
	t := s.T()
	r := &domainRepository{}

	err := r.checkCommonAndDataAndType(nil, "", nil)
	assert.EqualError(t, err, "code=500, message='db' cannot be nil")

	err = r.checkCommonAndDataAndType(s.DB, "12345", &model.Domain{})
	assert.EqualError(t, err, "code=500, message='Type' cannot be nil")

	err = r.checkCommonAndDataAndType(s.DB, "12345", &model.Domain{
		Type: pointy.Uint(model.DomainTypeIpa),
	})
	assert.NoError(t, err)
}

func (s *DomainRepositorySuite) TestCreateDomainToken() {
	var (
		key       []byte        = []byte("secret")
		validity  time.Duration = 2 * time.Hour
		testOrgID string        = "12345"
	)
	t := s.T()
	r := &domainRepository{}
	drt, err := r.CreateDomainToken(s.Ctx, key, validity, testOrgID, public.RhelIdm)
	assert.NoError(t, err)
	assert.Equal(t, drt.DomainType, public.RhelIdm)
	assert.NotEmpty(t, drt.DomainId)
	assert.Equal(
		t,
		drt.DomainId,
		domain_token.TokenDomainId(domain_token.DomainRegistrationToken(drt.DomainToken)),
	)
	assert.Greater(t, drt.ExpirationNS, uint64(time.Now().UnixNano()))
}

func (s *DomainRepositorySuite) TestPrepareUpdateUser() {
	var (
		value *string
		flag  *bool
		ok    bool
		v     interface{}
	)
	t := s.T()
	r := &domainRepository{}
	assert.Panics(t, func() {
		r.prepareUpdateUser(nil)
	})

	// title
	fields := r.prepareUpdateUser(&model.Domain{
		Title: pointy.String("My Domain Title"),
	})
	assert.Equal(t, 1, len(fields))
	v, ok = fields["title"]
	require.True(t, ok)
	require.NotNil(t, v)
	value, ok = v.(*string)
	require.True(t, ok)
	assert.Equal(t, "My Domain Title", *value)

	// description
	fields = r.prepareUpdateUser(&model.Domain{
		Description: pointy.String("My Domain Description"),
	})
	assert.Equal(t, 1, len(fields))
	v, ok = fields["description"]
	require.True(t, ok)
	require.NotNil(t, v)
	value, ok = v.(*string)
	require.True(t, ok)
	assert.Equal(t, "My Domain Description", *value)

	// auto_enrollment_enabled
	fields = r.prepareUpdateUser(&model.Domain{
		AutoEnrollmentEnabled: pointy.Bool(true),
	})
	assert.Equal(t, 1, len(fields))
	v, ok = fields["auto_enrollment_enabled"]
	require.True(t, ok)
	require.NotNil(t, v)
	flag, ok = v.(*bool)
	require.True(t, ok)
	assert.Equal(t, true, *flag)
}

func (s *DomainRepositorySuite) TestDeleteById() {
	var expectedErr error
	t := s.T()
	r := &domainRepository{}

	d := builder_model.NewDomain(builder_model.NewModel().WithID(1).Build()).Build()

	require.Panics(t, func() {
		_ = r.DeleteById(nil, "", model.NilUUID)
	})
	require.NoError(t, s.mock.ExpectationsWereMet())

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_ = r.DeleteById(context.Background(), "", model.NilUUID)
	})
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = fmt.Errorf("code=404, message=unknown domain '%s'", d.DomainUuid.String())
	test_sql.DeleteByID(1, s.mock, gorm.ErrRecordNotFound, d)
	err := r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = fmt.Errorf("invalid transaction")
	test_sql.DeleteByID(1, s.mock, gorm.ErrInvalidTransaction, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = fmt.Errorf("code=404, message=unknown domain '%s'", d.DomainUuid.String())
	test_sql.DeleteByID(2, s.mock, gorm.ErrRecordNotFound, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = fmt.Errorf("code=404, message=unknown domain '%s'", d.DomainUuid.String())
	test_sql.DeleteByID(3, s.mock, gorm.ErrRecordNotFound, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = gorm.ErrInvalidTransaction
	test_sql.DeleteByID(3, s.mock, gorm.ErrInvalidTransaction, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Success scenario
	expectedErr = nil
	test_sql.DeleteByID(3, s.mock, expectedErr, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.NoError(t, err)
	require.NoError(t, s.mock.ExpectationsWereMet())
}

func (s *DomainRepositorySuite) TestWrapErrNotFound() {
	var err error
	t := s.T()
	r := &domainRepository{}
	UUID := uuid.New()

	err = r.wrapErrNotFound(nil, UUID)
	require.NoError(t, err)

	err = r.wrapErrNotFound(gorm.ErrRecordNotFound, UUID)
	require.EqualError(t, err, fmt.Sprintf("code=404, message=unknown domain '%s'", UUID.String()))
}

func (s *DomainRepositorySuite) TestRegister() {
	var (
		err         error
		expectedErr error
	)
	t := s.T()
	r := &domainRepository{}
	realm := strings.ToUpper(helper.GenRandDomainName(3))
	domainID := uint(helper.GenRandNum(1, 2^63))
	gormModel := builder_model.NewModel().WithID(domainID).Build()
	d := builder_model.NewDomain(gormModel).Build()

	assert.Panics(t, func() {
		_ = r.Register(nil, d.OrgId, d)
	})
	require.NoError(t, s.mock.ExpectationsWereMet())

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_ = r.Register(context.Background(), d.OrgId, d)
	})
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = gorm.ErrDuplicatedKey
	test_sql.Register(1, s.mock, expectedErr, d)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, fmt.Sprintf("code=409, message=domain id '%s' is already registered.", d.DomainUuid))
	require.NoError(t, s.mock.ExpectationsWereMet())

	expectedErr = gorm.ErrInvalidField
	test_sql.Register(1, s.mock, expectedErr, d)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, "invalid field")
	require.NoError(t, s.mock.ExpectationsWereMet())

	d = builder_model.NewDomain(gormModel).
		WithIpaDomain(
			builder_model.NewIpaDomain().
				WithModel(gormModel).
				WithRealmName(&realm).
				WithRealmDomains(pq.StringArray{strings.ToLower(realm)}).
				WithServers([]model.IpaServer{
					builder_model.NewIpaServer(
						builder_model.NewModel().Build(),
					).WithIpaID(domainID).Build(),
				}).
				WithLocations([]model.IpaLocation{
					builder_model.NewIpaLocation(
						builder_model.NewModel().Build(),
					).WithIpaID(domainID).Build(),
				}).
				WithCaCerts([]model.IpaCert{
					builder_model.NewIpaCert(
						builder_model.NewModel().Build(),
						realm,
					).WithIpaID(domainID).Build(),
				}).
				Build(),
		).Build()
	test_sql.Register(1, s.mock, nil, d)
	test_sql.CreateIpaDomain(1, s.mock, gorm.ErrInvalidField, domainID, d.IpaDomain)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, "invalid field")
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Success case - FIXME Flaky test
	test_sql.Register(1, s.mock, nil, d)
	test_sql.CreateIpaDomain(4, s.mock, nil, domainID, d.IpaDomain)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.NoError(t, err)
	require.NoError(t, s.mock.ExpectationsWereMet())

	// IpaDomain is nil
	test_sql.Register(1, s.mock, nil, d)
	d.IpaDomain = nil
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, "code=500, message='IpaDomain' cannot be nil")
	require.NoError(t, s.mock.ExpectationsWereMet())
}

func (s *DomainRepositorySuite) TestDeleteByIdLogError() {
	t := s.T()
	r := &domainRepository{}

	domainID := uint(1)
	d := builder_model.NewDomain(builder_model.NewModel().WithID(domainID).Build()).Build()

	test_sql.DeleteByID(1, s.mock, gorm.ErrInvalidTransaction, d)
	err := r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, "invalid transaction")
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Check the log message
	assert.Contains(t, s.LogBuffer.String(), `level=ERROR msg="deleting domain when checking that the record exist"`)
}

func (s *DomainRepositorySuite) TestCheckList() {
	t := s.T()
	r := &domainRepository{}

	var (
		log *slog.Logger
		db  *gorm.DB
		err error
		ctx context.Context
	)

	ctx = context.TODO()
	ctx = app_context.CtxWithLog(ctx, slog.Default())
	ctx = app_context.CtxWithDB(ctx, s.DB)
	assert.NotPanics(t, func() {
		log, db, err = r.checkList(ctx, "", -1, -1)
	})
	require.NotNil(t, log)
	require.NotNil(t, db)
	require.EqualError(t, err, "'orgID' is empty")

	assert.NotPanics(t, func() {
		log, db, err = r.checkList(ctx, "12345", -1, -1)
	})
	require.NotNil(t, log)
	require.NotNil(t, db)
	require.EqualError(t, err, "'offset' is lower than 0")

	ctx = app_context.CtxWithDB(ctx, s.DB)
	assert.NotPanics(t, func() {
		log, db, err = r.checkList(ctx, "12345", 0, -1)
	})
	require.NotNil(t, log)
	require.NotNil(t, db)
	require.EqualError(t, err, "'limit' is lower than 0")

	assert.NotPanics(t, func() {
		log, db, err = r.checkList(ctx, "12345", 0, 10)
	})
	require.NoError(t, err)
}

func TestDomainRepositorySuite(t *testing.T) {
	suite.Run(t, new(DomainRepositorySuite))
}
