package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
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
	// testUuid := uuid.New()
	currentTime := time.Now()
	var (
		err  error
		data model.Ipa = model.Ipa{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			RealmName: pointy.String("MYDOMAIN.EXAMPLE"),
			CaCerts: []model.IpaCert{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					IpaID:        1,
					Nickname:     "MYDOMAIN.EXAMPLE IPA CA",
					Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
					Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
					SerialNumber: "1",
					NotBefore:    currentTime,
					NotAfter:     currentTime,
					Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
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
					RHSMId:              pointy.String("87353f5c-c05c-11ed-9a9b-482ae3863d30"),
					Location:            pointy.String("europe"),
					HCCEnrollmentServer: true,
					PKInitServer:        true,
					CaServer:            true,
				},
			},
			Locations: []model.IpaLocation{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					IpaID:       1,
					Name:        "boston",
					Description: pointy.String("Boston data center"),
				},
			},
			RealmDomains: []string{"mydomain.example"},
		}
		expectedError error
	)

	// Check nil
	err = s.repository.createIpaDomain(s.Log, s.DB, 1, nil)
	assert.EqualError(t, err, "code=500, message='data' cannot be nil")

	// Error on INSERT INTO "ipas"
	expectedError = fmt.Errorf(`error at INSERT INTO "ipas"`)
	test_sql.CreateIpaDomain(1, s.mock, expectedError, &data)
	err = s.repository.createIpaDomain(s.Log, s.DB, 1, &data)
	assert.EqualError(t, err, expectedError.Error())

	// Error on INSERT INTO "ipa_certs"
	expectedError = fmt.Errorf(`INSERT INTO "ipa_certs"`)
	test_sql.CreateIpaDomain(2, s.mock, expectedError, &data)
	err = s.repository.createIpaDomain(s.Log, s.DB, 1, &data)
	assert.EqualError(t, err, expectedError.Error())

	// Error on INSERT INTO "ipa_servers"
	expectedError = fmt.Errorf(`INSERT INTO "ipa_servers"`)
	test_sql.CreateIpaDomain(3, s.mock, expectedError, &data)
	err = s.repository.createIpaDomain(s.Log, s.DB, 1, &data)
	assert.EqualError(t, err, expectedError.Error())

	// Error on INSERT INTO "ipa_locations"
	expectedError = fmt.Errorf(`INSERT INTO "ipa_locations"`)
	test_sql.CreateIpaDomain(4, s.mock, expectedError, &data)
	err = s.repository.createIpaDomain(s.Log, s.DB, 1, &data)
	assert.EqualError(t, err, expectedError.Error())

	// Success scenario
	expectedError = nil
	test_sql.CreateIpaDomain(4, s.mock, nil, &data)
	err = s.repository.createIpaDomain(s.Log, s.DB, 1, &data)
	assert.NoError(t, err)
}

func (s *DomainRepositorySuite) TestUpdateErrors() {
	t := s.Suite.T()
	orgID := test.OrgId
	testUUID := test.DomainUUID
	domainID := uint(1)
	var (
		data *model.Domain = test.BuildDomainModel(test.OrgId)
		err  error
	)

	assert.Panics(t, func() {
		_ = s.repository.UpdateAgent(nil, "", nil)
	})

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_ = s.repository.UpdateAgent(context.Background(), "", nil)
	})

	err = s.repository.UpdateAgent(s.Ctx, "", nil)
	assert.EqualError(t, err, "'orgID' is empty")

	err = s.repository.UpdateAgent(s.Ctx, orgID, nil)
	assert.EqualError(t, err, "code=500, message='data' cannot be nil")

	expectedErr := fmt.Errorf("record not found")
	s.mock.MatchExpectationsInOrder(true)
	test_sql.FindByID(1, s.mock, nil, domainID, data)
	test_sql.FindIpaByID(1, s.mock, expectedErr, domainID, data)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.EqualError(t, err, "record not found")

	s.mock.MatchExpectationsInOrder(true)
	test_sql.FindByID(1, s.mock, nil, domainID, data)
	test_sql.FindIpaByID(4, s.mock, nil, domainID, data)
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "domains" SET "created_at"=$1,"updated_at"=$2,"org_id"=$3,"domain_uuid"=$4,"domain_name"=$5,"title"=$6,"description"=$7,"type"=$8,"auto_enrollment_enabled"=$9 WHERE (org_id = $10 AND domain_uuid = $11) AND "domains"."deleted_at" IS NULL AND "id" = $12`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),

			data.OrgId,
			testUUID.String(),
			data.DomainName,

			data.Title,
			data.Description,
			data.Type,
			data.AutoEnrollmentEnabled,

			data.OrgId,
			data.DomainUuid,
			data.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ipas" WHERE "ipas"."id" = $1`)).
		WithArgs(
			data.ID,
		).WillReturnResult(
		driver.RowsAffected(1),
	)
	test_sql.CreateIpaDomain(4, s.mock, nil, data.IpaDomain)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.NoError(t, err)
}

