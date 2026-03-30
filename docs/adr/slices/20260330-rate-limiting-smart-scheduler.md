# ADR Slice: Rate Limiting Smart Scheduler — GPU Math + Cool-Down

## Contexto
Com whisper-local (~12.7GB VRAM), qwen3.5:9b (~6.6GB VRAM) e Kokoro (~2-4GB) competindo na RTX 4090 24GB, o VRAM livre cai para ~0.5-1.5GB. É necessário um sistema de rate limiting inteligente para evitar OOM e manter latência consistente.

## Hardware Profile

| Recurso | Capacidade | Usado | Livre | Estratégia |
|---|---|---|---|---|
| VRAM RTX 4090 | 24.5 GB | 19.3 GB | **~5 GB** | Concurrency + timeout |
| RAM | 30 GB | 21 GB | 9 GB | CGroup caps |
| CPU 7900X | 24 threads | — | Alto | Concurrency pools |
| NVMe Gen5 | 4 TB | — | Alto | Sem limit local |

## Decisões

### 1. Ollama — Systemd Drop-in
```bash
# /etc/systemd/system/ollama.service.d/rate-limit.conf
[Service]
Environment="OLLAMA_NUM_PARALLEL=2"
Environment="OLLAMA_MAX_LOADED_MODELS=2"
Environment="OLLAMA_KEEP_ALIVE=10m"
```
```bash
sudo systemctl daemon-reload && sudo systemctl restart ollama
```

### 2. LiteLLM Cascade
```yaml
router_settings:
  routing_strategy: priority
  num_retries: 3
  allowed_fails: 2
  cooldown_time: 30      # 30s antes de retry do mesmo modelo
  timeout: 90
  retry_after: 2
```

### 3. Whisper
- 1 concurrent request (serializa)
- Timeout: 120s
- Não mantém VRAM entre requests (stateless inference)

### 4. Kokoro
- 2 concurrent synthesis
- Timeout: 30s por synthesis
- Stateless → sem memory pressure

## Cool-Down Strategy

Quando GPU util > 80% OU temp > 50°C:
1. **Throttle Ollama**: NUM_PARALLEL → 1 temporariamente
2. **Bloquear whisper**: fila de espera com timeout 180s
3. **LiteLLM routeia para Groq Free**: cloud enquanto GPU esfria

## Consequências
- **Positivo**: Zero OOM, latência previsível, VRAM nunca estoura
- **Negativo**: Throughput reduzido em burst (GPU saturada = respostas mais lentas)
- **Trade-off**: Local vs Cloud — sempre que GPU saturada, cloud entra como fallback
