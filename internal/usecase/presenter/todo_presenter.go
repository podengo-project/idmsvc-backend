package presenter

import (
	"fmt"

	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/presenter"
	"github.com/openlyinc/pointy"
	"github.com/rs/zerolog/log"
)

type todoPresenter struct{}

func NewTodoPresenter() presenter.TodoPresenter {
	return todoPresenter{}
}

func (p todoPresenter) Create(in *model.Todo, out *public.Todo) error {
	return p.FullUpdate(in, out)
}

func (p todoPresenter) List(prefix string, offset int64, count int32, data []model.Todo, output *public.ListTodo) error {
	if output == nil {
		return fmt.Errorf("'output' cannot be nil")
	}
	output.Meta = &public.Meta{
		Count: pointy.Int32(count),
	}
	output.Links = &public.Links{}
	if offset > 0 {
		output.Links.First = pointy.String(fmt.Sprintf(prefix+"/todo?offset=%d&count=%d", 0, count))
	}
	if offset-int64(count) < 0 {
		output.Links.Previous = pointy.String(fmt.Sprintf(prefix+"/todo?offset=%d&count=%d", 0, count))
	} else {
		output.Links.Previous = pointy.String(fmt.Sprintf(prefix+"/todo?offset=%d&count=%d", offset-int64(count), count))
	}
	output.Links.Next = pointy.String(fmt.Sprintf(prefix+"/todo?offset=%d&count=%d", offset+int64(count), count))
	// FIXME this is weird and I am not happy with this.
	//       I would need to find to modify the openapi spec to
	//       generate a structure more accessible.
	output.Data = &public.TodoArray{}
	*output.Data = []public.Todo{}
	slice := *output.Data
	var (
		outputItem public.Todo
		err        error
	)
	for _, item := range data {
		if err = p.Get(&item, &outputItem); err != nil {
			log.Error().Err(err)
		}
		slice = append(slice, outputItem)
	}
	*output.Data = slice
	return nil
}

func (p todoPresenter) Get(in *model.Todo, out *public.Todo) error {
	if in == nil {
		return fmt.Errorf("'in' cannot be nil")
	}
	if out == nil {
		return fmt.Errorf("'out' cannot be nil")
	}
	out.Id = pointy.Uint(in.ID)
	out.Title = pointy.String(*in.Title)
	out.Body = pointy.String(*in.Description)
	return nil
}

func (p todoPresenter) PartialUpdate(in *model.Todo, out *public.Todo) error {
	if in == nil {
		return fmt.Errorf("'in' cannot be nil")
	}
	if out == nil {
		return fmt.Errorf("'out' cannot be nil")
	}
	if in.ID == 0 {
		out.Id = nil
	} else {
		out.Id = pointy.Uint(in.ID)
	}
	if in.Title == nil {
		out.Title = nil
	} else {
		out.Title = pointy.String(*in.Title)
	}
	if in.Description == nil {
		out.Body = nil
	} else {
		out.Body = pointy.String(*in.Description)
	}
	return nil
}

func (p todoPresenter) FullUpdate(in *model.Todo, out *public.Todo) error {
	return p.PartialUpdate(in, out)
}
