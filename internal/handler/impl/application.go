package impl

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

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

// Application secrets
type appSecrets struct {
	domainRegKey []byte
}

type application struct {
	config    *config.Config
	secrets   appSecrets
	metrics   *metrics.Metrics
	domain    domainComponent
	host      hostComponent
	db        *gorm.DB
	inventory client.HostInventory
}

func NewHandler(config *config.Config, db *gorm.DB, m *metrics.Metrics, inventory client.HostInventory) handler.Application {
	var (
		err          error
		domainRegKey []byte
	)
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
	// TODO move unmarshal and verification to Viper?
	if domainRegKey, err = getSecretBytes(
		"app.domain_reg_key", config.Application.DomainRegTokenKey, 8,
	); err != nil {
		panic(err)
	}

	sec := appSecrets{
		domainRegKey: domainRegKey,
	}

	// Instantiate application
	return &application{
		config:    config,
		secrets:   sec,
		db:        db,
		metrics:   m,
		domain:    dc,
		host:      hc,
		inventory: inventory,
	}
}

// Convert and check secret (raw standard base64 string)
func getSecretBytes(name string, value string, minLength int) (data []byte, err error) {
	// ephemeral random key for testing and development
	if value == "random" {
		data = make([]byte, minLength)
		if _, err = rand.Read(data); err != nil {
			return nil, err
		}
		return data, nil
	}
	if data, err = base64.RawStdEncoding.DecodeString(value); err != nil {
		return nil, fmt.Errorf("Failed to decode std base64 secret '%s': %v", name, err)
	}
	if len(data) < minLength {
		return nil, fmt.Errorf("Secrets '%s' is too short, expected %d bytes.", name, minLength)
	}
	return data, nil
}
