package impl

import (
	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/handler"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/hmsidm/internal/interface/presenter"
	"github.com/hmsidm/internal/interface/repository"
	metrics "github.com/hmsidm/internal/metrics"
	usecase_interactor "github.com/hmsidm/internal/usecase/interactor"
	usecase_presenter "github.com/hmsidm/internal/usecase/presenter"
	usecase_repository "github.com/hmsidm/internal/usecase/repository"
	"gorm.io/gorm"
)

type todoComponent struct {
	interactor interactor.TodoInteractor
	repository repository.TodoRepository
	presenter  presenter.TodoPresenter
}

type application struct {
	config  *config.Config
	metrics *metrics.Metrics
	todo    todoComponent
	db      *gorm.DB
}

func NewHandler(config *config.Config, db *gorm.DB, m *metrics.Metrics) handler.Application {
	if config == nil {
		panic("config is nil")
	}
	if db == nil {
		panic("db is nil")
	}
	i := usecase_interactor.NewTodoInteractor()
	r := usecase_repository.NewTodoRepository()
	p := usecase_presenter.NewTodoPresenter()

	// Instantiate application
	return &application{
		config:  config,
		db:      db,
		metrics: m,
		todo: todoComponent{
			interactor: i,
			repository: r,
			presenter:  p,
		},
	}
}
