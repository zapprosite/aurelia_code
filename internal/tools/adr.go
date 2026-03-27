package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/agent"
)

// CreateADRHandler generates a new ADR file in docs/adr/
func CreateADRHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	title, err := requireStringArg(args, "title")
	if err != nil {
		return "", err
	}
	status := optionalStringArg(args, "status")
	if status == "" {
		status = "PROPOSED"
	}
	contextText := optionalStringArg(args, "context")
	decision := optionalStringArg(args, "decision")

	now := time.Now()
	dateStr := now.Format("20060102")
	
	// Clean title for filename
	cleanTitle := strings.ToLower(title)
	cleanTitle = strings.ReplaceAll(cleanTitle, " ", "-")
	
	workdir, _ := agent.WorkdirFromContext(ctx)
	if workdir == "" {
		workdir = "." 
	}

	adrDir := filepath.Join(workdir, "docs/adr")
	filename := fmt.Sprintf("%s-%s.md", dateStr, cleanTitle)
	fullPath := filepath.Join(adrDir, filename)

	if err := os.MkdirAll(adrDir, 0755); err != nil {
		return "", fmt.Errorf("falha ao criar diretorio docs/adr: %w", err)
	}

	content := fmt.Sprintf(`# ADR %s: %s

## Status
%s

## Contexto
%s

## Decisão
%s

## Consequências
- [A preencher pelo arquiteto]
`, dateStr, title, status, contextText, decision)

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("falha ao escrever ADR: %w", err)
	}

	return fmt.Sprintf("ADR criado com sucesso: docs/adr/%s", filename), nil
}
