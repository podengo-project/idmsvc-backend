package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "db-tool",
	Short: "idmsvc database tool",
	Long: `The db-tool handles database migration and updates. For example:

	db-tool new [migration-name]
	db-tool migrate up [steps]
	db-tool migrate down [steps]
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
