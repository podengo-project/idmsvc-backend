package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository *domainRepository
}

// https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
func (s *Suite) SetupTest() {
	var err error
	s.mock, s.DB, err = test.NewSqlMock(&gorm.Session{
		SkipHooks: true,
	})
	if err != nil {
		s.Suite.FailNow("Error calling gorm.Open: %s", err.Error())
		return
	}
	s.repository = &domainRepository{}
}

func (s *Suite) TestNewDomainRepository() {
	t := s.Suite.T()
	assert.NotPanics(t, func() {
		_ = NewDomainRepository()
	})
}

func (s *Suite) helperTestCreateIpaDomain(stage int, data *model.Ipa, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_domains","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
				WithArgs(
					data.Model.CreatedAt,
					data.Model.UpdatedAt,
					nil,

					data.RealmName,
					data.RealmDomains,
					data.ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(data.ID))
			}
		case 2:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_certs" ("created_at","updated_at","deleted_at","ipa_id","issuer","nickname","not_after","not_before","pem","serial_number","subject","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).
				WithArgs(
					data.CaCerts[0].CreatedAt,
					data.CaCerts[0].UpdatedAt,
					nil,

					data.CaCerts[0].IpaID,
					data.CaCerts[0].Issuer,
					data.CaCerts[0].Nickname,
					data.CaCerts[0].NotAfter,
					data.CaCerts[0].NotBefore,
					data.CaCerts[0].Pem,
					data.CaCerts[0].SerialNumber,
					data.CaCerts[0].Subject,
					data.CaCerts[0].ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(data.CaCerts[0].ID))
			}
		case 3:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_servers" ("created_at","updated_at","deleted_at","ipa_id","fqdn","rhsm_id","location","ca_server","hcc_enrollment_server","hcc_update_server","pk_init_server","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).
				WithArgs(
					data.Servers[0].CreatedAt,
					data.Servers[0].UpdatedAt,
					nil,

					data.Servers[0].IpaID,
					data.Servers[0].FQDN,
					data.Servers[0].RHSMId,
					data.Servers[0].Location,
					data.Servers[0].CaServer,
					data.Servers[0].HCCEnrollmentServer,
					data.Servers[0].HCCUpdateServer,
					data.Servers[0].PKInitServer,
					data.Servers[0].ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(data.Servers[0].ID))
			}
		case 4:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_locations" ("created_at","updated_at","deleted_at","ipa_id","name","description","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
				WithArgs(
					data.Locations[0].CreatedAt,
					data.Locations[0].UpdatedAt,
					nil,

					data.Locations[0].IpaID,
					data.Locations[0].Name,
					data.Locations[0].Description,
					data.Locations[0].ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(data.Locations[0].ID))
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *Suite) TestCreateIpaDomain() {
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
	err = s.repository.createIpaDomain(s.DB, 1, nil)
	assert.EqualError(t, err, "code=500, message='data' of type '*model.Ipa' cannot be nil")

	// Error on INSERT INTO "ipas"
	expectedError = fmt.Errorf(`error at INSERT INTO "ipas"`)
	s.helperTestCreateIpaDomain(1, &data, s.mock, expectedError)
	err = s.repository.createIpaDomain(s.DB, 1, &data)
	assert.EqualError(t, err, expectedError.Error())

	// Error on INSERT INTO "ipa_certs"
	expectedError = fmt.Errorf(`INSERT INTO "ipa_certs"`)
	s.helperTestCreateIpaDomain(2, &data, s.mock, expectedError)
	err = s.repository.createIpaDomain(s.DB, 1, &data)
	assert.EqualError(t, err, expectedError.Error())

	// Error on INSERT INTO "ipa_servers"
	expectedError = fmt.Errorf(`INSERT INTO "ipa_servers"`)
	s.helperTestCreateIpaDomain(3, &data, s.mock, expectedError)
	err = s.repository.createIpaDomain(s.DB, 1, &data)
	assert.EqualError(t, err, expectedError.Error())

	// Error on INSERT INTO "ipa_locations"
	expectedError = fmt.Errorf(`INSERT INTO "ipa_locations"`)
	s.helperTestCreateIpaDomain(4, &data, s.mock, expectedError)
	err = s.repository.createIpaDomain(s.DB, 1, &data)
	assert.EqualError(t, err, expectedError.Error())

	// Success scenario
	expectedError = nil
	s.helperTestCreateIpaDomain(4, &data, s.mock, nil)
	err = s.repository.createIpaDomain(s.DB, 1, &data)
	assert.NoError(t, err)
}

