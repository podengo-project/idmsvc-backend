package cmd

import (
	"github.com/spf13/cobra"
)

// jwkCmd represents the jwk command
var jwkCmd = &cobra.Command{
	Use:   "jwk",
	Short: "Hostconf JWK management",
}

func init() {
	rootCmd.AddCommand(jwkCmd)
}
