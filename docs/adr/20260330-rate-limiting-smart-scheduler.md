# ADR 20260330: Rate Limiting Inteligente — Hardware Math + Cool-Down

## Status
✅ Aceito

## Hardware Profile (will-zappro)

| Recurso | Capacidade | Usado | Livre | Strategy |
|---|---|---|---|---|
| **VRAM** | 24.5 GB | ~19.3 GB | **11.3 GB** | Resfriamento por concurrency + timeout |
| **RAM** | 30 GB | 21 GB | 2.4 GB | CGroup memory cap nos containers |
| **CPU** | 24 threads (7900X) | Variável | Alto |Concurrency pools por container |
| **NVMe Gen5** | 4 TB (seq: 14GB/s) | — | Alto | Sem rate limit local |

## GPU Math

```
VRAM Budget:
  whisper-local (faster-whisper-large-v3): ~12.7 GB  ← MAIOR CONSUMIDOR
  qwen3.5:9b (Q4_K_M):                          ~6.6 GB
  Overhead & KV cache:                          ~2 GB
  ─────────────────────────────────────────────
  Total:                                        ~21.3 GB
  Livre:                                        ~3.2 GB  ← CRÍTICO

  ⚠️ Com ambos ativos,只剩 3.2 GB livre.
  Isso não é suficiente para KV cache de burst.
```

## Estratégia de Rate Limiting

### 1. Ollama (qwen3.5:9b + nomic-embed-text)

**Config:** `~/.ollama/.env`
```
OLLAMA_NUM_PARALLEL=2       # 2 requests simultâneos (resfria GPU)
OLLAMA_MAX_LOADED_MODELS=2  # Mantém apenas 2 modelos carregados
OLLAMA_KEEP_ALIVE=10m       # Descarrega após 10min idle
```

### 2. Whisper (faster-whisper-server)

**Config:** `faster-whisper-server` já não carrega VRAM persistentemente — uma inference por vez. Rate limit por timeout.

```
Whisper timeout: 120s
Concurrent: 1 (serializa requisições)
```

### 3. Kokoro TTS

```
Concurrent: 2
Timeout: 30s por synthesis
```

### 4. LiteLLM Cascade

```yaml
router_settings:
  num_retries: 3
  allowed_fails: 2
  cooldown_time: 30          # 30s antes de retry do mesmo modelo
  timeout: 90                # Timeout por request
  retry_after: 2             # Backoff inicial 2s
```

### 5. Groq STT

```
Concurrent: 1 (STT é sequencial no pipeline)
Timeout: 60s
```

## Cool-Down Strategy (Resfriamento)

Quando GPU temp > 50°C OU utilization > 80%:

1. **Throttle Ollama**: `OLLAMA_NUM_PARALLEL=1` temporariamente
2. **Bloquear whisper**: fila de espera, timeout 180s
3. **LiteLLM routeia para cloud**: Groq Free enquanto GPU esfria

Monitoramento via `nvidia-smi` + Prometheus metrics.

## Consequências
- **Positivo**: GPU não estoura, latency consistente, zero OOM
- **Negativo**: Throughput reduzido quando ambos modelos ativos simultaneamente
