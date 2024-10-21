package cmd

import (
	"os"

	"log/slog"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var jwkRevokeCmd = &cobra.Command{
	Use:   "revoke [kid]",
	Short: "Revoke a JWK",
	Long:  `The revoke command marks a JWK as revoked.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		r := datastore.NewHostconfJwkDb(cfg, slog.Default())
		err := r.Revoke(args[0])
		if err != nil {
			slog.Error("Revoke failed", slog.String("error", err.Error()))
			os.Exit(2)
		} else {
			slog.Info("Done")
		}
	},
}

func init() {
	jwkCmd.AddCommand(jwkRevokeCmd)
}
