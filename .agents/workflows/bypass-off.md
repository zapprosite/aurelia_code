---
description: Desativa o modo de auto-aprovação, exigindo permissões manuais e revisões.
---
# Workflow: bypass-off

Passos:
1.  **Claude Code**: Atualizar `~/.claude/settings.json` para `"defaultMode": "planned"` (ou remover bypass).
2.  **Codex CLI**: Atualizar `~/snap/codex/current/config.toml` para `approval_policy = "always"`.
3.  **Antigravity IDE**: Atualizar `~/.config/Antigravity/User/settings.json` para `"antigravity.agent.dangerouslyBypassApprovals": false`.
4.  **Cortex/Workflow**: Executar `workflow-manage({ action: "setAutonomous", enabled: false })`.
5.  **Notificar**: Informar ao usuário que o modo de **Segurança Máxima** (Permissões Manuais) está ativo.
