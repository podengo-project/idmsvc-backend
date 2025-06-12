package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	builder_model "github.com/podengo-project/idmsvc-backend/internal/test/builder/model"
	test_sql "github.com/podengo-project/idmsvc-backend/internal/test/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.openly.dev/pointy"
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

	// context is nil
	assert.Panics(t, func() {
		_, _ = s.repository.MatchDomain(nil, &interactor.HostConfOptions{})
	})

	// db is nil
	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_, _ = s.repository.MatchDomain(context.Background(), &interactor.HostConfOptions{})
	})

	// options is nil
	domain, err := s.repository.MatchDomain(s.Ctx, nil)
	assert.Nil(t, domain)
	require.EqualError(t, err, "code=500, message='options' cannot be nil")

	// Error at Find
	test_sql.MatchDomain(1, s.mock, gorm.ErrInvalidTransaction, options, domains)
	domain, err = s.repository.MatchDomain(s.Ctx, options)
	assert.Nil(t, domain)
	require.EqualError(t, err, "invalid transaction")

	// Domains empty
	domainsEmpty := []model.Domain{}
	test_sql.MatchDomain(1, s.mock, nil, options, domainsEmpty)
	domain, err = s.repository.MatchDomain(s.Ctx, options)
	assert.Nil(t, domain)
	require.EqualError(t, err, "code=404, message=no matching domains")

	// More than 1 match
	domainsMoreThan1 := []model.Domain{
		domains[0],
		domains[0],
	}
	test_sql.MatchDomain(1, s.mock, nil, options, domainsMoreThan1)
	domain, err = s.repository.MatchDomain(s.Ctx, options)
	assert.Nil(t, domain)
	require.EqualError(t, err, "code=409, message=matched 2 domains, only one expected")

	// Success
	test_sql.MatchDomain(2, s.mock, nil, options, domains)
	domain, err = s.repository.MatchDomain(s.Ctx, options)
	assert.NotNil(t, domain)
	require.NoError(t, err)
}

func (s *SuiteHost) TestSignHostConfToken() {
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
	domain := *builder_model.NewDomain(builder_model.NewModel().WithID(id).Build()).
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
		).Build()

	// guard ctx is nil
	require.PanicsWithValue(t, "'ctx' is nil", func() {
		_, _ = s.repository.SignHostConfToken(nil, nil, nil, nil)
	})

	// guard options is nil
	ctx := app_context.CtxWithLog(context.Background(), slog.Default())
	token, err := s.repository.SignHostConfToken(ctx, nil, nil, nil)
	assert.Equal(t, "", token)
	require.EqualError(t, err, "code=500, message='options' cannot be nil")

	// guard domain is nil
	token, err = s.repository.SignHostConfToken(ctx, nil, options, nil)
	assert.Equal(t, "", token)
	require.EqualError(t, err, "code=500, message='domain' cannot be nil")

	// no signers available
	token, err = s.repository.SignHostConfToken(ctx, nil, options, &domain)
	assert.Equal(t, "", token)
	require.EqualError(t, err, "jws.Sign: no signers available. Specify an alogirthm and akey using jws.WithKey()")

	// no signers
	token, err = s.repository.SignHostConfToken(ctx, nil, options, &domain)
	assert.Equal(t, "", token)
	require.EqualError(t, err, "jws.Sign: no signers available. Specify an alogirthm and akey using jws.WithKey()")
}

func TestSuiteHost(t *testing.T) {
	suite.Run(t, new(SuiteHost))
}
