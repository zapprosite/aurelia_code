---
description: Entrega inteligente — staging seguro, commit semântico, push e criação de PR no GitHub.
---

# Fluxo de Entrega Inteligente (v2)

## Fase 1 — Diagnóstico
1. Obter branch atual: `git branch --show-current`.
   **Interromper** se for `main` ou `master`.
2. Checar estado: `git status --short`. Se não houver mudanças, informar e encerrar.
3. Exibir resumo: `git diff --stat HEAD`.

## Fase 2 — Staging Seguro
4. Adicionar rastreados e novos não-ignorados:
   ```bash
   git add -A
   ```
   O `.gitignore` é a proteção real contra `.env`, secrets e build artifacts.
5. Verificar staged: `git diff --cached --stat`.

## Fase 3 — Commit Semântico
6. Analisar o diff staged para detectar:
   - **Tipo**: `feat` / `fix` / `chore` / `refactor` / `docs` / `test`
   - **Escopo**: derivado dos paths alterados
     - `apps/api/` → `(api)`, `apps/frontend/` → `(web)`, `packages/ui/` → `(ui)`
     - Múltiplos pacotes → omitir escopo
7. Gerar mensagem clara e específica (não genérica):
   - ✅ `feat(api): add JWT refresh token endpoint`
   - ✅ `fix(web): resolve hydration mismatch on dashboard`
   - ❌ `feat: auto-ship changes from agent`
8. Commitar com co-authoria:
   ```bash
   git commit -m "[tipo(escopo)]: [descrição específica]

   Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
   ```

## Fase 4 — Push
9. Enviar com proteção contra sobrescrita acidental:
   ```bash
   git push --force-with-lease origin HEAD
   ```

## Fase 5 — Pull Request
10. Verificar se PR já existe: `gh pr view 2>/dev/null`.
11. Se não existir, criar PR com título e body derivados do commit:
    ```bash
    gh pr create \
      --title "[mensagem do commit principal]" \
      --body "## Summary\n- [mudanças detectadas]\n\n## Test plan\n- [ ] Smoke test\n- [ ] CI passing" \
      --base main
    ```
12. Exibir URL do PR criado.

## Pós-entrega
13. Mostrar status final: branch, commits enviados, URL do PR.
14. Sugerir: aguardar CI, solicitar review com `gh pr status`.
