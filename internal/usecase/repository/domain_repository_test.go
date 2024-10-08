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

func (s *DomainRepositorySuite) helperTestCreateIpaDomain(stage int, domainID uint, data *model.Ipa, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_domains","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					nil,

					data.RealmName,
					data.RealmDomains,
					domainID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(domainID))
			}
		case 2:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_certs" ("created_at","updated_at","deleted_at","ipa_id","issuer","nickname","not_after","not_before","pem","serial_number","subject") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					nil,

					domainID,
					data.CaCerts[0].Issuer,
					data.CaCerts[0].Nickname,
					data.CaCerts[0].NotAfter,
					data.CaCerts[0].NotBefore,
					data.CaCerts[0].Pem,
					data.CaCerts[0].SerialNumber,
					data.CaCerts[0].Subject,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(uint(2)))
			}
		case 3:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_servers" ("created_at","updated_at","deleted_at","ipa_id","fqdn","rhsm_id","location","ca_server","hcc_enrollment_server","hcc_update_server","pk_init_server") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					nil,

					domainID,
					data.Servers[0].FQDN,
					data.Servers[0].RHSMId,
					data.Servers[0].Location,
					data.Servers[0].CaServer,
					data.Servers[0].HCCEnrollmentServer,
					data.Servers[0].HCCUpdateServer,
					data.Servers[0].PKInitServer,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(uint(3)))
			}
		case 4:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_locations" ("created_at","updated_at","deleted_at","ipa_id","name","description") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					nil,

					domainID,
					data.Locations[0].Name,
					data.Locations[0].Description,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(uint(4)))
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *DomainRepositorySuite) TestCreateIpaDomain() {
	t := s.Suite.T()
	currentTime := time.Now()
	domainID := uint(1)
	domainName := "mydomain.example"
	realmName := strings.ToUpper(domainName)

	gormModel := builder_model.NewModel().
		WithCreatedAt(currentTime).
		WithUpdatedAt(currentTime).
		Build()

	var (
		err           error
		expectedError error
	)
	data := builder_model.NewIpaDomain().
		WithModel(
			builder_model.NewModel().
				WithID(domainID).
				WithCreatedAt(currentTime).
				WithUpdatedAt(currentTime).
				Build()).
		WithRealmName(pointy.String(realmName)).
		WithRealmDomains([]string{domainName}).
		WithCaCerts([]model.IpaCert{
			builder_model.NewIpaCert(gormModel, realmName).
				WithIpaID(domainID).
				Build(),
		}).
		WithServers([]model.IpaServer{
			builder_model.NewIpaServer(gormModel).
				WithIpaID(domainID).
				Build(),
		}).
		WithLocations([]model.IpaLocation{
			builder_model.NewIpaLocation(gormModel).
				WithIpaID(domainID).
				Build(),
		}).
		Build()
	require.NotNil(t, data)

	// Check nil
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, nil)
	require.EqualError(t, err, "code=500, message='data' cannot be nil")

	// Error on INSERT INTO "ipas"
	expectedError = fmt.Errorf(`error at INSERT INTO "ipas"`)
	s.helperTestCreateIpaDomain(1, domainID, data, s.mock, expectedError)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data)
	assert.EqualError(t, err, expectedError.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Error on INSERT INTO "ipa_certs"
	expectedError = fmt.Errorf(`error at INSERT INTO "ipa_certs"`)
	s.helperTestCreateIpaDomain(2, domainID, data, s.mock, expectedError)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data)
	assert.EqualError(t, err, expectedError.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Error on INSERT INTO "ipa_servers"
	expectedError = fmt.Errorf(`error at INSERT INTO "ipa_servers"`)
	s.helperTestCreateIpaDomain(3, domainID, data, s.mock, expectedError)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data)
	assert.EqualError(t, err, expectedError.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Error on INSERT INTO "ipa_locations"
	expectedError = fmt.Errorf(`error at INSERT INTO "ipa_locations"`)
	s.helperTestCreateIpaDomain(4, domainID, data, s.mock, expectedError)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data)
	assert.EqualError(t, err, expectedError.Error())
	require.NoError(t, s.mock.ExpectationsWereMet())

	// Success scenario
	expectedError = nil
	s.helperTestCreateIpaDomain(4, domainID, data, s.mock, nil)
	err = s.repository.createIpaDomain(s.Log, s.DB, domainID, data)
	assert.NoError(t, err)
	require.NoError(t, s.mock.ExpectationsWereMet())
}

