package interactor

import (
	"github.com/hmsidm/internal/api/event"
	"github.com/hmsidm/internal/domain/model"
)

type TodoEventInteractor interface {
	TodoCreated(in *event.TodoCreatedEvent, out *model.Todo) error
}
