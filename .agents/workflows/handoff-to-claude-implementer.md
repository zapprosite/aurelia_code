---
description: Transfere a execução de implementação para o subagente especializado do Claude.
---
# Workflow: handoff-to-claude-implementer

1. **Contexto**: Garante que o `implementation_plan.md` está aprovado e atualizado.
2. **Alvo**: Identifica os arquivos e a worktree de destino.
3. **Comando**: Invoca o Claude CLI com a persona de `implementer`: `claude --agent .claude/agents/implementer.md "Execute o plano de implementação em <path>"`.
4. **Monitoramento**: Aguarda a conclusão e revisa o diff gerado.
