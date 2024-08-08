package impl

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	app_middleware "github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/service"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
)

var (
	errRbacMockAwaitTimeout = errors.New("timeout awaiting rbac mock to be ready")
	errRbacMockUnknown      = errors.New("unknown error happened on rbac mock")
)

type MockPendo interface {
	GetBaseURL() string
	WaitAddress(timeout time.Duration) error
	// TODO Add here additional methods
	GetMetrics() pendo.SetMetadataRequest
}

type mockPendo struct {
	echo       *echo.Echo
	lock       sync.RWMutex
	context    context.Context
	cancelFunc context.CancelFunc
	address    string
	waitGroup  *sync.WaitGroup
	port       string
	appName    string
	// TODO Add here additional fields
	metrics pendo.SetMetadataRequest // TODO This will change as we remove uncertainty
}

func newPendoMockGuards(ctx context.Context, cfg *config.Config) {
	if ctx == nil {
		panic("'ctx' is nil")
	}
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if cfg.Clients.PendoBaseURL == "" {
		panic("'Config.Clients.PendoBaseURL' is an empty string")
	}
	if cfg.Clients.PendoAPIKey == "" {
		panic("'Config.Clients.PendoAPIKey' is an empty string")
	}
}

// NewPendoMock return a new rbac mock service for testing.
func NewPendoMock(ctx context.Context, cfg *config.Config) (service.ApplicationService, MockPendo) {
	var cancelFunc context.CancelFunc
	newPendoMockGuards(ctx, cfg)
	urlData, err := url.Parse(cfg.Clients.PendoBaseURL)
	if err != nil {
		panic(fmt.Sprintf("error parsing pendo client url: %s", err.Error()))
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
	m := &mockPendo{
		appName:    cfg.Application.Name,
		address:    address,
		echo:       e,
		context:    ctx,
		cancelFunc: cancelFunc,
		waitGroup:  &sync.WaitGroup{},
		lock:       sync.RWMutex{},
	}
	e.POST(fmt.Sprintf("%s/:kind/:group/access", urlData.Path), m.CreateMetadataAccountCustomValue)
	e.GET(fmt.Sprintf("%s/:kind/:group/access", urlData.Path), m.GetMetadataAccountCustomValue)
	return m, m
}

func (m *mockPendo) Start() error {
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

func (m *mockPendo) Stop() error {
	slog.Info("mock rback service stopping")
	defer m.waitGroup.Wait()
	m.cancelFunc()
	return nil
}

// WaitAddress is a naive implementation to await the rbac
// mock has an address assigned.
func (m *mockPendo) WaitAddress(timeout time.Duration) error {
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
func (m *mockPendo) GetBaseURL() string {
	addr := ""
	if m.echo.Listener != nil {
		addr = m.echo.Listener.Addr().String()
		if addr != "" {
			return fmt.Sprintf("http://%s/api/rbac/v1", addr)
		}
	}

	return addr
}

// GetMetrics return a copy of the internal metrics
func (m *mockPendo) GetMetrics() pendo.SetMetadataRequest {
	m.lock.RLock()
	defer m.lock.RUnlock()
	output := make(pendo.SetMetadataRequest, 0, len(m.metrics))
	for i := range m.metrics {
		output[i] = m.metrics[i]
	}
	return output
}
