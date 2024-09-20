package cmd

import (
	"catalog-go/database"
	"log"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("config").Value.String()
		loadConfig(path)

		if err := database.ConnectDatabase(); err != nil {
			log.Fatalf("[F] Failed to connect database: %v\n", err)
		}
		if err := database.MigrateDatabase(); err != nil {
			log.Fatalf("[F] Failed to migrate database: %v\n", err)
		}
		log.Printf("[I] Database migrated\n")
	},
}
