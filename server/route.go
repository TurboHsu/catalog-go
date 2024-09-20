package server

import (
	"catalog-go/config"
	"catalog-go/server/cat"
	"catalog-go/server/static"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(addr string, ctx context.Context) {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.CONFIG.Server.AllowOrigins
	r.Use(cors.New(corsConfig))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "catalog-go",
		})
	})

	cat.ConfigureRoute(r)
	static.ConfigureRoute(r)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[F] Failed to start server: %v\n", err)
		}
	}()

	log.Printf("[I] Server started at %s\n", addr)

	<-ctx.Done()

	log.Printf("[I] Gracefully shutting down server\n")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[F] Failed to shutdown server: %v\n", err)
	}
	shutdownCancel()
}
