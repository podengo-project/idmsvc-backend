package smoke

import (
	"net/http"
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	mock_rbac "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/rbac/impl"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
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

func (s *SuiteRbacPermission) prepareDomainIpa(t *testing.T) {
	var err error
	// As an admin user
	s.As(RBACAdmin, XRHIDUser)
	s.token, err = s.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, s.token)
	require.NotEqual(t, "", s.token.DomainToken)

	// This operation set AutoEnrollmentEnabled = False whatever
	// is the value we indicate here; we have to PATCH in a second
	// operation
	s.As(XRHIDSystem)
	s.domain, err = s.RegisterIpaDomain(s.token.DomainToken,
		builder_api.NewDomain("test.example").
			WithAutoEnrollmentEnabled(pointy.Bool(true)).
			WithDomainID(&s.token.DomainId).
			WithRhelIdm(builder_api.NewRhelIdmDomain("test.example").
				WithServers([]public.DomainIpaServer{}).
				AddServer(builder_api.NewDomainIpaServer("1.test.example").
					WithHccUpdateServer(true).
					WithSubscriptionManagerId(s.systemXRHID.Identity.System.CommonName).
					Build(),
				).Build(),
			).Build(),
	)
	require.NoError(t, err)
	require.NotNil(t, s.domain)

	// As an admin user
	s.As(RBACAdmin, XRHIDUser)
	s.domain, err = s.PatchDomain(
		s.domain.DomainId.String(),
		builder_api.NewUpdateDomainUserRequest().
			WithAutoEnrollmentEnabled(pointy.Bool(true)).
			Build())
	require.NoError(t, err)
	require.NotNil(t, s.domain)
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
	res, err := s.ReadDomainWithResponse(*s.domain.DomainId)
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

func (s *SuiteRbacPermission) commonRun(rbacProfile RBACProfile, testCases []TestCasePermission) {
	t := s.T()
	xrhidProfiles := []XRHIDProfile{XRHIDUser, XRHIDServiceAccount}
	for _, testCase := range testCases {
		for _, xrhidProfile := range xrhidProfiles {
			t.Logf("rbacProfile=%s, xrhidProfile=%s: %s", rbacProfile, xrhidProfile, testCase.Name)
			require.NotNil(t, testCase.Given)
			require.NotNil(t, testCase.Then)
			require.NotEqual(t, 0, testCase.Expected)

			testCase.Given(t)

			// As a role in SuperAdmin, Admin, ReadOnly, NoPerms
			// and a User or ServiceAccount identity
			s.As(rbacProfile, xrhidProfile)
			result := testCase.Then(t)
			assert.Equal(t, testCase.Expected, result)
		}
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
	}
	return testCases
}

func (s *SuiteRbacPermission) TestSuperAdminRole() {
	// This check the wildcards
	s.commonRun(mock_rbac.ProfileSuperAdmin, s.helperCommonAdmin())
}

func (s *SuiteRbacPermission) TestAdminRole() {
	// This check the Domains administrator role
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
			Name:     "Test idmsvc:domain:update for patching",
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
			Name:     "Test idmsvc:domain:update for patching",
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
	}
	s.commonRun(mock_rbac.ProfileDomainNoPerms, testCases)
}
