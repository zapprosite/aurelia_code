# ADR Slice: Onboarding Simplificado — Senior Direct

## Contexto
O onboarding do bot Telegram era verboso demais ("cockpit de comando", "cinematográfica", "SAP", emojis decorativos). O target é um onboarding direto e profissional que respeite o tempo do usuário sênior.

## Decisões

### Mensagens simplificadas (`internal/telegram/messages.go`)
```go
alreadyConfiguredMessage = "**Aurélia online.** Manda sua tarefa."
bootstrapWelcomeMessage = "**Aurélia** — selecione o perfil:"
bootstrapSuccessMessage = "Pronta. Manda sua tarefa — texto, áudio, imagem ou comando."
bootstrapFailureMessage = "Erro ao criar identidade. Tenta `/start` novamente."
bootstrapProfileMessage = "Perfil configurado. Qual é o seu nome e como prefere ser tratado?"
```

### Handler `/start` (`internal/telegram/bootstrap.go`)
- Botão `btn_status` registrado → `handleStatus` (antes: botão morto)
- Welcome em Markdown limpo: "**Aurélia online.** Manda sua tarefa — código, infra, análise, pesquisa."
- Botões sem emoji decorativo: "Dashboard" e "Status do sistema"

### System Prompt (`internal/telegram/input_pipeline.go`)
```go
const defaultSystemPrompt = `Você é Aurélia, engenheira sênior e assistente Jarvis do Will no Ubuntu Desktop.

REGRAS:
- Responda em português (BR), Markdown limpo. Direto: diagnóstico → solução → código.
- Tarefas técnicas: run_command primeiro, sempre.
- Pesquisa: web_search para dados externos — nunca suponha.
- Agendamentos: tools de scheduling direto, sem intermediário.

DESKTOP UBUNTU (DISPLAY=:1 + modo privilegiado ativo):
- Mouse/teclado: DISPLAY=:1 xdotool type/click/key
- Janelas: wmctrl -l / DISPLAY=:1 wmctrl -a
- Screenshot: DISPLAY=:1 scrot /tmp/screen.png && cat /tmp/screen.png
- Apps: DISPLAY=:1 xdg-open / DISPLAY=:1 gnome-terminal
- Notificação: DISPLAY=:1 notify-send "título" "msg"
- Clipboard: DISPLAY=:1 xclip -selection clipboard -i <<< "texto"`
```

### Identity Templates (`internal/telegram/bootstrap_config.go`)
- `name: "Aurélia"`, `role: "Engenheira Sênior — Homelab, Backend Go, Infra"`
- 5 regras concisas em vez de 6 regras verbosas e redundantes

## Consequências
- **Positivo**: Onboarding direto, sem marketing; identidade clara; ferramentas de desktop ativadas
- **Negativo**: Zero — mudança de comunicação apenas
