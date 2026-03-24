---
description: Inicia uma nova funcionalidade seguindo a governança de worktree.
---
# Workflow: start-feature

1. **Planejamento**: Antigravity/Gemini analisa `plan.md` e `AGENTS.md`, cria um novo plano em `.context/plans/feat-<nome>.md` com ADR em `docs/adr/`.
2. **Worktree**: Execute o comando para criar uma nova worktree: `git worktree add ../feat-<nome> main`.
3. **Iniciação**: Mova o contexto do plano para a nova worktree e notifique o usuário.
4. **Handoff**: Delegar a primeira tarefa de pesquisa para o Claude `researcher`.
