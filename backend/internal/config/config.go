// Package config читает env в типизированный конфиг (tech-stack §2: всё за env с дефолтом).
package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                  string // HTTP-порт
	DatabaseURL           string // DSN PostgreSQL
	BlobDir               string // локальное хранилище байтов изображений
	LogLevel              string // zerolog level
	IngestConcurrency     int    // параллелизм воркера обработки
	PhashThreshold        int    // макс. Hamming-расстояние для «почти-дубля»
	MatchPersistThreshold int    // минимальный risk_score, чтобы сохранить совпадение
	OTelEndpoint          string // пусто = трейсы выключены (observability §6)
}

func Load() Config {
	return Config{
		Port:                  env("PORT", "8080"),
		DatabaseURL:           env("DATABASE_URL", "postgres://roof:roof@localhost:5432/roof_dupes?sslmode=disable"),
		BlobDir:               env("BLOB_DIR", "/data/blobs"),
		LogLevel:              env("LOG_LEVEL", "info"),
		IngestConcurrency:     envInt("INGEST_CONCURRENCY", 4),
		PhashThreshold:        envInt("PHASH_THRESHOLD", 10),
		MatchPersistThreshold: envInt("MATCH_PERSIST_THRESHOLD", 31),
		OTelEndpoint:          env("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
