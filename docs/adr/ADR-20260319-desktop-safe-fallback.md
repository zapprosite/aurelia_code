---
description: Slice nonstop para endurecer o desktop fallback com click seguro, digitação segura e kill-switch.
status: in_progress
---

# ADR-20260319-desktop-safe-fallback

## Status

- Em execução

## Slice

- slug: desktop-safe-fallback
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: `docs/adr/taskmaster/ADR-20260319-desktop-safe-fallback.json`

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O repositório já tem baseline de `xdotool`, `wmctrl`, `scrot`, foco de janela e screenshot local. O que falta é impedir que desktop fallback vire automação cega.

## Decisão

Endurecer o desktop fallback com quatro guardrails:

- click seguro com screenshot antes/depois
- digitação segura com confirmação de foco
- orçamento máximo de passos
- kill-switch explícito e reversão rápida

## Escopo

- click seguro
- digitação segura
- limite de passos
- kill-switch
- evidência visual mínima

## Fora de escopo

- automação desktop irrestrita
- uso de desktop como caminho primário

## Arquivos afetados

- `internal/tools/`
- `scripts/deploy-desktop.sh`
- `docs/adr/taskmaster/ADR-20260319-desktop-safe-fallback.json`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/health`
- testes:
  - `go test ./internal/tools -count=1`
  - `go test ./... -count=1`
- scripts:
  - `bash ./scripts/deploy-desktop.sh`
  - `xdotool getwindowfocus getwindowname`
  - `scrot /tmp/aurelia-desktop-before.png`
- fallback:
  - abortar para browser/CLI
  - exigir confirmação humana

## Rollout

1. validar a janela/foco antes de qualquer ação
2. limitar número de passos
3. adicionar screenshot antes/depois
4. adicionar abort central

## Rollback

- desligar desktop fallback
- manter browser-first e CLI-first

## Evidência esperada

- click sem foco válido é recusado
- digitação sem confirmação de contexto é recusada
- limite de passos aborta a automação
- screenshots existem antes e depois

## Pendências / bloqueios

- falta costurar tudo em uma única política operacional de desktop
