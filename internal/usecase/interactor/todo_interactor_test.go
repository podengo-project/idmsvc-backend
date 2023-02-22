package interactor

import (
	"fmt"
	"testing"

	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTodoInteractor(t *testing.T) {
	var component interactor.TodoInteractor
	assert.NotPanics(t, func() {
		component = NewTodoInteractor()
	})
	assert.NotNil(t, component)
}

func TestCreate(t *testing.T) {
	type TestCaseGiven struct {
		Params *api_public.CreateTodoParams
		In     *api_public.Todo
		Out    *model.Todo
	}
	type TestCaseExpected struct {
		Err error
		Out *model.Todo
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "nil for params",
			Given: TestCaseGiven{
				Params: nil,
				In:     &api_public.Todo{},
				Out:    &model.Todo{},
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'params' cannot be nil"),
				Out: nil,
			},
		},
		{
			Name: "nil for In",
			Given: TestCaseGiven{
				Params: &api_public.CreateTodoParams{},
				In:     nil,
				Out:    &model.Todo{},
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'in' cannot be nil"),
				Out: nil,
			},
		},
		{
			Name: "nil for Out",
			Given: TestCaseGiven{
				Params: &api_public.CreateTodoParams{},
				In:     &api_public.Todo{},
				Out:    nil,
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'out' cannot be nil"),
				Out: nil,
			},
		},
		{
			Name: "success case",
			Given: TestCaseGiven{
				Params: &api_public.CreateTodoParams{},
				In: &api_public.Todo{
					Title: pointy.String("test title"),
					Body:  pointy.String("test body"),
				},
				Out: &model.Todo{},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &model.Todo{
					Title:       pointy.String("test title"),
					Description: pointy.String("test body"),
				},
			},
		},
	}
	component := NewTodoInteractor()
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		var out *model.Todo = nil
		if testCase.Given.Out != nil {
			out = &model.Todo{}
			*out = *testCase.Given.Out
		}
		result := component.Create(testCase.Given.Params, testCase.Given.In, out)
		if testCase.Expected.Err != nil {
			require.EqualError(t, result, testCase.Expected.Err.Error())
		} else {
			assert.NoError(t, result)
			require.NotNil(t, testCase.Expected.Out)

			assert.Equal(t, testCase.Expected.Out.Title, out.Title)
			assert.Equal(t, testCase.Expected.Out.Description, out.Description)
		}
	}
}
