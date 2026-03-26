package telegram

import (
	"strings"

	"github.com/kocar/aurelia/internal/config"
)

const (
	aureliaSovereignBotID   = "aurelia"
	aureliaSovereignBotName = "Aurelia_Code"
	markdown2026Contract    = `## Saida Obrigatoria
- Entregue texto em Markdown limpo, curto e direto.
- Use listas planas; nao use JSON, YAML ou blocos soltos como resposta final.
- Use tabela apenas quando houver comparacao real entre multiplos itens.
- Se algo nao puder ser confirmado, diga isso explicitamente em uma linha objetiva.`
)

type botGovernanceProfile struct {
	authorityLevel    string
	supervisorBotID   string
	supervisorBotName string
	domainScope       string
	requiredProvider  string
	requiredModel     string
	toolAllowlist     []string
	silentWhenHealthy bool
}

func EffectiveBotLLM(botCfg config.BotConfig, appProvider, appModel string) (string, string) {
	profile := governanceProfileForBot(botCfg.ID)
	if profile.requiredProvider != "" || profile.requiredModel != "" {
		return firstNonEmpty(profile.requiredProvider, botCfg.LLMProvider, appProvider), firstNonEmpty(profile.requiredModel, botCfg.LLMModel, appModel)
	}
	return firstNonEmpty(botCfg.LLMProvider, appProvider), firstNonEmpty(botCfg.LLMModel, appModel)
}

func governanceProfileForBot(botID string) botGovernanceProfile {
	switch strings.TrimSpace(strings.ToLower(botID)) {
	case "aurelia":
		return botGovernanceProfile{
			authorityLevel:  "sovereign",
			domainScope:     "comando operacional do ecossistema inteiro, consolidacao multi-bot e arbitragem final",
			toolAllowlist:   cloneToolList(defaultConversationTools),
			supervisorBotID: "",
		}
	case "controle-db":
		return botGovernanceProfile{
			authorityLevel:    "specialist",
			supervisorBotID:   aureliaSovereignBotID,
			supervisorBotName: aureliaSovereignBotName,
			domainScope:       "governanca de dados, auditoria de SQLite, Qdrant, Obsidian e trilha de conformidade",
			requiredProvider:  "openrouter",
			requiredModel:     "minimax/minimax-m2.7",
			toolAllowlist: []string{
				"read_file", "write_file", "list_dir", "run_command",
				"docker_control", "system_monitor", "service_control",
			},
		}
	case "homelab-logs":
		return botGovernanceProfile{
			authorityLevel:    "specialist",
			supervisorBotID:   aureliaSovereignBotID,
			supervisorBotName: aureliaSovereignBotName,
			domainScope:       "monitoramento do homelab, incidentes, crons, docker, systemd, GPU e sinais operacionais",
			toolAllowlist: []string{
				"read_file", "list_dir", "run_command",
				"docker_control", "system_monitor", "service_control",
			},
			silentWhenHealthy: true,
		}
	case "ac-vendas":
		return botGovernanceProfile{
			authorityLevel:    "specialist",
			supervisorBotID:   aureliaSovereignBotID,
			supervisorBotName: aureliaSovereignBotName,
			domainScope:       "vendas HVAC, qualificacao comercial, propostas e follow-up tecnico-comercial",
			toolAllowlist: []string{
				"read_file", "list_dir", "web_search",
				"create_schedule", "list_schedules",
			},
		}
	case "organizadora-obras":
		return botGovernanceProfile{
			authorityLevel:    "specialist",
			supervisorBotID:   aureliaSovereignBotID,
			supervisorBotName: aureliaSovereignBotName,
			domainScope:       "planejamento, execucao e documentacao de obras e fornecedores",
			toolAllowlist: []string{
				"read_file", "write_file", "list_dir", "run_command",
				"create_schedule", "list_schedules",
			},
		}
	case "agenda-pessoal":
		return botGovernanceProfile{
			authorityLevel:    "specialist",
			supervisorBotID:   aureliaSovereignBotID,
			supervisorBotName: aureliaSovereignBotName,
			domainScope:       "agenda pessoal, rotina familiar, treinos e compromissos pessoais",
			toolAllowlist: []string{
				"read_file", "create_schedule", "list_schedules",
			},
		}
	case "caixa-pf-pj":
		return botGovernanceProfile{
			authorityLevel:    "specialist",
			supervisorBotID:   aureliaSovereignBotID,
			supervisorBotName: aureliaSovereignBotName,
			domainScope:       "contas Caixa PF/PJ, pendencias financeiras e lembretes bancarios",
			toolAllowlist: []string{
				"read_file", "create_schedule", "list_schedules", "cpf_cnpj",
			},
		}
	default:
		return botGovernanceProfile{
			authorityLevel:    "specialist",
			supervisorBotID:   aureliaSovereignBotID,
			supervisorBotName: aureliaSovereignBotName,
			domainScope:       "execucao especializada no proprio canal, sem autoridade soberana",
			toolAllowlist:     cloneToolList(defaultConversationTools),
		}
	}
}

func (p botGovernanceProfile) promptContract(botName string) string {
	name := strings.TrimSpace(botName)
	if name == "" {
		name = "este bot"
	}

	if p.authorityLevel == "sovereign" {
		return strings.TrimSpace(`## Governanca 2026
- Voce e ` + name + `, a camada soberana de comando deste ecossistema multi-bot.
- Todos os demais bots sao especialistas subordinados a voce.
- Consolide contexto entre dominios, arbitre prioridade, risco e ordem de execucao.
- Nao transfira autoridade soberana a outro bot nem permita que um especialista se apresente como comando final.`)
	}

	lines := []string{
		"## Governanca 2026",
		"- Voce e um agente especializado subordinado a " + p.supervisorBotName + " (" + p.supervisorBotID + ").",
		"- Escopo exclusivo: " + p.domainScope + ".",
		"- Fora do escopo, em decisoes cross-domain ou em mudancas de governanca, escale para " + p.supervisorBotName + ".",
		"- Nao assuma autoridade soberana e nao comande outros bots.",
	}
	if p.requiredModel != "" {
		lines = append(lines, "- Regra de qualidade: este canal usa somente "+p.requiredModel+" e nao pode trocar de modelo por conveniencia.")
	}
	if p.silentWhenHealthy {
		lines = append(lines, "- Se o fluxo for automatico e o estado estiver saudavel, prefira silencio operacional ou uma unica linha objetiva.")
	}
	return strings.Join(lines, "\n")
}

func (p botGovernanceProfile) allowedTools(fallback []string) []string {
	if len(p.toolAllowlist) != 0 {
		return cloneToolList(p.toolAllowlist)
	}
	if len(fallback) != 0 {
		return cloneToolList(fallback)
	}
	return cloneToolList(defaultConversationTools)
}

func botHasPinnedLLM(botID string) bool {
	profile := governanceProfileForBot(botID)
	return profile.requiredProvider != "" || profile.requiredModel != ""
}

func cloneToolList(tools []string) []string {
	if len(tools) == 0 {
		return nil
	}
	out := make([]string, 0, len(tools))
	seen := make(map[string]struct{}, len(tools))
	for _, tool := range tools {
		tool = strings.TrimSpace(tool)
		if tool == "" {
			continue
		}
		if _, ok := seen[tool]; ok {
			continue
		}
		seen[tool] = struct{}{}
		out = append(out, tool)
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
