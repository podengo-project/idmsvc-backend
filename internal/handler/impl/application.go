package impl

import (
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
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

const (
	errXRHIDIsNil     = "failed because XRHID is nil"
	errUnserializing  = "failed to unserialize http api data"
	errInputAdapter   = "failed to translate the API request to business objects"
	errDBTXBegin      = "failed to begin database transaction"
	errDBNotFound     = "failed because a record not found in the database"
	errDBGeneralError = "failed on database operation"
	errDBTXCommit     = "failed to commit database transaction"
	errOutputAdapter  = "failed to translate the business object to the API response"
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
	pendo       client_pendo.Pendo
}

func guardNewHandler(cfg *config.Config, db *gorm.DB, m *metrics.Metrics, rbac client_rbac.Rbac, pendo client_pendo.Pendo) {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if db == nil {
		panic("'db' is nil")
	}
	if m == nil {
		panic("'m' is nil")
	}
	if rbac == nil {
		panic("'rbac' is nil")
	}
	if pendo == nil {
		panic("'pendo' is nil")
	}
}

func NewHandler(cfg *config.Config, db *gorm.DB, m *metrics.Metrics, rbac client_rbac.Rbac, pendo client_pendo.Pendo) handler.Application {
	dc := domainComponent{
		usecase_interactor.NewDomainInteractor(),
		usecase_repository.NewDomainRepository(),
		usecase_presenter.NewDomainPresenter(cfg),
	}
	hc := hostComponent{
		usecase_interactor.NewHostInteractor(),
		usecase_repository.NewHostRepository(),
		usecase_presenter.NewHostPresenter(cfg),
	}
	hcjc := hostconfJwkComponent{
		usecase_interactor.NewHostconfJwkInteractor(),
		usecase_repository.NewHostconfJwkRepository(cfg),
		usecase_presenter.NewHostconfJwkPresenter(cfg),
	}

	// Instantiate application
	return &application{
		config:      cfg,
		db:          db,
		metrics:     m,
		domain:      dc,
		host:        hc,
		hostconfjwk: hcjc,
		pendo:       pendo,
	}
}
