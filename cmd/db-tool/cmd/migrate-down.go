package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/spf13/cobra"
)

// downCmd represents the migrate down command
var downCmd = &cobra.Command{
	Use:   "down [steps]",
	Short: "Downgrade database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		steps, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "step %s is not an integer", args[0])
			os.Exit(2)
		}
		config := config.Get()
		err = datastore.MigrateDown(config, steps)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(downCmd)
}
