package perf

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	repository_impl "github.com/podengo-project/idmsvc-backend/internal/usecase/repository"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

type PrepData struct {
	cfg     *config.Config
	db      *gorm.DB
	logger  *slog.Logger
	context context.Context
}

func NewPrepData() *PrepData {
	cfg := config.Get()
	logger := slog.Default()
	db := datastore.NewDB(cfg)
	ctx := context.Background()
	ctx = app_context.CtxWithLog(ctx, logger)
	ctx = app_context.CtxWithDB(ctx, db)
	prepData := &PrepData{
		cfg:     cfg,
		db:      datastore.NewDB(cfg),
		logger:  slog.Default(),
		context: ctx,
	}

	return prepData
}

// Create_domains creates 1 enabled domain for each organizations.
//
// count: number of organizations
func CreateDomains(prepData PrepData, count, orgIDBase int) ([]*model.Domain, error) {
	domains := make([]*model.Domain, count)
	for i := 0; i < count; i++ {
		orgNumber := orgIDBase + i
		domain, err := createDomain(prepData.context, orgNumber)
		domains[i] = domain
		if err != nil {
			return domains, err
		}
	}
	return domains, nil
}

func DeleteDomains(prepData PrepData, domains []*model.Domain) error {
	domainRepository := repository_impl.NewDomainRepository()
	for _, domain := range domains {
		if domain == nil {
			continue
		}
		err := domainRepository.DeleteById(prepData.context, domain.OrgId, domain.DomainUuid)
		if err != nil {
			continue
		}
	}
	return nil
}

func createIPADomainData(orgNumber int) *model.Domain {
	orgID := strconv.Itoa(orgNumber)
	domainName := orgID + ".example.test"

	currentTime := time.Now()
	testIpaCert := model.IpaCert{
		IpaID:        uint(orgNumber),
		Nickname:     "IPA.TEST IPA CA",
		Issuer:       "CN=Certificate Authority,O=IPA.TEST",
		Subject:      "CN=Certificate Authority,O=IPA.TEST",
		SerialNumber: "1",
		NotBefore:    currentTime,
		NotAfter:     currentTime,
		Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
	}
	testIpaServer := model.IpaServer{
		IpaID:               uint(orgNumber),
		FQDN:                "server1.ipa.test",
		RHSMId:              pointy.String(uuid.NewString()),
		Location:            pointy.String("europe"),
		CaServer:            true,
		HCCEnrollmentServer: true,
		HCCUpdateServer:     true,
		PKInitServer:        true,
	}

	domain := &model.Domain{
		OrgId:                 orgID,
		DomainUuid:            uuid.New(),
		DomainName:            pointy.String(domainName),
		Title:                 pointy.String(domainName),
		Description:           pointy.String(""),
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain: &model.Ipa{
			RealmName:    pointy.String(""),
			CaCerts:      []model.IpaCert{testIpaCert},
			Servers:      []model.IpaServer{testIpaServer},
			Locations:    []model.IpaLocation{},
			RealmDomains: pq.StringArray{domainName},
		},
	}

	return domain
}

// createDomain directly creates a domain record in the database.
func createDomain(
	ctx context.Context,
	orgNumber int,
) (*model.Domain, error) {
	domainRepository := repository_impl.NewDomainRepository()
	domain := createIPADomainData(orgNumber)

	// Delete existing domains to ensure only one domain per organization
	// E.g. when re-running the test after some cleanup issues.
	domains, _, err := domainRepository.List(ctx, domain.OrgId, 0, 10)
	if err != nil {
		return nil, err
	}
	for _, d := range domains {
		err = domainRepository.DeleteById(ctx, domain.OrgId, d.DomainUuid)
		if err != nil {
			return nil, err
		}
	}

	err = domainRepository.Register(ctx, domain.OrgId, domain)
	if err != nil {
		return domain, err
	}
	return domain, nil
}
