package smoke

import (
	"context"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	mock_rbac "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/rbac/impl"
	client_idmsvc "github.com/podengo-project/idmsvc-backend/pkg/public"
	"github.com/redhatinsights/platform-go-middlewares/identity"
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
	Given func(*testing.T) int
	// Expected the status code, it will be <400 for an allowed
	// operation, 403 for an unauthorized operation
	Expected int
}

func (s *SuiteRbacPermission) addXRHIDHeader(xrhid *identity.XRHID) client_idmsvc.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Rh-Identity", header.EncodeXRHID(xrhid))
		return nil
	}
}

func (s *SuiteRbacPermission) doTestTokenCreate(t *testing.T) int {
	c, err := client_idmsvc.NewClientWithResponses("http://localhost:8000/api/idmsvc/v1")
	require.NoError(t, err)
	assert.NotNil(t, c)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := c.CreateDomainTokenWithBody(ctx,
		&client_idmsvc.CreateDomainTokenParams{
			XRhInsightsRequestId: pointy.String("test_permission_idmsvc_token_create"),
		},
		echo.MIMEApplicationJSON,
		http.NoBody,
		s.addXRHIDHeader(&s.UserXRHID),
	)
	return res.StatusCode
}

func (s *SuiteRbacPermission) doTestDomainIpaCreate(t *testing.T) int {
	return http.StatusNotImplemented
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

func (s *SuiteRbacPermission) commonAdminAndSuperAdminRole() {
	t := s.T()
	// Prepare the tests
	testCases := []TestCasePermission{
		{
			Name:     "Test idmsvc:token:create",
			Given:    s.doTestTokenCreate,
			Expected: http.StatusCreated,
		},
		// {
		// 	Name:     "Test idmsvc:domain:create",
		// 	Given:    s.doTestDomainIpaCreate,
		// 	Expected: http.StatusCreated,
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

	// Execute the test cases
	for _, testCase := range testCases {
		result := testCase.Given(t)
		assert.Equal(t, testCase.Expected, result)
	}
}

func (s *SuiteRbacPermission) TestSuperAdminRole() {
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileSuperAdmin])
	s.commonAdminAndSuperAdminRole()
}

func (s *SuiteRbacPermission) TestAdminRole() {
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileDomainAdmin])
	s.commonAdminAndSuperAdminRole()
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
