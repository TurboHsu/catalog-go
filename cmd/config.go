package cmd

import (
	"log"

	"catalog-go/config"
)

func loadConfig(path string) {
	if err := config.CONFIG.Load(path); err != nil {
		log.Fatalf("[F] Failed to load config: %v\n", err)
	}
}
