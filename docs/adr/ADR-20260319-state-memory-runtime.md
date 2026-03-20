---
description: Slice nonstop para persistir governor/breaker e consolidar a verdade local do voice plane.
status: accepted
---

# ADR-20260319-state-memory-runtime

## Status

- Aceito

## Slice

- slug: state-memory-runtime
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: docs/adr/taskmaster/ADR-20260319-state-memory-runtime.json

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O runtime já tinha budgets em memória no gateway e orçamento diário em arquivo no `voice processor`. Isso era suficiente para testes locais curtos, mas fraco para operação longa:

- restart do processo zerava breaker e contadores do gateway
- transcript de áudio só vivia no spool e nos mirrors remotos opcionais
- a verdade local do homelab ficava pouco consultável

## Decisão

Esta slice fecha o estado mínimo profissional do runtime local:

- `gateway_route_states` persiste requests, failures, breaker e modelo por rota em SQLite
- rollover diário reseta orçamento/estado para não travar budgets indefinidamente
- `voice_events` vira espelho local em SQLite para transcripts aceitos e rejeitados
- `Supabase` e `Qdrant` continuam como mirrors opcionais, não como única verdade

## Escopo

- store SQLite do gateway
- persistência do breaker/budgets do gateway
- mirror local de transcripts em SQLite
- defaults operacionais para `voice_reply_user_id/chat_id` a partir do usuário Telegram autorizado
- testes do store e do mirror local

## Fora de escopo

- consultas de analytics ou UI para `voice_events`
- persistência distribuída de breaker em rede
- reconciliação full `SQLite -> Supabase -> Qdrant`

## Arquivos afetados

- `internal/gateway/provider.go`
- `internal/gateway/store.go`
- `internal/gateway/provider_test.go`
- `internal/voice/sqlite_mirror.go`
- `internal/voice/sqlite_mirror_test.go`
- `internal/config/config.go`
- `cmd/aurelia/app.go`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/v1/router/status`
  - `curl -fsS http://127.0.0.1:8484/v1/voice/status`
- testes:
  - `go test ./internal/gateway ./internal/voice ./cmd/aurelia -count=1`
  - `go test ./... -count=1`
- scripts:
  - `sqlite3 ~/.aurelia/data/aurelia.db 'select count(*) from voice_events;'`
- fallback:
  - se o mirror local falhar, manter spool + mirrors remotos
  - se o store do gateway falhar, degradar para store em memória

## Rollout

1. validar suite local
2. levar o store/mirror para a worktree de deploy
3. provar que restart não zera o breaker/budget lane
4. provar `voice_events` no banco local após processamento real

## Rollback

- desligar o store persistente do gateway e manter memória do processo
- remover apenas o mirror SQLite e manter spool + Supabase/Qdrant

## Evidência esperada

- `go test ./... -count=1` verde
- `gateway_route_states` persistido em SQLite
- `voice_events` preenchido após transcript
- restart do processo preserva state load do gateway

## Pendências / bloqueios

- reconciliação mais rica entre `SQLite`, `Supabase` e `Qdrant` segue aberta para a próxima rodada

## Evidência registrada

- `gateway_route_states` validado no SQLite local com contadores e breaker state
- `voice_events` recebeu transcript real no runtime live
- `go test ./internal/gateway ./internal/voice ./cmd/aurelia -count=1` passou
