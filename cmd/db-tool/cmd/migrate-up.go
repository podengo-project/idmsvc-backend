package cmd

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/spf13/cobra"
)

// upCmd represents the migrate up command
var upCmd = &cobra.Command{
	Use:   "up [steps]",
	Short: "Upgrade database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		steps, err := strconv.Atoi(args[0])
		if err != nil {
			slog.Error("step is not an integer", slog.String("step", args[0]))
			os.Exit(2)
		}
		config := config.Get()
		err = datastore.MigrateUp(config, steps)
		if err != nil {
			slog.Error(err.Error())
			panic(err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(upCmd)
}
