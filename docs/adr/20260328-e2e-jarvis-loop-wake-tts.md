# ADR 20260328: End-to-End Jarvis Loop (Wake Word → TTS Response)

## Status
🟢 Proposto (P3 - Polish)

## Contexto
O Jarvis Loop completo é: **Wake Word → Audio Capture → STT → LLM (gemma3 27b) → TTS (Kokoro) → Audio Response**. Cada componente existe isoladamente, mas falta o fluxo integrado e o teste E2E.

## Decisões Arquiteturais

### 1. Fluxo Integrado

```
┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
│Wake Word│──▶│Capture  │──▶│  STT    │──▶│  LLM    │──▶│  TTS    │
│"Jarvis" │   │  Audio  │   │(Groq+Wh)│   │(gemma3) │   │(Kokoro)│
└─────────┘   └─────────┘   └─────────┘   └─────────┘   └─────────┘
                                                                │
                                                                ▼
                                                         ┌─────────┐
                                                         │  Play   │
                                                         │  Audio  │
                                                         └─────────┘
```

### 2. Conversation Context Preservation

O loop precisa manter contexto entre interações:

```go
// internal/voice/context.go
type ConversationContext struct {
    UserID        int64
    ChatID       int64
    History      []Message
    LastIntent   string
    LastEntities map[string]string
    TTL          time.Duration  // 5 minutos de timeout
}

const defaultContextTTL = 5 * time.Minute

func (c *ConversationContext) Add(user, assistant string) {
    c.History = append(c.History, Message{
        Role:      "user",
        Content:   user,
        Timestamp: time.Now(),
    }, Message{
        Role:      "assistant",
        Content:   assistant,
        Timestamp: time.Now(),
    })
}
```

### 3. TTS Streaming (Não Batch)

O TTS atual é batch (gera arquivo → envia). Para UX responsivo:

```go
// internal/tts/stream.go
type TTSStream struct {
    kokoro  *KokoroProvider
    buffer  *bytes.Buffer
    chunkSize int
}

func (t *TTSStream) Stream(ctx context.Context, text string) (<-chan []byte, error) {
    chunks := make(chan []byte, 10)

    go func() {
        defer close(chunks)

        // Chunk o texto em partes menores
        for _, chunk := range segmentText(text, 1200) {
            audio, err := t.kokoro.Synthesize(ctx, chunk)
            if err != nil {
                logger.Error("TTS chunk failed", "error", err)
                continue
            }
            select {
            case chunks <- audio:
            case <-ctx.Done():
                return
            }
        }
    }()

    return chunks, nil
}
```

### 4. Teste E2E

```go
// e2e/jarvis_loop_test.go
func TestJarvisLoop(t *testing.T) {
    // 1. Setup
    ctx := context.Background()
    cleanup := startTestServices(t, ctx)
    defer cleanup()

    // 2. Simula wake word
    audio := generateWakeWordAudio(t, "jarvis", "me mostre o clima")
    spoolPath := submitAudio(t, ctx, audio)

    // 3. Espera processamento
    transcript := waitForTranscript(t, ctx, spoolPath, 30*time.Second)
    require.Equal(t, "me mostre o clima", transcript)

    // 4. Verifica TTS output
    audioResponse := waitForAudio(t, ctx, 60*time.Second)
    require.NotEmpty(t, audioResponse)

    // 5. Valida formato
    validateAudioFormat(t, audioResponse, "mp3", 44100)
}
```

### 5. Conversation API

```go
// internal/voice/api.go
func (v *VoiceAPI) handleVoiceInput(ctx context.Context, req VoiceInputRequest) (*VoiceResponse, error) {
    // 1. Transcribe
    transcript, err := v.stt.Transcribe(ctx, req.Audio)
    if err != nil {
        return nil, fmt.Errorf("transcription: %w", err)
    }

    // 2. Get context
    convCtx := v.contextStore.Get(req.UserID, req.ChatID)
    defer v.contextStore.Set(req.UserID, req.ChatID, convCtx)

    // 3. LLM inference
    response, err := v.llm.Chat(ctx, convCtx.History, transcript)
    if err != nil {
        return nil, fmt.Errorf("llm: %w", err)
    }

    // 4. Update context
    convCtx.Add(transcript, response)

    // 5. TTS
    audioURL, err := v.tts.SynthesizeURL(ctx, response)
    if err != nil {
        return nil, fmt.Errorf("tts: %w", err)
    }

    return &VoiceResponse{
        Transcript: transcript,
        Response:   response,
        AudioURL:   audioURL,
    }, nil
}
```

## Consequências

### Positivas
- Loop completo funcionando: wake → response
- Contexto preservado entre interações
- UX responsivo com streaming TTS
- Testes E2E garantem confiança

### Negativas
- Latência total: ~2-5 segundos por interação
- Contexto em memória pode crescer indefinidamente
- TTS streaming requer client support (Telegram não suporta streaming)

### Trade-offs
- Batch vs Streaming TTS: Batch é mais simples mas menos responsivo
- Contexto em memória vs Redis: Memória é mais rápida mas não persiste

## Dependências
- ✅ `internal/voice/processor.go` (STT loop)
- ✅ `pkg/tts/openai_compatible.go` (Kokoro TTS)
- ✅ `internal/gateway/provider.go` (LLM gateway)
- ⚠️ `internal/voice/context.go` (NÃO EXISTE - contexto)
- ⚠️ `internal/tts/stream.go` (NÃO EXISTE - streaming)
- ❌ `e2e/jarvis_loop_test.go` (NÃO EXISTE)

## Referências
- [ADR-20260328-jarvis-voice-computer-use.md](./20260328-implementacao-jarvis-voice-e-computer-use.md)
- [ADR-20260328-tts-br-portuguese-industrialization.md](./20260328-tts-br-portuguese-industrialization.md)
- [internal/voice/processor.go](../../internal/voice/processor.go)
- [pkg/tts/openai_compatible.go](../../pkg/tts/openai_compatible.go)
- [internal/gateway/provider.go](../../internal/gateway/provider.go)

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
