# ADR Slice: STT Local-First Cascade — Whisper + Groq Fallback

## Contexto
O bot Telegram usava Groq (cloud) como transcrição primária de áudio. Para 50-80% de sovereignidade local, o Whisper local (faster-whisper-server) deve ser priorizado.

## Decisão

### Prioridade em `pkg/stt/factory.go`

```
Modo "local" / "faster-whisper":
  1. faster-whisper-server (:8020)  ← LOCAL
  2. Groq Cloud (whisper-large-v3) ← FALLBACK

Modo "" / "groq":
  1. Groq Cloud (whisper-large-v3)  ← PRIMARY
  2. faster-whisper-server (:8020)  ← FALLBACK
```

### Docker Compose
O container `whisper-local` já estava rodando (fedirz/faster-whisper-server:latest-cuda → porta 8020).

### TTS Parallel (saída de voz)
`deliverWithParallelTTS()` envia texto imediatamente → sintetiza áudio em goroutine → envia como voice note (opus/ogg). O sanitizador `sanitizeTextForSpeech()` remove markdown e expande símbolos.

## Consequências
- **Positivo**: 100% local para STT quando faster-whisper disponível; zero custo de API
- **Negativo**: Groq fallback consome API quota se faster-whisper falhar
- **Neutro**: Latência de transcrição pode variar (local ~2-5s vs cloud ~1-3s)
