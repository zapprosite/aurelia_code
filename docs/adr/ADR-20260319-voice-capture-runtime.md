---
description: Slice nonstop para integrar captura contínua de microfone ao voice plane já existente.
status: in_progress
---

# ADR-20260319-voice-capture-runtime

## Status

- Em execução

## Slice

- slug: voice-capture-runtime
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: `docs/adr/taskmaster/ADR-20260319-voice-capture-runtime.json`

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O repositório já tinha a primeira fase executável do voice plane:

- spool local de áudio
- processor com heartbeat, budget diário e fallback STT
- dispatch do transcript aceito para o mesmo fluxo do Telegram
- health e métricas do `voice_processor`

O gargalo restante era a captura contínua. Sem ela, o runtime dependia de `aurelia voice enqueue <arquivo>`, o que não serve como experiência JARVIS nem como base para um serviço de voz em background.

## Decisão

Esta slice introduz um **capture worker comandado por contrato**:

- `CaptureSource` abstrai a origem da captura
- `CommandCaptureSource` permite plugar um capturador externo
- `CaptureWorker` roda em loop, mantém heartbeat, health e route de status
- o worker enfileira clipes no spool já existente

Contrato do capturador externo:

- stdout vazio: nenhum clipe aprovado
- stdout JSON: clipe detectado e aprovado
- campos mínimos:
  - `detected`
  - `audio_file`
- campos opcionais:
  - `user_id`
  - `chat_id`
  - `requires_audio`
  - `source`
  - `delete_source_after`

Esse desenho permite conectar `openWakeWord + Silero VAD + ring buffer` depois, sem reescrever o plano de controle em Go.

## Escopo

- config do capture worker
- worker de captura em background
- `voice_capture` no `/health`
- `GET /v1/voice/capture/status`
- testes do contrato de captura

## Fora de escopo

- implementação local de `openWakeWord`
- implementação local de `Silero VAD`
- ring buffer real dentro do processo Go
- deploy em `/home/will/aurelia-24x7`
- TTS

## Arquivos afetados

- `internal/voice/capture.go`
- `internal/voice/capture_test.go`
- `internal/voice/metrics.go`
- `internal/config/config.go`
- `cmd/aurelia/app.go`
- `docs/adr/taskmaster/ADR-20260319-voice-capture-runtime.json`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/v1/voice/capture/status`
  - `curl -fsS http://127.0.0.1:8484/health`
- testes:
  - `go test ./internal/voice ./cmd/aurelia -count=1`
  - `go test ./... -count=1`
- scripts:
  - capturador externo futuro via `voice_capture_command`
- fallback:
  - `go run ./cmd/aurelia voice enqueue <arquivo>`
  - `voice_capture_enabled=false`

## Rollout

1. integrar o capture worker no repositório principal
2. validar unit + suite
3. plugar capturador real (`openWakeWord + Silero`) via `voice_capture_command`
4. só então portar a slice para a worktree de deploy

## Rollback

- `voice_capture_enabled=false`
- manter `voice processor` e `voice enqueue`
- remover apenas o worker de captura, sem tocar no spool/processador/STT já verdes

## Evidência esperada

- worker inicia e para junto do app
- `/v1/voice/capture/status` responde
- `voice_capture` entra no `/health`
- `CaptureOnce` enfileira clipe detectado no spool
- suite relevante passa

## Pendências / bloqueios

- ainda falta o capturador real com `openWakeWord + Silero`
- ainda falta o rollout da captura na worktree `/home/will/aurelia-24x7`
