# Arquitetura Soberana (Aurelia SOTA 2026)

## 1. Stack
```text
╔══════════════════════════════════════════════════════════════════════╗
║  AURÉLIA SOVEREIGN OS [v2026.04.01]                                  ║
╠══════════════════════════════════════════════════════════════════════╣
║  HOST: will-zappro        GPU: RTX 4090 [24GB]                       ║
║  STT: Whisper large-v3   [GPU] @ localhost:8020                      ║
║  TTS: Edge TTS (GRÁTIS)  pt-BR-ThalitaMultilingualNeural             ║
║  TTS: Kokoro-82M ONNX    [GPU] @ localhost:8012 (fallback)           ║
║  LLM: gemma3:27b-it-qat [Ollama] @ localhost:11434                   ║
║  VL: Qwen3.5-9B Vision  [Ollama] @ localhost:11434                   ║
╚══════════════════════════════════════════════════════════════════════╝
```

## 2. Fluxo de Mensagem (Texto & Roteamento)
1. **Telegram Retry** -> **Redis Deduplication** (Bloqueio).
2. **Payload Autorizado** -> **Porteiro IsSafe** (Valida intenção).
3. **Módulo de Agente** -> **Tier 1 LLM** (Raciocínio & Orquestração).
4. **Resposta Gerada** -> **Porteiro Polisher** (Markdown + Emoji + Caching).
5. **Saída Bruta** -> **Porteiro SecretMask** (Censura de chaves).
6. **Delivery** -> Usuário via Telegram API.

## 3. Fluxo de Áudio (Voz)
```text
[Áudio .wav do Telegram]
    │ POST :8020/v1/audio/transcriptions
    ▼
[Whisper large-v3 — STT GPU / Groq Fallback]
    │ texto
    ▼
[BGE-M3 :11434] → embed → [Qdrant :6333] → Contexto da Memória Semântica
    │
    ▼ POST :11434/api/chat
[Gemma 3 / Cloud LLM]
    │ resposta
    ▼ POST :8012/v1/audio/speech
[Kokoro TTS GPU / Edge TTS Microsoft]
    │
    ▼
[Áudio .ogg/mp3 via Telebot]
```
