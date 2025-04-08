package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(ctx context.Context, server *http.Server) error {
	const readHeaderTimeout = 5 * time.Second
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("Starting server at %s", server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal (context canceled)")
	case sig := <-stop:
		log.Printf("Received shutdown signal: %s", sig)
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), readHeaderTimeout)
	defer cancel()

	log.Println("Shutting down server...")

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("Server gracefully stopped")
	return nil
}
