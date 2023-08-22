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
	"github.com/podengo-project/idmsvc-backend/internal/usecase/client"
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

func main() {
	wg := &sync.WaitGroup{}
	cfg := config.Get()
	logger.InitLogger(cfg)
	db := datastore.NewDB(cfg)
	defer datastore.Close(db)

	ctx, cancel := startSignalHandler(context.Background())
	inventory := client.NewHostInventory(cfg)
	s := impl_service.NewApplication(ctx, wg, cfg, db, inventory)
	if e := s.Start(); e != nil {
		panic(e)
	}
	<-ctx.Done()
	defer cancel()
	s.Stop()
	wg.Wait()
}
