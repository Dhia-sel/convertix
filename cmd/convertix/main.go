package main

import (
	"log"
)

func main() {
	
	
	// TODO 1 : lire PORT (os.Getenv), valeur par défaut "8080"
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	// TODO 2 : construire le http.Server (Addr, Handler, 4 délais)
	httpServer := &http.Server{
		Addr:    ":" + PORT,
		Handler: api.NewHandler(),
	}
	// TODO 3 : lancer ListenAndServe dans une goroutine
	listenandServeErr := make(chan error, 1)
	go func() {
		log.Printf("Serveur démarré sur le port %s", PORT)
		listenandServeErr <- httpServer.ListenAndServe()
	}()
	// TODO 4 : attendre SIGINT ou SIGTERM (signal.NotifyContext)
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// TODO 5 : arrêter proprement avec srv.Shutdown et un contexte avec délai
	srvShutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}