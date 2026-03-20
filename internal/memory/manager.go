package memory

import (
	"context"
	"fmt"
)

// MemoryManager implementa a estratégia híbrida de Tri-Banco.
type MemoryManager struct {
	sqlite   *SQLiteProvider
	supabase *SupabaseProvider
	qdrant   *QdrantProvider
}

func NewMemoryManager(sq *SQLiteProvider, sb *SupabaseProvider, qd *QdrantProvider) *MemoryManager {
	return &MemoryManager{
		sqlite:   sq,
		supabase: sb,
		qdrant:   qd,
	}
}

// GetContext busca o contexto consolidado das três camadas.
func (m *MemoryManager) GetContext(ctx context.Context, sessionID string, role string) (string, error) {
	// 1. Busca estado quente (SQLite)
	hot, err := m.sqlite.GetRecent(ctx, sessionID)
	// ... (Lógica de consolidação real aqui)
	return "Contexto Híbrido: SQLite + Supabase + Qdrant", nil
}

// Push ingere memórias de forma assíncrona para as camadas permanentes.
func (m *MemoryManager) Push(ctx context.Context, sessionID string, content string) error {
	// Escrita imediata no SQLite
	m.sqlite.Save(ctx, sessionID, content)

	// Ingestão assíncrona nas camadas Supabase e Qdrant (Passive Ingestion)
	go func() {
		m.supabase.StoreVector(context.Background(), sessionID, content)
		m.qdrant.LogExperience(context.Background(), sessionID, content)
	}()

	return nil
}

// Placeholders para os providers Reais
type SQLiteProvider struct{}
func (s *SQLiteProvider) GetRecent(ctx context.Context, sid string) ([]Message, error) { return nil, nil }
func (s *SQLiteProvider) Save(ctx context.Context, sid string, c string) {}

type SupabaseProvider struct{}
func (s *SupabaseProvider) GetHistory(ctx context.Context, sid string, limit int) ([]Message, error) { return nil, nil }
func (s *SupabaseProvider) StoreVector(ctx context.Context, sid string, c string) {}

type QdrantProvider struct{}
func (q *QdrantProvider) LogExperience(ctx context.Context, sid string, c string) {}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
