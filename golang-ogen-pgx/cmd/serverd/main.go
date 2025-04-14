package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize logger
	// logger.Setup()

	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	// dbPool, err := database.NewPgxPool(ctx)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to connect to database")
	// }
	// defer dbPool.Close()

	// Create repository
	// itemRepo := repository.NewItemRepository(dbPool)

	// Create service
	// svc := service.NewService(itemRepo)

	// Create API handler
	// handler, err := api.NewServer(svc)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to create API server")
	// }

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Create HTTP server
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		// Handler: handler,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Hello, World!"}`))
		}),
	}

	// Start server in a goroutine
	go func() {
		log.Info().Str("port", port).Msg("Starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Info().Msg("Shutting down server...")

	// Create a deadline to wait for
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited properly")
}
