package smoke

import (
	"net/http"
	"testing"

	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	mock_rbac "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/rbac/impl"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteRbacPermission struct {
	SuiteBase
	token  *public.DomainRegTokenResponse
	domain *public.Domain
}

type TestCasePermission struct {
	Name string
	// Given represent a function that launch the operation
	Given func(*testing.T)
	Then  func(*testing.T) int
	// Expected the status code, it will be <400 for an allowed
	// operation, 403 for an unauthorized operation
	Expected int
}

func (s *SuiteRbacPermission) prepareNoop(t *testing.T) {
	// It is empty on porpose
}

func (s *SuiteRbacPermission) doTestTokenCreate(t *testing.T) int {
	res, err := s.CreateTokenWithResponse()
	require.NoError(t, err)
	assert.NotNil(t, res)
	return res.StatusCode
}

func (s *SuiteRbacPermission) prepareDomainIpaCreate(t *testing.T) {
	token, err := s.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, token)
	require.NotEqual(t, "", token.DomainToken)
	s.token = token
}

func (s *SuiteRbacPermission) prepareDomainIpa(t *testing.T) {
	var err error
	s.token, err = s.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, s.token)
	require.NotEqual(t, "", s.token.DomainToken)

	// This operation set AutoEnrollmentEnabled = False whatever
	// is the value we indicate here; we have to PATCH in a second
	// operation
	s.domain, err = s.RegisterIpaDomain(s.token.DomainToken,
		builder_api.NewDomain("test.example").
			WithAutoEnrollmentEnabled(pointy.Bool(true)).
			WithDomainID(&s.token.DomainId).
			WithRhelIdm(builder_api.NewRhelIdmDomain("test.example").
				WithServers([]public.DomainIpaServer{}).
				AddServer(builder_api.NewDomainIpaServer("1.test.example").
					WithHccUpdateServer(true).
					WithSubscriptionManagerId(s.SystemXRHID.Identity.System.CommonName).
					Build(),
				).Build(),
			).Build(),
	)
	require.NoError(t, err)
	require.NotNil(t, s.domain)

	s.domain, err = s.PatchDomain(
		s.domain.DomainId.String(),
		builder_api.NewUpdateDomainUserRequest().
			WithAutoEnrollmentEnabled(pointy.Bool(true)).
			Build())
	require.NoError(t, err)
	require.NotNil(t, s.domain)
}

func (s *SuiteRbacPermission) doTestDomainIpaCreate(t *testing.T) int {
	res, err := s.RegisterIpaDomainWithResponse(s.token.DomainToken,
		builder_api.NewDomain("test.example").
			WithDomainID(&s.token.DomainId).
			WithRhelIdm(builder_api.NewRhelIdmDomain("test.example").
				AddServer(builder_api.NewDomainIpaServer("1.test.example").
					WithHccUpdateServer(true).
					WithSubscriptionManagerId(s.SystemXRHID.Identity.System.CommonName).
					Build()).
				Build(),
			).Build(),
	)
	require.NoError(t, err)
	require.NotNil(t, res)
	return res.StatusCode
}

func (s *SuiteRbacPermission) doTestDomainIpaUpdate(t *testing.T) int {
	subscriptionManagerID := s.domain.RhelIdm.Servers[0].SubscriptionManagerId.String()
	domainID := s.domain.DomainId.String()
	res, err := s.UpdateDomainWithResponse(
		domainID,
		builder_api.NewUpdateDomainAgent("test.example").
			WithHCCUpdate(true).
			WithDomainRhelIdm(*builder_api.NewRhelIdmDomain("test.example").
				WithServers([]public.DomainIpaServer{}).
				AddServer(
					builder_api.NewDomainIpaServer("1.test.example").
						WithHccUpdateServer(true).
						WithSubscriptionManagerId(subscriptionManagerID).
						Build(),
				).Build(),
			).WithSubscriptionManagerID(subscriptionManagerID).
			Build(),
	)
	require.NoError(t, err)
	require.NotNil(t, res)
	return res.StatusCode
}

func (s *SuiteRbacPermission) doTestDomainIpaPatch(t *testing.T) int {
	res, err := s.PatchDomainWithResponse(
		s.domain.DomainId.String(),
		builder_api.NewUpdateDomainUserRequest().
			WithAutoEnrollmentEnabled(pointy.Bool(true)).
			WithTitle(pointy.String("new title")).
			Build(),
	)
	require.NoError(t, err)
	require.NotNil(t, res)
	return res.StatusCode
}

func (s *SuiteRbacPermission) doTestDomainIpaRead(t *testing.T) int {
	res, err := s.UserReadDomainWithResponse(*s.domain.DomainId)
	require.NoError(t, err)
	require.NotNil(t, res)
	return res.StatusCode
}

func (s *SuiteRbacPermission) doTestDomainIpaDelete(t *testing.T) int {
	res, err := s.DeleteDomainWithResponse(*s.domain.DomainId)
	require.NoError(t, err)
	require.NotNil(t, res)
	return res.StatusCode
}

func (s *SuiteRbacPermission) doTestDomainList(t *testing.T) int {
	res, err := s.ListDomainWithResponse(0, 10)
	require.NoError(t, err)
	require.NotNil(t, res)
	return res.StatusCode
}

func (s *SuiteRbacPermission) doTestReadSigningKeys(t *testing.T) int {
	res, err := s.SystemSigningKeysWithResponse()
	require.NoError(t, err)
	return res.StatusCode
}

func (s *SuiteRbacPermission) doTestSystemReadDomain(t *testing.T) int {
	res, err := s.SystemReadDomainWithResponse(*s.domain.DomainId)
	require.NoError(t, err)
	return res.StatusCode
}

