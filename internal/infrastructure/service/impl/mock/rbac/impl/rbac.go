package impl

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	app_middleware "github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/service"
	"gopkg.in/yaml.v3"
)

var (
	errRbacMockAwaitTimeout = errors.New("timeout awaiting rbac mock to be ready")
	errRbacMockUnknown      = errors.New("unknown error happened on rbac mock")
)

//go:embed super-admin.yaml
var rbacProfileSuperAdmin []byte

//go:embed domain-admin.yaml
var rbacProfileAdmin []byte

//go:embed domain-readonly.yaml
var rbacProfileReadOnly []byte

//go:embed domain-none.yaml
var rbacProfileNoPerms []byte

//go:embed custom.yaml
var rbacProfileCustom []byte

const (
	ProfileSuperAdmin     = "super-admin"
	ProfileDomainAdmin    = "domain-admin"
	ProfileDomainReadOnly = "domain-readonly"
	ProfileDomainNoPerms  = "domain-no-perms"
	ProfileCustom         = "custom"
)

var Profiles map[string][]string = map[string][]string{
	ProfileSuperAdmin:     LoadProfile(rbacProfileSuperAdmin),
	ProfileDomainAdmin:    LoadProfile(rbacProfileAdmin),
	ProfileDomainReadOnly: LoadProfile(rbacProfileReadOnly),
	ProfileDomainNoPerms:  LoadProfile(rbacProfileNoPerms),
	ProfileCustom:         LoadProfile(rbacProfileCustom),
}

// https://consoledot.pages.redhat.com/docs/dev/services/rbac.html#_retrieve_and_handle_access_list
type Page struct {
	Meta  map[string]any    `json:"meta"`
	Links map[string]string `json:"links"`
	Data  []Permission      `json:"data"`
}

type Permission struct {
	Permission          string `json:"permission"`
	ResourceDefinitions []any  `json:"resourceDefinitions"`
}

type MockRbac interface {
	SetPermissions(data []string)
	GetBaseURL() string
	WaitAddress(timeout time.Duration) error
}

type mockRbac struct {
	echo       *echo.Echo
	lock       sync.Mutex
	context    context.Context
	cancelFunc context.CancelFunc
	address    string
	waitGroup  *sync.WaitGroup
	port       string
	appName    string
	data       []Permission
}

// LoadProfile unmarshall a yaml content with a list
// of permissions to let to externalize the static
// contents.
// Return the list of strings, but currently no checks
// are made currently.
func LoadProfile(data []byte) []string {
	result := []string{}
	err := yaml.Unmarshal(data, &result)
	if err != nil {
		panic(err.Error())
	}
	// TODO The data are not being validated:
	// - every string is a tuple of 3 items separated by ':'.
	// - first item must match the service id, for instance 'idmsvc'.
	// - second item must match the allowed resources by the service.
	// - third item must match a valid verb.
	return result
}

func newRbacMockGuards(ctx context.Context, cfg *config.Config) {
	if ctx == nil {
		panic("ctx is nil")
	}
	if cfg == nil {
		panic("cfg is nil")
	}
	if !cfg.Application.EnableRBAC {
		panic("Config.Application.EnableRBAC is false")
	}
	if cfg.Clients.RbacBaseURL == "" {
		panic("Config.Clients.RbacBaseURL is an empty string")
	}
}

// NewRbacMock return a new rbac mock service for testing.
func NewRbacMock(ctx context.Context, cfg *config.Config) (service.ApplicationService, MockRbac) {
	var (
		cancelFunc  context.CancelFunc
		profileName string
	)
	urlData, err := url.Parse(cfg.Clients.RbacBaseURL)
	if err != nil {
		panic(fmt.Sprintf("error parsing rbac client url: %s", err.Error()))
	}
	address := fmt.Sprintf("%s:%s", urlData.Hostname(), urlData.Port())
	ctx, cancelFunc = context.WithCancel(ctx)
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(app_middleware.ContextLogConfig(&app_middleware.LogConfig{}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		// Request logger values for middleware.RequestLoggerValues
		LogError:  true,
		LogMethod: true,
		LogStatus: true,
		LogURI:    true,
		// Forwards error to the global error handler, so it can decide
		// appropriate status code.
		HandleError:   true,
		LogValuesFunc: logger.MiddlewareLogValues,
	}))
	m := &mockRbac{
		appName:    cfg.Application.Name,
		address:    address,
		echo:       e,
		context:    ctx,
		cancelFunc: cancelFunc,
		waitGroup:  &sync.WaitGroup{},
		lock:       sync.Mutex{},
	}
	e.GET(fmt.Sprintf("%s/access", urlData.Path), m.accessHandler)
	profileName = os.Getenv("APP_CLIENTS_RBAC_PROFILE")
	if profileName == "" {
		profileName = ProfileDomainAdmin
	}
	profileData, ok := Profiles[profileName]
	if !ok {
		slog.Error("not found", "profile_name", profileName)
		panic("rbac mock profile not found")
	}
	m.SetPermissions(profileData)
	return m, m
}

func (m *mockRbac) Start() error {
	m.echo.HideBanner = true
	m.echo.Debug = false
	m.echo.HidePort = false
	m.waitGroup.Add(2)
	go func() {
		defer m.waitGroup.Done()
		slog.Info("mock rbac service starting")
		if err := m.echo.Start(m.address); err != nil {
			if err != http.ErrServerClosed {
				slog.Error(err.Error())
			} else {
				slog.Info("Service rbac mock closed")
			}
			return
		}
	}()
	go func() {
		defer m.waitGroup.Done()
		defer m.cancelFunc()
		<-m.context.Done()
		if err := m.echo.Shutdown(m.context); err != nil {
			slog.Error(err.Error())
			return
		}
	}()
	return nil
}

func (m *mockRbac) Stop() error {
	slog.Info("mock rback service stopping")
	defer m.waitGroup.Wait()
	m.cancelFunc()
	return nil
}

// SetPermissions allow to dynamically assign the list
// of permissions that will return the mock service
// when it is reached out for accessin the acl list.
// data contains the information to be returned.
func (m *mockRbac) SetPermissions(data []string) {
	newData := make([]Permission, len(data))
	for i := range data {
		newData[i] = Permission{
			Permission:          data[i],
			ResourceDefinitions: []any{},
		}
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data = newData
}

// WaitAddress is a naive implementation to await the rbac
// mock has an address assigned.
func (m *mockRbac) WaitAddress(timeout time.Duration) error {
	isListening := false
	deadline := time.Now().Add(timeout)
	for time.Now().Compare(deadline) < 0 {
		if m.echo.Listener != nil && m.echo.Listener.Addr().String() != "" {
			slog.Info(fmt.Sprintf("rbac mock listening at: %s", m.echo.Listener.Addr().String()))
			isListening = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !isListening {
		if time.Now().Compare(deadline) >= 0 {
			return errRbacMockAwaitTimeout
		} else {
			return errRbacMockUnknown
		}
	}
	return nil
}

// GetBaseURL retrieve the base URL to reach out the
// rbac mock.
// Return empty string if the listener is not yet assigned;
// see WaitAddress method.
func (m *mockRbac) GetBaseURL() string {
	addr := ""
	if m.echo.Listener != nil {
		addr = m.echo.Listener.Addr().String()
		if addr != "" {
			return fmt.Sprintf("http://%s/api/rbac/v1", addr)
		}
	}

	return addr
}
