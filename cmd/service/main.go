package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	impl_service "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/rbac"
	client_pendo "github.com/podengo-project/idmsvc-backend/internal/usecase/client/pendo"
	client_rbac "github.com/podengo-project/idmsvc-backend/internal/usecase/client/rbac"
)

const component = "service"

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
	logger.LogBuildInfo(component)
	cfg := config.Get()
	logger.InitLogger(cfg, component)
	defer logger.DoneLogger()

	db := datastore.NewDB(cfg)
	defer datastore.Close(db)

	ctx, cancel := startSignalHandler(context.Background())
	rbac := initRbacWrapper(ctx, cfg)
	pendo := client_pendo.NewClient(cfg)
	s := impl_service.NewApplication(ctx, wg, cfg, db, rbac, pendo)
	if e := s.Start(); e != nil {
		panic(e)
	}
	<-ctx.Done()
	defer cancel()
	s.Stop()
	wg.Wait()
}
