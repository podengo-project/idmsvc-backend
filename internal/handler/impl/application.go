package impl

import (
	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/handler"
	"github.com/hmsidm/internal/interface/client"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/hmsidm/internal/interface/presenter"
	"github.com/hmsidm/internal/interface/repository"
	metrics "github.com/hmsidm/internal/metrics"
	usecase_interactor "github.com/hmsidm/internal/usecase/interactor"
	usecase_presenter "github.com/hmsidm/internal/usecase/presenter"
	usecase_repository "github.com/hmsidm/internal/usecase/repository"
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
