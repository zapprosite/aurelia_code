# ADR Slice: Voxtral STT GPU — Substitui Whisper CPU por Voxtral Mini 3B

## Status
Proposto

## Contexto
Whisper LARGE v3 está rodando em **CPU** (~12GB RAM, não VRAM), causando latência de 10-30s por transcrição. Kokoro TTS também em CPU (3-5x real-time, aceitável mas não ideal). A pesquisa mostra que o Kokoro **não tem voz PT-BR nativa**.

## Decisão

Substituir Whisper CPU pelo **Voxtral Mini 3B** rodando em GPU via vLLM Docker.

### Por que Voxtral Mini 3B
- WER: **4% no FLEURS** — melhor que Whisper Large v3 (~4.4%)
- VRAM: **~9.5GB bf16** — CABE junto com qwen3.5:9b (~9.2GB) = 18.7GB / 24GB
- Idiomas: **13 idiomas** incluindo **Português nativo**
- Velocidade: **3x mais rápido** que ElevenLabs Scribe v2
- Acesso: **Open weights** no Hugging Face (vLLM ou Transformers)
- Release: Fevereiro 2026

### VRAM Budget (RTX 4090 24GB)

| Serviço | VRAM | Acumulado | Status |
|---|---|---|---|
| Ollama (qwen3.5:9b) | ~9.2 GB | 9.2 GB | ✅ GPU quente |
| Voxtral Mini 3B (vLLM) | ~9.5 GB | 18.7 GB | ✅ GPU sob demanda |
| Kokoro TTS | ~0 GB | 18.7 GB | CPU (ONNX, 3-5x RT) |
| **TOTAL** | | **~18.7 GB** | **5.3 GB livre** |

## Implementação

### Voxtral Mini 3B via vLLM Docker

```yaml
# docker-compose.yml
  voxtral-stt:
    container_name: aurelia-voxtral
    image: ghcr.io/vllm/vllm-openai:latest
    ports:
      - "127.0.0.1:8030:8000"
    volumes:
      - ./data/voxtral:/root/.cache/huggingface
    environment:
      - NVIDIA_VISIBLE_DEVICES=all
      - CUDA_VISIBLE_DEVICES=0
      - VLLM_WORKER_MULTIPROC_METHOD=spawn
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
    restart: always
    healthcheck:
      test: ["CMD-SHELL", "python3 -c \"import httpx; r = httpx.get('http://localhost:8000/health', timeout=5); exit(0 if r.status_code==200 else 1)\""]
      interval: 30s
      timeout: 15s
      retries: 5
      start_period: 120s
    command: >
      --model mistralai/Voxtral-Mini-3B-2507
      --tokenizer_mode mistral
      --config_format mistrasl
      --load_format mistral
      --task transcribe
      --max_model_len 131072
      --gpu_memory_utilization 0.85
```

### STT Factory (`pkg/stt/factory.go`)

```go
// Prioridade: voxtral (GPU) → groq (cloud)
func NewTranscriber(provider, groqAPIKey, baseURL, model, language string) (Transcriber, error) {
    switch provider {
    case "voxtral", "mistral":
        t := NewVoxtralTranscriber(baseURL, language)
        if t.IsAvailable() {
            return t, nil
        }
        // Fallback Groq
        if groqAPIKey != "" {
            t2 := NewGroqTranscriber(groqAPIKey, baseURL, model, language)
            if t2.IsAvailable() {
                return t2, nil
            }
        }
        return nil, fmt.Errorf("voxtral unavailable e groq não configurado")
    case "groq":
        t := NewGroqTranscriber(groqAPIKey, baseURL, model, language)
        if t.IsAvailable() {
            return t, nil
        }
        // Fallback voxtral
        t2 := NewVoxtralTranscriber(baseURL, language)
        if t2.IsAvailable() {
            return t2, nil
        }
        return nil, fmt.Errorf("groq e voxtral indisponíveis")
    case "", "local":
        // voxtral primeiro (GPU), groq como backup
        t := NewVoxtralTranscriber(baseURL, language)
        if t.IsAvailable() {
            return t, nil
        }
        if groqAPIKey != "" {
            t2 := NewGroqTranscriber(groqAPIKey, baseURL, model, language)
            if t2.IsAvailable() {
                return t2, nil
            }
        }
        return nil, fmt.Errorf("voxtral e groq indisponíveis")
    default:
        return nil, fmt.Errorf("provider %q não suportado", provider)
    }
}
```

### VoxtralTranscriber (`pkg/stt/voxtral.go`)

```go
type VoxtralTranscriber struct {
    baseURL   string
    language string
}

func NewVoxtralTranscriber(baseURL, language string) *VoxtralTranscriber {
    if baseURL == "" {
        baseURL = "http://localhost:8030"
    }
    return &VoxtralTranscriber{baseURL: baseURL, language: language}
}

func (v *VoxtralTranscriber) Transcribe(ctx context.Context, audio io.Reader, format string) (string, error) {
    // vLLM OpenAI-compatible: POST /v1/audio/transcriptions
    // multipart/form-data com field "file"
    // Model: mistralai/Voxtral-Mini-3B-2507
    // Response: {"text": "..."}
}

func (v *VoxtralTranscriber) IsAvailable() bool {
    // GET /health → 200 OK
}
```

## Remoção (Legacy Prune)

| Componente | Ação |
|---|---|
| `whisper-local` container | **DELETE** — não precisa mais do faster-whisper-server |
| `fedirz/faster-whisper-server:latest-cuda` image | **DELETE** do docker |
| `pkg/stt/local.go` (faster-whisper) | **DELETE** ou **DEPRECATE** |
| STT factory `faster-whisper` case | **REMOVE** ou marcar como deprecated |
| whisper volumes em compose | **DELETE** |

## Consequências
- **Positivo**: STT em GPU, latência ~3x menor, Português nativo, VRAM gerenciada (18.7GB/24GB)
- **Negativo**: Container vLLM de ~15GB; primeira carga do modelo pode demorar 2-5min
- **Risco**: vLLM pode conflitar com Ollama se ambos tentarem alocar toda GPU — mitigado com `--gpu_memory_utilization 0.85`