func (s *DomainRepositorySuite) TestUpdateErrors() {
	t := s.Suite.T()
	orgID := test.OrgId
	domainID := uint(1)
	data := test.BuildDomainModel(test.OrgId)
	var err error

	assert.Panics(t, func() {
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

	s.mock.MatchExpectationsInOrder(true)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.EqualError(t, err, "data.Model.ID cannot be 0")
	require.NoError(t, s.mock.ExpectationsWereMet())

	s.mock.MatchExpectationsInOrder(true)
	data.Model.ID = domainID
	s.helperTestFindByID(1, domainID, data, s.mock, nil)
	helperTestFindIpaByID(1, domainID, data, s.mock, fmt.Errorf("record not found"))
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.EqualError(t, err, "record not found")
	require.NoError(t, s.mock.ExpectationsWereMet())
}

func helperTestUpdateDomain(stage int, domainID uint, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectedQuery := mock.ExpectExec(regexp.QuoteMeta(`UPDATE "domains" SET "created_at"=$1,"updated_at"=$2,"org_id"=$3,"domain_uuid"=$4,"domain_name"=$5,"title"=$6,"description"=$7,"type"=$8,"auto_enrollment_enabled"=$9 WHERE (org_id = $10 AND domain_uuid = $11) AND "domains"."deleted_at" IS NULL AND "id" = $12`)).
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

					data.OrgId,
					data.DomainUuid,
					domainID,
				)
			if i == stage && expectedErr != nil {
				expectedQuery.WillReturnError(expectedErr)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(1, 1))
			}

		case 2:
			expectedExec := mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ipas" WHERE "ipas"."id" = $1`)).
				WithArgs(
					domainID,
				)
			if i == stage && expectedErr != nil {
				expectedExec.WillReturnError(expectedErr)
			} else {
				expectedExec.WillReturnResult(
					driver.RowsAffected(1),
				)
			}

		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *DomainRepositorySuite) TestUpdateAgentSuccess() {
	var err error
	t := s.Suite.T()
	domainID := uint(1)
	orgID := test.OrgId
	currentTime := time.Now()
	domainName := helper.GenRandDomainName(2)
	realmName := strings.ToUpper(domainName)
	gormModel := builder_model.NewModel().
		WithID(domainID).
		WithCreatedAt(currentTime).
		WithUpdatedAt(currentTime).
		Build()
	gormModelEmptyID := builder_model.NewModel().
		WithID(0).
		WithCreatedAt(currentTime).
		WithUpdatedAt(currentTime).
		Build()
	autoEnrollment := true
	data := builder_model.NewDomain(gormModel).
		WithDomainUUID(test.DomainUUID).
		WithOrgID(orgID).
		WithTitle(&domainName).
		WithDescription(pointy.String("Long description for " + domainName)).
		WithAutoEnrollmentEnabled(&autoEnrollment).
		WithIpaDomain(builder_model.NewIpaDomain().
			WithModel(builder_model.NewModel().
				WithID(domainID).
				WithCreatedAt(currentTime).
				WithUpdatedAt(currentTime).
				Build()).
			WithRealmDomains(pq.StringArray{domainName}).
			WithCaCerts([]model.IpaCert{
				builder_model.NewIpaCert(gormModelEmptyID, realmName).
					WithIpaID(domainID).
					Build(),
			}).
			WithRealmName(&realmName).
			WithServers([]model.IpaServer{
				builder_model.NewIpaServer(gormModelEmptyID).
					WithIpaID(domainID).
					WithFQDN(helper.GenRandFQDNWithDomain(domainName)).
					Build(),
			}).
			WithLocations([]model.IpaLocation{
				builder_model.NewIpaLocation(gormModelEmptyID).
					WithIpaID(domainID).
					Build(),
			}).
			Build(),
		).Build()

	s.mock.MatchExpectationsInOrder(false)
	s.helperTestFindByID(1, domainID, data, s.mock, nil)
	helperTestFindIpaByID(4, domainID, data, s.mock, nil)
	helperTestUpdateDomain(2, domainID, data, s.mock, nil)
	s.helperTestCreateIpaDomain(4, domainID, data.IpaDomain, s.mock, nil)
	err = s.repository.UpdateAgent(s.Ctx, orgID, data)
	require.NoError(t, err)
	require.NoError(t, s.mock.ExpectationsWereMet())
}

