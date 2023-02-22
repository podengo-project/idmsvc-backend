package interactor

import (
	"fmt"

	"github.com/hmsidm/internal/api/event"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
)

type todoEventInteractor struct {
}

func NewTodoEvent() interactor.TodoEventInteractor {
	return &todoEventInteractor{}
}

func (i *todoEventInteractor) TodoCreated(in *event.TodoCreatedEvent, out *model.Todo) error {
	return fmt.Errorf("Not implemented")
}
