---
description: Slice nonstop para fechar login guiado seguro no plano de browser-use da Aurelia.
status: in_progress
---

# ADR-20260319-browser-safe-login

## Status

- Em execução

## Slice

- slug: browser-safe-login
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: `docs/adr/taskmaster/ADR-20260319-browser-safe-login.json`

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O browser-use já está validado no repositório com baseline de Playwright, screenshot e DevTools em loopback. O que ainda falta é o caminho de login seguro, para que a Aurelia consiga conduzir autenticação sem virar automação cega nem vazar segredo operacional.

Essa slice precisa proteger:

- credenciais
- 2FA
- domínios permitidos
- logs e screenshots
- handoff humano quando a etapa se tornar sensível

## Decisão

Implementar um `guided login flow` com contrato explícito:

- allowlist de domínio/host
- checkpoints por etapa (`start`, `username`, `password`, `2fa`, `success`, `abort`)
- pausa obrigatória antes de senha e 2FA
- captura visual apenas antes/depois das etapas sensíveis
- logs sem segredos
- rollback claro para handoff humano

## Escopo

- fluxo guiado de login no browser
- validação de domínio
- pause/handoff em etapas sensíveis
- artefatos seguros de screenshot
- contrato de sucesso/abort

## Fora de escopo

- login via desktop fallback
- armazenamento de credenciais
- bypass de CAPTCHA

## Arquivos afetados

- `internal/agent/`
- `internal/tools/`
- `internal/telegram/`
- `scripts/jarvis-playwright-smoke.mjs`
- `docs/adr/taskmaster/ADR-20260319-browser-safe-login.json`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/health`
  - `curl -fsS http://127.0.0.1:8484/v1/router/status`
- testes:
  - `go test ./internal/agent ./internal/tools ./internal/telegram -count=1`
  - `go test ./... -count=1`
- scripts:
  - `node ./scripts/jarvis-playwright-smoke.mjs`
- fallback:
  - converter o login para `handoff_required`
  - parar antes de senha/2FA e delegar ao humano

## Rollout

1. definir contrato do login guiado
2. aplicar gates por domínio e etapa
3. adicionar artefatos seguros e abort explícito
4. validar smoke browser
5. validar handoff antes de passos sensíveis

## Rollback

- desabilitar fluxo de login guiado
- manter apenas navegação/read-only
- delegar login para humano ou Antigravity

## Evidência esperada

- domínio fora da allowlist é recusado
- senha/2FA nunca aparece em log
- o fluxo entra em pausa antes do trecho sensível
- o operador consegue retomar com contexto suficiente

## Pendências / bloqueios

- falta definir o shape final do contrato de pausa/retomada entre browser e operador
