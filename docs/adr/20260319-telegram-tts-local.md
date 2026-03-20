---
title: Telegram TTS local via voice-proxy
status: accepted
date: 2026-03-19
---

# Contexto

A Aurelia já aceitava áudio no Telegram e transcrevia com Groq, mas a resposta ainda caía em texto.
O homelab já expõe um TTS local em `http://127.0.0.1:8011` via `voice-proxy`, compatível com `/v1/audio/speech`.

# Decisão

Adotar TTS local no Telegram com os defaults abaixo:

- provider: `openai_compatible`
- base URL: `http://127.0.0.1:8011`
- model: `chatterbox`
- voice: `Aurelia.wav` (PT-BR, sweet & educated)
- format: `opus`
- speed: `1.0`

Quando `requiresAudio=true`, a Aurelia:

1. sanitiza markdown para fala
2. sintetiza no `voice-proxy`
3. envia `Voice` no Telegram quando o formato for `opus`
4. faz fallback para texto se TTS ou envio falharem

# Arquivos afetados

- `pkg/tts/openai_compatible.go`
- `internal/telegram/output.go`
- `internal/telegram/bot.go`
- `internal/telegram/input_pipeline.go`
- `internal/config/config.go`

# Provas

- `go test ./... -count=1`
- `voice-proxy` respondeu `POST /v1/audio/speech 200`
- o runtime live registrou `voice_events.accepted=1` e `requires_tts=1`

# Consequências

- o caminho de voz no Telegram deixa de depender de mock de edge-tts
- a resposta em áudio permanece local-first
- falha de TTS não quebra a conversa; só degrada para texto
