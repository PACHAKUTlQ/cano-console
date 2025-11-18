package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	forceUpdate := flag.Bool("update", false, "Force download of all assets, even if they exist.")
	flag.Parse()

	log.SetFlags(0)

	if err := UpdateAssets(*forceUpdate); err != nil {
		log.Fatalf("❌ Failed to update assets: %v", err)
	}

	server, url := StartLocalServer()
	log.Printf("🛰️  Local server started at: %s", url)

	// Graceful shutdown setup
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		log.Println("\n🔌 Signal received, shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("⚠️ Error during server shutdown: %v", err)
		}
	}()

	if err := LaunchSandboxedBrowser(url); err != nil {
		log.Fatalf("❌ Failed to launch browser: %v", err)
	}

	log.Println("🚪 Browser closed. Exiting.")
}