func (s *Suite) helperTestUpdateAgent(stage int, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	if stage == 0 {
		return
	}
	if stage < 0 {
		panic("'stage' cannot be lower than 0")
	}
	if stage > 6 {
		panic("'stage' cannot be greater than 6")
	}

	s.mock.MatchExpectationsInOrder(true)
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectExec := s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "domains" SET "created_at"=$1,"updated_at"=$2,"org_id"=$3,"domain_uuid"=$4,"domain_name"=$5,"title"=$6,"description"=$7,"type"=$8,"auto_enrollment_enabled"=$9 WHERE "domains"."deleted_at" IS NULL AND "id" = $10`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),

					data.OrgId,
					data.DomainUuid,
					data.DomainName,

					data.Title,
					data.Description,
					data.Type,
					data.AutoEnrollmentEnabled,

					data.ID,
				)
			if i == stage && expectedErr != nil {
				expectExec.WillReturnError(expectedErr)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			}
		case 2:
			expectExec := s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ipas" WHERE "ipas"."id" = $1`)).
				WithArgs(
					data.ID,
				).WillReturnResult(
				driver.RowsAffected(1),
			)
			if i == stage && expectedErr != nil {
				expectExec.WillReturnError(expectedErr)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			}
		case 3:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_domains","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),

					data.IpaDomain.RealmName,
					data.IpaDomain.RealmDomains,
					data.IpaDomain.ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(data.ID))
			}
		case 4:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_certs" ("created_at","updated_at","deleted_at","ipa_id","issuer","nickname","not_after","not_before","pem","serial_number","subject","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) ON CONFLICT ("id") DO UPDATE SET "ipa_id"="excluded"."ipa_id" RETURNING "id"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),

					data.IpaDomain.CaCerts[0].IpaID,

					data.IpaDomain.CaCerts[0].Issuer,
					data.IpaDomain.CaCerts[0].Nickname,
					data.IpaDomain.CaCerts[0].NotAfter,
					data.IpaDomain.CaCerts[0].NotBefore,
					data.IpaDomain.CaCerts[0].Pem,
					data.IpaDomain.CaCerts[0].SerialNumber,
					data.IpaDomain.CaCerts[0].Subject,
					data.IpaDomain.CaCerts[0].ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(data.IpaDomain.CaCerts[0].ID))
			}
		case 5:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_servers" ("created_at","updated_at","deleted_at","ipa_id","fqdn","rhsm_id","location","ca_server","hcc_enrollment_server","hcc_update_server","pk_init_server","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) ON CONFLICT ("id") DO UPDATE SET "ipa_id"="excluded"."ipa_id" RETURNING "id"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),

					data.IpaDomain.Servers[0].IpaID,

					data.IpaDomain.Servers[0].FQDN,
					data.IpaDomain.Servers[0].RHSMId,
					data.IpaDomain.Servers[0].Location,
					data.IpaDomain.Servers[0].CaServer,
					data.IpaDomain.Servers[0].HCCEnrollmentServer,
					data.IpaDomain.Servers[0].HCCUpdateServer,
					data.IpaDomain.Servers[0].PKInitServer,

					data.IpaDomain.Servers[0].ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(data.IpaDomain.Servers[0].ID))
			}
		case 6:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_locations" ("created_at","updated_at","deleted_at","ipa_id","name","description","id") VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT ("id") DO UPDATE SET "ipa_id"="excluded"."ipa_id" RETURNING "id"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),

					data.IpaDomain.Locations[0].IpaID,
					data.IpaDomain.Locations[0].Name,
					data.IpaDomain.Locations[0].Description,

					data.IpaDomain.Locations[0].ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(data.IpaDomain.Locations[0].ID))
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *Suite) TestUpdateAgent() {
	t := s.Suite.T()
	orgID := test.OrgId
	var (
		data *model.Domain = test.BuildDomainModel(test.OrgId)
		err  error
	)

	err = s.repository.UpdateAgent(nil, "", nil)
	assert.EqualError(t, err, "code=500, message='db' cannot be nil")

	err = s.repository.UpdateAgent(s.DB, "", nil)
	assert.EqualError(t, err, "'orgID' is empty")

	err = s.repository.UpdateAgent(s.DB, orgID, nil)
	assert.EqualError(t, err, "code=500, message='data' cannot be nil")

	s.helperTestUpdateAgent(6, data, s.mock, nil)
	err = s.repository.UpdateAgent(s.DB, orgID, data)
	require.NoError(t, err)
}

