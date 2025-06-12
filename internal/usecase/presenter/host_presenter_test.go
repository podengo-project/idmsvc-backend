package presenter

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

func TestNewHostPresenter(t *testing.T) {
	assert.Panics(t, func() {
		NewHostPresenter(nil)
	})
	cfg := config.Config{}
	assert.NotPanics(t, func() {
		NewHostPresenter(&cfg)
	})
}

func TestHostConf(t *testing.T) {
	currentTime := time.Now()
	testDomainID := &uuid.UUID{}
	*testDomainID = uuid.MustParse("188a62fc-0720-11ee-9dfd-482ae3863d30")
	testRHSMID := pointy.String("fe106208-dd32-11ed-aa87-482ae3863d30")
	testDomain := "ipa.test"
	testRealm := "IPA.TEST"
	testIpaCert := model.IpaCert{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		IpaID:        1,
		Nickname:     "IPA.TEST IPA CA",
		Issuer:       "CN=Certificate Authority,O=IPA.TEST",
		Subject:      "CN=Certificate Authority,O=IPA.TEST",
		SerialNumber: "1",
		NotBefore:    currentTime,
		NotAfter:     currentTime,
		Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
	}
	testIpaServer := model.IpaServer{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: gorm.DeletedAt{},
		},
		IpaID:               1,
		FQDN:                "server1.ipa.test",
		RHSMId:              testRHSMID,
		Location:            pointy.String("europe"),
		CaServer:            true,
		HCCEnrollmentServer: true,
		HCCUpdateServer:     true,
		PKInitServer:        true,
	}

	type TestCaseGiven struct {
		Input  *model.Domain
		Token  string
		Output *public.HostConfResponse
	}
	type TestCaseExpected struct {
		Err    error
		Output *public.HostConfResponse
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
				Err:    internal_errors.NilArgError("domain"),
				Output: nil,
			},
		},
		{
			Name: "Unsupported domain type",
			Given: TestCaseGiven{
				Input: &model.Domain{
					Model:                 gorm.Model{ID: 1},
					OrgId:                 "12345",
					DomainUuid:            *testDomainID,
					DomainName:            pointy.String(testDomain),
					Type:                  pointy.Uint(model.DomainTypeUndefined),
					AutoEnrollmentEnabled: pointy.Bool(true),
				},
			},
			Expected: TestCaseExpected{
				Err:    fmt.Errorf("domain 'ipa.test' has unsupported domain type ''"),
				Output: nil,
			},
		},
		{
			Name: "Missing CA certs",
			Given: TestCaseGiven{
				Input: &model.Domain{
					Model:                 gorm.Model{ID: 1},
					OrgId:                 "12345",
					DomainUuid:            *testDomainID,
					DomainName:            pointy.String(testDomain),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String(testRealm),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
						RealmDomains: pq.StringArray{testDomain},
					},
				},
			},
			Expected: TestCaseExpected{
				Err:    fmt.Errorf("domain 'ipa.test' has no CA certificates"),
				Output: nil,
			},
		},
		{
			Name: "Missing servers",
			Given: TestCaseGiven{
				Input: &model.Domain{
					Model:                 gorm.Model{ID: 1},
					OrgId:                 "12345",
					DomainUuid:            *testDomainID,
					DomainName:            pointy.String(testDomain),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String(testRealm),
						CaCerts:      []model.IpaCert{testIpaCert},
						Servers:      []model.IpaServer{},
						RealmDomains: pq.StringArray{testDomain},
					},
				},
			},
			Expected: TestCaseExpected{
				Err:    fmt.Errorf("domain 'ipa.test' has no enrollment servers"),
				Output: nil,
			},
		},
		{
			Name: "Success",
			Given: TestCaseGiven{
				Input: &model.Domain{
					Model:                 gorm.Model{ID: 1},
					OrgId:                 "12345",
					DomainUuid:            *testDomainID,
					DomainName:            pointy.String(testDomain),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String(testRealm),
						CaCerts:      []model.IpaCert{testIpaCert},
						Servers:      []model.IpaServer{testIpaServer},
						RealmDomains: pq.StringArray{testDomain},
					},
				},
				Token: "token",
			},
			Expected: TestCaseExpected{
				Err: nil,
				Output: &public.HostConfResponse{
					AutoEnrollmentEnabled: true,
					DomainType:            model.DomainTypeIpaString,
					DomainId:              *testDomainID,
					DomainName:            testDomain,
					RhelIdm: public.HostConfIpa{
						Cabundle: testIpaCert.Pem,
						EnrollmentServers: []public.HostConfIpaServer{
							{Fqdn: testIpaServer.FQDN, Location: pointy.String("europe")},
						},
						RealmName: testRealm,
					},
					Token: pointy.String("token"),
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		obj := &hostPresenter{cfg: test.GetTestConfig()}
		output, err := obj.HostConf(testCase.Given.Input, testCase.Given.Token)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Err.Error(), err.Error())
			assert.Nil(t, output)
		} else {
			assert.NoError(t, err)
			assert.Equal(t,
				testCase.Expected.Output.DomainId,
				output.DomainId)
			assert.Equal(t,
				testCase.Expected.Output.DomainName,
				output.DomainName)
			assert.Equal(t,
				testCase.Expected.Output.DomainType,
				output.DomainType)
			assert.Equal(t,
				testCase.Expected.Output.AutoEnrollmentEnabled,
				output.AutoEnrollmentEnabled)
			assert.Equal(t,
				testCase.Expected.Output.RhelIdm.RealmName,
				output.RhelIdm.RealmName)
			assert.Equal(t,
				testCase.Expected.Output.RhelIdm.Cabundle,
				output.RhelIdm.Cabundle)
			assert.Equal(t,
				testCase.Expected.Output.RhelIdm.EnrollmentServers,
				output.RhelIdm.EnrollmentServers)
			assert.Equal(t,
				testCase.Expected.Output.Token,
				output.Token)
		}
	}
}
