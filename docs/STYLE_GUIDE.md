# Aurélia Go Style Guide

> Regras de estilo Go para o ecossistema Aurélia. Atualizar este ficheiro quando novas regras de implementação forem adotadas.

## 1. Error Handling

```go
// ✅ Sempre trate erros — nunca ignore com _
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething: %w", err)
}

// ✅ Erros estruturados com contexto
return nil, fmt.Errorf("process audio chunk %d/%d: %w", i, total, err)

// ❌ Nunca ignore erros com _
data, _ := os.ReadFile(path) // PRUNE neste padrão
```

## 2. Logging

```go
// ✅ Usar slog do pacote internal/observability
import "github.com/kocar/aurelia/internal/observability"

// ✅ Níveis corretos
slog.Info("component_action", "key", value)
slog.Warn("degraded_mode", "reason", reason)
slog.Error("operation_failed", "error", err)

// ❌ Não usar log.Printf
log.Printf("something happened: %v", err)
```

## 3. Configuração

```go
// ✅ Usar internal/config — nunca hardcoded
import "github.com/kocar/aurelia/internal/config"

val := cfg.GetString("FIELD_NAME")  // dobra env var + JSON

// ❌ Não hardcoded valores
const MaxRetries = 3 // PRUNE — use config
```

## 4. Context

```go
// ✅ Sempre passe ctx como primeiro argumento
func Process(ctx context.Context, input string) error

// ✅ Check context cancellation em loops longos
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // continue
}
```

## 5. Naming Conventions

```go
// ✅ PascalCase para tipos e interfaces
type AudioBuffer struct { ... }
type Synthesizer interface { ... }

// ✅ snake_case para funções, variáveis, campos
func new_audio_buffer() { ... }
chunk_size := 4096

// ✅ Acrónimos em caps: URL, ID, API, HTTP
userID string
baseURL string

// ❌ Não use camelCase em variáveis Go
// userId, baseUrl, HttpClient — PRUNE estes padrões
```

## 6. Tests

```go
// ✅ Ficheiros de teste no mesmo package com sufixo _test.go
// internal/audio/buffer_test.go

// ✅ Table-driven tests para casos múltiplos
func TestBuffer(t *testing.T) {
    tests := []struct {
        name  string
        input int
        want  error
    }{
        {"valid", 1024, nil},
        {"negative", -1, ErrBufferSize},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewBuffer(tt.input)
            if err != tt.want {
                t.Errorf("NewBuffer() = %v, want %v", err, tt.want)
            }
        })
    }
}

// ✅ Mock de interfaces, não de implementações concretas
// Use interfaces para facilitar teste (STT, TTS, LLM providers)
```

## 7. Structured Errors

```go
// ✅ Erros descritivos com stack preservado
return nil, fmt.Errorf("audio/synthesize: chunk %d/%d: %w", i, n, err)

// ✅ Errors são valores — defina Sentinel errors para packages
var ErrBufferFull = errors.New("buffer full")
var ErrTimeout = errors.New("timeout")

// ❌ Não retorne errors.New em loops — cacheie
// ❌ Não ignore erros silently
```

## 8. Imports

```go
// ✅ Imports agrupados: stdlib → external → internal
import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/mark3labs/mcp-go/mcp"

    "github.com/kocar/aurelia/internal/config"
    "github.com/kocar/aurelia/internal/observability"
)

// ✅ Use alias apenas quando necessário
import (
    llm "github.com/kocar/aurelia/pkg/llm"  // evitar conflito de nome
)
```

## 9. Concurrency

```go
// ✅ Context para cancellation — goroutines devem morrer com o request
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

// ✅ Channel ownership: quem cria fecha
ch := make(chan Result)
go producer(ctx, ch)
for r := range ch { ... }

// ❌ Não feche channels de outros
// ❌ Não escreva em channels depois de close
```

## 10. Documentation

```go
// ✅ Doc comments em todas as funções exportadas
// Synthesize converts text to audio using the configured TTS provider.
// Returns raw WAV bytes. Chunking is handled internally.
func (s *Synthesizer) Synthesize(ctx context.Context, text string) ([]byte, error)

// ✅ Explicar decisões não-óbvias
// Chunking strategy: paragraphs → newlines → sentences → commas → words
// This order prioritizes natural phrasing over raw character limits.
```

## 11. No Placeholder Code

```go
// ❌ NUNCA cometa TODO/FIXME/HACK
// Se não está implementado: abra issue ou use stub com erro claro
func (t *LocalTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
    return "", fmt.Errorf("local transcription not implemented: use groq or openai provider")
}

// ✅ Stubs devem ser úteis — retornar erro descritivo
// Não deixar código commented-out sem explicação
```

## 12. Package Structure

```
internal/
  ├── agent/          # domain logic
  ├── audio/          # audio pipeline
  ├── config/         # configuration (single source of truth)
  ├── observability/   # logging + metrics
  ├── skill/          # skill system
  └── ...
pkg/
  ├── llm/            # LLM adapters (OpenAI, OpenRouter, Groq)
  ├── stt/            # STT adapters (Whisper, Groq)
  └── tts/            # TTS adapters (Kokoro, OpenAI)
```

---

*Última atualização: 2026-03-31 — criado como parte do cleanup MCP/A2A 2026*
