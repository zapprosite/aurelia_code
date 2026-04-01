package memory

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/observability"
)

// QdrantUpserter define a interface para upsert no Qdrant
type QdrantUpserter interface {
	UpsertKnowledge(ctx context.Context, text string) error
}

// LLMClient define a interface para geração de resumos
type LLMClient interface {
	GenerateSummary(ctx context.Context, text string) (string, error)
}

// DreamConsolidator analisa conversas antigas e as consolida em conhecimento longo.
type DreamConsolidator struct {
	db       *sql.DB
	cfg      *config.AppConfig
	qdrant   QdrantUpserter
	llm      LLMClient
	mu       sync.Mutex
}

// NewDreamConsolidator cria a estrutura.
func NewDreamConsolidator(db *sql.DB, cfg *config.AppConfig, qdrant QdrantUpserter, llm LLMClient) *DreamConsolidator {
	return &DreamConsolidator{
		db:       db,
		cfg:      cfg,
		qdrant:   qdrant,
		llm:      llm,
	}
}

// Start inicia o processo em uma goroutine se habilitado e bloqueia múltiplas runs.
func (dc *DreamConsolidator) Start(ctx context.Context) {
	if dc.cfg == nil || !dc.cfg.Features.DreamEnabled {
		return
	}
	go dc.runLoop(ctx)
}

func (dc *DreamConsolidator) runLoop(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dc.consolidate(ctx)
		}
	}
}

func (dc *DreamConsolidator) consolidate(ctx context.Context) {
	logger := observability.Logger("memory.dream")
	
	if !dc.mu.TryLock() {
		logger.Debug("consolidation already running, skipping")
		return
	}
	defer dc.mu.Unlock()

	if !dc.cfg.Features.DreamEnabled {
		return
	}

	// Gate 1: 24h desde o último run persisting in SQLite
	var lastRun time.Time
	err := dc.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(created_at), '2000-01-01') FROM system_events WHERE event_type = 'dream_run'").Scan(&lastRun)
	if err != nil && err != sql.ErrNoRows {
		logger.Warn("dream failed to get last run", slog.Any("err", err))
		return
	}

	if time.Since(lastRun) < 24*time.Hour {
		return
	}

	// Gate 2: 5+ sessões novas desde o último run
	var sessionCount int
	err = dc.db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT conversation_id) FROM messages WHERE created_at > ?", lastRun).Scan(&sessionCount)
	if err != nil {
		logger.Warn("dream failed to count new sessions", slog.Any("err", err))
		return
	}

	if sessionCount < 5 {
		return
	}

	// Ação: Busca mensagens antigas > 7 dias
	limitDate := time.Now().Add(-7 * 24 * time.Hour)
	rows, err := dc.db.QueryContext(ctx, "SELECT content FROM messages WHERE created_at < ? ORDER BY created_at ASC LIMIT 1000", limitDate)
	if err != nil {
		logger.Warn("dream failed to fetch old messages", slog.Any("err", err))
		return
	}
	defer rows.Close()

	var allText string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err == nil {
			allText += content + "\n"
		}
	}

	// Gera resumo via LLM (ops-cron alias simulado aqui pelo LLMClient com context com timeout e sem stream)
	if allText != "" && dc.llm != nil {
		summaryCtx, cancel := context.WithTimeout(ctx, 3*time.Minute)
		summary, err := dc.llm.GenerateSummary(summaryCtx, allText)
		cancel()
		if err == nil && dc.qdrant != nil {
			// Upsert no Qdrant como knowledge
			_ = dc.qdrant.UpsertKnowledge(ctx, summary)
		}
	}

	// Persiste timestamp do run
	_, _ = dc.db.ExecContext(ctx, "INSERT INTO system_events (event_type, created_at) VALUES ('dream_run', ?)", time.Now())
}
