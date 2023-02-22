package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/repository"
	"github.com/hmsidm/internal/test"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository.TodoRepository
}

func (s *Suite) SetupSuite() {
	var err error
	s.mock, s.DB, err = test.NewSqlMock()
	if err != nil {
		s.Suite.FailNow("Error calling gorm.Open: %s", err.Error())
		return
	}
	s.repository = NewTodoRepository()
}

func (s *Suite) TearDownSuite() {
}

func (s *Suite) TestCreate() {
	t := s.Suite.T()
	currentTime := time.Now()
	var data model.Todo = model.Todo{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		Title:       pointy.String("Todo Title"),
		Description: pointy.String("Todo Description"),
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todos" ("created_at","updated_at","deleted_at","title","description","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(data.CreatedAt, data.UpdatedAt, data.DeletedAt, data.Title, data.Description, data.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(data.ID))

	err := s.repository.Create(s.DB, &data)
	require.NoError(t, err)
}

func (s *Suite) TestCreateErrors() {
	t := s.Suite.T()
	currentTime := time.Now()
	var (
		data model.Todo = model.Todo{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			Title:       pointy.String("Todo Title"),
			Description: pointy.String("Todo Description"),
		}
		err error
	)

	err = s.repository.Create(nil, nil)
	assert.Error(t, err)

	err = s.repository.Create(s.DB, nil)
	assert.Error(t, err)

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todos" ("created_at","updated_at","deleted_at","title","description","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(data.CreatedAt, data.UpdatedAt, data.DeletedAt, data.Title, data.Description, data.ID).
		WillReturnError(fmt.Errorf("an error happened"))
	err = s.repository.Create(s.DB, &data)
	assert.Error(t, err)
	assert.Equal(t, "an error happened", err.Error())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
