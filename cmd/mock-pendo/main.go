package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	mock_pendo_impl "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/pendo/impl"
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
	ctx, cancel := startSignalHandler(context.Background())
	defer cancel()

	cfg := config.Get()
	logger.InitLogger(cfg)
	srvPendo, mockPendo := mock_pendo_impl.NewPendoMock(ctx, cfg)
	_ = mockPendo

	if err := srvPendo.Start(); err != nil {
		panic(err)
	}
	<-ctx.Done()
	if err := srvPendo.Stop(); err != nil {
		panic(err)
	}
}
