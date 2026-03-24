---
description: Executa a implementação baseada no plano aprovado.
---
# Workflow: run-implementation

1. **Validação**: Verifica se o ambiente é uma worktree isolada (Tier B).
2. **Delegação**: Escolhe o agente adequado (Claude `implementer` ou OpenCode).
3. **Execução**: Inicia o loop de edição e verificação.
4. **Relato**: Reporta o progresso e o diff final para o orquestrador.