func (s *DomainRepositorySuite) helperTestUpdateIpaDomain(stage int, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	if stage == 0 {
		return
	}
	if stage < 0 {
		panic("'stage' cannot be lower than 0")
	}
	if stage > 5 {
		panic("'stage' cannot be greater than 5")
	}

	mock.MatchExpectationsInOrder(true)
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			data.IpaDomain.Model.ID = data.Model.ID
			expectExec := mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ipas" WHERE "ipas"."id" = $1`)).
				WithArgs(
					data.IpaDomain.Model.ID,
				)
			if i == stage && expectedErr != nil {
				expectExec.WillReturnError(expectedErr)
			} else {
				expectExec.WillReturnResult(
					driver.RowsAffected(1),
				)
			}
		case 2:
			data.IpaDomain.Model.ID = data.Model.ID
			expectExec := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_domains","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					nil,

					data.IpaDomain.RealmName,
					data.IpaDomain.RealmDomains,
					sqlmock.AnyArg(),
				)
			if i == stage && expectedErr != nil {
				expectExec.WillReturnError(expectedErr)
			} else {
				expectExec.WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(data.Model.ID))
			}
		case 3:
			for j := range data.IpaDomain.CaCerts {
				data.IpaDomain.CaCerts[j].Model.ID = 0
				expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_certs" ("created_at","updated_at","deleted_at","ipa_id","issuer","nickname","not_after","not_before","pem","serial_number","subject") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						nil,

						data.Model.ID,

						data.IpaDomain.CaCerts[j].Issuer,
						data.IpaDomain.CaCerts[j].Nickname,
						data.IpaDomain.CaCerts[j].NotAfter,
						data.IpaDomain.CaCerts[j].NotBefore,
						data.IpaDomain.CaCerts[j].Pem,
						data.IpaDomain.CaCerts[j].SerialNumber,
						data.IpaDomain.CaCerts[j].Subject,
					)
				if i == stage && expectedErr != nil {
					expectQuery.WillReturnError(expectedErr)
				} else {
					expectQuery.WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(uint(10 + j)))
				}
			}
		case 4:
			for j := range data.IpaDomain.Servers {
				data.IpaDomain.Servers[j].Model.ID = 0
				expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_servers" ("created_at","updated_at","deleted_at","ipa_id","fqdn","rhsm_id","location","ca_server","hcc_enrollment_server","hcc_update_server","pk_init_server") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						nil,

						data.Model.ID,

						data.IpaDomain.Servers[j].FQDN,
						data.IpaDomain.Servers[j].RHSMId,
						data.IpaDomain.Servers[j].Location,
						data.IpaDomain.Servers[j].CaServer,
						data.IpaDomain.Servers[j].HCCEnrollmentServer,
						data.IpaDomain.Servers[j].HCCUpdateServer,
						data.IpaDomain.Servers[j].PKInitServer,
					)
				if i == stage && expectedErr != nil {
					expectQuery.WillReturnError(expectedErr)
				} else {
					expectQuery.WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(uint(20 + j)))
				}
			}
		case 5:
			for j := range data.IpaDomain.Locations {
				data.IpaDomain.Locations[j].Model.ID = 0
				expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_locations" ("created_at","updated_at","deleted_at","ipa_id","name","description") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
					WithArgs(
						data.IpaDomain.Locations[j].Model.CreatedAt,
						data.IpaDomain.Locations[j].Model.UpdatedAt,
						data.IpaDomain.Locations[j].Model.DeletedAt,

						data.Model.ID,
						data.IpaDomain.Locations[j].Name,
						data.IpaDomain.Locations[j].Description,
					)
				if i == stage && expectedErr != nil {
					expectQuery.WillReturnError(expectedErr)
				} else {
					expectQuery.WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(uint(30 + j)))
				}
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *DomainRepositorySuite) helperTestUpdateUser(stage int, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
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
				s.helperTestFindByID(1, 1, data, mock, expectedErr)
			} else {
				s.helperTestFindByID(1, 1, data, mock, nil)
				helperTestFindIpaByID(4, 1, data, mock, nil)
			}
		case 2: // Update
			expectExec := mock.ExpectExec(regexp.QuoteMeta(`UPDATE "domains" SET "auto_enrollment_enabled"=$1,"description"=$2,"title"=$3 WHERE (org_id = $4 AND domain_uuid = $5) AND "domains"."deleted_at"`)).
				WithArgs(
					data.AutoEnrollmentEnabled,
					data.Description,
					data.Title,

					data.OrgId,
					data.DomainUuid,
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

func (s *DomainRepositorySuite) TestUpdateIpaDomain() {
	var (
		err error
	)
	t := s.Suite.T()
	orgID := "11111"
	currentTime := time.Now()
	domainID := uint(1)
	domainName := helper.GenRandDomainName(2)
	realmName := strings.ToUpper(domainName)
	gormModel := builder_model.NewModel().
		WithID(domainID).
		WithCreatedAt(currentTime).
		WithUpdatedAt(currentTime).
		Build()
	gormModelEmptyID := builder_model.NewModel().
		WithID(0).
		WithCreatedAt(currentTime).
		WithUpdatedAt(currentTime).
		Build()
	data := builder_model.NewDomain(gormModel).
		WithDomainUUID(test.DomainUUID).
		WithOrgID(orgID).
		WithIpaDomain(builder_model.NewIpaDomain().
			WithModel(builder_model.NewModel().
				WithID(domainID).
				WithCreatedAt(currentTime).
				WithUpdatedAt(currentTime).
				Build()).
			WithRealmDomains(pq.StringArray{domainName}).
			WithCaCerts([]model.IpaCert{
				builder_model.NewIpaCert(gormModelEmptyID, realmName).
					WithIpaID(domainID).
					Build(),
			}).
			WithRealmName(&realmName).
			WithServers([]model.IpaServer{
				builder_model.NewIpaServer(gormModelEmptyID).
					WithIpaID(domainID).
					WithFQDN(helper.GenRandFQDNWithDomain(domainName)).
					Build(),
			}).
			WithLocations([]model.IpaLocation{
				builder_model.NewIpaLocation(gormModelEmptyID).
					WithIpaID(domainID).
					Build(),
			}).
			Build(),
		).Build()
	dataEmptyID := builder_model.NewDomain(gormModel).
		WithDomainUUID(test.DomainUUID).
		WithOrgID(orgID).
		WithIpaDomain(builder_model.NewIpaDomain().
			WithModel(gormModelEmptyID).
			Build()).
		Build()

	// Wrong arguments: log is nil
	expectedError := internal_errors.NilArgError("log")
	s.helperTestUpdateIpaDomain(0, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(nil, nil, "", nil)
	require.EqualError(t, err, expectedError.Error())

	// Wrong arguments: db is nil
	expectedError = internal_errors.NilArgError("db")
	s.helperTestUpdateIpaDomain(0, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), nil, "", nil)
	require.EqualError(t, err, expectedError.Error())

	// Wrong arguments: orgID is empty
	expectedError = internal_errors.EmptyArgError("orgID")
	s.helperTestUpdateIpaDomain(0, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), s.DB, "", nil)
	require.EqualError(t, err, expectedError.Error())

	// Wrong arguments: data is nil
	expectedError = internal_errors.NilArgError("data")
	s.helperTestUpdateIpaDomain(0, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), s.DB, orgID, nil)
	require.EqualError(t, err, expectedError.Error())

	// Wrong arguments: data.Model.ID is empty
	expectedError = internal_errors.EmptyArgError("data.Model.ID")
	s.helperTestUpdateIpaDomain(0, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), s.DB, orgID, dataEmptyID.IpaDomain)
	require.EqualError(t, err, expectedError.Error())

	// database error at DELETE FROM 'ipas'
	expectedError = fmt.Errorf("database error at DELETE FROM 'ipas'")
	s.helperTestUpdateIpaDomain(1, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), s.DB, orgID, data.IpaDomain)
	require.EqualError(t, err, expectedError.Error())

	// database error at INSERT INTO 'ipas'
	expectedError = fmt.Errorf("database error at INSERT INTO 'ipas'")
	s.helperTestUpdateIpaDomain(2, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), s.DB, orgID, data.IpaDomain)
	require.EqualError(t, err, expectedError.Error())

	// database error at INSERT INTO 'ipa_certs'
	expectedError = fmt.Errorf("database error at INSERT INTO 'ipa_certs'")
	s.helperTestUpdateIpaDomain(3, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), s.DB, orgID, data.IpaDomain)
	require.EqualError(t, err, expectedError.Error())

	// database error at INSERT INTO 'ipa_servers'
	expectedError = fmt.Errorf("database error at INSERT INTO 'ipa_servers'")
	s.helperTestUpdateIpaDomain(4, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), s.DB, orgID, data.IpaDomain)
	require.EqualError(t, err, expectedError.Error())

	// database error at INSERT INTO 'ipa_locations'
	expectedError = fmt.Errorf("database error at INSERT INTO 'ipa_locations'")
	s.helperTestUpdateIpaDomain(5, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), s.DB, orgID, data.IpaDomain)
	require.EqualError(t, err, expectedError.Error())

	// Success scenario
	expectedError = nil
	s.helperTestUpdateIpaDomain(5, data, s.mock, expectedError)
	err = s.repository.updateIpaDomain(slog.Default(), s.DB, orgID, data.IpaDomain)
	require.NoError(t, err)
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

