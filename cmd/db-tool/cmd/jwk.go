package cmd

import (
	"github.com/spf13/cobra"
)

// jwkCmd represents the jwk command
var jwkCmd = &cobra.Command{
	Use:   "jwk",
	Short: "JSON Web Key management",
}

func init() {
	rootCmd.AddCommand(jwkCmd)
}
