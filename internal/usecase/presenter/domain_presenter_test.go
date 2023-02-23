package presenter

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewTodoPresenter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewDomainPresenter()
	})
}

type mynewerror struct{}

func (e *mynewerror) Error() string {
	return "mynewerror"
}

func TestDomainPresenterGet(t *testing.T) {
	testUuid := uuid.New()
	type TestCaseGiven struct {
		Input  *model.Domain
		Output *public.ReadDomainResponse
	}
	type TestCaseExpected struct {
		Err    error
		Output *public.ReadDomainResponse
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
				Input: nil,
			},
			Expected: TestCaseExpected{
				Err:    fmt.Errorf("'domain' cannot be nil"),
				Output: nil,
			},
		},
		{
			Name: "Success case",
			Given: TestCaseGiven{
				Input: &model.Domain{
					Model:                 gorm.Model{ID: 1},
					OrgId:                 "12345",
					DomainUuid:            testUuid,
					DomainName:            pointy.String("domain.example"),
					DomainType:            pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
					IpaDomain: &model.Ipa{
						RealmName:  pointy.String("DOMAIN.EXAMPLE"),
						CaList:     pointy.String(""),
						ServerList: pointy.String("server1.domain.example,server2.domain.example"),
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Output: &public.ReadDomainResponse{
					AutoEnrollmentEnabled: pointy.Bool(true),
					DomainUuid:            pointy.String(testUuid.String()),
					DomainName:            pointy.String("domain.example"),
					DomainType:            pointy.String(model.DomainTypeString(model.DomainTypeIpa)),
					Ipa: &public.ReadDomainIpa{
						RealmName:  pointy.String("DOMAIN.EXAMPLE"),
						CaList:     "",
						ServerList: &[]string{"server1.domain.example", "server2.domain.example"},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		obj := NewDomainPresenter()
		output, err := obj.Get(testCase.Given.Input)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Err.Error(), err.Error())
			assert.Nil(t, output)
		} else {
			assert.NoError(t, err)
			assert.Equal(t,
				*testCase.Expected.Output.DomainUuid,
				*output.DomainUuid)
			assert.Equal(t,
				*testCase.Expected.Output.DomainName,
				*output.DomainName)
			assert.Equal(t,
				*testCase.Expected.Output.DomainType,
				*output.DomainType)
			assert.Equal(t,
				*testCase.Expected.Output.AutoEnrollmentEnabled,
				*output.AutoEnrollmentEnabled)
			assert.Equal(t,
				*testCase.Expected.Output.Ipa.RealmName,
				*output.Ipa.RealmName)
			assert.Equal(t,
				testCase.Expected.Output.Ipa.CaList,
				output.Ipa.CaList)
			assert.Equal(t,
				testCase.Expected.Output.Ipa.ServerList,
				output.Ipa.ServerList)
		}
	}
}
