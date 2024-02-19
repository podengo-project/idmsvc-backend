package cmd

import (
	"os"

	"golang.org/x/exp/slog"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var jwkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List JWKs",
	Long:  `The list command shows all JWKs in the database.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logger.LogBuildInfo("db-tool")
		cfg := config.Get()
		logger.InitLogger(cfg)
		r := datastore.NewHostconfJwkDb(cfg)
		err := r.ListKeys()
		if err != nil {
			slog.Error("List failed", slog.String("error", err.Error()))
			os.Exit(2)
		} else {
			slog.Info("Done")
		}
	},
}

func init() {
	jwkCmd.AddCommand(jwkListCmd)
}
