package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"example.com/field-expedition-ledger/internal/config"
	"example.com/field-expedition-ledger/internal/httpapi"
	"example.com/field-expedition-ledger/internal/store"
)

func main() {
	settings := config.Load()
	repository := store.NewMemoryRepository()
	if settings.Seed {
		if err := store.SeedDemo(context.Background(), repository); err != nil {
			log.Fatalf("seed demo data: %v", err)
		}
	}
	server := &http.Server{
		Addr:              settings.Address,
		Handler:           httpapi.NewHandler(repository),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("field expedition ledger listening on %s", settings.Address)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}
