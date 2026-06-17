// Package logger настраивает zerolog (observability §5: JSON, уровень через LOG_LEVEL, без PII).
package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(level string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return zerolog.New(os.Stdout).Level(lvl).With().Timestamp().Str("service", "roof-dupes").Logger()
}
