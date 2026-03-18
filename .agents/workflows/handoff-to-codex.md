---
description: Transfere a execução de implementação para o Codex CLI.
---
# Workflow: handoff-to-codex

1. **Setup**: Garante que o Codex está operando dentro da worktree correta.
2. **Comando**: Invoca o Codex CLI: `codex "Implemente a função X seguindo o contrato Y em <path>"`.
3. **Auto-Approve**: Permite modo auto (Tier B) se dentro da worktree.
4. **Sincronia**: Atualiza o plano de execução com o resultado do Codex.
