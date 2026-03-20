package telegram

import (
	"encoding/json"
	"fmt"
	"strings"
)

type antigravityTaskKind string

const (
	antigravityTaskResearch antigravityTaskKind = "research"
	antigravityTaskConfig   antigravityTaskKind = "config"
	antigravityTaskDiff     antigravityTaskKind = "diff"
)

type antigravityDelegation struct {
	Kind    antigravityTaskKind
	Prompt  string
	Request antigravityHandoffRequest
}

type antigravityHandoffRequest struct {
	Version      string              `json:"version"`
	Kind         antigravityTaskKind `json:"kind"`
	Workspace    string              `json:"workspace"`
	TaskClass    string              `json:"task_class"`
	UserText     string              `json:"user_text"`
	Constraints  []string            `json:"constraints"`
	Deliverables []string            `json:"deliverables"`
}

type antigravityHandoffResult struct {
	Status       string   `json:"status"`
	Summary      string   `json:"summary"`
	ProposedDiff string   `json:"proposed_diff,omitempty"`
	Commands     []string `json:"commands,omitempty"`
	Validation   []string `json:"validation,omitempty"`
	ResidualRisk string   `json:"residual_risk"`
}

var antigravityHighRiskTerms = []string{
	"senha",
	"secret",
	"token",
	"api key",
	"api_key",
	"deploy",
	"produção",
	"producao",
	"firewall",
	"tailscale",
	"cloudflare",
	"merge",
	"rebase",
	"migrate",
	"migration",
	"drop ",
	" rm ",
	"delete ",
}

var antigravityDiffTerms = []string{
	"review diff",
	"revisar diff",
	"revise diff",
	"analisa diff",
	"analisar diff",
	"veja o diff",
}

var antigravityConfigTerms = []string{
	"config",
	"configur",
	"ajuste",
	"ajustar",
	"json",
	"yaml",
	"toml",
	"monta um curl",
	"monte um curl",
	"monta curl",
	"localiza",
	"localize",
	"flag",
	"env ",
}

var antigravityResearchTerms = []string{
	"pesquise",
	"pesquisa",
	"estude",
	"brainstorm",
	"compare",
	"explique",
	"resuma",
	"investigue",
}

func maybeBuildAntigravityDelegationPrompt(text string) *antigravityDelegation {
	normalized := normalizeAntigravityText(text)
	if normalized == "" {
		return nil
	}
	if len(normalized) > 500 {
		return nil
	}
	for _, term := range antigravityHighRiskTerms {
		if strings.Contains(normalized, term) {
			return nil
		}
	}

	kind := classifyAntigravityTaskKind(normalized)
	if kind == "" {
		return nil
	}

	req := buildAntigravityHandoffRequest(kind, text)
	return &antigravityDelegation{
		Kind:    kind,
		Prompt:  buildAntigravityPrompt(req),
		Request: req,
	}
}

func normalizeAntigravityText(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	text = strings.ReplaceAll(text, "\n", " ")
	return strings.Join(strings.Fields(text), " ")
}

func classifyAntigravityTaskKind(text string) antigravityTaskKind {
	for _, term := range antigravityDiffTerms {
		if strings.Contains(text, term) {
			return antigravityTaskDiff
		}
	}
	for _, term := range antigravityConfigTerms {
		if strings.Contains(text, term) {
			return antigravityTaskConfig
		}
	}
	for _, term := range antigravityResearchTerms {
		if strings.Contains(text, term) {
			return antigravityTaskResearch
		}
	}
	return ""
}

func buildAntigravityHandoffRequest(kind antigravityTaskKind, userText string) antigravityHandoffRequest {
	return antigravityHandoffRequest{
		Version:   "2026-03-19",
		Kind:      kind,
		Workspace: "/home/will/aurelia",
		TaskClass: "light",
		UserText:  strings.TrimSpace(userText),
		Constraints: []string{
			"nao tocar em secrets",
			"nao fazer deploy",
			"nao declarar sucesso sem prova",
			"manter a mudanca pequena e reversivel",
		},
		Deliverables: []string{
			"diff proposto ou passos exatos",
			"comandos de validacao",
			"risco residual em uma linha",
		},
	}
}

