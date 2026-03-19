package telegram

import (
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
	Kind   antigravityTaskKind
	Prompt string
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

	return &antigravityDelegation{
		Kind:   kind,
		Prompt: buildAntigravityPrompt(kind, text),
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

func buildAntigravityPrompt(kind antigravityTaskKind, userText string) string {
	extra := promptBodyForKind(kind)
	return fmt.Sprintf("Classifiquei esta tarefa como `light` para o chat do Antigravity.\n\n```text\nVoce esta no chat leve do Antigravity para o workspace /home/will/aurelia.\n\nTarefa:\n- %s\n\nRestricoes:\n- nao tocar em secrets\n- nao fazer deploy\n- nao declarar sucesso sem prova\n- manter a mudanca pequena e reversivel\n\nEntregue:\n- diff proposto ou passos exatos\n- comandos de validacao\n- risco residual em uma linha\n\n%s\n```\n", strings.TrimSpace(userText), extra)
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
