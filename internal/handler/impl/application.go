package impl

import (
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	client_inventory "github.com/podengo-project/idmsvc-backend/internal/interface/client/inventory"
	client_pendo "github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
	client_rbac "github.com/podengo-project/idmsvc-backend/internal/interface/client/rbac"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/interface/presenter"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	metrics "github.com/podengo-project/idmsvc-backend/internal/metrics"
	usecase_interactor "github.com/podengo-project/idmsvc-backend/internal/usecase/interactor"
	usecase_presenter "github.com/podengo-project/idmsvc-backend/internal/usecase/presenter"
	usecase_repository "github.com/podengo-project/idmsvc-backend/internal/usecase/repository"
	"gorm.io/gorm"
)

type domainComponent struct {
	interactor interactor.DomainInteractor
	repository repository.DomainRepository
	presenter  presenter.DomainPresenter
}

type hostComponent struct {
	interactor interactor.HostInteractor
	repository repository.HostRepository
	presenter  presenter.HostPresenter
}

type hostconfJwkComponent struct {
	interactor interactor.HostconfJwkInteractor
	repository repository.HostconfJwkRepository
	presenter  presenter.HostconfJwkPresenter
}

type application struct {
	config      *config.Config
	metrics     *metrics.Metrics
	domain      domainComponent
	host        hostComponent
	hostconfjwk hostconfJwkComponent
	db          *gorm.DB
	inventory   client_inventory.HostInventory
	pendo       client_pendo.Pendo
}

func NewHandler(config *config.Config, db *gorm.DB, m *metrics.Metrics, inventory client_inventory.HostInventory, rbac client_rbac.Rbac, pendo client_pendo.Pendo) handler.Application {
	if config == nil {
		panic("config is nil")
	}
	if db == nil {
		panic("db is nil")
	}
	dc := domainComponent{
		usecase_interactor.NewDomainInteractor(),
		usecase_repository.NewDomainRepository(),
		usecase_presenter.NewDomainPresenter(config),
	}
	hc := hostComponent{
		usecase_interactor.NewHostInteractor(),
		usecase_repository.NewHostRepository(),
		usecase_presenter.NewHostPresenter(config),
	}
	hcjc := hostconfJwkComponent{
		usecase_interactor.NewHostconfJwkInteractor(),
		usecase_repository.NewHostconfJwkRepository(config),
		usecase_presenter.NewHostconfJwkPresenter(config),
	}

	// Instantiate application
	return &application{
		config:      config,
		db:          db,
		metrics:     m,
		domain:      dc,
		host:        hc,
		hostconfjwk: hcjc,
		inventory:   inventory,
		pendo:       pendo,
	}
}
