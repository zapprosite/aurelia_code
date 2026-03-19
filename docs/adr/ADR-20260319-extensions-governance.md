---
description: Slice nonstop para fechar a política de extensões sem contaminar o core do runtime.
status: accepted
---

# ADR-20260319-extensions-governance

## Status

- Aceito

## Slice

- slug: extensions-governance
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: docs/adr/taskmaster/ADR-20260319-extensions-governance.json

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O projeto vinha discutindo extensões para Chrome e Antigravity como aceleradores. O risco era deixar tooling opcional virar dependência implícita do runtime.

## Decisão

- nenhuma extensão é parte do core do runtime
- o core continua em `Go + CLI + agent-browser + Playwright + systemd`
- extensões só entram como perfil opcional do operador humano
- toda extensão precisa de classificação:
  - `core`: proibido nesta fase
  - `nice-to-have`: permitido com rollback simples
  - `risky`: documentado, mas não recomendado

## Escopo

- matriz de extensões recomendadas
- separação `core` vs `nice-to-have`
- política de rollback

## Fora de escopo

- instalar extensões automaticamente
- depender de extensão para login, browser-use, Antigravity ou deploy

## Arquivos afetados

- `docs/extension_matrix_20260319.md`
- `docs/adr/PENDING-SLICES-20260319.md`
- `plan.md`

## Simulações e smoke previstos

- curl:
  - n/a
- testes:
  - revisão documental
- scripts:
  - n/a
- fallback:
  - remover a extensão do perfil isolado
  - voltar ao fluxo base com DevTools/Playwright

## Rollout

1. manter o core livre de extensões
2. documentar candidatos
3. avaliar manualmente só depois do runtime ficar verde no deploy

## Rollback

- excluir extensão do perfil dedicado do Chrome
- manter o perfil principal intacto

## Evidência esperada

- matriz documental clara
- nenhuma dependência do runtime em extensão

## Pendências / bloqueios

- a avaliação prática de extensões só faz sentido depois do fechamento do deploy gateway+voice
