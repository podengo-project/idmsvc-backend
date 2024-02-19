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
var jwkPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge expired JWK",
	Long:  `The purge command removes all expired JWK from the database.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logger.LogBuildInfo("db-tool")
		cfg := config.Get()
		logger.InitLogger(cfg)
		r := datastore.NewHostconfJwkDb(cfg)
		err := r.Purge()
		if err != nil {
			slog.Error("Purge failed", slog.String("error", err.Error()))
			os.Exit(2)
		} else {
			slog.Info("Done")
		}
	},
}

func init() {
	jwkCmd.AddCommand(jwkPurgeCmd)
}
