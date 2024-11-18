// Package main is the entry point for the pendo-client application.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	pendo "github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
	client_pendo "github.com/podengo-project/idmsvc-backend/internal/usecase/client/pendo"
	"github.com/spf13/cobra"
)

const component = "pendo-client"

func initPendo() (cfg *config.Config, ctx context.Context, pendoClient pendo.Pendo) {
	// Initialize dependencies
	logger.LogBuildInfo(component)
	cfg = config.Get()
	cfg.Logging.Level = "debug"
	logger.InitLogger(cfg, component)

	ctx = context.Background()
	ctx = app_context.CtxWithLog(ctx, slog.Default())

	log := app_context.LogFromCtx(ctx)

	// Initialize the pendo client
	log.Info("Starting pendo-client")
	log.Info(fmt.Sprintf("Pendo Base URL: %s", cfg.Clients.PendoBaseURL))
	log.Info(fmt.Sprintf("Pendo API Key: %s", cfg.Clients.PendoAPIKey))
	log.Info(fmt.Sprintf("Pendo Track Event Key: %s", cfg.Clients.PendoTrackEventKey))
	pendoClient = client_pendo.NewClient(cfg)

	return cfg, ctx, pendoClient
}

func sendPendoTrackEvent(ctx context.Context, client pendo.Pendo, event, accountID, visitorID string) {
	log := app_context.LogFromCtx(ctx)
	track := pendo.TrackRequest{
		AccountID: accountID,
		Type:      "track",
		Event:     event,
		VisitorID: visitorID,
		Timestamp: time.Now().UTC().UnixMilli(),
	}

	log.Debug(fmt.Sprintf("%v", track))

	err := client.SendTrackEvent(ctx, &track)
	if err != nil {
		log.Warn(err.Error())
	}
}

func createRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pendo-client",
		Short: "Pendo Client",
		Long: `pendo-client is a CLI tool to interact with the Pendo API.:

	export CLIENTS_PENDO_BASE_URL="https://app.pendo.io"
	export CLIENTS_PENDO_API_KEY="your-api-key"
	export CLIENTS_PENDO_TRACK_EVENT_KEY="your-track-event-secret-key"
	pendo-client track [name] [accountID] [visitorID]
`,
	}
}

func createTrackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "track [name] [accountID] [visitorID]",
		Short: "Create new Track Event",
		Args:  cobra.ExactArgs(3),
		Run: func(_ *cobra.Command, args []string) {
			_, ctx, pendoClient := initPendo()
			sendPendoTrackEvent(ctx, pendoClient, args[0], args[1], args[2])
		},
	}
}

func initCLI() *cobra.Command {
	rootCmd := createRootCmd()
	trackCmd := createTrackCmd()

	rootCmd.AddCommand(trackCmd)
	return rootCmd
}

func main() {
	rootCmd := initCLI()
	defer logger.DoneLogger()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