func (s *SuiteRbacPermission) commonRun(profile string, testCases []TestCasePermission) {
	t := s.T()

	for _, testCase := range testCases {
		t.Logf("profile=%s: %s", profile, testCase.Name)
		require.NotNil(t, testCase.Given)
		require.NotNil(t, testCase.Then)
		require.NotEqual(t, 0, testCase.Expected)

		s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileSuperAdmin])
		testCase.Given(t)
		s.RbacMock.SetPermissions(mock_rbac.Profiles[profile])
		result := testCase.Then(t)
		assert.Equal(t, testCase.Expected, result)
	}
}

func (s *SuiteRbacPermission) helperCommonAdmin() []TestCasePermission {
	testCases := []TestCasePermission{
		{
			Name:     "Test idmsvc:token:create",
			Given:    s.prepareNoop,
			Then:     s.doTestTokenCreate,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test User Update idmsvc:domain:update",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaPatch,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test idmsvc:domain:read",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaRead,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test idmsvc:domain:delete",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaDelete,
			Expected: http.StatusNoContent,
		},
		{
			Name:     "Test idmsvc:domain:list",
			Given:    s.prepareNoop,
			Then:     s.doTestDomainList,
			Expected: http.StatusOK,
		},
		// System requests are identified by its certificate
		{
			Name:     "Test Agent register domain",
			Given:    s.prepareDomainIpaCreate,
			Then:     s.doTestDomainIpaCreate,
			Expected: http.StatusCreated,
		},
		{
			Name:     "Test Update Agent",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaUpdate,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test Agent Read domain",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestSystemReadDomain,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test Read SigningKeys",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestReadSigningKeys,
			Expected: http.StatusOK,
		},
	}
	return testCases
}

func (s *SuiteRbacPermission) TestSuperAdminRole() {
	s.commonRun(mock_rbac.ProfileSuperAdmin, s.helperCommonAdmin())
}

func (s *SuiteRbacPermission) TestAdminRole() {
	s.commonRun(mock_rbac.ProfileDomainAdmin, s.helperCommonAdmin())
}

func (s *SuiteRbacPermission) TestReadPermission() {
	testCases := []TestCasePermission{
		{
			Name:     "Test idmsvc:token:create",
			Given:    s.prepareNoop,
			Then:     s.doTestTokenCreate,
			Expected: http.StatusUnauthorized,
		},
		{
			Name:     "Test Update User idmsvc:domain:update",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaPatch,
			Expected: http.StatusUnauthorized,
		},
		{
			Name:     "Test idmsvc:domain:read",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaRead,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test idmsvc:domain:delete",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaDelete,
			Expected: http.StatusUnauthorized,
		},
		{
			Name:     "Test idmsvc:domain:list",
			Given:    s.prepareNoop,
			Then:     s.doTestDomainList,
			Expected: http.StatusOK,
		},
		// System requests are identified by its certificate
		{
			Name:     "Test Agent register domain",
			Given:    s.prepareDomainIpaCreate,
			Then:     s.doTestDomainIpaCreate,
			Expected: http.StatusCreated,
		},
		{
			Name:     "Test Update Agent",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaUpdate,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test Agent Read domain",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestSystemReadDomain,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test Read SigningKeys",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestReadSigningKeys,
			Expected: http.StatusOK,
		},
	}
	s.commonRun(mock_rbac.ProfileDomainReadOnly, testCases)
}

func (s *SuiteRbacPermission) TestNoPermission() {
	testCases := []TestCasePermission{
		{
			Name:     "Test idmsvc:token:create",
			Given:    s.prepareNoop,
			Then:     s.doTestTokenCreate,
			Expected: http.StatusUnauthorized,
		},
		{
			Name:     "Test Update User idmsvc:domain:update",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaPatch,
			Expected: http.StatusUnauthorized,
		},
		{
			Name:     "Test idmsvc:domain:read",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaRead,
			Expected: http.StatusUnauthorized,
		},
		{
			Name:     "Test idmsvc:domain:delete",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaDelete,
			Expected: http.StatusUnauthorized,
		},
		{
			Name:     "Test idmsvc:domain:list",
			Given:    s.prepareNoop,
			Then:     s.doTestDomainList,
			Expected: http.StatusUnauthorized,
		},
		// System requests are identified by its certificate
		{
			Name:     "Test Agent register domain",
			Given:    s.prepareDomainIpaCreate,
			Then:     s.doTestDomainIpaCreate,
			Expected: http.StatusCreated,
		},
		{
			Name:     "Test Update Agent",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestDomainIpaUpdate,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test Agent Read domain",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestSystemReadDomain,
			Expected: http.StatusOK,
		},
		{
			Name:     "Test Read SigningKeys",
			Given:    s.prepareDomainIpa,
			Then:     s.doTestReadSigningKeys,
			Expected: http.StatusOK,
		},
	}
	s.commonRun(mock_rbac.ProfileDomainNoPerms, testCases)
}

func (s *SuiteRbacPermission) TestHostConfExecute() {
	// This is executed on their own test because
	// one verification is that only one domain match
	// for the current organization, for the specified
	// criteria.
	t := s.T()
	s.prepareDomainIpa(t)
	domainType := public.RhelIdm
	res, err := s.SystemHostConfWithResponse(
		s.domain.RhelIdm.Servers[0].SubscriptionManagerId.String(),
		"client."+s.domain.DomainName,
		builder_api.NewHostConf().
			WithDomainName(pointy.String(s.domain.DomainName)).
			WithDomainType(&domainType).
			Build())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}
