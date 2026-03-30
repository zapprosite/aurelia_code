# ADR Slice: Kokoro TTS GPU Local — 100% Sovereign Audio

## Contexto
O Kokoro TTS permite voz PT-BR local sem custo de API. A meta é 50-80% das respostas sintetizadas localmente na RTX 4090, com fallback opcional para cloud se VRAM estiver saturada.

## Decisão

### Container Docker
```yaml
# docker-compose.yml
kokoro-tts:
  image: ghcr.io/remsky/kokoro-fastapi-gpu:latest
  ports:
    - "8012:8880"   # porta host → porta interna do container
  environment:
    - NVIDIA_VISIBLE_DEVICES=all
    - CUDA_VISIBLE_DEVICES=0
```

**Nota**: O container sobe em CPU mode por padrão (mesmo a imagem GPU) — GPU acceleration depende do driver CUDA disponível no host. A imagem `ghcr.io/remsky/kokoro-fastapi-gpu` detecta automaticamente se CUDA está presente.

### API Endpoint
- Health: `GET http://localhost:8012/health` → `{"status":"healthy"}`
- Speech: `POST http://localhost:8012/v1/audio/speech` (OpenAI-compatible)
- Vozeiro: `pt-br_isabela` (voz PT-BR feminina)
- Formato: `opus` → `.ogg` (voice note Telegram)

### Integração Go
`pkg/tts/openai_compatible.go` já faz POST para `/v1/audio/speech` com formato OpenAI-compatible. Config default: `http://127.0.0.1:8012` (local).

### VRAM Budget
| Componente | VRAM |
|---|---|
| whisper-local | ~12.7 GB |
| qwen3.5:9b | ~6.6 GB |
| Kokoro (GPU) | ~2-4 GB |
| Overhead | ~2 GB |
| **Total** | ~23-25 GB / 24.5 GB |
| **Livre** | ~0.5-1.5 GB |

**⚠️ Quando Kokoro GPU + Whisper + qwen3.5 rodam simultaneamente, VRAM fica crítica.**

## Consequências
- **Positivo**: Voz local PT-BR sem custo de API; latência menor para TTS
- **Negativo**: VRAMcompetition com Whisper e qwen3.5; container precisa de warm-up (~15s)
- **Mitigação**: Rate limiting com `OLLAMA_NUM_PARALLEL=2` + cooldown 30s
