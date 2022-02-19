package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"messenger-backend/database"
	"messenger-backend/models"

	"github.com/jackc/pgx/v4/log/zapadapter"
)

func main() {
	pgxLogLevel, err := database.LogLevelFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	pgPool, err := database.NewPgxPool(context.Background(), "",
		zapadapter.NewLogger(database.GetLogger()), pgxLogLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()

	go database.InitialSetup(pgPool)

	server := &http.Server{Addr: ":8080", Handler: server(models.NewService(
		&database.DB{Postgres: pgPool},
	))}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sig
		log.Println("Shutting down server...")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
