---
description: Realiza uma auditoria profunda do contexto local antes de agir.
---
# Workflow: inspect-context

1. **Enumeração**: Listar arquivos em `.context/docs/` e `.context/plans/`.
2. **Autoridade**: Ler `AGENTS.md` e validar se as regras estão atualizadas.
3. **Estado**: Verificar se há uma worktree ativa e qual o status do plano atual.
4. **Relatório**: Consolidar um resumo do estado atual para o usuário.