func (s *DomainRepositorySuite) helperTestFindByID(stage int, domainID uint, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL ORDER BY "domains"."id" LIMIT $3`)).
				WithArgs(
					data.OrgId,
					data.DomainUuid,
					1,
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
						domainID,
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
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *DomainRepositorySuite) TestFindByID() {
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
	assert.Panics(t, func() {
		_, _ = r.FindByID(nil, "", uuid.Nil)
	})

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_, _ = r.FindByID(context.Background(), "", uuid.Nil)
	})

	// Check path when an error hapens into the sql statement
	expectedErr = fmt.Errorf(`error at SELECT * FROM "domains"`)
	s.helperTestFindByID(1, 1, data, s.mock, expectedErr)
	domain, err = r.FindByID(s.Ctx, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, domain)

	// Check path when a domain type is NULL
	expectedErr = internal_errors.NilArgError("Type")
	s.helperTestFindByID(1, 1, dataTypeNil, s.mock, nil)
	domain, err = r.FindByID(s.Ctx, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, domain)

	// Check for 'ipas' record not found
	expectedErr = gorm.ErrRecordNotFound
	s.helperTestFindByID(1, 1, data, s.mock, nil)
	helperTestFindIpaByID(1, 1, data, s.mock, expectedErr)
	domain, err = r.FindByID(s.Ctx, data.OrgId, data.DomainUuid)
	require.NoError(t, s.mock.ExpectationsWereMet())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, domain)

	// Successful scenario
	expectedErr = nil
	s.helperTestFindByID(1, 1, data, s.mock, nil)
	helperTestFindIpaByID(4, 1, data, s.mock, expectedErr)
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
		s.helperTestUpdateUser(testCase.Given.Stage, data, s.mock, testCase.Expected)

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
		require.NoError(t, s.mock.ExpectationsWereMet())
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

func (s *DomainRepositorySuite) helperTestDeleteById(stage int, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL ORDER BY "domains"."id" LIMIT $3`)).
				WithArgs(
					data.OrgId,
					data.DomainUuid,
					1,
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
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL LIMIT $3`)).
				WithArgs(
					data.OrgId,
					data.DomainUuid,
					1,
				)
			if i == stage && expectedErr != nil {
				if expectedErr == gorm.ErrRecordNotFound {
					expectQuery.WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(int64(0)))
				} else {
					expectQuery.WillReturnError(expectedErr)
				}
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(int64(1)))
			}
		case 3:
			expectQuery := s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."id" = $3`)).
				WithArgs(
					data.OrgId,
					data.DomainUuid,
					data.ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnResult(driver.RowsAffected(1))
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *DomainRepositorySuite) TestDeleteById() {
	t := s.T()
	r := &domainRepository{}

	d := builder_model.NewDomain(builder_model.NewModel().WithID(1).Build()).Build()

	require.Panics(t, func() {
		_ = r.DeleteById(nil, "", model.NilUUID)
	})

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_ = r.DeleteById(context.Background(), "", model.NilUUID)
	})

	s.helperTestDeleteById(1, d, s.mock, gorm.ErrRecordNotFound)
	err := r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, fmt.Sprintf("code=404, message=unknown domain '%s'", d.DomainUuid.String()))

	s.helperTestDeleteById(1, d, s.mock, gorm.ErrInvalidTransaction)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, "invalid transaction")

	s.helperTestDeleteById(2, d, s.mock, gorm.ErrRecordNotFound)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, fmt.Sprintf("code=404, message=unknown domain '%s'", d.DomainUuid.String()))

	s.helperTestDeleteById(3, d, s.mock, gorm.ErrRecordNotFound)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, fmt.Sprintf("code=404, message=unknown domain '%s'", d.DomainUuid.String()))

	s.helperTestDeleteById(3, d, s.mock, gorm.ErrInvalidTransaction)
	err = r.DeleteById(s.Ctx, d.OrgId, d.DomainUuid)
	require.EqualError(t, err, gorm.ErrInvalidTransaction.Error())

	// Success scenario
	s.helperTestDeleteById(3, d, s.mock, nil)
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

