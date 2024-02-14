package impl

import (
	"encoding/json"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk"
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

type hostconfJwkComponent struct {
	interactor interactor.HostconfJwkInteractor
	repository repository.HostconfJwkRepository
	presenter  presenter.HostconfJwkPresenter
}

// Application secrets
type hostConfKeys struct {
	// TODO: store JWKs in database
	signingKeys []jwk.Key
	publicKeys  []string
}

type application struct {
	config      *config.Config
	jwks        *hostConfKeys
	metrics     *metrics.Metrics
	domain      domainComponent
	host        hostComponent
	hostconfjwk hostconfJwkComponent
	db          *gorm.DB
	inventory   client.HostInventory
}

func NewHandler(config *config.Config, db *gorm.DB, m *metrics.Metrics, inventory client.HostInventory) handler.Application {
	var (
		err error
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
	hcjwkc := hostconfJwkComponent{
		usecase_interactor.NewHostconfJwkInteractor(),
		usecase_repository.NewHostconfJwkRepository(config),
		usecase_presenter.NewHostconfJwkPresenter(config),
	}

	jwks, err := getJwks()
	if err != nil {
		panic(err)
	}
	// Instantiate application
	return &application{
		config:      config,
		jwks:        jwks,
		db:          db,
		metrics:     m,
		domain:      dc,
		host:        hc,
		hostconfjwk: hcjwkc,
		inventory:   inventory,
	}
}

// Generate ephemeral JWKs
func getJwks() (k *hostConfKeys, err error) {
	// TODO: temporary hack
	var (
		priv jwk.Key
		pub  jwk.Key
		pubs []byte
	)
	k = &hostConfKeys{}
	expiration := time.Now().Add(90 * 24 * time.Hour)
	if priv, err = hostconf_jwk.GeneratePrivateJWK(expiration); err != nil {
		return nil, err
	}
	k.signingKeys = []jwk.Key{priv}

	if pub, err = priv.PublicKey(); err != nil {
		return nil, err
	}
	if pubs, err = json.Marshal(pub); err != nil {
		return nil, err
	}
	k.publicKeys = []string{string(pubs)}

	return k, nil

}
