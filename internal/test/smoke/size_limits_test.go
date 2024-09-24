package smoke

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type sizeLimitConfig struct {
	IdleTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	HeaderSizeLimit int
	BodySizeLimit   int
}

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteLimitSize struct {
	SizeLimitConfig map[string]sizeLimitConfig
	SuiteBase
}

func generateConfig(customSizeLimitConfig sizeLimitConfig) *config.Config {
	cfg := &config.Config{}
	_ = config.Load(cfg)
	if customSizeLimitConfig.IdleTimeout != 0 {
		cfg.Application.IdleTimeout = customSizeLimitConfig.IdleTimeout
	}
	if customSizeLimitConfig.ReadTimeout != 0 {
		cfg.Application.ReadTimeout = customSizeLimitConfig.ReadTimeout
	}
	if customSizeLimitConfig.WriteTimeout != 0 {
		cfg.Application.WriteTimeout = customSizeLimitConfig.WriteTimeout
	}
	if customSizeLimitConfig.HeaderSizeLimit != 0 {
		cfg.Application.SizeLimitRequestHeader = customSizeLimitConfig.HeaderSizeLimit
	}
	if customSizeLimitConfig.BodySizeLimit != 0 {
		cfg.Application.SizeLimitRequestBody = customSizeLimitConfig.BodySizeLimit
	}
	return cfg
}

func (s *SuiteLimitSize) SetupTest() {
	testName := s.T().Name()
	s.SuiteBase.Config = generateConfig(s.SizeLimitConfig[testName])
	s.SuiteBase.SetupTest()
}

func (s *SuiteLimitSize) TearDownTest() {
	s.SuiteBase.TearDownTest()
}

func (s *SuiteLimitSize) TestReadTimeout() {
	t := s.T()
	requestID := "test_read_timeout_with_create_token"

	// Use a pipe that force the timeout
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		// Wait a time bigger than the ReadTimeout
		time.Sleep(s.Config.Application.ReadTimeout + 1*time.Second)
	}()

	// Create a HTTP client and sent the request
	client := &http.Client{}

	url := s.DefaultPublicBaseURL() + "/domains/token"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, pr)
	require.NoError(t, err)
	require.NotNil(t, req)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	s.As(XRHIDUser)
	s.addXRHIDHeader(&req.Header, s.currentXRHID)
	s.addRequestID(&req.Header, requestID)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	errCurrentBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	errCurrent := &public.ErrorResponse{}
	err = json.Unmarshal(errCurrentBytes, errCurrent)
	require.NoError(t, err)
	require.NotNil(t, errCurrent.Errors)
	require.NotEmpty(t, *errCurrent.Errors)
	require.Equal(t, strconv.Itoa(http.StatusBadRequest), (*errCurrent.Errors)[0].Status)
	regExp := `request body has an error: reading failed: read tcp (127\.0\.0\.1|\[::1\]):(\d+)->(127\.0\.0\.1|\[::1\]):(\d+): (i/o timeout)`
	require.Regexp(t, regExp, (*errCurrent.Errors)[0].Title)
}

func (s *SuiteLimitSize) TestIdleTimeout() {
	t := s.T()
	// TODO Pending test to check idle timeout
	//      Initial idea was to do a request with
	//      Connection: Keep-Alive
	//      And await the time for idle timeout
	//      and check in some way if the connection
	//      was still open; but I have not found
	//      yet a way to check this; it was observed
	//      that even sending the Connection: Keep-Alive
	//      the reponse was not sending back the header
	//      with any value (indicating if the connection
	//      is immeditly closed or similar).
	t.SkipNow()
}

func (s *SuiteLimitSize) TestWriteTimeout() {
	t := s.T()
	// TODO Skiping so far, but open for ideas about check this behavior
	t.SkipNow()
}

func (s *SuiteLimitSize) TestLimitSizeRequestHeader() {
	t := s.T()

	// Create a HTTP client and sent the request
	client := &http.Client{}

	url := s.DefaultPublicBaseURL() + "/domains/token"
	reqBody, err := json.Marshal(&public.DomainRegTokenRequest{
		DomainType: public.RhelIdm,
	})
	require.NoError(t, err)
	require.NotEmpty(t, reqBody)
	reqBodyReader := bytes.NewReader(reqBody)
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		url,
		reqBodyReader,
	)
	require.NoError(t, err)
	require.NotNil(t, req)

	// Header limit is tricky to check, because even if we set
	// the value of 1 byte, the internal server set that size
	// plus 4096 bytes, see link below:
	// https://cs.opensource.google/go/go/+/master:src/net/http/server.go;l=926
	// so we have to go over that threshold on the headers to
	// get an 431 too large error response for the headers.
	headerSize := 4096 + s.Config.Application.SizeLimitRequestHeader
	tooBigHeader := strings.Repeat("-", headerSize)
	req.Header.Add("X-Padding", tooBigHeader)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusRequestHeaderFieldsTooLarge, resp.StatusCode)
}

func (s *SuiteLimitSize) TestLimitSizeRequestBody() {
	requestID := "test_limit_size_request_body"
	t := s.T()

	// Create a HTTP client and sent the request
	client := &http.Client{}

	hdr := http.Header{}
	url := s.DefaultPublicBaseURL() + "/domains/token"
	s.addRequestID(&hdr, requestID)
	reqBody, err := json.Marshal(&public.DomainRegTokenRequest{
		DomainType: public.RhelIdm,
	})
	require.NoError(t, err)
	require.NotEmpty(t, reqBody)
	reqBodyReader := bytes.NewReader(reqBody)
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		url,
		reqBodyReader,
	)
	require.NoError(t, err)
	require.NotNil(t, req)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for key, value := range hdr {
		req.Header.Set(key, strings.Join(value, "; "))
	}
	// Override
	s.currentXRHID = &s.userXRHID
	s.addXRHIDHeader(&req.Header, s.currentXRHID)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
}

func TestSuiteLimitSize(t *testing.T) {
	suite.Run(t, &SuiteLimitSize{
		SizeLimitConfig: map[string]sizeLimitConfig{
			"TestSuiteLimitSize/TestLimitSizeRequestHeader": {
				HeaderSizeLimit: 5,
			},
			"TestSuiteLimitSize/TestLimitSizeRequestBody": {
				BodySizeLimit: 5,
			},
			"TestSuiteLimitSize/TestReadTimeout": {
				ReadTimeout: time.Duration(1 * time.Millisecond),
			},
			"TestSuiteLimitSize/TestWriteTimeout": {
				WriteTimeout: time.Duration(1 * time.Millisecond),
			},
			"TestSuiteLimitSize/TestIdleTimeout": {
				IdleTimeout: time.Duration(1 * time.Millisecond),
			},
		},
	})
}
