package smoke

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/service"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"github.com/podengo-project/idmsvc-backend/internal/usecase/client"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	service_impl "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl"
)

// SuiteBase represents the base Suite to be used for smoke tests, this
// start the services before run the smoke tests.
// TODO the smoke tests cannot be executed in parallel yet, an alternative
// for them would be to use specific http and metrics service in one side,
// and to use a specific OrgID per test by using a generator for it which
// would provide data partition between the tests.
type SuiteBase struct {
	suite.Suite
	cfg         *config.Config
	OrgID       string
	UserXRHID   identity.XRHID
	SystemXRHID identity.XRHID

	cancel context.CancelFunc
	svc    service.ApplicationService
	wg     *sync.WaitGroup
	db     *gorm.DB
}

// SetupTest start the services and await until they are ready
// for being used.
func (s *SuiteBase) SetupTest() {
	s.cfg = config.Get()
	s.wg = &sync.WaitGroup{}
	logger.InitLogger(s.cfg)
	s.db = datastore.NewDB(s.cfg)

	ctx, cancel := StartSignalHandler(context.Background())
	s.cancel = cancel
	inventory := client.NewHostInventory(s.cfg)
	s.svc = service_impl.NewApplication(ctx, s.wg, s.cfg, s.db, inventory)
	go func() {
		if e := s.svc.Start(); e != nil {
			panic(e)
		}
	}()
	s.OrgID = strconv.Itoa(int(builder_helper.GenRandNum(1, 99999999)))
	s.UserXRHID = builder_api.NewUserXRHID().WithOrgID(s.OrgID).WithActive(true).Build()
	s.SystemXRHID = builder_api.NewSystemXRHID().WithOrgID(s.OrgID).Build()
	s.WaitReady(s.cfg)
}

// TearDownTest Stop the services in an ordered way before every
// smoke test executed.
func (s *SuiteBase) TearDownTest() {
	TearDownSignalHandler()
	defer datastore.Close(s.db)
	defer s.cancel()
	s.svc.Stop()
	s.wg.Wait()
}

// WaitReady poll the ready healthcheck until the response is http.StatusOK
// cfg is the current configuration to use for the application.
func (s *SuiteBase) WaitReady(cfg *config.Config) {
	if cfg == nil {
		panic("cfg is nil")
	}
	header := http.Header{}
	for i := 0; i < 300; i++ {
		resp, err := s.DoRequest(
			http.MethodGet,
			s.DefaultPrivateBaseURL()+"/readyz",
			header,
			nil,
		)
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	panic("WaitReady didn't return after 30 seconds checking for it")
}

// CreateToken is a helper function to request a token to the API for registration
// for a rhel-idm domain using the OrgID assigned to the unit test.
// Return the token response or error.
func (s *SuiteBase) CreateToken() (*public.DomainRegTokenResponse, error) {
	var headers http.Header = http.Header{}
	var resp *http.Response
	var err error

	url := s.DefaultPublicBaseURL() + "/domains/token"

	headers.Add("X-Rh-Insights-Request-Id", "get_token")
	headers.Add("X-Rh-Identity", header.EncodeXRHID(&s.UserXRHID))
	if resp, err = s.DoRequest(
		http.MethodPost,
		url,
		headers,
		&public.DomainRegTokenRequest{
			DomainType: "rhel-idm",
		},
	); err != nil {
		return nil, fmt.Errorf("failure when POST %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failure when POST %s: expected '%d' but got '%d'", url, http.StatusOK, resp.StatusCode)
	}
	var data []byte
	data, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failure when reading body for POST %s because an empty response", url)
	}
	var token *public.DomainRegTokenResponse = &public.DomainRegTokenResponse{}
	err = json.Unmarshal(data, token)
	if err != nil {
		return nil, fmt.Errorf("failure when unmarshalling the information for POST %s", url)
	}
	return token, nil
}