func (s *DomainRepositorySuite) helperTestRegister(stage int, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "domains" ("created_at","updated_at","deleted_at","org_id","domain_uuid","domain_name","title","description","type","auto_enrollment_enabled","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
				WithArgs(
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

					data.Model.ID,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(data.Model.ID))
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *DomainRepositorySuite) TestRegisterSuccess() {
	var err error
	t := s.T()
	r := &domainRepository{}
	realm := strings.ToUpper(helper.GenRandDomainName(3))
	domainID := uint(helper.GenRandNum(1, 2^63-4))
	gormModel := builder_model.NewModel().WithID(domainID).Build()
	gormModelServer := builder_model.NewModel().WithID(domainID + 1).Build()
	gormModelLocation := builder_model.NewModel().WithID(domainID + 2).Build()
	gormModelCaCert := builder_model.NewModel().WithID(domainID + 3).Build()
	d := builder_model.NewDomain(gormModel).
		WithIpaDomain(
			builder_model.NewIpaDomain().
				WithModel(gormModel).
				WithRealmName(&realm).
				WithRealmDomains(pq.StringArray{strings.ToLower(realm)}).
				WithServers([]model.IpaServer{
					builder_model.NewIpaServer(gormModelServer).WithIpaID(domainID).Build(),
				}).
				WithLocations([]model.IpaLocation{
					builder_model.NewIpaLocation(gormModelLocation).WithIpaID(domainID).Build(),
				}).
				WithCaCerts([]model.IpaCert{
					builder_model.NewIpaCert(
						gormModelCaCert,
						realm,
					).WithIpaID(domainID).Build(),
				}).
				Build(),
		).Build()
	s.helperTestRegister(1, d, s.mock, nil)
	s.helperTestCreateIpaDomain(4, domainID, d.IpaDomain, s.mock, nil)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.NoError(t, err)
}

func (s *DomainRepositorySuite) TestRegisterFailure() {
	var err error
	t := s.T()
	r := &domainRepository{}
	realm := strings.ToUpper(helper.GenRandDomainName(3))
	domainID := uint(helper.GenRandNum(1, 2^63))
	currentTime := time.Now()
	gormModel := builder_model.NewModel().
		WithID(domainID).
		WithCreatedAt(currentTime).
		WithUpdatedAt(currentTime).
		Build()
	d := builder_model.NewDomain(gormModel).Build()

	assert.Panics(t, func() {
		_ = r.Register(nil, d.OrgId, d)
	})

	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_ = r.Register(context.Background(), d.OrgId, d)
	})

	s.helperTestRegister(1, d, s.mock, gorm.ErrDuplicatedKey)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, fmt.Sprintf("code=409, message=domain id '%s' is already registered.", d.DomainUuid))

	s.helperTestRegister(1, d, s.mock, gorm.ErrInvalidField)
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
						builder_model.NewModel().
							WithID(2).
							WithCreatedAt(currentTime).
							WithUpdatedAt(currentTime).
							Build(),
					).
						WithIpaID(domainID).
						Build(),
				}).
				WithLocations([]model.IpaLocation{
					builder_model.NewIpaLocation(
						builder_model.NewModel().
							WithID(3).
							WithCreatedAt(currentTime).
							WithUpdatedAt(currentTime).
							Build(),
					).
						WithIpaID(domainID).
						Build(),
				}).
				WithCaCerts([]model.IpaCert{
					builder_model.NewIpaCert(
						builder_model.NewModel().
							WithID(4).
							WithCreatedAt(currentTime).
							WithUpdatedAt(currentTime).
							Build(),
						realm,
					).
						WithIpaID(domainID).
						Build(),
				}).
				Build(),
		).Build()
	s.helperTestRegister(1, d, s.mock, nil)
	s.helperTestCreateIpaDomain(1, domainID, d.IpaDomain, s.mock, gorm.ErrInvalidField)
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, "invalid field")
	require.NoError(t, s.mock.ExpectationsWereMet())

	// IpaDomain is nil
	s.helperTestRegister(1, d, s.mock, nil)
	d.IpaDomain = nil
	err = r.Register(s.Ctx, d.OrgId, d)
	require.EqualError(t, err, "code=500, message='IpaDomain' cannot be nil")
	require.NoError(t, s.mock.ExpectationsWereMet())
}

