# ADR Slice: Jarvis Local Pinned — Audio + Vision 100% Local RTX 4090

## Status
✅ Implementado (2026-03-30)

## Contexto
STT (Whisper) e TTS (Kokoro) devem rodar 24/7 como padrão local na RTX 4090, sem dependência de nuvem. Audio + Vision juntos consomem ~14.7GB VRAM, deixando ~9.8GB para o qwen3.5:9b + overhead.

## VRAM Budget (RTX 4090 24GB)

| Serviço | VRAM | Acumulado | Notas |
|---|---|---|---|
| Whisper large-v3 | ~12.7 GB | 12.7 GB | STT — maior consumidor |
| Kokoro TTS | ~2.0 GB | 14.7 GB | TTS GPU |
| qwen3.5:9b | ~6.6 GB | 21.3 GB | Se Kokoro unload → carrega |
| Overhead/scratch | ~3.3 GB | 24.6 GB | Near capacity! |
| **TOTAL** | | **~21.3 GB** | Com Whisper + qwen3.5 + Kokoro |

> ⚠️ Cuidado: Whisper (12.7GB) + Kokoro (2GB) + qwen3.5 (6.6GB) = 21.3GB — muito próximo do limite. Se Whisper não estiver em uso, Kokoro pode carregar e qwen3.5 fica com mais headroom.

## Padrão pinned

### STT — Whisper Local (faster-whisper-server)
```
Container: fedirz/faster-whisper-server:latest-cuda
Porta:     127.0.0.1:8020 → :8000 (bind localhost, exposto apenas local)
VRAM:      ~12.7 GB
Health:    GET /health → "OK"
Modelo:    Systran/faster-whisper-large-v3
Lingua:    pt (prioridade)
Fallback:  Groq Whisper (cloud, gratuito)
Restart:   always
GPU:       nvidia reservation (1 device)
```

### TTS — Kokoro Local (GPU)
```
Container: ghcr.io/remsky/kokoro-fastapi-gpu:latest
Porta:    localhost:8012 → :8880
VRAM:     ~2 GB
Health:    GET /health → {"status":"healthy"}
Voice:    pt-br_isabela (feminino PT-BR)
Format:   opus → .ogg (Telegram voice note)
Restart:  always
GPU:       nvidia reservation (1 device)
```

### Vision — qwen3.5:9b VL (Ollama local)
```
Modelo:     ollama/qwen3.5:9b
VRAM:       ~6.6 GB
Timeout:    45s
extra_body: {think: false} — garante content no response
```

## Fallback Cascade

```
Audio/Vision Request
    │
    ▼
[Whisper :8020] ── fail (3x) ──→ [Groq Whisper cloud]
    │
    ▼
[qwen3.5:9b via LiteLLM]
    │
    ├── Ollama local (rpm=1, least-busy prefere)
    ├── fallback → qwen-3.6-plus-preview:free (OpenRouter)
    ├── fallback → minimax-2.5:free (OpenRouter)
    ├── fallback → groq/llama-3.3-70b-versatile
    └── fallback → minimax-m2.7 (OpenRouter paid)
    │
    ▼
[Kokoro :8012] ── fail (3x) ──→ [Texto puro — sem TTS]
```

## Endurecimento Docker

### Health checks
```yaml
# smart-router — wget instalado via apk (container alpine)
healthcheck:
  test: ["CMD-SHELL", "wget -qO- http://localhost:4000/health || exit 1"]

# Kokoro — healthcheck já funcional
healthcheck:
  test: ["CMD-SHELL", "curl -f http://localhost:8880/health || exit 1"]

# Whisper — healthcheck funcional via curl :8020
healthcheck:
  test: ["CMD-SHELL", "curl -f http://localhost:8000/health || exit 1"]
```

### Restart policies
```yaml
restart: always   # todos os serviços críticos
```

## Consequências
- **Positivo**: Audio + Vision 100% local, zero custo de API, VRAM gerenciada
- **Negativo**: VRAM no limite (~21.3GB / 24.5GB) — monitoramento obrigatório
- **Risco**: Se Whisper + Kokoro + qwen3.5 todos carregados simultaneamente, OOM
