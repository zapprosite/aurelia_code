---
description: Inicia uma nova funcionalidade seguindo a governança de worktree.
---
# Workflow: start-feature

1. **Planejamento**: Antigravity/Gemini analisa o PRD.md e cria um novo plano em `.context/plans/feat-<nome>.md`.
2. **Worktree**: Execute o comando para criar uma nova worktree: `git worktree add ../feat-<nome> main`.
3. **Iniciação**: Mova o contexto do plano para a nova worktree e notifique o usuário.
4. **Handoff**: Delegar a primeira tarefa de pesquisa para o Claude `researcher`.
