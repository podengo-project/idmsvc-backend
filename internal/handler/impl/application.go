package impl

import (
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client"
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

type application struct {
	config    *config.Config
	metrics   *metrics.Metrics
	domain    domainComponent
	host      hostComponent
	db        *gorm.DB
	inventory client.HostInventory
}

func NewHandler(config *config.Config, db *gorm.DB, m *metrics.Metrics, inventory client.HostInventory) handler.Application {
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

	// Instantiate application
	return &application{
		config:    config,
		db:        db,
		metrics:   m,
		domain:    dc,
		host:      hc,
		inventory: inventory,
	}
}
