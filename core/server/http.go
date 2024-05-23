package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"kiit-lab-engine/core/db"
	"kiit-lab-engine/routes"
)

func StartServer() error {
	// Load configuration from .env file
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// Set Gin mode based on configuration
	ginMode := viper.GetString("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize DB client
	dbClient := db.NewClient()
	if err := dbClient.Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbClient.Disconnect()

	// Initialize Gin router
	r := gin.Default()
	r.SetTrustedProxies(nil)
	routes.InitRoutes(r, dbClient)

	// x default
	port := viper.GetString("PORT")
	if port == "" {
		port = "8421"
	}

	// Create the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// Create a context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// Create a context with a timeout to allow ongoing requests to finish
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
	return nil
}
