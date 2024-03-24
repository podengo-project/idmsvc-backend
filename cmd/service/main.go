package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	impl_service "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl"
	mock_rbac_impl "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/rbac/impl"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/rbac"
	client_inventory "github.com/podengo-project/idmsvc-backend/internal/usecase/client/inventory"
	client_rbac "github.com/podengo-project/idmsvc-backend/internal/usecase/client/rbac"
	"golang.org/x/exp/slog"
)

func startSignalHandler(c context.Context) (context.Context, context.CancelFunc) {
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

// initRbacWrapper initialize the client wrapper to communicate with
// rbac microservice.
func initRbacWrapper(ctx context.Context, cfg *config.Config) rbac.Rbac {
	if env, ok := os.LookupEnv("ENV_NAME"); ok && env == "local" {
		if cfg.Application.EnableRBAC && cfg.Clients.RbacBaseURL != "" {
			// Only for local environment and if it is enabled
			srvRbac, mockRbac := mock_rbac_impl.NewRbacMock(ctx, cfg)
			srvRbac.Start()
			if err := mockRbac.WaitAddress(3 * time.Second); err != nil {
				panic(err.Error())
			}
			slog.Info("rbac mock listening", "rbac_mock_base_url", mockRbac.GetBaseURL())
		} else {
			slog.Warn("local env running with no rbac mock")
		}
	}

	// Initialize the rbac client wrapper
	rbacClient, err := client_rbac.NewClient("idmsvc", client_rbac.WithBaseURL(cfg.Clients.RbacBaseURL))
	if err != nil {
		panic(err)
	}
	rbac := client_rbac.New(cfg.Clients.RbacBaseURL, rbacClient)
	return rbac
}

func main() {
	wg := &sync.WaitGroup{}
	logger.LogBuildInfo("idmscv-backend")
	cfg := config.Get()
	logger.InitLogger(cfg)
	db := datastore.NewDB(cfg)
	defer datastore.Close(db)

	ctx, cancel := startSignalHandler(context.Background())
	inventory := client_inventory.NewHostInventory(cfg)
	rbac := initRbacWrapper(ctx, cfg)
	s := impl_service.NewApplication(ctx, wg, cfg, db, inventory, rbac)
	if e := s.Start(); e != nil {
		panic(e)
	}
	<-ctx.Done()
	defer cancel()
	s.Stop()
	wg.Wait()
}
