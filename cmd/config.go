package cmd

import (
	"log"

	"catalog-go/config"

	"github.com/spf13/cobra"
)

var updateConfigSchema = &cobra.Command{
	Use:   "update-config-schema",
	Short: "Update config schema",
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("config").Value.String()
		err := config.CONFIG.Load(path)
		if err != nil {
			log.Fatalf("[F] Failed to load config: %v\n", err)
		}
		err = config.CONFIG.Save(path)
		if err != nil {
			log.Fatalf("[F] Failed to save config: %v\n", err)
		}
		log.Printf("[I] Config schema updated\n")
	},
}

func loadConfig(path string) {
	if err := config.CONFIG.Load(path); err != nil {
		log.Fatalf("[F] Failed to load config: %v\n", err)
	}
}
