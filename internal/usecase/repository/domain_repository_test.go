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
	"github.com/hmsidm/internal/domain/model"
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
	repository *domainRepository
}

func (s *Suite) SetupSuite() {
	var err error
	s.mock, s.DB, err = test.NewSqlMock(&gorm.Session{SkipHooks: true})
	if err != nil {
		s.Suite.FailNow("Error calling gorm.Open: %s", err.Error())
		return
	}
	s.repository = &domainRepository{}
}

func (s *Suite) TearDownSuite() {
}

func (s *Suite) TestNewDomainRepository() {
	t := s.Suite.T()
	assert.NotPanics(t, func() {
		_ = NewDomainRepository()
	})
}

func (s *Suite) TestCreate() {
	orgID := "12345"
	token := uuid.NewString()
	tokenExpiration := &time.Time{}
	*tokenExpiration = time.Now().Add(model.DefaultTokenExpiration())
	testUUID := uuid.New()
	t := s.Suite.T()
	currentTime := time.Now()
	var data model.Domain = model.Domain{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		OrgId:                 orgID,
		DomainUuid:            testUUID,
		DomainName:            nil,
		Title:                 pointy.String("My Domain Example Title"),
		Description:           pointy.String("My Domain Example Description"),
		Type:                  pointy.Uint(model.DomainTypeIpa),
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

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "domains" ("created_at","updated_at","deleted_at","org_id","domain_uuid","domain_name","title","description","type","auto_enrollment_enabled","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(
			data.CreatedAt,
			data.UpdatedAt,
			nil,

			orgID,
			data.DomainUuid,
			data.DomainName,
			data.Title,
			data.Description,
			data.Type,
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

	err := s.repository.Create(s.DB, orgID, &data)
	require.NoError(t, err)
}

func (s *Suite) TestCreateErrors() {
	orgID := "12345"
	testUUID := uuid.New()
	t := s.Suite.T()
	currentTime := time.Now()
	var (
		data model.Domain = model.Domain{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			OrgId:                 orgID,
			DomainUuid:            testUUID,
			DomainName:            nil,
			Title:                 pointy.String("My domain test title"),
			Description:           pointy.String("My domain test description"),
			Type:                  pointy.Uint(model.DomainTypeIpa),
			AutoEnrollmentEnabled: pointy.Bool(true),
			IpaDomain: &model.Ipa{
				CaCerts: []model.IpaCert{},
				Servers: []model.IpaServer{},
			},
		}
		domainTypeIsNil model.Domain = model.Domain{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			OrgId:                 orgID,
			DomainUuid:            testUUID,
			DomainName:            nil,
			Title:                 pointy.String("My domain test title"),
			Description:           pointy.String("My domain test description"),
			Type:                  nil,
			AutoEnrollmentEnabled: pointy.Bool(true),
			IpaDomain:             nil,
		}
		ipaDomainTypeIsNotValid model.Domain = model.Domain{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			OrgId:                 orgID,
			DomainUuid:            testUUID,
			DomainName:            nil,
			Title:                 pointy.String("My domain test title"),
			Description:           pointy.String("My domain test description"),
			Type:                  pointy.Uint(1000),
			AutoEnrollmentEnabled: pointy.Bool(true),
			IpaDomain:             nil,
		}
		err error
	)

	err = s.repository.Create(nil, "", nil)
	assert.Error(t, err)

	err = s.repository.Create(s.DB, "", nil)
	assert.Error(t, err)

	err = s.repository.Create(s.DB, orgID, nil)
	assert.Error(t, err)

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "domains" ("created_at","updated_at","deleted_at","org_id","domain_uuid","domain_name","title","description","type","auto_enrollment_enabled","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(
			data.CreatedAt,
			data.UpdatedAt,
			nil,
			orgID,
			data.DomainUuid,
			nil,
			data.Title,
			data.Description,
			data.Type,
			data.AutoEnrollmentEnabled,
			data.ID,
		).
		WillReturnError(fmt.Errorf("an error happened"))
	err = s.repository.Create(s.DB, orgID, &data)
	assert.Error(t, err)
	assert.Equal(t, "an error happened", err.Error())

	err = s.repository.Create(s.DB, orgID, &domainTypeIsNil)
	assert.EqualError(t, err, "'Type' is nil")

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "domains" ("created_at","updated_at","deleted_at","org_id","domain_uuid","domain_name","title","description","type","auto_enrollment_enabled","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(
			ipaDomainTypeIsNotValid.CreatedAt,
			ipaDomainTypeIsNotValid.UpdatedAt,
			nil,
			ipaDomainTypeIsNotValid.OrgId,
			ipaDomainTypeIsNotValid.DomainUuid,
			ipaDomainTypeIsNotValid.DomainName,
			ipaDomainTypeIsNotValid.Title,
			ipaDomainTypeIsNotValid.Description,
			ipaDomainTypeIsNotValid.Type,
			ipaDomainTypeIsNotValid.AutoEnrollmentEnabled,
			ipaDomainTypeIsNotValid.ID,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(uint(1)))
	err = s.repository.Create(s.DB, orgID, &ipaDomainTypeIsNotValid)
	assert.EqualError(t, err, "'Type' is invalid")
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
			RealmName:       pointy.String("MYDOMAIN.EXAMPLE"),
			Token:           nil,
			TokenExpiration: nil,
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
			RealmDomains: []string{"mydomain.example"},
		}
	)

	// Check nil
	err = s.repository.createIpaDomain(s.DB, 1, nil)
	assert.EqualError(t, err, "'data' of type '*model.Ipa' is nil")

	//
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_domains","token","token_expiration","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
		WithArgs(
			data.Model.CreatedAt,
			data.Model.UpdatedAt,
			nil,

			data.RealmName,
			data.RealmDomains,
			data.Token,
			sqlmock.AnyArg(),
			data.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(data.ID))

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_certs" ("created_at","updated_at","deleted_at","ipa_id","issuer","nickname","not_valid_after","not_valid_before","pem","serial_number","subject","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).
		WithArgs(
			data.CaCerts[0].CreatedAt,
			data.CaCerts[0].UpdatedAt,
			nil,

			data.CaCerts[0].IpaID,
			data.CaCerts[0].Issuer,
			data.CaCerts[0].Nickname,
			data.CaCerts[0].NotValidAfter,
			data.CaCerts[0].NotValidBefore,
			data.CaCerts[0].Pem,
			data.CaCerts[0].SerialNumber,
			data.CaCerts[0].Subject,
			data.CaCerts[0].ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(data.CaCerts[0].ID))

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_servers" ("created_at","updated_at","deleted_at","ipa_id","fqdn","rhsm_id","ca_server","hcc_enrollment_server","hcc_update_server","pk_init_server","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(
			data.Servers[0].CreatedAt,
			data.Servers[0].UpdatedAt,
			nil,

			data.Servers[0].IpaID,
			data.Servers[0].FQDN,
			data.Servers[0].RHSMId,
			data.Servers[0].CaServer,
			data.Servers[0].HCCEnrollmentServer,
			data.Servers[0].HCCUpdateServer,
			data.Servers[0].PKInitServer,
			data.Servers[0].ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(data.Servers[0].ID))

	err = s.repository.createIpaDomain(s.DB, 1, &data)
	assert.NoError(t, err)
}

