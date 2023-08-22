package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/router"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type metricsService struct {
	context   context.Context
	cancel    context.CancelFunc
	waitGroup *sync.WaitGroup
	config    *config.Config

	echo   *echo.Echo
	logger zerolog.Logger
}

func NewMetrics(ctx context.Context, wg *sync.WaitGroup, config *config.Config, app handler.Application) service.ApplicationService {
	if config == nil {
		panic("config is nil")
	}
	if wg == nil {
		panic("wg is nil")
	}

	routerConfig := router.RouterConfig{
		Version:     "v1.0",
		MetricsPath: config.Metrics.Path,
		PrivatePath: "/private",
		Handlers:    app,
	}

	result := &metricsService{}
	result.context, result.cancel = context.WithCancel(ctx)
	result.waitGroup = wg
	result.logger = zerolog.New(os.Stderr)
	result.config = config

	result.echo = router.NewRouterForMetrics(
		echo.New(),
		routerConfig,
	)
	result.echo.HideBanner = true

	result.echo.Pre(middleware.RemoveTrailingSlash())
	if config.Logging.Level == "debug" {
		result.echo.Use(middleware.Logger())
	}
	result.echo.Use(middleware.Recover())

	if result.config.Logging.Level == "debug" {
		routes := result.echo.Routes()
		log.Debug().Msg("Printing metrics routes")
		for idx, route := range routes {
			if data, err := json.Marshal(route); err == nil {
				log.Debug().Int("idx", idx).RawJSON("route", []byte(data)).Send()
			}
		}
	}

	return result
}

func (srv *metricsService) Start() error {
	srv.waitGroup.Add(2)
	go func() {
		defer srv.waitGroup.Done()
		srvAddress := fmt.Sprintf(":%d", srv.config.Metrics.Port)
		log.Debug().Msgf("metrics: srvAddress=%s", srvAddress)
		if err := srv.echo.Start(srvAddress); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Failed to start metricsService")
		}
	}()

	go func() {
		defer srv.waitGroup.Done()
		defer srv.cancel()
		<-srv.context.Done()
		srv.logger.Info().Msg("Shutting down metricsService")
		if err := srv.echo.Shutdown(context.Background()); err != nil {
			srv.logger.Error().Err(err).Msg("error shuttingdown metricsService")
		}
	}()

	return nil
}

func (srv *metricsService) Stop() error {
	srv.cancel()
	return nil
}
