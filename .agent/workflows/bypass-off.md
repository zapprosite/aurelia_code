---
description: Desativa o modo de auto-aprovação, exigindo permissões manuais e revisões.
---
# Workflow: bypass-off

Passos:
1.  **Claude Code**: Atualizar `~/.claude/settings.json` para `"defaultMode": "planned"` (ou remover bypass).
2.  **Antigravity IDE**: Atualizar `~/.config/Antigravity/User/settings.json` para `"antigravity.agent.dangerouslyBypassApprovals": false`.
3.  **Cortex/Workflow**: Executar `workflow-manage({ action: "setAutonomous", enabled: false })`.
4.  **Notificar**: Informar ao usuário que o modo de **Segurança Máxima** (Permissões Manuais) está ativo.
