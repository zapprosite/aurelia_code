# PLANO: Jarvis Tutor - Escuta Tudo 24/7

## AUDITORIA COMPLETA

| Componente | Status | Host:Porta | Modelo |
|------------|--------|------------|---------|
| LiteLLM | ✅ Healthy | localhost:4000 | qwen3.5 (tier 0) |
| Ollama | ✅ Ok | localhost:11434 | qwen3.5, qwen2.5, nomic |
| Kokoro TTS | ✅ Ok | localhost:8880 | pt-br female |
| Groq STT | ✅ | api.groq.com | whisper-large-v3 |

## ARQUITETURA

```
┌─────────────────────────────────────────────────────────────┐
│                    JARVIS TUTOR                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Mic (PulseAudio)  ──┐                                 │
│                       │                                 │
│                       ▼                                 │
│              ┌──────────────┐                          │
│              │ Audio Capture │  (parec → WAV)           │
│              └──────┬───────┘                          │
│                     │ WAV                               │
│                     ▼                                   │
│              ┌──────────────┐                          │
│              │  Groq STT    │  (whisper-large-v3)     │
│              └──────┬───────┘                          │
│                     │ texto                            │
│                     ▼                                   │
│              ┌──────────────┐                          │
│              │ LiteLLM     │  (qwen3.5)            │
│              │  Gateway    │                          │
│              └──────┬───────┘                          │
│                     │ resposta                         │
│                     ▼                                  │
│              ┌──────────────┐                          │
│              │ Kokoro TTS  │  (pt-br female)          │
│              │ localhost:8880│                          │
│              └──────┬───────┘                          │
│                     │ audio                           │
│                     ▼                                  │
│              ┌──────────────┐                          │
│              │  Speakers   │  (paplay/PulseAudio)       │
│              └──────────────┘                          │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## IMPLEMENTAÇÃO

### 1. Audio Capture (parec + sox)
### 2. STT (Groq API → whisper)
### 3. LLM (LiteLLM → qwen3.5)
### 4. TTS (Kokoro → localhost:8880)
### 5. Playback (paplay)

## EXECUÇÃO

1. Build com código limpo
2. Teste manual
3. Service systemd

