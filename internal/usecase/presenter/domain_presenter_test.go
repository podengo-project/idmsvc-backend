package presenter

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

func TestNewTodoPresenter(t *testing.T) {
	assert.Panics(t, func() {
		NewDomainPresenter(nil)
	})

	assert.NotPanics(t, func() {
		NewDomainPresenter(test.GetTestConfig())
	})
}

func TestGet(t *testing.T) {
	testUUID := &uuid.UUID{}
	*testUUID = uuid.MustParse("188a62fc-0720-11ee-9dfd-482ae3863d30")
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
				Err:    internal_errors.NilArgError("domain"),
				Output: nil,
			},
		},
		{
			Name: "Success case",
			Given: TestCaseGiven{
				Input: &model.Domain{
					Model:                 gorm.Model{ID: 1},
					OrgId:                 "12345",
					DomainUuid:            *testUUID,
					DomainName:            pointy.String("domain.example"),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String("DOMAIN.EXAMPLE"),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
						RealmDomains: pq.StringArray{"domain.example"},
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Output: &public.ReadDomainResponse{
					AutoEnrollmentEnabled: pointy.Bool(true),
					DomainId:              testUUID,
					DomainName:            "domain.example",
					DomainType:            public.DomainType(model.DomainTypeString(model.DomainTypeIpa)),
					RhelIdm: &public.DomainIpa{
						RealmName:    "DOMAIN.EXAMPLE",
						CaCerts:      []public.Certificate{},
						Servers:      []public.DomainIpaServer{},
						RealmDomains: []string{"domain.example"},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		obj := &domainPresenter{cfg: test.GetTestConfig()}
		output, err := obj.Get(testCase.Given.Input)
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
				testCase.Expected.Output.RhelIdm.CaCerts,
				output.RhelIdm.CaCerts)
			assert.Equal(t,
				testCase.Expected.Output.RhelIdm.Servers,
				output.RhelIdm.Servers)
		}
	}
}

func TestFillRhelmIdmCertsPanics(t *testing.T) {
	var err error
	p := &domainPresenter{cfg: test.GetTestConfig()}

	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(nil, nil)
	})

	domain := &model.Domain{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(nil, domain)
	})

	domain.IpaDomain = &model.Ipa{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(nil, domain)
	})

	output := &public.Domain{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(output, domain)
	})

	domain.IpaDomain.CaCerts = []model.IpaCert{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(output, domain)
	})

	// Minimal success
	output.RhelIdm = &public.DomainIpa{}
	p.fillRhelIdmCerts(output, domain)
	assert.NoError(t, err)
}

