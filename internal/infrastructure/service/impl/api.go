package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/router"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/service"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type apiService struct {
	context   context.Context
	cancel    context.CancelFunc
	waitGroup *sync.WaitGroup
	config    *config.Config

	echo   *echo.Echo
	logger zerolog.Logger
}

func NewApi(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, app handler.Application, metrics *metrics.Metrics) service.ApplicationService {
	if cfg == nil {
		panic("config is nil")
	}
	if wg == nil {
		panic("wg is nil")
	}

	result := &apiService{}
	result.context, result.cancel = context.WithCancel(ctx)
	result.waitGroup = wg
	result.logger = logger.NewApiLogger(cfg)
	result.config = cfg
	routerConfig := router.RouterConfig{
		Version:            "1.0",
		PublicPath:         "/api/idmsvc",
		PrivatePath:        "/private",
		Handlers:           app,
		Metrics:            metrics,
		EnableAPIValidator: cfg.Application.ValidateAPI,
	}
	if cfg.Application.AcceptXRHFakeIdentity {
		routerConfig.IsFakeEnabled = true
	}
	result.echo = router.NewRouterWithConfig(
		echo.New(),
		routerConfig,
	)
	result.echo.HideBanner = true
	if result.config.Logging.Level == "debug" || result.config.Logging.Level == "trace" {
		result.echo.Debug = true
	}
	if result.config.Logging.Level == "debug" || result.config.Logging.Level == "trace" {
		routes := result.echo.Routes()
		log.Debug().Msg("Printing routes")
		for idx, route := range routes {
			if data, err := json.Marshal(route); err == nil {
				log.Debug().Int("idx", idx).RawJSON("route", []byte(data)).Send()
			}
		}
	}

	return result
}

func (srv *apiService) Start() error {
	srv.waitGroup.Add(2)
	go func() {
		defer srv.waitGroup.Done()
		srvAddress := fmt.Sprintf(":%d", srv.config.Web.Port)
		log.Debug().Msgf("srvAddress=%s", srvAddress)
		if err := srv.echo.Start(srvAddress); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Failed to start server")
		}
	}()

	go func() {
		defer srv.waitGroup.Done()
		defer srv.cancel()
		<-srv.context.Done()
		srv.logger.Info().Msg("Shutting down apiService")
		if err := srv.echo.Shutdown(context.Background()); err != nil {
			srv.logger.Error().Err(err).Msg("error shuttingdown apiService")
		}
	}()

	return nil
}

func (srv *apiService) Stop() error {
	srv.cancel()
	return nil
}
