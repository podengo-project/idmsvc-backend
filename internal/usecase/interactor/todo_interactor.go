package interactor

import (
	"fmt"

	"github.com/hmsidm/internal/api/public"
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/openlyinc/pointy"
)

type TodoInteractor struct{}

func NewTodoInteractor() interactor.TodoInteractor {
	return TodoInteractor{}
}

func (i TodoInteractor) Create(params *api_public.CreateTodoParams, in *api_public.Todo, out *model.Todo) error {
	if params == nil {
		return fmt.Errorf("'params' cannot be nil")
	}
	if in == nil {
		return fmt.Errorf("'in' cannot be nil")
	}
	if out == nil {
		return fmt.Errorf("'out' cannot be nil")
	}
	out.Title = pointy.String(*in.Title)
	out.Description = pointy.String(*in.Body)
	return nil
}

func (i TodoInteractor) PartialUpdate(id public.Id, params *api_public.PartialUpdateTodoParams, in *api_public.Todo, out *model.Todo) error {
	if id <= 0 {
		return fmt.Errorf("'id' should be a positive int64")
	}
	if params == nil {
		return fmt.Errorf("'params' cannot be nil")
	}
	if in == nil {
		return fmt.Errorf("'in' cannot be nil")
	}
	if out == nil {
		return fmt.Errorf("'out' cannot be nil")
	}
	out.Model.ID = id
	if in.Title != nil {
		out.Title = pointy.String(*in.Title)
	}
	if in.Body != nil {
		out.Description = pointy.String(*in.Body)
	}
	return nil
}

func (i TodoInteractor) FullUpdate(id public.Id, params *api_public.UpdateTodoParams, in *api_public.Todo, out *model.Todo) error {
	if id <= 0 {
		return fmt.Errorf("'id' should be a positive int64")
	}
	if params == nil {
		return fmt.Errorf("'params' cannot be nil")
	}
	if in == nil {
		return fmt.Errorf("'in' cannot be nil")
	}
	if out == nil {
		return fmt.Errorf("'out' cannot be nil")
	}
	out.ID = id
	out.Title = pointy.String(*in.Title)
	out.Description = pointy.String(*in.Body)
	return nil
}

func (i TodoInteractor) Delete(id public.Id, params *api_public.DeleteTodoParams, out *uint) error {
	if params == nil {
		return fmt.Errorf("'params' cannot be nil")
	}
	if out == nil {
		return fmt.Errorf("'out' cannot be nil")
	}
	*out = id
	return nil
}

func (i TodoInteractor) List(params *api_public.ListTodosParams, offset *int64, limit *int32) error {
	if params == nil {
		return fmt.Errorf("'in' cannot be nil")
	}
	if offset == nil {
		return fmt.Errorf("'offset' cannot be nil")
	}
	if limit == nil {
		return fmt.Errorf("'limit' cannot be nil")
	}
	*offset = *params.Offset
	*limit = *params.Limit
	return nil
}

func (i TodoInteractor) GetById(in *api_public.Id, id *uint) error {
	if in == nil {
		return fmt.Errorf("'in' cannot be nil")
	}
	if id == nil {
		return fmt.Errorf("'id' cannot be nil")
	}
	*id = *in
	return nil
}
