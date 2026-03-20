---
description: Slice nonstop para provar no ambiente live o caminho completo de wake word até resposta.
status: in_progress
---

# ADR-20260319-voice-e2e-proof-live

## Status

- Em execução

## Slice

- slug: voice-e2e-proof-live
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: `docs/adr/taskmaster/ADR-20260319-voice-e2e-proof-live.json`

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O voice plane já existe no runtime:

- capture worker
- wake word + VAD
- spool
- processor
- Groq STT
- TTS no Telegram

O gap real restante é a prova humana E2E no live: falar com o headset, passar pelo wake word e receber resposta válida.

## Decisão

Tratar isso como slice própria de prova operacional, com foco em:

- headset ALSA correto
- wake phrase positiva
- transcript aceito
- resposta final útil
- evidência de banco, logs e health

## Escopo

- prova humana do voice plane
- validação live em `/home/will/aurelia-24x7`
- evidência em `voice_events`
- governança de fallback

## Fora de escopo

- troca de provider STT
- nova arquitetura de voz

## Arquivos afetados

- `scripts/voice-capture-openwakeword.sh`
- `scripts/voice-capture-openwakeword.py`
- `scripts/voice-capture-smoke.sh`
- `internal/voice/`
- `docs/adr/taskmaster/ADR-20260319-voice-e2e-proof-live.json`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/health`
  - `curl -fsS http://127.0.0.1:8484/v1/voice/capture/status`
  - `curl -fsS http://127.0.0.1:8484/v1/voice/status`
- testes:
  - `go test ./internal/voice ./cmd/aurelia -count=1`
  - `go test ./... -count=1`
- scripts:
  - `bash ./scripts/voice-capture-smoke.sh`
  - `bash ./scripts/health-check.sh`
- fallback:
  - `go run ./cmd/aurelia voice enqueue <arquivo>`
  - `voice_capture_enabled=false`

## Rollout

1. validar headset e dispositivo ALSA
2. validar silêncio sem falso positivo
3. validar wake phrase positiva
4. provar transcript aceito
5. provar resposta útil

## Rollback

- desligar `voice_capture`
- manter `voice enqueue`
- manter Telegram voice input como caminho manual

## Evidência esperada

- `voice_capture=ok` e `voice_processor=ok`
- `voice_events` registra wake positivo
- transcript aceito chega ao pipeline
- resposta volta ao usuário

## Pendências / bloqueios

- depende de prova humana com headset real no ambiente live
