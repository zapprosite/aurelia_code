---
title: Voz pronta da Aurelia via Gemini TTS
status: accepted
date: 2026-03-19
---

# Contexto

O runtime já possuía:

- STT via Groq
- TTS local no Telegram via `voice-proxy`
- chave Gemini validada no host

O usuário pediu uma voz pronta, feminina e profissional em PT-BR, sem depender de clonagem nem de nova conta de áudio.

# Decisão

Adicionar uma lane `tts_provider=gemini` ao runtime, usando:

- model: `gemini-2.5-flash-preview-tts`
- voice: `Sulafat`
- format: `wav`

`Sulafat` foi escolhida como default da voz pronta da Aurelia porque o catálogo oficial a descreve como `Warm`, o que se aproxima melhor do alvo:

- feminina
- doce
- calma
- acolhedora
- profissional

O fallback operacional continua sendo o `voice-proxy` local.

# Arquivos afetados

- `pkg/tts/gemini.go`
- `pkg/tts/gemini_test.go`
- `internal/telegram/bot.go`
- `internal/config/config.go`
- `scripts/gemini-tts-smoke.sh`
- `docs/aurelia_voice_profile_20260319.md`

# Provas

- `go test ./pkg/tts ./internal/telegram ./internal/config -count=1`
- `go test ./... -count=1`
- `bash ./scripts/gemini-tts-smoke.sh`

# Consequências

- a Aurelia ganha uma voz pronta de qualidade sem depender de clonagem
- a clonagem autorizada continua em slice separada
- se Gemini falhar, o caminho local segue disponível
