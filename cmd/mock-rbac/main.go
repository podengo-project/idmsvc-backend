package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	mock_rbac_impl "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/rbac/impl"
)

const component = "mock-rbac"

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
	logger.InitLogger(cfg, component)
	defer logger.DoneLogger()

	if cfg.Clients.RbacBaseURL == "" {
		panic("'RbacBaseURL' is empty")
	}
	srvRbac, mockRbac := mock_rbac_impl.NewRbacMock(ctx, cfg)

	profileName := os.Getenv("APP_CLIENTS_RBAC_PROFILE")
	if profileName == "" {
		profileName = mock_rbac_impl.ProfileDomainAdmin
	}
	if profileData, ok := mock_rbac_impl.Profiles[profileName]; ok {
		mockRbac.SetPermissions(profileData)
	} else {
		slog.Error("not found", "profile_name", profileName)
		panic("rbac mock profile not found")
	}

	if err := srvRbac.Start(); err != nil {
		panic(err)
	}
	<-ctx.Done()
	if err := srvRbac.Stop(); err != nil {
		panic(err)
	}
}
