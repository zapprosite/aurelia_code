---
description: Processo final de revisão de código, conformidade e merge na main.
---
# Workflow: review-merge

1. **Auditoria**: Invoca o Claude `reviewer` para validar o código e a conformidade com `AGENTS.md`.
2. **Teste**: Executa a base de testes automatizados na worktree.
3. **Aprovação**: Solicita autorização humana obrigatória (Tier C).
4. **Merge**: Realiza o merge da branch/worktree para a `main` e remove a worktree temporária.
