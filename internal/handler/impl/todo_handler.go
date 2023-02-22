package impl

import (
	"net/http"

	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// List all Todos
// (GET /todo)
func (a *application) ListTodos(ctx echo.Context, params public.ListTodosParams) error {
	var (
		err    error
		data   []model.Todo
		output public.ListTodo
		tx     *gorm.DB
	)
	// TODO A call to an internal validator could be here to check public.ListTodosParams
	offset := int64(0)
	count := int32(0)
	if err = a.todo.interactor.List(&params, &offset, &count); err != nil {
		return err
	}
	tx = a.db.Begin()
	if data, err = a.todo.repository.FindAll(tx, offset, count); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	// TODO Read prefix from configuration
	if err = a.todo.presenter.List("/api/hmsidm/v1", offset, count, data, &output); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, output)
}

// Return a Todo resource
// (GET /todo/{id})
func (a *application) GetTodo(ctx echo.Context, id public.Id, params public.GetTodoParams) error {
	var (
		err    error
		data   model.Todo
		output public.Todo
		itemId uint
		tx     *gorm.DB
	)

	if err = a.todo.interactor.GetById(&id, &itemId); err != nil {
		return err
	}
	tx = a.db.Begin()
	if data, err = a.todo.repository.FindById(tx, itemId); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	if err = a.todo.presenter.Get(&data, &output); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, output)
}

// Modify an existing Todo
// (PATCH /todo/{id})
func (a *application) PartialUpdateTodo(ctx echo.Context, id public.Id, params public.PartialUpdateTodoParams) error {
	var (
		err    error
		data   *model.Todo
		output public.Todo
		input  public.Todo
		tx     *gorm.DB
	)

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	data = &model.Todo{}
	if err = a.todo.interactor.PartialUpdate(id, &params, &input, data); err != nil {
		tx.Rollback()
		return err
	}
	tx = a.db.Begin()
	if *data, err = a.todo.repository.PartialUpdate(tx, data); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	if err = a.todo.presenter.PartialUpdate(data, &output); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, output)
}

// Replace an existing Todo
// (PUT /todo/{id})
func (a *application) UpdateTodo(ctx echo.Context, id public.Id, params public.UpdateTodoParams) error {
	var (
		err    error
		data   *model.Todo
		output public.Todo
		input  public.Todo
		tx     *gorm.DB
	)

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	data = &model.Todo{}
	if err = a.todo.interactor.FullUpdate(id, &params, &input, data); err != nil {
		return err
	}
	tx = a.db.Begin()
	if *data, err = a.todo.repository.Update(tx, data); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	if err = a.todo.presenter.FullUpdate(data, &output); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, output)
}

// Create a Todo resource
// (POST /todo)
func (a *application) CreateTodo(ctx echo.Context, params public.CreateTodoParams) error {
	var (
		err    error
		input  public.Todo
		data   model.Todo
		output public.Todo
		tx     *gorm.DB
	)

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	if err = a.todo.interactor.Create(&params, &input, &data); err != nil {
		return err
	}
	tx = a.db.Begin()
	if err = a.todo.repository.Create(tx, &data); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	if err = a.todo.presenter.Create(&data, &output); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, output)
}

// Delete a Todo resource
// (POST /todo)
func (a *application) DeleteTodo(ctx echo.Context, id uint, params public.DeleteTodoParams) error {
	var err error
	var tx *gorm.DB
	tx = a.db.Begin()
	if err = a.todo.repository.DeleteById(tx, id); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return ctx.NoContent(http.StatusNoContent)
}