func (s *DomainRepositorySuite) TestUpdateIpaDomain() {
	var (
		err error
	)
	t := s.Suite.T()
	currentTime := time.Now()
	orgID := "11111"
	domainId := uuid.MustParse("3bccb88e-dd25-11ed-99e0-482ae3863d30")
	subscriptionManagerID := pointy.String("fe106208-dd32-11ed-aa87-482ae3863d30")
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
						ID:        2,
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
						ID:        3,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
						DeletedAt: gorm.DeletedAt{},
					},
					IpaID:               1,
					FQDN:                "server1.mydomain.example",
					RHSMId:              subscriptionManagerID,
					Location:            pointy.String("europe"),
					CaServer:            true,
					HCCEnrollmentServer: true,
					HCCUpdateServer:     true,
					PKInitServer:        true,
				},
			},
			Locations: []model.IpaLocation{
				{
					Model: gorm.Model{
						ID:        4,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
						DeletedAt: gorm.DeletedAt{},
					},
					Name:        "boston",
					Description: pointy.String("Boston data center"),
					IpaID:       1,
				},
			},
			RealmDomains: pq.StringArray{"mydomain.example"},
		},
	}

	type TestCaseGiven struct {
		Stage  int
		DB     *gorm.DB
		Domain *model.Domain
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected error
	}

	testCases := []TestCase{
		{
			Name: "Wrong arguments: db is nil",
			Given: TestCaseGiven{
				Stage:  0,
				DB:     nil,
				Domain: nil,
			},
			Expected: internal_errors.NilArgError("db"),
		},
		{
			Name: "Wrong arguments: data is nil",
			Given: TestCaseGiven{
				Stage:  0,
				DB:     s.DB,
				Domain: nil,
			},
			Expected: internal_errors.NilArgError("data"),
		},
		{
			Name: "database error at DELETE FROM 'ipas'",
			Given: TestCaseGiven{
				Stage:  1,
				DB:     s.DB,
				Domain: &data,
			},
			Expected: fmt.Errorf("database error at DELETE FROM 'ipas'"),
		},
		{
			Name: "database error at INSERT INTO 'ipas'",
			Given: TestCaseGiven{
				Stage:  2,
				DB:     s.DB,
				Domain: &data,
			},
			Expected: fmt.Errorf("database error at INSERT INTO 'ipas'"),
		},
		{
			Name: "database error at INSERT INTO 'ipa_certs'",
			Given: TestCaseGiven{
				Stage:  3,
				DB:     s.DB,
				Domain: &data,
			},
			Expected: fmt.Errorf("database error at INSERT INTO 'ipa_certs'"),
		},
		{
			Name: "database error at INSERT INTO 'ipa_servers'",
			Given: TestCaseGiven{
				Stage:  4,
				DB:     s.DB,
				Domain: &data,
			},
			Expected: fmt.Errorf("database error at INSERT INTO 'ipa_servers'"),
		},
		{
			Name: "database error at INSERT INTO 'ipa_locations'",
			Given: TestCaseGiven{
				Stage:  5,
				DB:     s.DB,
				Domain: &data,
			},
			Expected: fmt.Errorf("database error at INSERT INTO 'ipa_locations'"),
		},
		{
			Name: "Success scenario",
			Given: TestCaseGiven{
				Stage:  5,
				DB:     s.DB,
				Domain: &data,
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)

		// Prepare the db mock
		test_sql.UpdateIpaDomain(testCase.Given.Stage, s.mock, testCase.Expected, &data)

		// Run for error or success
		if testCase.Given.Domain != nil {
			err = s.repository.updateIpaDomain(s.Log, testCase.Given.DB, testCase.Given.Domain.IpaDomain)
		} else {
			err = s.repository.updateIpaDomain(s.Log, testCase.Given.DB, nil)
		}

		// Check expectations for error and success scenario
		if testCase.Expected != nil {
			assert.EqualError(t, err, testCase.Expected.Error())
		} else {
			assert.NoError(t, err)
		}
	}
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
		err error
	)
	t := s.Suite.T()
	data := test.BuildDomainModel(test.OrgId)

	type TestCaseGiven struct {
		Stage  int
		DB     *gorm.DB
		Domain *model.Domain
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected error
	}

	testCases := []TestCase{
		{
			Name: "database error at FindByID",
			Given: TestCaseGiven{
				Stage:  1,
				DB:     s.DB,
				Domain: data,
			},
			Expected: fmt.Errorf("database error at FindByID"),
		},
		{
			Name: "database error at UPDATE INTO 'domains'",
			Given: TestCaseGiven{
				Stage:  2,
				DB:     s.DB,
				Domain: data,
			},
			Expected: fmt.Errorf("database error at UPDATE INTO 'domains'"),
		},
		{
			Name: "successful scenario",
			Given: TestCaseGiven{
				Stage:  2,
				DB:     s.DB,
				Domain: data,
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)

		// Prepare the db mock
		test_sql.UpdateUser(testCase.Given.Stage, s.mock, testCase.Expected, data)

		// Run for error or success
		c := app_context.CtxWithLog(app_context.CtxWithDB(context.Background(), s.DB), slog.Default())
		if testCase.Given.Domain != nil {
			err = s.repository.UpdateUser(c, test.OrgId, testCase.Given.Domain)
		} else {
			err = s.repository.UpdateUser(c, "", nil)
		}

		// Check expectations for error and success scenario
		if testCase.Expected != nil {
			assert.EqualError(t, err, testCase.Expected.Error())
		} else {
			assert.NoError(t, err)
		}
	}
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

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_ = r.DeleteById(context.Background(), "", model.NilUUID)
	})

	expectedErr = fmt.Errorf("code=404, message=unknown domain '%s'", d.DomainUuid.String())
	test_sql.DeleteByID(1, s.mock, gorm.ErrRecordNotFound, d)
	err := r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())

	expectedErr = fmt.Errorf("invalid transaction")
	test_sql.DeleteByID(1, s.mock, gorm.ErrInvalidTransaction, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())

	expectedErr = fmt.Errorf("code=404, message=unknown domain '%s'", d.DomainUuid.String())
	test_sql.DeleteByID(2, s.mock, gorm.ErrRecordNotFound, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())

	expectedErr = fmt.Errorf("code=404, message=unknown domain '%s'", d.DomainUuid.String())
	test_sql.DeleteByID(3, s.mock, gorm.ErrRecordNotFound, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())

	expectedErr = gorm.ErrInvalidTransaction
	test_sql.DeleteByID(3, s.mock, gorm.ErrInvalidTransaction, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, expectedErr.Error())

	// Success scenario
	expectedErr = nil
	test_sql.DeleteByID(3, s.mock, expectedErr, d)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.NoError(t, err)
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
	id := uint(helper.GenRandNum(1, 2^63))
	gormModel := builder_model.NewModel().WithID(id).Build()
	d := builder_model.NewDomain(gormModel).Build()

	assert.Panics(t, func() {
		_ = r.Register(nil, d.OrgId, d)
	})

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_ = r.Register(context.Background(), d.OrgId, d)
	})

	expectedErr = gorm.ErrDuplicatedKey
	test_sql.Register(1, d, s.mock, expectedErr)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, fmt.Sprintf("code=409, message=domain id '%s' is already registered.", d.DomainUuid))

	expectedErr = gorm.ErrInvalidField
	test_sql.Register(1, d, s.mock, expectedErr)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, "invalid field")

	d = builder_model.NewDomain(gormModel).
		WithIpaDomain(
			builder_model.NewIpaDomain().
				WithModel(gormModel).
				WithRealmName(&realm).
				WithRealmDomains(pq.StringArray{strings.ToLower(realm)}).
				WithServers([]model.IpaServer{
					builder_model.NewIpaServer(
						builder_model.NewModel().Build(),
					).WithIpaID(id).Build(),
				}).
				WithLocations([]model.IpaLocation{
					builder_model.NewIpaLocation(
						builder_model.NewModel().Build(),
					).WithIpaID(id).Build(),
				}).
				WithCaCerts([]model.IpaCert{
					builder_model.NewIpaCert(
						builder_model.NewModel().Build(),
						realm,
					).WithIpaID(id).Build(),
				}).
				Build(),
		).Build()
	test_sql.Register(1, d, s.mock, nil)
	test_sql.CreateIpaDomain(1, s.mock, gorm.ErrInvalidField, d.IpaDomain)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, "invalid field")

	// Success case - FIXME Flaky test
	test_sql.Register(1, d, s.mock, nil)
	test_sql.CreateIpaDomain(4, s.mock, nil, d.IpaDomain)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.NoError(t, err)

	// IpaDomain is nil
	test_sql.Register(1, d, s.mock, nil)
	d.IpaDomain = nil
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, "code=500, message='IpaDomain' cannot be nil")
}

func (s *DomainRepositorySuite) TestDeleteByIdLogError() {
	t := s.T()
	r := &domainRepository{}

	d := builder_model.NewDomain(builder_model.NewModel().WithID(1).Build()).Build()

	test_sql.DeleteByID(1, s.mock, gorm.ErrInvalidTransaction, d)
	err := r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, "invalid transaction")

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
