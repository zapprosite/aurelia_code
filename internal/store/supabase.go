package store

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/kocar/aurelia/internal/config"
)

// SupabaseStore implementa a camada L3 (Galactic) de memória.
type SupabaseStore struct {
	url     string
	apiKey  string
	enabled bool
	client  *http.Client
}

func NewSupabaseStore(cfg *config.AppConfig) *SupabaseStore {
	return &SupabaseStore{
		url:     cfg.SupabaseURL,
		apiKey:  os.Getenv("SUPABASE_SERVICE_ROLE_KEY"), // Deve estar no .env
		enabled: cfg.SupabaseEnabled && cfg.SupabaseURL != "",
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SupabaseStore) SyncFact(ctx context.Context, fact map[string]any) error {
	if !s.enabled {
		return nil
	}
	// Logica de upsert via PostgREST do Supabase
	return nil
}

func (s *SupabaseStore) SyncNote(ctx context.Context, note map[string]any) error {
	if !s.enabled {
		return nil
	}
	return nil
}
