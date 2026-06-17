package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

//go:embed all:migrations
var migrationsFS embed.FS

// Store — тонкий слой репозиториев поверх PostgreSQL (tech-stack §2: один файл ≈ одна сущность).
type Store struct {
	db  *sql.DB
	log zerolog.Logger
}

func Open(dsn string, log zerolog.Logger) (*Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	return &Store{db: db, log: log}, nil
}

// WaitReady пингует БД с ретраями — postgres может стартовать дольше backend.
func (s *Store) WaitReady(attempts int) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = s.db.Ping(); err == nil {
			return nil
		}
		s.log.Info().Int("attempt", i+1).Msg("ожидание PostgreSQL...")
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("postgres недоступен: %w", err)
}

func (s *Store) Close() error { return s.db.Close() }

// Migrate применяет forward-only миграции из embed.FS, отслеживая schema_migrations.
func (s *Store) Migrate() error {
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
		return err
	}
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	for _, name := range files {
		var exists bool
		if err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, name).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		sqlBytes, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if _, err := s.db.Exec(`INSERT INTO schema_migrations(version) VALUES ($1)`, name); err != nil {
			return err
		}
		s.log.Info().Str("migration", name).Msg("migration applied")
	}
	return nil
}
