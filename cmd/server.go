package cmd

import (
	"catalog-go/config"
	"catalog-go/database"
	"catalog-go/server"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

func runServer(cmd *cobra.Command, args []string) {
	path := cmd.Flag("config").Value.String()
	loadConfig(path)
	if err := database.ConnectDatabase(); err != nil {
		log.Fatalf("[F] Failed to connect database: %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()

		// Clean up
		err := database.CloseDatabase()
		if err != nil {
			log.Fatalf("[F] Failed to close database: %v\n", err)
		}
	}()

	// Start all services
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run(config.CONFIG.Server.Listen, ctx)
	}()

	wg.Wait()

}
