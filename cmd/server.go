package cmd

import (
	"catalog-go/config"
	"catalog-go/database"
	"catalog-go/receiver/telegram"
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
	registerGin(&wg, ctx)
	registerTelegramBot(&wg, ctx)

	wg.Wait()

}

func registerGin(wg *sync.WaitGroup, ctx context.Context) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run(config.CONFIG.Server.Listen, ctx)
	}()
}

func registerTelegramBot(wg *sync.WaitGroup, ctx context.Context) {
	if !config.CONFIG.Receiver.TelegramBot.Enable {
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		telegram.Run(config.CONFIG.Receiver.TelegramBot.Token, ctx)
	}()
}