// RegisterIpaDomain is a helper function to register a domain with the API
// for a rhel-idm domain using the OrgID assigned to the unit test.
// Return the token response or error.
func (s *SuiteBase) RegisterIpaDomain(domain *public.Domain) (*public.Domain, error) {
	var headers http.Header = http.Header{}
	var resp *http.Response
	var err error
	var token *public.DomainRegTokenResponse = nil

	if token, err = s.CreateToken(); err != nil {
		s.FailNow("Error creating a token for registering a rhel-idm domain", "%s", err.Error())
	}

	url := s.DefaultPublicBaseURL() + "/domains"
	headers.Add("X-Rh-Idm-Registration-Token", token.DomainToken)
	headers.Add("X-Rh-Insights-Request-Id", "get_token")
	headers.Add("X-Rh-Identity", header.EncodeXRHID(&s.SystemXRHID))
	headers.Add("X-Rh-Idm-Version", header.EncodeXRHIDMVersion(&header.XRHIDMVersion{
		IPAHCCVersion:      "1.0.0",
		IPAVersion:         "4.19.0",
		OSReleaseID:        "9.3",
		OSReleaseVersionID: "redhat-9.3",
	}))
	if resp, err = s.DoRequest(
		http.MethodPost,
		url,
		headers,
		domain,
	); err != nil {
		return nil, fmt.Errorf("failure when %s %s: %w", http.MethodPost, url, err)
	}
	// FIXME This should be http.StatusCreated instead of StatusOK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failure when POST %s: expected '%d' but got '%d'", url, http.StatusOK, resp.StatusCode)
	}
	var data []byte
	data, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failure when reading body for POST %s because an empty response", url)
	}
	var createdDomain *public.Domain = &public.Domain{}
	err = json.Unmarshal(data, createdDomain)
	if err != nil {
		return nil, fmt.Errorf("failure when unmarshalling the information for %s %s", http.MethodPost, url)
	}
	return createdDomain, nil
}

// RunTestCase run test for one specific testcase
func (s *SuiteBase) RunTestCase(testCase *TestCase) {
	t := s.T()

	var (
		body []byte
		resp *http.Response
		err  error
	)

	// GIVEN testCase
	bodyCount := 0
	if testCase.Given.Body != nil {
		bodyCount++
	}
	if testCase.Given.BodyBytes != nil {
		bodyCount++
	}
	if bodyCount > 1 {
		t.Errorf("Given Body and BodyBytes are exclusive between them.")
	}
	bodyCount = 0
	if testCase.Expected.Body != nil {
		bodyCount++
	}
	if testCase.Expected.BodyFunc != nil {
		bodyCount++
	}
	if testCase.Expected.BodyBytes != nil {
		bodyCount++
	}
	if bodyCount > 1 {
		t.Errorf("Expected Body, BodyFunc and BodyBytes are exclusive between them.")
	}

	// WHEN
	resp, err = s.DoRequest(testCase.Given.Method, testCase.Given.URL, testCase.Given.Header, testCase.Given.Body)

	// THEN

	// Check no error
	require.NoError(t, err)
	if resp != nil {
		body, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)
	}

	// Check response status code
	require.Equal(t, testCase.Expected.StatusCode, resp.StatusCode)

	// Check response headers
	t.Log("Checking response headers")
	for key := range testCase.Expected.Header {
		expectedValue := fmt.Sprintf("%s: %s", key, testCase.Expected.Header.Get(key))
		currentValue := fmt.Sprintf("%s: %s", key, resp.Header.Get(key))
		assert.Equal(t, expectedValue, currentValue)
	}

	// Check response body
	if bodyCount == 0 && len(body) == 0 {
		return
	}
	if testCase.Expected.Body != nil {
		assert.Equal(t, testCase.Expected.Body, body)
	}
	if testCase.Expected.BodyFunc != nil {
		assert.True(t, testCase.Expected.BodyFunc(t, body))
	}
	if testCase.Expected.BodyBytes != nil {
		assert.Equal(t, testCase.Expected.BodyBytes, body)
	}
}

// RunTestCases run a slice of test cases.
// testCases is the list of test cases to be executed.
func (s *SuiteBase) RunTestCases(testCases []TestCase) {
	t := s.T()
	for i := range testCases {
		t.Log(testCases[i].Name)
		s.RunTestCase(&testCases[i])
	}
}

// DefaultPublicBaseURL retrieve the public base endpoint URL.
// Return for the URL for the current configuration.
func (s *SuiteBase) DefaultPublicBaseURL() string {
	return fmt.Sprintf("http://localhost:%d/api/idmsvc/v1", s.cfg.Web.Port)
}