func (s *Suite) TestUpdateErrors() {
	t := s.Suite.T()
	orgID := "11111"
	testUUID := uuid.New()
	currentTime := time.Now()
	var (
		data model.Domain = model.Domain{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			OrgId:                 orgID,
			DomainUuid:            testUUID,
			DomainName:            pointy.String("domain.example"),
			Title:                 pointy.String("My domain test title"),
			Description:           pointy.String("My domain test description"),
			Type:                  pointy.Uint(model.DomainTypeIpa),
			AutoEnrollmentEnabled: pointy.Bool(true),
			IpaDomain: &model.Ipa{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				},
				CaCerts: []model.IpaCert{},
				Servers: []model.IpaServer{},
			},
		}
		err error
	)

	err = s.repository.Update(nil, "", nil)
	assert.EqualError(t, err, "'db' is nil")

	err = s.repository.Update(s.DB, "", nil)
	assert.EqualError(t, err, "'orgId' is empty")

	err = s.repository.Update(s.DB, orgID, nil)
	assert.EqualError(t, err, "'data' is nil")

	s.mock.MatchExpectationsInOrder(true)
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL ORDER BY "domains"."id" LIMIT 1`)).
		WithArgs(
			orgID,
			testUUID,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"org_id", "domain_uuid", "domain_name",
			"title", "description", "type", "auto_enrollment_enabled",
		}).AddRow(
			data.Model.ID,
			data.Model.CreatedAt,
			data.Model.UpdatedAt,
			nil,

			data.OrgId,
			testUUID,
			*data.DomainName,
			*data.Title,
			*data.Description,
			*data.Type,
			*data.AutoEnrollmentEnabled,
		))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipas" WHERE "ipas"."deleted_at" IS NULL AND "ipas"."id" = $1 ORDER BY "ipas"."id" LIMIT 1`)).
		WithArgs(
			data.Model.ID,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at",
				"realm_name", "realm_names", "token", "token_expiration",
			}))
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "domains" SET "created_at"=$1,"updated_at"=$2,"org_id"=$3,"domain_uuid"=$4,"domain_name"=$5,"title"=$6,"description"=$7,"type"=$8,"auto_enrollment_enabled"=$9 WHERE org_id = $10 AND "domains"."deleted_at" IS NULL AND "id" = $11`)).
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
			data.ID,
		).
		WillReturnError(fmt.Errorf("An error"))
	err = s.repository.Update(s.DB, orgID, &data)
	assert.EqualError(t, err, "An error")

	s.mock.MatchExpectationsInOrder(true)
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL ORDER BY "domains"."id" LIMIT 1`)).
		WithArgs(
			orgID,
			testUUID,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"org_id", "domain_uuid", "domain_name",
			"title", "description", "type", "auto_enrollment_enabled",
		}).AddRow(
			data.Model.ID,
			data.Model.CreatedAt,
			data.Model.UpdatedAt,
			nil,

			data.OrgId,
			testUUID,
			*data.DomainName,
			*data.Title,
			*data.Description,
			*data.Type,
			*data.AutoEnrollmentEnabled,
		))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipas" WHERE "ipas"."deleted_at" IS NULL AND "ipas"."id" = $1 ORDER BY "ipas"."id" LIMIT 1`)).
		WithArgs(
			data.Model.ID,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at",
				"realm_name", "realm_names", "token", "token_expiration",
			}))
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "domains" SET "created_at"=$1,"updated_at"=$2,"org_id"=$3,"domain_uuid"=$4,"domain_name"=$5,"title"=$6,"description"=$7,"type"=$8,"auto_enrollment_enabled"=$9 WHERE org_id = $10 AND "domains"."deleted_at" IS NULL AND "id" = $11`)).
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
			data.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ipas" WHERE id = $1 AND "ipas"."id" = $2`)).
		WithArgs(
			data.ID,
			data.ID,
		).WillReturnResult(
		driver.RowsAffected(1),
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_domains","token","token_expiration","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,

			data.IpaDomain.RealmName,
			data.IpaDomain.RealmDomains,
			data.IpaDomain.Token,
			data.IpaDomain.TokenExpiration,
			data.IpaDomain.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow("1"),
		)
	err = s.repository.Update(s.DB, orgID, &data)
	assert.NoError(t, err)
}