func (s *Suite) helperTestUpdateIpaDomain(stage int, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	if stage == 0 {
		return
	}
	if stage < 0 {
		panic("'stage' cannot be lower than 0")
	}
	if stage > 5 {
		panic("'stage' cannot be greater than 5")
	}

	s.mock.MatchExpectationsInOrder(true)
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectExec := s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ipas" WHERE "ipas"."id" = $1`)).
				WithArgs(
					data.Model.ID,
				)
			if i == stage && expectedErr != nil {
				expectExec.WillReturnError(expectedErr)
			} else {
				expectExec.WillReturnResult(
					driver.RowsAffected(1),
				)
			}
		case 2:
			expectExec := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_domains","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
				WithArgs(
					data.Model.CreatedAt,
					data.Model.UpdatedAt,
					data.Model.DeletedAt,

					data.IpaDomain.RealmName,
					data.IpaDomain.RealmDomains,
					data.ID,
				)
			if i == stage && expectedErr != nil {
				expectExec.WillReturnError(expectedErr)
			} else {
				expectExec.WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(data.ID))
			}
		case 3:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_certs" ("created_at","updated_at","deleted_at","ipa_id","issuer","nickname","not_after","not_before","pem","serial_number","subject","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT ("id") DO UPDATE RETURNING "id"`)).
				WithArgs(
					data.IpaDomain.CaCerts[0].Model.CreatedAt,
					data.IpaDomain.CaCerts[0].Model.UpdatedAt,
					data.IpaDomain.CaCerts[0].Model.DeletedAt,

					data.IpaDomain.CaCerts[0].IpaID,

					data.IpaDomain.CaCerts[0].Issuer,
					data.IpaDomain.CaCerts[0].Nickname,
					data.IpaDomain.CaCerts[0].NotAfter,
					data.IpaDomain.CaCerts[0].NotBefore,
					data.IpaDomain.CaCerts[0].Pem,
					data.IpaDomain.CaCerts[0].SerialNumber,
					data.IpaDomain.CaCerts[0].Subject,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(data.IpaDomain.CaCerts[0].ID))
			}
		case 4:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_servers" ("created_at","updated_at","deleted_at","ipa_id","fqdn","rhsm_id","location","ca_server","hcc_enrollment_server","hcc_update_server","pk_init_server","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT ("id") DO UPDATE RETURNING "id"`)).
				WithArgs(
					data.IpaDomain.Servers[0].Model.CreatedAt,
					data.IpaDomain.Servers[0].Model.UpdatedAt,
					data.IpaDomain.Servers[0].Model.DeletedAt,

					data.IpaDomain.Servers[0].IpaID,

					data.IpaDomain.Servers[0].FQDN,
					data.IpaDomain.Servers[0].RHSMId,
					data.IpaDomain.Servers[0].Location,
					data.IpaDomain.Servers[0].CaServer,
					data.IpaDomain.Servers[0].HCCEnrollmentServer,
					data.IpaDomain.Servers[0].HCCUpdateServer,
					data.IpaDomain.Servers[0].PKInitServer,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(data.IpaDomain.Servers[0].ID))
			}
		case 5:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_locations" ("created_at","updated_at","deleted_at","ipa_id","name","description","id") VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT ("id") DO UPDATE RETURNING "id"`)).
				WithArgs(
					data.IpaDomain.Locations[0].Model.CreatedAt,
					data.IpaDomain.Locations[0].Model.UpdatedAt,
					data.IpaDomain.Locations[0].Model.DeletedAt,

					data.IpaDomain.Locations[0].IpaID,
					data.IpaDomain.Locations[0].Name,
					data.IpaDomain.Locations[0].Description,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(data.IpaDomain.Locations[0].ID))
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *Suite) TestUpdateIpaDomain() {
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
			Expected: internal_errors.NilArgError("data' of type '*model.Ipa"),
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
		s.helperTestUpdateIpaDomain(testCase.Given.Stage, &data, s.mock, testCase.Expected)

		// Run for error or success
		if testCase.Given.Domain != nil {
			err = s.repository.updateIpaDomain(testCase.Given.DB, testCase.Given.Domain.IpaDomain)
		} else {
			err = s.repository.updateIpaDomain(testCase.Given.DB, nil)
		}

		// Check expectations for error and success scenario
		if testCase.Expected != nil {
			assert.Error(t, err, testCase.Expected.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func (s *Suite) TestList() {
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

	// db is nil
	output, count, err := r.List(nil, "", -1, -1)
	assert.EqualError(t, err, "code=500, message='db' cannot be nil")
	assert.Equal(t, int64(0), count)
	assert.Nil(t, output)

	// orgID is empty
	output, count, err = r.List(s.DB, "", -1, -1)
	assert.EqualError(t, err, "'orgID' is empty")
	assert.Equal(t, int64(0), count)
	assert.Nil(t, output)

	// offset is lower than 0
	output, count, err = r.List(s.DB, orgID, -1, -1)
	assert.EqualError(t, err, "'offset' is lower than 0")
	assert.Equal(t, int64(0), count)
	assert.Nil(t, output)

	// limit is lower than 0
	offset := 0
	output, count, err = r.List(s.DB, orgID, offset, -1)
	assert.EqualError(t, err, "'limit' is lower than 0")
	assert.Equal(t, int64(0), count)
	assert.Nil(t, output)

	// Return error
	limit := 5
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "domains" WHERE org_id = $1`)).
		WithArgs(orgID).
		WillReturnError(fmt.Errorf("an error happened"))
	output, count, err = r.List(s.DB, orgID, offset, limit)
	assert.EqualError(t, err, "an error happened")
	assert.Equal(t, int64(0), count)
	assert.Nil(t, output)

	// Success case
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "domains" WHERE org_id = $1`)).
		WithArgs(orgID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE org_id = $1 AND "domains"."deleted_at" IS NULL LIMIT 5`)).
		WithArgs(orgID).
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
	output, count, err = r.List(s.DB, orgID, offset, limit)
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

func (s *Suite) helperTestFindByID(stage int, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL ORDER BY "domains"."id" LIMIT 1`)).
				WithArgs(
					data.OrgId,
					data.DomainUuid,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				autoenrollment := false
				if data.AutoEnrollmentEnabled != nil {
					autoenrollment = *data.AutoEnrollmentEnabled
				}
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deletet_at",

					"org_id", "domain_uuid", "domain_name",
					"title", "description", "type",
					"auto_enrollment_enabled",
				}).
					AddRow(
						data.ID,
						data.CreatedAt,
						data.UpdatedAt,
						nil,

						data.OrgId,
						data.DomainUuid,
						data.DomainName,
						data.Title,
						data.Description,
						data.Type,
						autoenrollment,
					))
			}
		case 2:
			expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipas" WHERE id = $1 AND "ipas"."deleted_at" IS NULL ORDER BY "ipas"."id" LIMIT 1`)).
				WithArgs(1)
			if i == stage && expectedErr != nil {
				expectedQuery.WillReturnError(expectedErr)
			} else {
				expectedQuery.WillReturnRows(sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deletet_at",

					"realm_name", "realm_domains",
				}).AddRow(
					data.Model.ID,
					data.Model.CreatedAt,
					data.Model.UpdatedAt,
					nil,

					data.IpaDomain.RealmName,
					data.IpaDomain.RealmDomains,
				))
			}
		case 3:
			expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_certs" WHERE "ipa_certs"."ipa_id" = $1 AND "ipa_certs"."deleted_at" IS NULL`)).
				WithArgs(1)
			if i == stage && expectedErr != nil {
				expectedQuery.WillReturnError(expectedErr)
			} else {
				expectedQuery.
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "created_at", "updated_at", "deletet_at",

						"ipa_id", "issuer", "nickname",
						"not_after", "not_before", "serial_number",
						"subject", "pem",
					}).
						AddRow(
							data.IpaDomain.CaCerts[0].Model.ID,
							data.IpaDomain.CaCerts[0].Model.CreatedAt,
							data.IpaDomain.CaCerts[0].Model.UpdatedAt,
							data.IpaDomain.CaCerts[0].Model.DeletedAt,

							data.IpaDomain.CaCerts[0].IpaID,
							data.IpaDomain.CaCerts[0].Issuer,
							data.IpaDomain.CaCerts[0].Nickname,
							data.IpaDomain.CaCerts[0].NotAfter,
							data.IpaDomain.CaCerts[0].NotBefore,
							data.IpaDomain.CaCerts[0].SerialNumber,
							data.IpaDomain.CaCerts[0].Subject,
							data.IpaDomain.CaCerts[0].Pem,
						))
			}
		case 4:
			expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_locations" WHERE "ipa_locations"."ipa_id" = $1 AND "ipa_locations"."deleted_at" IS NULL`)).
				WithArgs(1)
			if i == stage && expectedErr != nil {
				expectedQuery.WillReturnError(expectedErr)
			} else {
				expectedQuery.WillReturnRows(sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deletet_at",

					"ipa_id",
					"name", "description",
				}).
					AddRow(
						data.IpaDomain.Locations[0].Model.ID,
						data.IpaDomain.Locations[0].Model.CreatedAt,
						data.IpaDomain.Locations[0].Model.UpdatedAt,
						data.IpaDomain.Locations[0].Model.DeletedAt,

						data.IpaDomain.Locations[0].IpaID,
						data.IpaDomain.Locations[0].Name,
						data.IpaDomain.Locations[0].Description,
					))
			}
		case 5:
			expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_servers" WHERE "ipa_servers"."ipa_id" = $1 AND "ipa_servers"."deleted_at" IS NULL`)).
				WithArgs(1)
			if i == stage && expectedErr != nil {
				expectedQuery.WillReturnError(expectedErr)
			} else {
				expectedQuery.WillReturnRows(sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deletet_at",

					"ipa_id", "fqdn", "rhsm_id", "location",
					"ca_server", "hcc_enrollment_server", "hcc_update_server",
					"pk_init_server",
				}).
					AddRow(
						data.IpaDomain.Servers[0].Model.ID,
						data.IpaDomain.Servers[0].Model.CreatedAt,
						data.IpaDomain.Servers[0].Model.UpdatedAt,
						data.IpaDomain.Servers[0].Model.DeletedAt,

						data.IpaDomain.Servers[0].IpaID,
						data.IpaDomain.Servers[0].FQDN,
						data.IpaDomain.Servers[0].RHSMId,
						data.IpaDomain.Servers[0].Location,
						data.IpaDomain.Servers[0].CaServer,
						data.IpaDomain.Servers[0].HCCEnrollmentServer,
						data.IpaDomain.Servers[0].HCCUpdateServer,
						data.IpaDomain.Servers[0].PKInitServer,
					))
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *Suite) TestFindByID() {
	t := s.T()
	r := &domainRepository{}
	s.mock.MatchExpectationsInOrder(true)

	currentTime := time.Now()
	notBefore := currentTime
	notAfter := currentTime.Add(356 * 24 * time.Hour)
	domainID := uuid.MustParse("c5d2c9d0-2b2f-11ee-8ec5-482ae3863d30")
	data := &model.Domain{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		OrgId:                 "12345",
		DomainUuid:            domainID,
		DomainName:            pointy.String("mydomain.example"),
		Title:                 pointy.String("My Example Domain"),
		Description:           pointy.String("My long description for my example domain"),
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain: &model.Ipa{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			RealmName:    pointy.String("MYDOMAIN.EXAMPLE"),
			RealmDomains: pq.StringArray{"mydomain.example"},
			CaCerts: []model.IpaCert{
				{
					Model: gorm.Model{
						ID:        2,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					IpaID:        1,
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
						ID:        3,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					IpaID:               1,
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
						ID:        4,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					IpaID:       1,
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
	domain, err = r.FindByID(nil, "", uuid.Nil)
	assert.EqualError(t, err, "code=500, message='db' cannot be nil")
	assert.Nil(t, domain)

	// Check path when an error hapens into the sql statement
	expectedErr = fmt.Errorf(`error at SELECT * FROM "domains"`)
	s.helperTestFindByID(1, data, s.mock, expectedErr)
	domain, err = r.FindByID(s.DB, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, domain)

	// Check path when a domain type is NULL
	expectedErr = internal_errors.NilArgError("Type")
	s.helperTestFindByID(1, dataTypeNil, s.mock, nil)
	domain, err = r.FindByID(s.DB, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, domain)

	// Check for 'ipas' record not found
	expectedErr = gorm.ErrRecordNotFound
	s.helperTestFindByID(2, data, s.mock, expectedErr)
	domain, err = r.FindByID(s.DB, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, domain)

	// Successful scenario
	expectedErr = nil
	s.helperTestFindByID(5, data, s.mock, expectedErr)
	domain, err = r.FindByID(s.DB, data.OrgId, data.DomainUuid)
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

func (s *Suite) helperTestUpdateUser(stage int, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	if stage == 0 {
		return
	}
	if stage < 0 {
		panic("'stage' cannot be lower than 0")
	}
	if stage > 2 {
		panic("'stage' cannot be greater than 3")
	}

	s.mock.MatchExpectationsInOrder(true)
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			if i == stage && expectedErr != nil {
				s.helperTestFindByID(1, data, mock, expectedErr)
			} else {
				s.helperTestFindByID(5, data, mock, nil)
			}
		case 2: // Update
			expectExec := mock.ExpectExec(regexp.QuoteMeta(`UPDATE "domains" SET "auto_enrollment_enabled"=$1,"description"=$2,"title"=$3 WHERE (org_id = $4 AND domain_uuid = $5) AND "domains"."deleted_at" IS NULL AND "id" = $6`)).
				WithArgs(
					data.AutoEnrollmentEnabled,
					data.Description,
					data.Title,

					data.OrgId,
					data.DomainUuid,
					data.ID,
				)
			if i == stage && expectedErr != nil {
				expectExec.WillReturnError(expectedErr)
			} else {
				expectExec.WillReturnResult(
					driver.RowsAffected(1))
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *Suite) TestUpdateUser() {
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
			Name: "Wrong arguments: db is nil",
			Given: TestCaseGiven{
				Stage:  0,
				DB:     nil,
				Domain: nil,
			},
			Expected: internal_errors.NilArgError("db"),
		},
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
		s.helperTestUpdateUser(testCase.Given.Stage, data, s.mock, testCase.Expected)

		// Run for error or success
		if testCase.Given.Domain != nil {
			err = s.repository.UpdateUser(testCase.Given.DB, test.OrgId, testCase.Given.Domain)
		} else {
			err = s.repository.UpdateUser(testCase.Given.DB, "", nil)
		}

		// Check expectations for error and success scenario
		if testCase.Expected != nil {
			assert.Error(t, err, testCase.Expected.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

// ---------------- Test for private methods ---------------------

func (s *Suite) TestCheckCommon() {
	t := s.T()
	r := &domainRepository{}

	err := r.checkCommon(nil, "")
	assert.EqualError(t, err, "code=500, message='db' cannot be nil")

	err = r.checkCommon(s.DB, "")
	assert.EqualError(t, err, "'orgID' is empty")

	err = r.checkCommon(s.DB, "12345")
	assert.NoError(t, err)
}

func (s *Suite) TestCheckCommonAndUUID() {
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

func (s *Suite) TestCheckCommonAndData() {
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

func (s *Suite) TestCheckCommonAndDataAndType() {
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

func (s *Suite) TestCreateDomainToken() {
	var (
		key       []byte        = []byte("secret")
		validity  time.Duration = 2 * time.Hour
		testOrgID string        = "12345"
	)
	t := s.T()
	r := &domainRepository{}
	drt, err := r.CreateDomainToken(key, validity, testOrgID, public.RhelIdm)
	assert.NoError(t, err)
	assert.Equal(t, drt.DomainType, public.RhelIdm)
	assert.NotEmpty(t, drt.DomainId)
	assert.Equal(
		t,
		drt.DomainId,
		token.TokenDomainId(token.DomainRegistrationToken(drt.DomainToken)),
	)
	assert.Greater(t, drt.ExpirationNS, uint64(time.Now().UnixNano()))
}

func (s *Suite) TestPrepareUpdateUser() {
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

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
