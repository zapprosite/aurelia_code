---
description: Slice nonstop para fechar handoff fim a fim do Antigravity sem retrabalho nem perda de contexto.
status: in_progress
---

# ADR-20260319-antigravity-handoff-e2e

## Status

- Em execução

## Slice

- slug: antigravity-handoff-e2e
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: `docs/adr/taskmaster/ADR-20260319-antigravity-handoff-e2e.json`

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O runtime já gera prompts automáticos para tarefas `light` e já tem skill do Antigravity instalada. O gap restante é o E2E completo: da classificação até a volta do resultado com evidência suficiente para a Aurélia continuar sem duplicar trabalho.

## Decisão

Fechar um handoff estruturado com:

- prompt de ida padronizado
- payload mínimo obrigatório
- resultado de volta estruturado
- evidência e próximos passos
- classificação clara entre `approved`, `revise`, `blocked`

## Escopo

- contrato de ida
- contrato de volta
- payload mínimo
- menos retrabalho entre chat leve e runtime
- prova de E2E

## Fora de escopo

- substituir o executor principal pelo Antigravity
- múltiplos chats concorrentes

## Arquivos afetados

- `internal/telegram/antigravity_prompt.go`
- `internal/telegram/antigravity_prompt_test.go`
- `internal/telegram/input_pipeline.go`
- `docs/antigravity_gemini_operator_blueprint.md`
- `docs/adr/taskmaster/ADR-20260319-antigravity-handoff-e2e.json`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/v1/router/status`
- testes:
  - `go test ./internal/telegram ./internal/agent -count=1`
  - `go test ./... -count=1`
- scripts:
  - smoke manual de tarefa `light` no Telegram
- fallback:
  - responder com handoff estruturado em texto
  - marcar `handoff_required=true`

## Rollout

1. fixar payload mínimo de ida
2. fixar shape de volta
3. adicionar guarda contra retrabalho
4. provar um ciclo completo no Telegram/Antigravity

## Rollback

- manter geração do prompt de ida
- desabilitar consumo automático do retorno
- exigir revisão humana do handoff

## Evidência esperada

- tarefa `light` vira prompt consistente
- resposta do Antigravity volta sem perder contexto
- Aurélia consegue continuar sem pedir a mesma informação de novo

## Pendências / bloqueios

- ainda falta prova humana fim a fim com um caso real no cockpit do Antigravity

## Progresso registrado

- o handoff de ida agora já nasce com contrato estruturado em JSON
- o prompt exige resposta final em `approved|revise|blocked`
- o parser da volta já existe no runtime para consumo futuro
- o pipeline do Telegram agora já consome o retorno estruturado quando ele é colado no chat
- a resposta volta formatada como handoff final utilizável pela Aurélia
- `go test ./internal/telegram ./internal/agent -count=1` passou
- `go test ./... -count=1` passou
