package smoke

import (
	"net/http"
	"testing"
)

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteRbacPermission struct {
	SuiteBase
}

type TestCasePermission struct {
	Name string
	// Given represent a function that launch the operation
	Given func(t *testing.T) int
	// Expected the status code, it will be <400 for an allowed
	// operation, 403 for an unauthorized operation
	Expected int
}

func (s *SuiteRbacPermission) doTestTokenCreate(t *testing.T) int {
	return http.StatusNotImplemented
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

func (s *SuiteRbacPermission) TestEditPermission() {
	// Prepare the tests
	testCases := []TestCasePermission{
		{
			Name:     "Test idmsvc:token:create",
			Given:    doTestTokenCreate,
			Expected: http.StatusCreated,
		},
		{
			Name:     "Test idmsvc:domain:create",
			Given:    doTestTokenCreate,
			Expected: http.StatusCreated,
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}

func (s *SuiteRbacPermission) TestReadPermission() {
}

func (s *SuiteRbacPermission) TestNoPermission() {
}
