// Команда server — HTTP-бэкенд системы поиска дублей (Go-монолит, tech-stack §2).
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/roof/dupes/internal/config"
	"github.com/roof/dupes/internal/ext"
	"github.com/roof/dupes/internal/logger"
	"github.com/roof/dupes/internal/pipeline"
	"github.com/roof/dupes/internal/scoring"
	"github.com/roof/dupes/internal/storage"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	version := os.Getenv("GIT_TAG")
	if version == "" {
		version = "dev"
	}

	store, err := storage.Open(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal().Err(err).Msg("open db")
	}
	defer store.Close()
	if err := store.WaitReady(30); err != nil {
		log.Fatal().Err(err).Msg("db not ready")
	}
	if err := store.Migrate(); err != nil {
		log.Fatal().Err(err).Msg("migrate")
	}
	log.Info().Msg("migrations applied")

	score := scoring.Config{PhashThreshold: cfg.PhashThreshold}
	pipe := pipeline.New(store, score, cfg.MatchPersistThreshold, cfg.IngestConcurrency, cfg.BlobDir, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pipe.Start(ctx)

	handlers := &ext.Handlers{Store: store, Pipe: pipe, BlobDir: cfg.BlobDir, Version: version, Log: log}
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           ext.Router(handlers, log),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("backend слушает")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info().Msg("останавливаюсь...")
	cancel()
	shCtx, shCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shCancel()
	_ = srv.Shutdown(shCtx)
}
