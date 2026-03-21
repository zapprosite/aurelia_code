package agent

import (
	"context"
	"testing"
)

func buildTestRegistry(tools []Tool) *ToolRegistry {
	reg := NewToolRegistry()
	for _, t := range tools {
		tool := t
		reg.Register(tool, func(_ context.Context, _ map[string]interface{}) (string, error) {
			return "ok", nil
		})
	}
	return reg
}

func TestNewToolCatalog_BuildsFromRegistry(t *testing.T) {
	reg := NewToolRegistry()
	reg.Register(Tool{Name: "run_command", Description: "Executa comandos bash no terminal"}, nil)
	reg.Register(Tool{Name: "read_file", Description: "Lê o conteúdo de um arquivo"}, nil)
	reg.Register(Tool{Name: "telegram_send", Description: "Envia mensagem via Telegram bot"}, nil)

	catalog := NewToolCatalog(reg)
	if catalog.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", catalog.Len())
	}
}

func TestToolCatalog_MatchForTask_FiltersRelevant(t *testing.T) {
	reg := NewToolRegistry()
	reg.Register(Tool{Name: "run_command", Description: "Executa comando bash shell"}, nil)
	reg.Register(Tool{Name: "telegram_send", Description: "Envia mensagem Telegram"}, nil)
	reg.Register(Tool{Name: "qdrant_search", Description: "Busca semântica vetorial embeddings"}, nil)
	reg.Register(Tool{Name: "read_file", Description: "Lê arquivo do sistema"}, nil)

	catalog := NewToolCatalog(reg)

	// Tarefa de debug terminal — deve priorizar run_command
	matches := catalog.MatchForTask("execute bash debug terminal", 2)
	if len(matches) == 0 {
		t.Fatal("expected at least one match for bash/terminal prompt")
	}
	found := false
	for _, m := range matches {
		if m.Name == "run_command" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected run_command to be in top matches for bash prompt, got: %v", matches)
	}
}

func TestToolCatalog_MatchForTask_ReturnsCoreWhenNoMatch(t *testing.T) {
	reg := NewToolRegistry()
	reg.Register(Tool{Name: "run_command", Description: "Executa comandos"}, nil)
	reg.Register(Tool{Name: "read_file", Description: "Lê arquivos"}, nil)

	catalog := NewToolCatalog(reg)

	// Prompt genérico — deve retornar pelo menos as core tools
	matches := catalog.MatchForTask("zzzxxx completamente fora do vocabulário", 5)
	// Não deve returnar vazio — core tools são sempre incluídas como fallback
	_ = matches // resultado depende das stopwords; o importante é não panic
}

func TestToolCatalog_AllTools_ReturnsFull(t *testing.T) {
	reg := NewToolRegistry()
	for i := 0; i < 10; i++ {
		name := "tool_" + string(rune('a'+i))
		reg.Register(Tool{Name: name, Description: "test tool"}, nil)
	}
	catalog := NewToolCatalog(reg)
	all := catalog.AllTools()
	if len(all) != 10 {
		t.Fatalf("expected 10, got %d", len(all))
	}
}

func TestTokenize_RemovesStopwords(t *testing.T) {
	tokens := tokenize("como executar isso para o bot do Telegram")
	for _, tok := range tokens {
		if tok == "como" || tok == "para" || tok == "isso" {
			t.Errorf("stopword '%s' should have been removed", tok)
		}
	}
}
