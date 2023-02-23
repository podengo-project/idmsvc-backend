package presenter

import (
	"fmt"
	"testing"

	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewTodoPresenter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewTodoPresenter()
	})
}

type mynewerror struct{}

func (e *mynewerror) Error() string {
	return "mynewerror"
}

func TestTodoPresenterGet(t *testing.T) {
	type TestCaseGiven struct {
		Input  *model.Todo
		Output *public.Todo
	}
	type TestCaseExpected struct {
		Err    error
		Output *public.Todo
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "error when 'in' is nil",
			Given: TestCaseGiven{
				Input:  nil,
				Output: nil,
			},
			Expected: TestCaseExpected{
				Err:    fmt.Errorf("'in' cannot be nil"),
				Output: nil,
			},
		},
		{
			Name: "error when 'out'' is nil",
			Given: TestCaseGiven{
				Input: &model.Todo{
					Model:       gorm.Model{ID: 1},
					Title:       pointy.String("mytitle"),
					Description: pointy.String("mydescription"),
				},
				Output: nil,
			},
			Expected: TestCaseExpected{
				Err:    fmt.Errorf("'out' cannot be nil"),
				Output: nil,
			},
		},
		{
			Name: "error when 'out'' is nil",
			Given: TestCaseGiven{
				Input: &model.Todo{
					Model:       gorm.Model{ID: 1},
					Title:       pointy.String("mytitle"),
					Description: pointy.String("mydescription"),
				},
				Output: &public.Todo{},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Output: &public.Todo{
					Id:    pointy.Uint(1),
					Title: pointy.String("mytitle"),
					Body:  pointy.String("mydescription"),
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		obj := NewTodoPresenter()
		err := obj.Get(testCase.Given.Input, testCase.Given.Output)
		if testCase.Expected.Err != nil {
			assert.Error(t, err)
			assert.Equal(t, testCase.Expected.Err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, *testCase.Expected.Output.Id, *testCase.Given.Output.Id)
			assert.Equal(t, *testCase.Expected.Output.Title, *testCase.Given.Output.Title)
			assert.Equal(t, *testCase.Expected.Output.Body, *testCase.Given.Output.Body)
		}
	}
}
