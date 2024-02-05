package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"github.com/stretchr/testify/assert"
)

type SuiteHost struct {
	SuiteBase
	repository *hostRepository
}

// https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
func (s *SuiteHost) SetupTest() {
	s.SuiteBase.SetupTest()
	s.repository = &hostRepository{}
}

func (s *SuiteHost) TestNewHostRepository() {
	t := s.Suite.T()
	assert.NotPanics(t, func() {
		_ = NewHostRepository()
	})
}
