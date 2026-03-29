# ADR 20260328: Otimização de GPU Budget (Whisper Medium + Groq STT)

## Status
🟡 Proposto (P2)

## Contexto
O ADR `20260328-multimodal-gpu-optimization` identifica que a carga combinada de Gemma3-27B (Ollama), Whisper-Large-v3 (STT) e Kodoro (TTS) excede os 24GB de VRAM da RTX 4090. Precisamos implementar o plano de alívio.

## Decisões Arquiteturais

### 1. STT Strategy: Groq Primary + Local Fallback

```go
// internal/voice/stt.go
package voice

type STTConfig struct {
    Primary   Transcriber  // Groq API (cloud)
    Fallback  Transcriber  // Whisper local (VRAM)
}

// Prioridade: Groq > Local
// Se Groq retorna 429 ou está down → fallback automático
```

```go
// Configuração em app.json
{
  "stt": {
    "primary": "groq",
    "fallback": "whisper-local",
    "groq_model": "whisper-large-v3",
    "local_model": "medium",  // Mudança de large-v3 para medium
    "language": "pt-BR"
  }
}
```

### 2. Whisper Model: Large → Medium

**Motivação**: Economia de ~2GB VRAM

| Modelo | VRAM | WER (pt-BR estimado) | Speed |
|--------|------|----------------------|-------|
| `large-v3` | ~3GB | ~5% | 1x |
| `medium` | ~1.5GB | ~7% | ~1.3x |
| `small` | ~500MB | ~10% | ~2x |

**Trade-off aceito**: Perda marginal de precisão (~2% WER) compensada pela estabilidade.

```bash
# Download do modelo menor
ollama pull whisper:medium

#Remover o large se não usar
# ollama rm whisper:large-v3
```

### 3. Hybrid-Cloud Fallback Chain

```go
// internal/voice/processor.go - transcribeWithBudget()
func (p *Processor) transcribeWithBudget(ctx context.Context, audioPath string) (string, error) {
    // 1. Tenta Groq (primário)
    if p.primary != nil && p.primary.IsAvailable() {
        transcript, err := p.primary.Transcribe(ctx, audioPath)
        if err == nil {
            return transcript, nil
        }
        // 2. Se Groq falhou, tenta local
        if p.fallback != nil && p.fallback.IsAvailable() {
            p.recordFallback("groq_fail")
            return p.fallback.Transcribe(ctx, audioPath)
        }
        return "", err
    }

    // 3. Groq down, vai direto para local
    if p.fallback != nil && p.fallback.IsAvailable() {
        p.recordFallback("groq_unavailable")
        return p.fallback.Transcribe(ctx, audioPath)
    }

    return "", fmt.Errorf("no STT available")
}
```

### 4. Silent Fallback UX

O ADR original pede "UX silencioso": tentar segundo provedor sem incomodar o usuário:

```go
// internal/telegram/messages.go
func (p *TelegramBot) handleTranscriptionError(ctx context.Context, err error, audioPath string) {
    // Não expõe erro ao usuário se fallback funcionar
    if transcript, fallbackErr := p.fallbackSTT.Transcribe(ctx, audioPath); fallbackErr == nil {
        p.processTranscript(ctx, transcript)
        // Log interno, mas usuário não precisa saber
        logger.Info("STT fallback successful", "error", err.Error())
        return
    }

    // Só expõe se ambos falharem
    p.sendError(ctx, "Erro de Transcrição. Tente novamente.")
}
```

### 5. VRAM Monitoring

```go
// internal/gpu/monitor.go
type GPUMonitor struct {
    memoryUsed   prometheus.Gauge
    memoryTotal  prometheus.Gauge
    temperature  prometheus.Gauge
}

func (m *GPUMonitor) CheckCapacity(requiredGB float64) bool {
    used, total, _ := m.GetMemoryInfo()
    available := total - used
    return available >= requiredGB
}

// Thresholds para alertas
const (
    VRAMWarningThreshold = 20.0  // GB
    VRAMCriticalThreshold = 22.0 // GB
)
```

## Consequências

### Positivas
- Fim de `CUDA Out of Memory` em transcriptions
- Groq é extremamente rápido (~300ms por áudio)
-Fallback local garante operação offline
- Preserva VRAM para Gemma3 27B

### Negativas
- Groq requer API key e tem rate limits
- Modelo medium tem WER ligeiramente maior
- Dependência de internet para Groq (primário)

### Trade-offs
- Precisão vs. Estabilidade: medium é "bom o suficiente"
- Custo vs. Soberania: Groq é cloud, mas fallback é local

## Dependências
- ⚠️ `internal/voice/processor.go` (já existe, precisa de refactor)
- ⚠️ `internal/telegram/messages.go` (precisa de error handling refactored)
- ❌ `groq` Go SDK ou REST calls
- ❌ `ollama pull whisper:medium`
- ❌ Environment var: `GROQ_API_KEY`

## Referências
- [ADR-20260328-multimodal-gpu-optimization.md](./20260328-multimodal-gpu-optimization.md)
- [ADR-20260328-tts-br-portuguese-industrialization.md](./20260328-tts-br-portuguese-industrialization.md)
- [internal/voice/processor.go](../../internal/voice/processor.go)
- [internal/telegram/messages.go](../../internal/telegram/messages.go)

## Links Obrigatórios
- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)

---
**Data**: 2026-03-28
**Status**: Proposto
**Autor**: Claude (Principal Engineer)
**Slice**: feature/neon-sentinel
**Progress**: 0%
