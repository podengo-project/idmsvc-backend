package impl

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/interface/presenter"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	metrics "github.com/podengo-project/idmsvc-backend/internal/metrics"
	usecase_interactor "github.com/podengo-project/idmsvc-backend/internal/usecase/interactor"
	usecase_presenter "github.com/podengo-project/idmsvc-backend/internal/usecase/presenter"
	usecase_repository "github.com/podengo-project/idmsvc-backend/internal/usecase/repository"
	"golang.org/x/crypto/hkdf"
	"gorm.io/gorm"
)

const (
	Salt             = "idmsvc-backend"
	DomainRegKeyInfo = "domain registration key"
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
	// TODO: store JWKs in database
	signingKeys []jwk.Key
	publicKeys  []string
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

	sec, err := getAppSecret(config)
	if err != nil {
		panic(err)
	}
	// Instantiate application
	return &application{
		config:    config,
		secrets:   *sec,
		db:        db,
		metrics:   m,
		domain:    dc,
		host:      hc,
		inventory: inventory,
	}
}

// Parse main secret and get sub secrets
func getAppSecret(config *config.Config) (sec *appSecrets, err error) {
	const mainSecretLength = 16
	// get / create main secret
	var secret []byte
	if config.Application.MainSecret == "random" {
		secret = make([]byte, mainSecretLength)
		if _, err = rand.Read(secret); err != nil {
			return nil, err
		}
	} else {
		if secret, err = base64.RawURLEncoding.DecodeString(config.Application.MainSecret); err != nil {
			return nil, fmt.Errorf("Failed to main secret: %v", err)
		}
		if len(secret) < mainSecretLength {
			return nil, fmt.Errorf("Master secret is too short, expected at least %d bytes.", mainSecretLength)
		}
	}

	// extract PRK from main secret
	var hash = sha256.New
	prk := hkdf.Extract(hash, secret, []byte(Salt))

	sec = &appSecrets{}
	sec.domainRegKey, err = hkdfExpand(hash, prk, []byte(DomainRegKeyInfo), 32)
	if err != nil {
		return nil, err
	}

	// TODO: temporary hack
	var (
		priv jwk.Key
		pub  jwk.Key
		pubs []byte
	)
	expiration := time.Now().Add(90 * 24 * time.Hour)
	if priv, err = token.GeneratePrivateJWK(expiration); err != nil {
		return nil, err
	}
	sec.signingKeys = []jwk.Key{priv}

	if pub, err = priv.PublicKey(); err != nil {
		return nil, err
	}
	if pubs, err = json.Marshal(pub); err != nil {
		return nil, err
	}
	sec.publicKeys = []string{string(pubs)}

	return sec, nil

}

// expand pseudo random key with HKDF
func hkdfExpand(hash func() hash.Hash, prk []byte, info []byte, length int) (secret []byte, err error) {
	reader := hkdf.Expand(hash, prk, info)
	secret = make([]byte, length)
	if _, err := io.ReadFull(reader, secret); err != nil {
		return nil, err
	}
	return secret, err
}
