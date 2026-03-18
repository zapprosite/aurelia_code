---
description: Ativa o modo de auto-aprovação (Full Auto) em todos os agentes e no IDE.
---
# Workflow: bypass-on

Passos:
1.  **Claude Code**: Atualizar `~/.claude/settings.json` para `"defaultMode": "bypassPermissions"`.
2.  **Codex CLI**: Atualizar `~/snap/codex/current/config.toml` para `approval_policy = "never"`.
3.  **Antigravity IDE**: Atualizar `~/.config/Antigravity/User/settings.json` para `"antigravity.agent.dangerouslyBypassApprovals": true`.
4.  **Cortex/Workflow**: Executar `workflow-manage({ action: "setAutonomous", enabled: true })`.
5.  **Notificar**: Informar ao usuário que o modo **Full Auto** está ativo em todas as camadas.
