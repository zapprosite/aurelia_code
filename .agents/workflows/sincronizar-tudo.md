---
description: Sincronização Git Padrão Sênior (Commit & Push)
---

# /sincronizar-tudo

Automatiza commit e push seguindo as convenções do monorepo.

## Passos

1. Verificar status: `git status --short`
   - Se limpo, informar "nothing to sync" e encerrar.

2. Adicionar tudo não-ignorado:
   ```bash
   git add -A
   ```
   O `.gitignore` protege secrets e build artifacts.

3. Analisar diff staged e gerar mensagem semântica real:
   - Detectar tipo: `feat` / `fix` / `chore` / `refactor` / `docs`
   - Detectar escopo pelo path (`api`, `web`, `ui`, `core`)
   - ✅ `chore(claude): add turbo workflow`
   - ❌ `feat: sincronização automática de workspace`
   ```bash
   git commit -m "[tipo(escopo)]: [descrição específica]

   Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
   ```

4. Push:
   ```bash
   git push --force-with-lease origin HEAD
   ```

5. Exibir confirmação com branch e hash do commit.
