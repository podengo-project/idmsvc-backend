package interactor

import (
	"github.com/hmsidm/internal/api/public"
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
)

type TodoInteractor interface {
	Create(params *api_public.CreateTodoParams, body *api_public.Todo, out *model.Todo) error
	PartialUpdate(id public.Id, params *api_public.PartialUpdateTodoParams, in *api_public.Todo, out *model.Todo) error
	FullUpdate(id public.Id, params *api_public.UpdateTodoParams, in *api_public.Todo, out *model.Todo) error
	Delete(id public.Id, params *api_public.DeleteTodoParams, out *uint) error
	List(params *api_public.ListTodosParams, offset *int64, limit *int32) error
	GetById(params *api_public.Id, out *uint) error
}