func TestFillRhelmIdmCerts(t *testing.T) {
	NotBefore := time.Now()
	NotAfter := NotBefore.Add(time.Hour * 24)
	type TestCaseGiven struct {
		To   *public.Domain
		From *model.Domain
	}
	type TestCaseExpected struct {
		Err error
		To  *public.Domain
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "Full success copy",
			Given: TestCaseGiven{
				To: &public.Domain{
					DomainType: public.RhelIdm,
					RhelIdm:    &public.DomainIpa{},
				},
				From: &model.Domain{
					Type: pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						CaCerts: []model.IpaCert{
							{
								Nickname:     "MYDOMAIN.EXAMPLE.IPA CA",
								Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
								NotBefore:    NotBefore,
								NotAfter:     NotAfter,
								SerialNumber: "1",
								Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
								Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				To: &public.Domain{
					DomainType: public.RhelIdm,
					RhelIdm: &public.DomainIpa{
						CaCerts: []public.Certificate{
							{
								Nickname:     "MYDOMAIN.EXAMPLE.IPA CA",
								Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
								NotBefore:    NotBefore,
								NotAfter:     NotAfter,
								SerialNumber: "1",
								Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
								Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		// Instantiate directly to access the private methods
		p := &domainPresenter{cfg: test.GetTestConfig()}
		if testCase.Expected.Err != nil {
			// assert.EqualError(t, err, testCase.Expected.Err.Error())
			assert.Panics(t, func() {
				p.fillRhelIdmCerts(testCase.Given.To, testCase.Given.From)
			})
		} else {
			// assert.NoError(t, err)
			assert.NotPanics(t, func() {
				p.fillRhelIdmCerts(testCase.Given.To, testCase.Given.From)
			})
			require.NotNil(t, testCase.Expected.To)
			require.NotNil(t, testCase.Expected.To.RhelIdm)
			require.NotNil(t, testCase.Expected.To.RhelIdm.CaCerts)
			require.NotNil(t, testCase.Given.To)
			require.NotNil(t, testCase.Given.To.RhelIdm)
			require.NotNil(t, testCase.Given.To.RhelIdm.CaCerts)
			assert.Equal(t, testCase.Expected.To.RhelIdm.CaCerts, testCase.Given.To.RhelIdm.CaCerts)
		}
	}
}

func TestRegister(t *testing.T) {
	testUUID := uuid.MustParse("ebac2444-e51b-11ed-a7f5-482ae3863d30")
	testDomainName := "mydomain.example"
	testTitle := pointy.String("My Example Domain Title")
	testModel := model.Domain{
		DomainUuid:            testUUID,
		DomainName:            pointy.String(testDomainName),
		Title:                 testTitle,
		Description:           pointy.String("My Example Domain Description"),
		OrgId:                 "12345",
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain: &model.Ipa{
			RealmName:    pointy.String(testDomainName),
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []model.IpaCert{},
			Servers:      []model.IpaServer{},
			Locations:    []model.IpaLocation{},
		},
	}
	testExpected := public.Domain{
		DomainId:              &testUUID,
		DomainName:            testDomainName,
		Title:                 testTitle,
		Description:           pointy.String("My Example Domain Description"),
		AutoEnrollmentEnabled: pointy.Bool(true),
		DomainType:            public.RhelIdm,
		RhelIdm: &public.DomainIpa{
			RealmName:    testDomainName,
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []public.Certificate{},
			Servers:      []public.DomainIpaServer{},
			Locations:    []public.Location{},
		},
	}

	p := &domainPresenter{cfg: test.GetTestConfig()}

	domain, err := p.Register(nil)
	assert.EqualError(t, err, "code=500, message='domain' cannot be nil")
	assert.Nil(t, domain)

	domain, err = p.Register(&testModel)
	assert.NoError(t, err)
	assert.Equal(t, testExpected, *domain)
}

func TestUpdate(t *testing.T) {
	testUUIDString := pointy.String("ebac2444-e51b-11ed-a7f5-482ae3863d30")
	testUUID := uuid.MustParse(*testUUIDString)
	testTitle := pointy.String("My Example Domain Title")
	testDescription := pointy.String("My Example Domain Description")
	testDomainName := "mydomain.example"
	testModel := model.Domain{
		DomainUuid:            testUUID,
		DomainName:            pointy.String(testDomainName),
		Title:                 testTitle,
		Description:           testDescription,
		OrgId:                 "12345",
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain: &model.Ipa{
			RealmName:    pointy.String(testDomainName),
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []model.IpaCert{},
			Servers:      []model.IpaServer{},
		},
	}
	testExpected := public.Domain{
		DomainId:              &testUUID,
		DomainName:            testDomainName,
		Title:                 testTitle,
		Description:           testDescription,
		AutoEnrollmentEnabled: pointy.Bool(true),
		DomainType:            public.RhelIdm,
		RhelIdm: &public.DomainIpa{
			RealmName:    testDomainName,
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []public.Certificate{},
			Servers:      []public.DomainIpaServer{},
			Locations:    []public.Location{},
		},
	}

	p := &domainPresenter{cfg: test.GetTestConfig()}

	domain, err := p.UpdateAgent(nil)
	assert.EqualError(t, err, "code=500, message='domain' cannot be nil")
	assert.Nil(t, domain)

	domain, err = p.UpdateAgent(&testModel)
	assert.NoError(t, err, "")
	assert.Equal(t, testExpected, *domain)

	domain, err = p.UpdateUser(nil)
	assert.EqualError(t, err, "code=500, message='domain' cannot be nil")
	assert.Nil(t, domain)

	domain, err = p.UpdateUser(&testModel)
	assert.NoError(t, err, "")
	assert.Equal(t, testExpected, *domain)
}

func TestList(t *testing.T) {
	// https://consoledot.pages.redhat.com/docs/dev/developer-references/rest/pagination.html
	testOrgID := "12345"
	testUUID1 := uuid.MustParse("5427c3d6-eaa1-11ed-99da-482ae3863d30")
	testUUID2 := uuid.MustParse("5ae8e844-eaa1-11ed-8f71-482ae3863d30")
	p := &domainPresenter{cfg: test.GetTestConfig()}

	// offset lower than 0
	count := int64(5)
	offset := -1
	limit := -1
	output, err := p.List(count, offset, limit, nil)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'offset' is lower than 0")

	// limit lower than 0
	offset = 5
	output, err = p.List(count, offset, limit, nil)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'limit' is lower than 0")

	// Offset is higher or equal to count
	limit = 10
	output, err = p.List(count, offset, limit, nil)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'offset' is higher or equal to 'count'")

	// set default limit
	offset = 0
	limit = 0
	output, err = p.List(count, offset, limit, nil)
	assert.NotNil(t, output)

	assert.Equal(t, count, output.Meta.Count)
	assert.Equal(t, p.cfg.Application.PaginationDefaultLimit, output.Meta.Limit)
	assert.Equal(t, offset, output.Meta.Offset)

	require.NotNil(t, output.Links.First)
	assert.Equal(t, p.buildPaginationLink(0, p.cfg.Application.PaginationDefaultLimit), *output.Links.First)
	assert.Nil(t, output.Links.Previous)
	assert.Nil(t, output.Links.Next)
	require.NotNil(t, output.Links.Last)
	assert.Equal(t, p.buildPaginationLink(0, p.cfg.Application.PaginationDefaultLimit), *output.Links.Last)

	// set max limit  paginationMaxLimit
	limit = p.cfg.Application.PaginationMaxLimit + 1
	output, err = p.List(count, offset, limit, nil)
	assert.NotNil(t, output)

	assert.Equal(t, count, output.Meta.Count)
	assert.Equal(t, p.cfg.Application.PaginationMaxLimit, output.Meta.Limit)
	assert.Equal(t, offset, output.Meta.Offset)

	require.NotNil(t, output.Links.First)
	assert.Equal(t, p.buildPaginationLink(0, p.cfg.Application.PaginationMaxLimit), *output.Links.First)
	assert.Nil(t, output.Links.Previous)
	assert.Nil(t, output.Links.Next)
	require.NotNil(t, output.Links.Last)
	assert.Equal(t, p.buildPaginationLink(0, p.cfg.Application.PaginationMaxLimit), *output.Links.Last)

	// domain slice is nil return empty list
	count = int64(0)
	offset = 0
	limit = 2
	expected := public.ListDomainsResponseSchema{
		Meta: public.PaginationMeta{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		},
		Links: public.PaginationLinks{
			First: pointy.String(p.buildPaginationLink(offset, limit)),
			Last:  pointy.String(p.buildPaginationLink(offset, limit)),
		},
		Data: []public.ListDomainsData{},
	}
	output, err = p.List(count, offset, limit, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, *output)

	// domain slice fill the list
	data := []model.Domain{
		{
			OrgId:                 testOrgID,
			DomainUuid:            testUUID1,
			DomainName:            pointy.String("mydomain1.example"),
			AutoEnrollmentEnabled: pointy.Bool(true),
			Title:                 pointy.String("mydomain1 example title"),
			Description:           pointy.String("mydomain1.example located in Boston"),
			Type:                  pointy.Uint(model.DomainTypeIpa),
		},
		{
			OrgId:                 testOrgID,
			DomainUuid:            testUUID2,
			DomainName:            pointy.String("mydomain2.example"),
			AutoEnrollmentEnabled: nil,
			Title:                 pointy.String("mydomain2 example title"),
			Description:           pointy.String("mydomain2.example located in Brno"),
			Type:                  pointy.Uint(model.DomainTypeIpa),
		},
	}
	count = int64(len(data))
	offset = 0
	limit = 2
	expected = public.ListDomainsResponseSchema{
		Meta: public.PaginationMeta{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		},
		Links: public.PaginationLinks{
			First: pointy.String(p.buildPaginationLink(offset, limit)),
			Last:  pointy.String(p.buildPaginationLink(offset, limit)),
		},
		Data: []public.ListDomainsData{
			{
				AutoEnrollmentEnabled: true,
				DomainType:            public.RhelIdm,
				DomainName:            "mydomain1.example",
				DomainId:              testUUID1,
				Title:                 "mydomain1 example title",
				Description:           "mydomain1.example located in Boston",
			},
			{
				AutoEnrollmentEnabled: false,
				DomainType:            public.RhelIdm,
				DomainName:            "mydomain2.example",
				DomainId:              testUUID2,
				Title:                 "mydomain2 example title",
				Description:           "mydomain2.example located in Brno",
			},
		},
	}
	output, err = p.List(count, offset, limit, data)
	assert.NoError(t, err)
	assert.Equal(t, &expected, output)
}

func TestCreateDomainToken(t *testing.T) {
	tok := &repository.DomainRegToken{
		DomainId:     uuid.New(),
		DomainToken:  "",
		DomainType:   public.RhelIdm,
		ExpirationNS: 0,
	}
	cfg := test.GetTestConfig()
	p := &domainPresenter{cfg: cfg}
	newTok, err := p.CreateDomainToken(tok)
	assert.Equal(t, tok.DomainToken, newTok.DomainToken)
	assert.Equal(t, tok.DomainId, newTok.DomainId)
	assert.Equal(t, tok.DomainType, newTok.DomainType)
	assert.NoError(t, err)
}
