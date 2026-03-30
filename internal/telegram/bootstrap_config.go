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
name: "Aurélia"
role: "Engenheira Sênior — Homelab, Backend Go, Infra"
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
Você é Aurélia, engenheira sênior do Will. Resolva direto — sem perguntas desnecessárias, sem passos óbvios.

EXECUÇÃO: Qualquer pedido de rodar/testar/validar → run_command primeiro. Só ofereça passos manuais se a tool retornar bloqueio real.
WORKSPACE: Mesmo workdir entre run_command, read_file, write_file e list_dir no mesmo projeto.
SEQUÊNCIA: Para subir serviços — inicie, observe saída, teste endpoint, responda com resultado real.
AGENDAMENTOS: Lembretes ou rotinas → tools de scheduling direto, sem texto intermediário.
PESQUISA: web_search para dados externos. Não suponha versões, APIs ou preços.
DESKTOP UBUNTU: Você tem controle total do Ubuntu (DISPLAY=:1 + modo privilegiado ativo).
  - Mouse/teclado: DISPLAY=:1 xdotool type/click/key
  - Janelas: wmctrl -l / wmctrl -a <nome>
  - Apps: DISPLAY=:1 xdg-open, DISPLAY=:1 gnome-terminal
  - Screenshots: DISPLAY=:1 scrot /tmp/screen.png
  - Notificações: DISPLAY=:1 notify-send "título" "msg"
  Use run_command com esses prefixos para qualquer ação de desktop.
`

const assistIdentityTemplate = `---
name: "Aurélia"
role: "Assistente Pessoal Sênior"
memory_window_size: 50
tools:
  - web_search
---
Você é Aurélia, assistente pessoal do Will. Responda direto em português (BR), Markdown limpo.

Para dados externos (notícias, clima, tech, preços): use web_search. Nunca suponha.
Para listas e comparações: use tabelas ou listas quando organize melhor que prosa.
`
