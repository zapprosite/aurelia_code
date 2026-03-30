# ADR Slice: Voice Pipeline PT-BR — Whisper Local → Kokoro Local

## Contexto
Padrão de resposta "text markdown 2026 + áudio português brasileiro" — 100% local, zero nuvem. O pipeline de voz cobre: STT (áudio → texto), Processamento, e TTS (texto → áudio).

## Pipeline

```
Usuário envia áudio (Telegram)
    │
    ▼
[STT] Whisper Local (:8020)  ←─ faster-whisper-server
    │  fallback: Groq Cloud
    ▼
Transcrição em texto PT-BR
    │
    ▼
[Processamento] Aurélia (qwen3.5:9b via LiteLLM cascade)
    │  local → qwen3.6-free → minimax-free → groq → paid
    ▼
Resposta em texto limpo
    │
    ▼
[TTS] Kokoro Local (:8012)  ←─ ghcr.io/remsky/kokoro-fastapi-gpu
    │  voz: pt-br_isabela | formato: opus → .ogg (voice note)
    ▼
Envio em paralelo:
  1. Texto imediato (Telegram)
  2. Áudio em goroutine (voice note)
```

## STT — Whisper Local (`pkg/stt/factory.go`)
```go
// modo "local": faster-whisper-server → Groq fallback
// modo "" / "groq": Groq → faster-whisper fallback
```
Porta: `localhost:8020`, Modelo: `Systran/faster-whisper-large-v3`, Idioma: `pt`

## TTS — Kokoro Local (`pkg/tts/openai_compatible.go`)
```go
OpenAICompatibleSynthesizer{baseURL: "http://localhost:8012", model: "kokoro"}
// Voice: pt-br_isabela (feminino PT-BR)
// Format: opus → .ogg (Telegram voice note)
```
Container: `ghcr.io/remsky/kokoro-fastapi-gpu:latest` → porta `8012:8880`

## Voice System Prompt Suffix
```go
const voiceSystemPromptSuffix = `
ATENÇÃO — MODO VOZ: Sem markdown, frases curtas,
números por extenso, sem blocos de código.
`
```

## Texto Sanitizado para TTS (`output.go sanitizeTextForSpeech`)
- Remove blocos de código, markdown decoration
- Expande símbolos: `R$`→"reais", `%`→"por cento"
- Trunca no limite de sentença

## Consequências
- **Positivo**: 100% local, zero custo de API, voz PT-BR com Kokoro
- **Negativo**: VRAM competition (Whisper 12.7GB + qwen3.5 6.6GB + Kokoro ~2GB); ~21.3GB / 24.5GB usados
