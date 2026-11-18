// Command convertix démarre l'API de conversion d'unités.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dhia-sel/convertix/internal/api"
)

const (
	defaultPort     = "8080"
	shutdownTimeout = 10 * time.Second
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	if err := run(context.Background(), ln); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, ln net.Listener) error {
	srv := &http.Server{
		Handler:           api.NewHandler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("serveur démarré sur %s", ln.Addr())
		serverErr <- srv.Serve(ln)
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("serveur : %w", err)
	case <-signalCtx.Done():
	}

	stop()
	log.Println("arrêt en cours...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("arrêt : %w", err)
	}

	log.Println("serveur arrêté proprement")
	return nil
}
