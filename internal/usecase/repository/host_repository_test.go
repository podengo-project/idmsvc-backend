package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	builder_model "github.com/podengo-project/idmsvc-backend/internal/test/builder/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SuiteHost struct {
	SuiteBase
	repository *hostRepository
}

// https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
func (s *SuiteHost) SetupTest() {
	s.SuiteBase.SetupTest()
	s.repository = &hostRepository{}
}

func (s *SuiteHost) TestNewHostRepository() {
	t := s.Suite.T()
	assert.NotPanics(t, func() {
		_ = NewHostRepository()
	})
}

func (s *SuiteHost) helperTestMatchDomain(stage int, options *interactor.HostConfOptions, domains []model.Domain, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT "domains"."id","domains"."created_at","domains"."updated_at","domains"."deleted_at","domains"."org_id","domains"."domain_uuid","domains"."domain_name","domains"."title","domains"."description","domains"."type","domains"."auto_enrollment_enabled" FROM "domains" left join ipas on domains.id = ipas.id WHERE domains.org_id = $1 AND domains.domain_uuid = $2 AND domains.domain_name = $3 AND domains.type = $4 AND "domains"."deleted_at" IS NULL`)).
				WithArgs(
					options.OrgId,
					options.DomainId,
					options.DomainName,
					model.DomainTypeUint((string)(*options.DomainType)),
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deletet_at",

					"org_id", "domain_uuid", "domain_name",
					"title", "description", "type",
					"auto_enrollment_enabled",
				})
				for j := range domains {
					rows.AddRow(
						domains[j].ID,
						domains[j].CreatedAt,
						domains[j].UpdatedAt,
						domains[j].DeletedAt,

						domains[j].OrgId,
						domains[j].DomainUuid,
						domains[j].DomainName,
						domains[j].Title,
						domains[j].Description,
						domains[j].Type,
						domains[j].AutoEnrollmentEnabled,
					)
				}
				expectQuery = expectQuery.WillReturnRows(rows)
			}
		case 2:
			if len(domains) == 0 {
				helperTestFindByIDIpa(1, &domains[0], mock, expectedErr)
			}
			helperTestFindByIDIpa(4, &domains[0], mock, expectedErr)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *SuiteHost) TestMatchDomain() {
	t := s.Suite.T()
	domainName := helper.GenRandDomainName(2)
	fqdn := helper.GenRandFQDNWithDomain(domainName)
	domainId := uuid.New()
	options := &interactor.HostConfOptions{
		OrgId:       "12345",
		CommonName:  domainId,
		InventoryId: uuid.MustParse("9db10f12-c421-11ee-8c1c-482ae3863d30"),
		Fqdn:        fqdn,
		DomainId:    &domainId,
		DomainName:  &domainName,
		DomainType:  (*api_public.DomainType)(pointy.String(model.DomainTypeIpaString)),
	}
	id := uint(helper.GenRandNum(0, 2^63))
	realm := strings.ToUpper(domainName)
	domains := []model.Domain{
		*builder_model.NewDomain(builder_model.NewModel().WithID(id).Build()).
			WithOrgID(options.OrgId).
			WithAutoEnrollmentEnabled(pointy.Bool(true)).
			WithIpaDomain(
				builder_model.NewIpaDomain().
					WithModel(builder_model.NewModel().WithID(id).Build()).
					WithRealmName(&realm).
					WithRealmDomains(pq.StringArray{domainName}).
					WithServers([]model.IpaServer{
						builder_model.NewIpaServer(
							builder_model.NewModel().Build(),
						).WithIpaID(id).
							WithFQDN(fqdn).
							Build(),
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
			).Build(),
		*builder_model.NewDomain(builder_model.NewModel().Build()).
			WithOrgID(options.OrgId).
			// Avoid flaky tests by disabling the second domain
			WithAutoEnrollmentEnabled(pointy.Bool(false)).
			Build(),
	}

	// db is nil
	domain, err := s.repository.MatchDomain(nil, &interactor.HostConfOptions{})
	assert.Nil(t, domain)
	require.EqualError(t, err, "code=500, message='db' cannot be nil")

	// options is nil
	domain, err = s.repository.MatchDomain(s.DB, nil)
	assert.Nil(t, domain)
	require.EqualError(t, err, "code=500, message='options' cannot be nil")

	// Error at Find
	s.helperTestMatchDomain(1, options, domains, s.mock, gorm.ErrInvalidTransaction)
	domain, err = s.repository.MatchDomain(s.DB, options)
	assert.Nil(t, domain)
	require.EqualError(t, err, "invalid transaction")

	// Domains empty
	domainsEmpty := []model.Domain{}
	s.helperTestMatchDomain(1, options, domainsEmpty, s.mock, nil)
	domain, err = s.repository.MatchDomain(s.DB, options)
	assert.Nil(t, domain)
	require.EqualError(t, err, "code=404, message=no matching domains")

	// More than 1 match
	domainsMoreThan1 := []model.Domain{
		domains[0],
		domains[0],
	}
	s.helperTestMatchDomain(1, options, domainsMoreThan1, s.mock, nil)
	domain, err = s.repository.MatchDomain(s.DB, options)
	assert.Nil(t, domain)
	require.EqualError(t, err, "code=409, message=matched 2 domains, only one expected")

	// Success
	s.helperTestMatchDomain(2, options, domains, s.mock, nil)
	domain, err = s.repository.MatchDomain(s.DB, options)
	assert.NotNil(t, domain)
	require.NoError(t, err)
}

func TestSuiteHost(t *testing.T) {
	suite.Run(t, new(SuiteHost))
}
