//go:build integration

package gateway

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
)

func ollamaURLForTest(t *testing.T) string {
	t.Helper()
	u := os.Getenv("OLLAMA_URL")
	if u == "" {
		u = "http://127.0.0.1:11434"
	}
	return u
}

func TestProvider_LocalLane_GenerateContent(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: -short skips")
	}

	cfg := &config.AppConfig{
		OllamaURL:    ollamaURLForTest(t),
		LLMProvider:  "ollama",
		LLMModel:     "gemma3:12b",
		DBPath:       t.TempDir() + "/gateway_integration.db",
		GroqSoftCapDaily: 800,
		GroqHardCapDaily: 1200,
	}

	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := p.GenerateContent(ctx,
		"Você é um assistente de teste. Responda em uma frase curta.",
		[]agent.Message{{Role: "user", Content: "diga apenas: OK"}},
		nil,
	)
	if err != nil {
		t.Fatalf("GenerateContent() error = %v (Ollama com gemma3:12b rodando em %s?)", err, cfg.OllamaURL)
	}
	if resp.Content == "" {
		t.Fatalf("esperava resposta não vazia")
	}

	t.Logf("Gateway local lane OK: response=%q tokens_in=%d tokens_out=%d",
		truncate(resp.Content, 80), resp.InputTokens, resp.OutputTokens)
}

func TestProvider_StatusSnapshot_LocalBreakerClosed(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: -short skips")
	}

	cfg := &config.AppConfig{
		OllamaURL: ollamaURLForTest(t),
		DBPath:    t.TempDir() + "/gateway_snap.db",
		GroqSoftCapDaily: 800,
		GroqHardCapDaily: 1200,
	}

	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}
	defer p.Close()

	snap := p.StatusSnapshot()
	if snap.PrimaryLane == "" {
		t.Fatalf("esperava PrimaryLane preenchido")
	}
	// Circuit breakers locais devem estar fechados em provider novo
	for key, state := range snap.Routes {
		if strings.HasPrefix(key, "local") && state.BreakerState == "open" {
			t.Fatalf("lane %q com breaker open em provider recém criado", key)
		}
	}
	t.Logf("StatusSnapshot OK: primary=%s routes=%d", snap.PrimaryLane, len(snap.Routes))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