// DefaultPrivateBaseURL retrieve the private base endpoint URL.
// Return for the URL for the current configuration.
func (s *SuiteBase) DefaultPrivateBaseURL() string {
	return fmt.Sprintf("http://localhost:%d/private", s.cfg.Web.Port)
}

// DoRequest execute a http request against a url using headers and the body specified.
// method is the HTTP method to use for the request.
// url is the to reach out.
// header represents the request headers.
// body is any golang type to be marshalled or a Reader interface (TODO future).
// Return the http.Response object and nil when the endpoint is reached out,
// or nil and a non error when some non API situation happens trying to reach
// out the endpoint.
func (s *SuiteBase) DoRequest(method string, url string, header http.Header, body any) (*http.Response, error) {
	var reader io.Reader = nil
	client := &http.Client{}

	if body != nil {
		// TODO add type assert so if a Reader interface
		// is provided, it will be used directly; this will
		// allow to wrong body to check for integration tests
		_body, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		if len(_body) > 0 {
			reader = bytes.NewReader(_body)
		}
	} else {
		reader = bytes.NewBufferString("")
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Set(key, strings.Join(value, "; "))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type BodyFunc func(t *testing.T, body []byte) bool

// TestCaseGiven represents the requirements for the smoke test to implement.
type TestCaseGiven struct {
	// Method represents the http method for the request.
	Method string
	// URL represents the url for the route to test.
	URL string
	// Header represents the set of header of the request.
	Header http.Header
	// Body represents a golang type to be marshalled before send the request;
	// this field exclude the BodyBytes field.
	Body any
	// BodyBytes represents a specific buffer for the request body; this
	// field exlude the Body field. This works for bad formed json documents,
	// and other scenarios where Body does not fit.
	BodyBytes []byte
}

// TestCaseExpect represents the expected response for a smoke test.
type TestCaseExpect struct {
	// StatusCode represents the http status code expected.
	StatusCode int
	// Header represents the expected http response headers.
	Header http.Header
	// Body represent an API type struct that after marshall should match the
	// returned response; this could be a situation, because the order of the
	// properties could not match. It is useful only when the property order
	// is deterministic, else use BodyFunc.
	Body any
	// BodyBytes represent a specific bytes returned on the expectations.
	BodyBytes []byte
	// BodyFunc represent a custom function that will return nil or error
	// to check some specifc body unserialized. This option exclude Body and
	// BodyBytes and is useful when we want expectations based on a
	// valid json document, but it is not a perfect fit of the Body.
	BodyFunc BodyFunc
}

// TestCase represents a test case for the smoke test
type TestCase struct {
	// Name represents a string to be printed out which will be displayed
	// in case of a failure happens.
	Name string
	// Given represents the given specification for the test case.
	Given TestCaseGiven
	// Expected represents the expected result for the operations.
	Expected TestCaseExpect
}

// StartSignalHandler set up the signal handler. This method MUST NOT
// be called several times, as that make no deterministic which signal
// handler will receive the call.
// c is the golang context associated, if it is nil a new background
// context is used.
// Return the cancel context generated that will called on exit and
// the cancel function associted to the context.
// See: https://pkg.go.dev/os/signal
func StartSignalHandler(c context.Context) (context.Context, context.CancelFunc) {
	if c == nil {
		c = context.Background()
	}
	ctx, cancel := context.WithCancel(c)
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-exit
		cancel()
	}()
	return ctx, cancel
}

// TearDownSignalHandler reset the signal handlers
func TearDownSignalHandler() {
	signal.Reset(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
}

func TestSuite(t *testing.T) {
	// TODO Add here your test suites
	suite.Run(t, new(SuiteTokenCreate))
	// suite.Run(t, new(SuiteRegisterDomain))
	// suite.Run(t, new(SuiteDomainUpdateUser))
	// suite.Run(t, new(SuiteDomainUpdateAgent))
	// suite.Run(t, new(SuiteDomainRead))
	// suite.Run(t, new(SuiteListDomains))
	// suite.Run(t, new(SuiteDomainDelete))
}
