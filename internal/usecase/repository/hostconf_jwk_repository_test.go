package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type HostConfJwkRepositorySuite struct {
	suite.Suite
	db         *gorm.DB
	mock       sqlmock.Sqlmock
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
	cfg := test.GetTestConfig()
	s.repository = NewHostconfJwkRepository(cfg)
}

func (s *HostConfJwkRepositorySuite) TestNewHostconfJwkRepository() {
	t := s.Suite.T()
	assert.NotNil(t, s.repository)
	assert.NotNil(t, s.repository.(*hostconfJwkRepository).config)
	assert.Panics(t, func() {
		NewHostconfJwkRepository(nil)
	})
}

func (s *HostConfJwkRepositorySuite) TestCreateJWK() {
	t := s.Suite.T()
	model, err := s.repository.CreateJWK(s.db)
	assert.Nil(t, model)
	assert.EqualError(t, err, notImplementedError.Error())
}

func (s *HostConfJwkRepositorySuite) TestRevokeJWK() {
	t := s.Suite.T()
	kid := "keyid"
	model, err := s.repository.RevokeJWK(s.db, kid)
	assert.Nil(t, model)
	assert.EqualError(t, err, notImplementedError.Error())
}

func (s *HostConfJwkRepositorySuite) TestListJWKs() {
	t := s.Suite.T()
	models, err := s.repository.ListJWKs(s.db)
	assert.Nil(t, models)
	assert.EqualError(t, err, notImplementedError.Error())
}

func (s *HostConfJwkRepositorySuite) TestGetPublicKeyArray() {
	t := s.Suite.T()
	keys, err := s.repository.GetPublicKeyArray(s.db)
	assert.Nil(t, keys)
	assert.EqualError(t, err, notImplementedError.Error())
}

func (s *HostConfJwkRepositorySuite) TestGetPrivateSigningKeys() {
	t := s.Suite.T()
	keys, err := s.repository.GetPrivateSigningKeys(s.db)
	assert.Nil(t, keys)
	assert.EqualError(t, err, notImplementedError.Error())
}

func TestHostConfJwkRepositorySuite(t *testing.T) {
	suite.Run(t, new(HostConfJwkRepositorySuite))
}