func helperTestFindIpaByID(stage int, domainID uint, data *model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipas" WHERE id = $1 AND "ipas"."deleted_at" IS NULL ORDER BY "ipas"."id" LIMIT $2`)).
				WithArgs(
					domainID,
					1,
				)
			if i == stage && expectedErr != nil {
				expectedQuery.WillReturnError(expectedErr)
			} else {
				expectedQuery.WillReturnRows(sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deletet_at",

					"realm_name", "realm_domains",
				}).AddRow(
					domainID,
					data.Model.CreatedAt,
					data.Model.UpdatedAt,
					data.Model.DeletedAt,

					data.IpaDomain.RealmName,
					data.IpaDomain.RealmDomains,
				))
			}
		case 2:
			expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_certs" WHERE "ipa_certs"."ipa_id" = $1 AND "ipa_certs"."deleted_at" IS NULL`)).
				WithArgs(domainID)
			if i == stage && expectedErr != nil {
				expectedQuery.WillReturnError(expectedErr)
			} else {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deletet_at",

					"ipa_id", "issuer", "nickname",
					"not_after", "not_before", "serial_number",
					"subject", "pem",
				})
				for j := range data.IpaDomain.CaCerts {
					rows.AddRow(
						domainID+uint(j)+1,
						data.IpaDomain.CaCerts[j].Model.CreatedAt,
						data.IpaDomain.CaCerts[j].Model.UpdatedAt,
						data.IpaDomain.CaCerts[j].Model.DeletedAt,

						domainID,
						data.IpaDomain.CaCerts[j].Issuer,
						data.IpaDomain.CaCerts[j].Nickname,
						data.IpaDomain.CaCerts[j].NotAfter,
						data.IpaDomain.CaCerts[j].NotBefore,
						data.IpaDomain.CaCerts[j].SerialNumber,
						data.IpaDomain.CaCerts[j].Subject,
						data.IpaDomain.CaCerts[j].Pem,
					)
				}
				expectedQuery.WillReturnRows(rows)
			}
		case 3:
			expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_locations" WHERE "ipa_locations"."ipa_id" = $1 AND "ipa_locations"."deleted_at" IS NULL`)).
				WithArgs(domainID)
			if i == stage && expectedErr != nil {
				expectedQuery.WillReturnError(expectedErr)
			} else {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deletet_at",

					"ipa_id",
					"name", "description",
				})
				for j := range data.IpaDomain.Locations {
					rows.AddRow(
						domainID+uint(j)+1,
						data.IpaDomain.Locations[j].Model.CreatedAt,
						data.IpaDomain.Locations[j].Model.UpdatedAt,
						data.IpaDomain.Locations[j].Model.DeletedAt,

						data.IpaDomain.Locations[j].IpaID,
						data.IpaDomain.Locations[j].Name,
						data.IpaDomain.Locations[j].Description,
					)
				}
				expectedQuery.WillReturnRows(rows)
			}
		case 4:
			expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_servers" WHERE "ipa_servers"."ipa_id" = $1 AND "ipa_servers"."deleted_at" IS NULL`)).
				WithArgs(domainID)
			if i == stage && expectedErr != nil {
				expectedQuery.WillReturnError(expectedErr)
			} else {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deletet_at",

					"ipa_id", "fqdn", "rhsm_id", "location",
					"ca_server", "hcc_enrollment_server", "hcc_update_server",
					"pk_init_server",
				})
				for j := range data.IpaDomain.Servers {
					rows.AddRow(
						domainID+uint(j)+1,
						data.IpaDomain.Servers[j].Model.CreatedAt,
						data.IpaDomain.Servers[j].Model.UpdatedAt,
						data.IpaDomain.Servers[j].Model.DeletedAt,

						data.IpaDomain.Servers[j].IpaID,
						data.IpaDomain.Servers[j].FQDN,
						data.IpaDomain.Servers[j].RHSMId,
						data.IpaDomain.Servers[j].Location,
						data.IpaDomain.Servers[j].CaServer,
						data.IpaDomain.Servers[j].HCCEnrollmentServer,
						data.IpaDomain.Servers[j].HCCUpdateServer,
						data.IpaDomain.Servers[j].PKInitServer,
					)
				}
				expectedQuery.WillReturnRows(rows)
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *DomainRepositorySuite) TestDeleteByIdLogError() {
	t := s.T()
	r := &domainRepository{}

	d := builder_model.NewDomain(builder_model.NewModel().WithID(1).Build()).Build()

	s.helperTestDeleteById(1, d, s.mock, gorm.ErrInvalidTransaction)
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
