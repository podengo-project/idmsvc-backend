package smoke

import (
	"net/http"
	"testing"

	mock_rbac "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/rbac/impl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteRbacPermission struct {
	SuiteBase
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

func (s *SuiteRbacPermission) prepareDomainIpaCreate(t *testing.T, state map[string]any) {
	token, err := s.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, token)
	require.NotEqual(t, "", token.DomainToken)
}

func (s *SuiteRbacPermission) doTestDomainIpaCreate(t *testing.T) int {
	res, err := s.CreateTokenWithResponse()
	require.NoError(t, err)
	require.NotNil(t, res)
	return res.StatusCode
}

func (s *SuiteRbacPermission) doTestDomainIpaUpdate(t *testing.T) int {
	return http.StatusNotImplemented
}

func (s *SuiteRbacPermission) doTestDomainIpaPatch(t *testing.T) int {
	return http.StatusNotImplemented
}

func (s *SuiteRbacPermission) doTestDomainIpaRead(t *testing.T) int {
	return http.StatusNotImplemented
}

func (s *SuiteRbacPermission) doTestDomainList(t *testing.T) int {
	return http.StatusNotImplemented
}

func (s *SuiteRbacPermission) doTestDomainIpaDelete(t *testing.T) int {
	return http.StatusNotImplemented
}

func (s *SuiteRbacPermission) doTestHostConfExecute(t *testing.T) int {
	return http.StatusNotImplemented
}

func (s *SuiteRbacPermission) doTestJWKExecute(t *testing.T) int {
	return http.StatusNotImplemented
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
		// {
		// 	Name:  "Test idmsvc:token:create",
		// 	Given: s.prepareNoop,
		// 	Then:  s.doTestTokenCreate,
		// 	// TODO Probably this will be change to http.StatusCreated
		// 	Expected: http.StatusOK,
		// },
		// {
		// 	Name:  "Test idmsvc:domain:create",
		// 	Given: s.prepareDomainIpaCreate,
		// 	Then:  s.doTestDomainIpaCreate,
		// 	// TODO Refactor to http.StatusCreated
		// 	Expected: http.StatusOK,
		// },
		// {
		// 	Name:     "Test idmsvc:domain:update",
		// 	Given:    s.doTestDomainIpaUpdate,
		// 	Expected: http.StatusOK,
		// },
		// {
		// 	Name:     "Test idmsvc:domain:update",
		// 	Given:    s.doTestDomainIpaPatch,
		// 	Expected: http.StatusOK,
		// },
		// {
		// 	Name:     "Test idmsvc:domain:read",
		// 	Given:    s.doTestDomainIpaRead,
		// 	Expected: http.StatusOK,
		// },
		// {
		// 	Name:     "Test idmsvc:domain:delete",
		// 	Given:    s.doTestDomainIpaDelete,
		// 	Expected: http.StatusNoContent,
		// },
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
	t := s.T()
	t.Log("TestReadPermission is not implemented")
	// TODO Add the test set for the read only permission
}

func (s *SuiteRbacPermission) TestNoPermission() {
	t := s.T()
	t.Log("TestReadPermission is not implemented")
	// TODO Add the test set for an empty set of permission
}
