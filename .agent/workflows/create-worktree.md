---
description: Cria uma nova worktree isolada para desenvolvimento seguro.
---
# Workflow: create-worktree

1. **Parâmetros**: Solicita o nome da branch/feature.
2. **Criação**: Executa `git worktree add ../<nome> main`.
3. **Setup**: Copia arquivos de contexto necessários para a nova worktree.
4. **Link**: Adiciona um link no plano principal de `.context/plans/`.