func buildAntigravityPrompt(req antigravityHandoffRequest) string {
	extra := promptBodyForKind(req.Kind)
	requestJSON, _ := json.MarshalIndent(req, "", "  ")
	return fmt.Sprintf("Classifiquei esta tarefa como `light` para o chat do Antigravity.\n\n```text\nVoce esta no chat leve do Antigravity para o workspace /home/will/aurelia.\n\nHandoff request:\n%s\n\n%s\n\nResposta obrigatoria em JSON:\n{\n  \"status\": \"approved|revise|blocked\",\n  \"summary\": \"\",\n  \"proposed_diff\": \"\",\n  \"commands\": [\"\"],\n  \"validation\": [\"\"],\n  \"residual_risk\": \"\"\n}\n```\n", string(requestJSON), extra)
}

func promptBodyForKind(kind antigravityTaskKind) string {
	switch kind {
	case antigravityTaskConfig:
		return "Instrucao adicional:\n- Revise a configuracao alvo e proponha a menor mudanca possivel.\n- Nao reestruture o projeto.\n- Se houver ambiguidade, liste no maximo 2 opcoes e recomende 1."
	case antigravityTaskDiff:
		return "Instrucao adicional:\n- Leia o diff ou a mudanca alvo e diga o que mudou.\n- Aponte o risco principal.\n- Diga a validacao minima necessaria."
	case antigravityTaskResearch:
		return "Instrucao adicional:\n- Pesquise apenas o necessario para responder ao ponto tecnico.\n- Nao escreva tutorial longo.\n- Responda com conclusao, prova e proximo comando ou diff."
	default:
		return ""
	}
}

func parseAntigravityHandoffResult(text string) (*antigravityHandoffResult, error) {
	trimmed := extractJSONPayload(strings.TrimSpace(text))
	if trimmed == "" {
		return nil, fmt.Errorf("empty handoff result")
	}
	var result antigravityHandoffResult
	if err := json.Unmarshal([]byte(trimmed), &result); err != nil {
		return nil, fmt.Errorf("decode handoff result: %w", err)
	}
	switch result.Status {
	case "approved", "revise", "blocked":
	default:
		return nil, fmt.Errorf("invalid handoff status %q", result.Status)
	}
	if strings.TrimSpace(result.Summary) == "" {
		return nil, fmt.Errorf("handoff summary is required")
	}
	if strings.TrimSpace(result.ResidualRisk) == "" {
		return nil, fmt.Errorf("handoff residual risk is required")
	}
	return &result, nil
}

func maybeParseAntigravityHandoffResult(text string) *antigravityHandoffResult {
	result, err := parseAntigravityHandoffResult(text)
	if err != nil {
		return nil
	}
	return result
}

func formatAntigravityHandoffResult(result *antigravityHandoffResult) string {
	if result == nil {
		return ""
	}
	var b strings.Builder
	switch result.Status {
	case "approved":
		b.WriteString("Handoff do Antigravity: aprovado.\n\n")
	case "revise":
		b.WriteString("Handoff do Antigravity: precisa de revisao.\n\n")
	case "blocked":
		b.WriteString("Handoff do Antigravity: bloqueado.\n\n")
	}
	b.WriteString(result.Summary)
	if len(result.Commands) > 0 {
		b.WriteString("\n\nComandos sugeridos:\n")
		for _, cmd := range result.Commands {
			cmd = strings.TrimSpace(cmd)
			if cmd == "" {
				continue
			}
			b.WriteString("- `")
			b.WriteString(cmd)
			b.WriteString("`\n")
		}
	}
	if len(result.Validation) > 0 {
		b.WriteString("\nValidacao minima:\n")
		for _, item := range result.Validation {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			b.WriteString("- ")
			b.WriteString(item)
			b.WriteString("\n")
		}
	}
	if strings.TrimSpace(result.ProposedDiff) != "" {
		b.WriteString("\nDiff proposto:\n```diff\n")
		b.WriteString(strings.TrimSpace(result.ProposedDiff))
		b.WriteString("\n```")
	}
	if strings.TrimSpace(result.ResidualRisk) != "" {
		b.WriteString("\n\nRisco residual: ")
		b.WriteString(strings.TrimSpace(result.ResidualRisk))
	}
	return strings.TrimSpace(b.String())
}

func extractJSONPayload(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) >= 3 && strings.HasPrefix(lines[0], "```") && strings.TrimSpace(lines[len(lines)-1]) == "```" {
			return strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))
		}
	}
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start >= 0 && end > start {
		return strings.TrimSpace(text[start : end+1])
	}
	return text
}
