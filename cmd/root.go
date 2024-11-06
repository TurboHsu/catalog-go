package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "catalog-go",
	Short: "Cat logger supercharged",
	Long: `Cat logger supercharged with cat power`,
	Run: runServer,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "config.toml", "config file path")

	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(updateConfigSchema)
}
