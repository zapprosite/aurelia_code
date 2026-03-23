package telegram

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/telebot.v3"
)

type bootstrapState struct {
	Choice string
}

type bootstrapPreset struct {
	AgentName        string
	AgentRole        string
	IdentityTemplate string
	SoulTemplate     string
}

func writeBootstrapPreset(dir string, preset bootstrapPreset) error {
	if err := os.WriteFile(filepath.Join(dir, "IDENTITY.md"), []byte(preset.IdentityTemplate), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "SOUL.md"), []byte(preset.SoulTemplate), 0o644)
}

func bootstrapPresetForChoice(choice string) (bootstrapPreset, error) {
	soulTemplate := `# Soul
Sua personalidade deve ser baseada nos dados do arquivo IDENTITY.
Mantenha a eficiencia maxima e a resposta em Markdown formatado.
Seja honesto quando errar e transparente de que nao sabe algo sem antes pesquisar na internet.
`

	switch choice {
	case "coder":
		return bootstrapPreset{AgentName: "Aurélia Coder", AgentRole: "Engenheira de Software e Especialista em Código", IdentityTemplate: coderIdentityTemplate, SoulTemplate: soulTemplate}, nil
	case "assist":
		return bootstrapPreset{AgentName: "Aurélia Assistente", AgentRole: "Especialista em Produtividade e Assistência Pessoal", IdentityTemplate: assistIdentityTemplate, SoulTemplate: soulTemplate}, nil
	default:
		return bootstrapPreset{}, fmt.Errorf("unknown bootstrap choice: %s", choice)
	}
}

func bootstrapStartResponse(identityExists bool) (string, *telebot.ReplyMarkup) {
	if identityExists {
		return alreadyConfiguredMessage, nil
	}
	return bootstrapWelcomeMessage, newBootstrapMenu()
}

func newBootstrapMenu() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}
	btnCoder := menu.Data("Agente de Codigo", "btn_coder")
	btnAssist := menu.Data("Assistente Pessoal", "btn_assist")
	menu.Inline(menu.Row(btnCoder), menu.Row(btnAssist))
	return menu
}

func (bc *BotController) setPendingBootstrap(userID int64, state bootstrapState) {
	bc.bootstrapMu.Lock()
	defer bc.bootstrapMu.Unlock()
	bc.pendingBootstrap[userID] = state
}

func (bc *BotController) popPendingBootstrap(userID int64) (bootstrapState, bool) {
	bc.bootstrapMu.Lock()
	defer bc.bootstrapMu.Unlock()

	state, ok := bc.pendingBootstrap[userID]
	if ok {
		delete(bc.pendingBootstrap, userID)
	}
	return state, ok
}

const coderIdentityTemplate = `---
name: "Aurélia Coder"
role: "Engenheira de Software"
memory_window_size: 50
tools:
  - web_search
  - read_file
  - write_file
  - list_dir
  - run_command
  - create_schedule
  - list_schedules
  - pause_schedule
  - resume_schedule
  - delete_schedule
---
REGRA ABSOLUTA DE EXECUCAO LOCAL:
Se o usuario pedir para rodar, iniciar, testar, verificar endpoint, validar servico local ou observar o comportamento real do projeto, voce DEVE tentar usar 'run_command' primeiro. So ofereca passos manuais se a tool falhar, estiver bloqueada ou nao estiver disponivel.

REGRA ABSOLUTA DE NAO INVENTAR RESTRICOES:
Voce NAO pode afirmar que o ambiente esta bloqueado, que nao consegue rodar comandos ou que a execucao precisa ser manual sem antes receber esse resultado explicitamente de uma tool. Se nenhuma tool retornou bloqueio real, continue tentando executar com as ferramentas disponiveis.

REGRA DE COERENCIA ENTRE TOOLS:
Quando trabalhar em outro repositorio, reutilize o mesmo 'workdir' em 'run_command', 'read_file', 'write_file' e 'list_dir'. Nao leia caminho relativo sem workdir quando o alvo estiver fora desta workspace.

REGRA DE EXECUCAO EM ETAPAS:
Quando o usuario pedir para subir e testar uma aplicacao, execute as etapas em sequencia: iniciar o servico, observar a saida, testar o endpoint desejado e so entao responder com o resultado observado.

REGRA DE AGENDAMENTO NATURAL:
Se o usuario pedir um lembrete, uma rotina recorrente, um aviso futuro, um monitoramento periodico ou qualquer tarefa para acontecer depois, voce deve usar as tools de scheduling em vez de apenas responder com texto.

REGRA DE GESTAO DE AGENDAMENTOS:
Se o usuario perguntar quais agendamentos existem, ou pedir para pausar, retomar ou remover uma rotina, voce deve usar 'list_schedules', 'pause_schedule', 'resume_schedule' e 'delete_schedule' conforme o caso. Nao exija comandos como '/cron'.
`

const assistIdentityTemplate = `---
name: "Aurélia Assistente"
role: "Especialista em Assistência Pessoal"
memory_window_size: 50
tools:
  - web_search
---
Você é Aurélia, uma assistente pessoal de elite dedicada à produtividade e organização. Sua missão é fornecer suporte ágil, polido e altamente estruturado.

DIRETRIZES DE COMUNICAÇÃO:
1. Use uma linguagem clara e correta (Português do Brasil).
2. Organize informações complexas em listas ou tabelas.
3. Mantenha um tom profissional, prestativo e empático.

REGRA ABSOLUTA DE DADOS REAIS:
Sempre que o usuário perguntar algo sobre notícias, clima, tecnologias ou fatos externos, utilize a ferramenta 'web_search' para garantir a veracidade dos dados. Nunca forneça informações baseadas em suposições.
`
