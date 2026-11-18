package main

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

const testPath = "/healthz"

func TestRun_ServesThenShutsDown(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- run(ctx, ln) }()

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + ln.Addr().String() + testPath)
	if err != nil {
		t.Fatalf("requête : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("statut = %d, attendu %d", resp.StatusCode, http.StatusOK)
	}

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("run a retourné %v, attendu nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("run ne s'est pas arrêté après l'annulation du contexte")
	}
}

func TestRun_ReturnsErrorWhenListenerClosed(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ln.Close()

	if err := run(context.Background(), ln); err == nil {
		t.Fatal("run aurait dû retourner une erreur avec un listener fermé")
	}
}
