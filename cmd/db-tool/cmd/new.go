package cmd

import (
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create new DB migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := datastore.CreateMigrationFile(args[0])
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
