package cmd

import (
	"log/slog"
	"os"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var jwkRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh or create JWKs",
	Long:  `The refresh command ensures that the database contains valid JWKs.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		r := datastore.NewHostconfJwkDb(cfg, slog.Default())
		err := r.Refresh()
		if err != nil {
			slog.Error("Refresh failed", slog.String("error", err.Error()))
			os.Exit(2)
		} else {
			slog.Info("Done")
		}
	},
}

func init() {
	jwkCmd.AddCommand(jwkRefreshCmd)
}
