package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func run() int {
	var configPath string
	flag.StringVar(&configPath, "config", "config/publicapi.yaml", "Path to config file")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Starting publicapi service...")

	if err := RunApp(ctx, configPath); err != nil {
		log.Printf("Service failed: %v", err)
		return 1
	}

	log.Println("Service stopped gracefully")
	return 0
}

func main() {
	os.Exit(run())
}
