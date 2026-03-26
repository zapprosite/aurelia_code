// Package store provides connections to external persistent stores (Supabase/Postgres).
// SQLite remains the primary runtime store — this package is for long-lived,
// queryable data that benefits from Postgres (knowledge_items, component_status).
package store

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SupabaseStore holds a pgx connection pool to a Supabase (Postgres) instance.
type SupabaseStore struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// Connect creates a new connection pool. Caller must call Close() when done.
func Connect(ctx context.Context, dsn string, logger *slog.Logger) (*SupabaseStore, error) {
	if dsn == "" {
		return nil, fmt.Errorf("supabase: DSN is required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("supabase: connect: %w", err)
	}
	s := &SupabaseStore{pool: pool, logger: logger}
	if err := s.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("supabase: ping after connect: %w", err)
	}
	logger.Info("supabase connected", slog.String("dsn_prefix", dsnPrefix(dsn)))
	return s, nil
}

// Ping verifies the connection is alive.
func (s *SupabaseStore) Ping(ctx context.Context) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire: %w", err)
	}
	defer conn.Release()
	return conn.Conn().Ping(ctx)
}

// Close releases all pool connections.
func (s *SupabaseStore) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

// Pool returns the underlying pgxpool for advanced usage.
func (s *SupabaseStore) Pool() *pgxpool.Pool {
	return s.pool
}

// Migrate runs the initial schema creation (idempotent).
func (s *SupabaseStore) Migrate(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, schemaSQL)
	if err != nil {
		return fmt.Errorf("supabase: migrate: %w", err)
	}
	s.logger.Info("supabase schema OK")
	return nil
}

func dsnPrefix(dsn string) string {
	if len(dsn) > 30 {
		return dsn[:30] + "..."
	}
	return dsn
}
