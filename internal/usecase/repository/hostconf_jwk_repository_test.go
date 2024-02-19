package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type HostConfJwkRepositorySuite struct {
	suite.Suite
	db         *gorm.DB
	mock       sqlmock.Sqlmock
	cfg        *config.Config
	repository repository.HostconfJwkRepository
}

// https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
func (s *HostConfJwkRepositorySuite) SetupTest() {
	var err error
	s.mock, s.db, err = test.NewSqlMock(&gorm.Session{
		SkipHooks: true,
	})
	if err != nil {
		s.Suite.FailNow("Error calling gorm.Open: %s", err.Error())
		return
	}
	s.mock.MatchExpectationsInOrder(true)
	s.cfg = test.GetTestConfig()
	s.repository = NewHostconfJwkRepository(s.cfg)
}

func (s *HostConfJwkRepositorySuite) newHostconfJWK() (hcjwk *model.HostconfJwk) {
	t := s.Suite.T()

	expiresAt := time.
		Now().
		Add(s.cfg.Application.HostconfJwkValidity).
		Truncate(time.Second)
	hcjwk, err := model.NewHostconfJwk(s.cfg.Secrets, expiresAt)
	require.Nil(t, err)
	return hcjwk
}

func (s *HostConfJwkRepositorySuite) TestNewHostconfJwkRepository() {
	t := s.Suite.T()
	assert.NotNil(t, s.repository)
	assert.NotNil(t, s.repository.(*hostconfJwkRepository).config)
	assert.Panics(t, func() {
		NewHostconfJwkRepository(nil)
	})
}

func (s *HostConfJwkRepositorySuite) TestInsertJWK() {
	t := s.Suite.T()

	hcjwk := s.newHostconfJWK()
	err := s.repository.InsertJWK(s.db, hcjwk)
	// TODO mock
	assert.Error(t, err)
}

func (s *HostConfJwkRepositorySuite) TestRevokeJWK() {
	t := s.Suite.T()

	hcjwk := s.newHostconfJWK()
	model, err := s.repository.RevokeJWK(s.db, hcjwk.KeyId)
	assert.Nil(t, model)
	assert.Error(t, err)
}

func (s *HostConfJwkRepositorySuite) TestListJWKs() {
	t := s.Suite.T()
	models, err := s.repository.ListJWKs(s.db)
	assert.Nil(t, models)
	assert.Error(t, err)
}

func (s *HostConfJwkRepositorySuite) TestGetPublicKeyArray() {
	t := s.Suite.T()
	keys, revokedKids, err := s.repository.GetPublicKeyArray(s.db)
	assert.Nil(t, keys)
	assert.Nil(t, revokedKids)
	assert.Error(t, err)
}

func (s *HostConfJwkRepositorySuite) TestGetPrivateSigningKeys() {
	t := s.Suite.T()
	keys, err := s.repository.GetPrivateSigningKeys(s.db)
	assert.Nil(t, keys)
	assert.Error(t, err)
}

func TestHostConfJwkRepositorySuite(t *testing.T) {
	suite.Run(t, new(HostConfJwkRepositorySuite))
}
