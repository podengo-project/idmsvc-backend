package cmd

import (
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database up or down",
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