func (s *Suite) TestRhelIdmClearToken() {
	var (
		err error
	)
	t := s.Suite.T()
	currentTime := time.Now()
	mismatchOrgID := "22222"
	orgID := "11111"
	uuidString := "3bccb88e-dd25-11ed-99e0-482ae3863d30"
	subscriptionManagerID := "fe106208-dd32-11ed-aa87-482ae3863d30"
	data := model.Domain{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: gorm.DeletedAt{},
		},
		OrgId:                 orgID,
		DomainUuid:            uuid.MustParse(uuidString),
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
					IpaID:          1,
					Issuer:         "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
					Nickname:       "MYDOMAIN.EXAMPLE IPA CA",
					NotValidAfter:  currentTime.Add(24 * time.Hour),
					NotValidBefore: currentTime,
					SerialNumber:   "1",
					Subject:        "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
					Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----",
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
					RHSMId:              subscriptionManagerID,
					CaServer:            true,
					HCCEnrollmentServer: true,
					HCCUpdateServer:     true,
					PKInitServer:        true,
				},
			},
			RealmDomains: pq.StringArray{"mydomain.example"},
		},
	}

	err = s.repository.RhelIdmClearToken(nil, "", "")
	assert.EqualError(t, err, "'db' is nil")

	err = s.repository.RhelIdmClearToken(s.DB, "", "")
	assert.EqualError(t, err, "'orgId' is empty")

	err = s.repository.RhelIdmClearToken(s.DB, orgID, "")
	assert.EqualError(t, err, "'uuid' is empty")

	s.mock.MatchExpectationsInOrder(true)
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL ORDER BY "domains"."id" LIMIT 1`)).
		WithArgs(
			orgID,
			uuidString,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at",
				"org_id", "domain_uuid", "domain_name", "title",
				"description", "type", "auto_enrollment_enabled",
			}).
				AddRow(
					1,
					data.CreatedAt,
					data.UpdatedAt,
					nil,

					mismatchOrgID,
					data.DomainUuid,
					*data.DomainName,
					*data.Title,
					*data.Description,
					*data.Type,
					*data.AutoEnrollmentEnabled,
				))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipas" WHERE "ipas"."deleted_at" IS NULL AND "ipas"."id" = $1 ORDER BY "ipas"."id" LIMIT 1`)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"realm_name", "realm_names", "token", "token_expiration",
		}).
			AddRow(
				1,
				currentTime,
				currentTime,
				nil,

				data.IpaDomain.RealmName,
				data.IpaDomain.RealmDomains,
				data.IpaDomain.Token,
				data.IpaDomain.TokenExpiration,
			))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_certs" WHERE "ipa_certs"."id" = $1 AND "ipa_certs"."deleted_at" IS NULL`)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"realm_name", "realm_names", "token", "token_expiration",
		}).
			AddRow(
				1,
				currentTime,
				currentTime,
				nil,

				data.IpaDomain.RealmName,
				data.IpaDomain.RealmDomains,
				data.IpaDomain.Token,
				data.IpaDomain.TokenExpiration,
			))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_servers" WHERE "ipa_servers"."id" = $1 AND "ipa_servers"."deleted_at" IS NULL`)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"ipa_id", "fqdn", "rhsm_id",
			"ca_server", "hcc_enrollment_server", "hcc_update_server",
			"pk_init_server",
		}))
	err = s.repository.RhelIdmClearToken(s.DB, orgID, uuidString)
	assert.EqualError(t, err, "'OrgId' mistmatch")

	// Success scenario
	s.mock.MatchExpectationsInOrder(true)
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL ORDER BY "domains"."id" LIMIT 1`)).
		WithArgs(
			orgID,
			uuidString,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at",
				"org_id", "domain_uuid", "domain_name", "title",
				"description", "type", "auto_enrollment_enabled",
			}).
				AddRow(
					1,
					data.CreatedAt,
					data.UpdatedAt,
					nil,

					orgID,
					data.DomainUuid,
					*data.DomainName,
					*data.Title,
					*data.Description,
					*data.Type,
					*data.AutoEnrollmentEnabled,
				))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipas" WHERE "ipas"."deleted_at" IS NULL AND "ipas"."id" = $1 ORDER BY "ipas"."id" LIMIT 1`)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"realm_name", "realm_names", "token", "token_expiration",
		}).
			AddRow(
				1,
				currentTime,
				currentTime,
				nil,

				data.IpaDomain.RealmName,
				data.IpaDomain.RealmDomains,
				data.IpaDomain.Token,
				data.IpaDomain.TokenExpiration,
			))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_certs" WHERE "ipa_certs"."id" = $1 AND "ipa_certs"."deleted_at" IS NULL`)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"realm_name", "realm_names", "token", "token_expiration",
		}).
			AddRow(
				1,
				currentTime,
				currentTime,
				nil,

				data.IpaDomain.RealmName,
				data.IpaDomain.RealmDomains,
				data.IpaDomain.Token,
				data.IpaDomain.TokenExpiration,
			))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_servers" WHERE "ipa_servers"."id" = $1 AND "ipa_servers"."deleted_at" IS NULL`)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"ipa_id", "fqdn", "rhsm_id",
			"ca_server", "hcc_enrollment_server", "hcc_update_server",
			"pk_init_server",
		}))
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ipas" SET "token"=$1 WHERE id = $2`)).
		WithArgs(
			nil,
			data.IpaDomain.ID,
		).
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ipas" SET "token_expiration"=$1 WHERE id = $2`)).
		WithArgs(
			nil,
			data.IpaDomain.ID,
		).
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	err = s.repository.RhelIdmClearToken(s.DB, orgID, uuidString)
	assert.NoError(t, err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
