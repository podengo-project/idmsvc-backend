package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh or create JWKs",
	Long:  `The refresh command ensures that the database contains valid JWKs.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print not implemented warning, but don't fail.
		fmt.Println("JWK refresh is not implemented yet.")
		os.Exit(0)
	},
}

func init() {
	jwkCmd.AddCommand(refreshCmd)
}
