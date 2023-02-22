package presenter

import (
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
)

type TodoPresenter interface {
	Create(in *model.Todo, out *public.Todo) error
	List(prefix string, offset int64, count int32, data []model.Todo, out *public.ListTodo) error
	Get(in *model.Todo, out *public.Todo) error
	PartialUpdate(in *model.Todo, out *public.Todo) error
	FullUpdate(in *model.Todo, out *public.Todo) error
}
