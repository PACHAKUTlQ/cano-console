package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func StartLocalServer() (*http.Server, string) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatalf("❌ Could not start listener: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", port)

	handler := http.FileServer(http.Dir(appDir))
	server := &http.Server{Handler: handler}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	return server, url
}
