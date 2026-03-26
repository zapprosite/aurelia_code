package agent

import (
	"context"
	"strings"
	"testing"
)

type captureRuntimeProvider struct {
	systemPrompt string
	tools        []Tool
}

func (c *captureRuntimeProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	c.systemPrompt = systemPrompt
	c.tools = append([]Tool(nil), tools...)
	return &ModelResponse{Content: "ok"}, nil
}

type scopedAssemblerStub struct {
	query string
	botID string
}

func (s *scopedAssemblerStub) AssembleContext(ctx context.Context, query string) string {
	s.query = query
	return "fallback context"
}

func (s *scopedAssemblerStub) AssembleContextForBot(ctx context.Context, botID, query string) string {
	s.botID = botID
	s.query = query
	return "scoped context"
}

func TestLoopRun_InjectsMemoryContextForScopedBot(t *testing.T) {
	t.Parallel()

	provider := &captureRuntimeProvider{}
	registry := NewToolRegistry()
	registry.Register(Tool{Name: "run_command", Description: "runs commands"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})

	assembler := &scopedAssemblerStub{}
	loop := NewLoop(provider, registry, 1).WithMemoryAssembler(assembler)
	ctx := WithBotContext(context.Background(), "controle-db")

	_, _, err := loop.Run(ctx, "base prompt", []Message{{Role: "user", Content: "organize o banco"}}, []string{"run_command"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if assembler.botID != "controle-db" {
		t.Fatalf("expected scoped assembler botID controle-db, got %q", assembler.botID)
	}
	if assembler.query != "organize o banco" {
		t.Fatalf("expected assembler query to match user message, got %q", assembler.query)
	}
	if provider.systemPrompt == "" || !strings.Contains(provider.systemPrompt, "# MEMORY CONTEXT") || !strings.Contains(provider.systemPrompt, "scoped context") {
		t.Fatalf("expected memory block in system prompt, got %q", provider.systemPrompt)
	}
}

func TestLoopRun_UsesToolCatalogWhenAllowedToolsAreEmpty(t *testing.T) {
	t.Parallel()

	provider := &captureRuntimeProvider{}
	registry := NewToolRegistry()
	registry.Register(Tool{Name: "create_schedule", Description: "cria cron schedule para backup e rotina"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "docker_control", Description: "gerencia containers docker"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})

	loop := NewLoop(provider, registry, 1).WithToolCatalog(NewToolCatalog(registry), 1)

	_, _, err := loop.Run(context.Background(), "base prompt", []Message{{Role: "user", Content: "crie um schedule cron para backup"}}, nil)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(provider.tools) != 1 {
		t.Fatalf("expected 1 tool from catalog, got %d", len(provider.tools))
	}
	if provider.tools[0].Name != "create_schedule" {
		t.Fatalf("expected create_schedule, got %q", provider.tools[0].Name)
	}
}
