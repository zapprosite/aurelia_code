# ADR-20260331-DUP: Correção de Respostas Duplicadas

**Status:** ✅ Implementado
**Data:** 2026-03-31
**Problema:** Bot `aurelia_code` enviava texto + nota de voz separada = "duas respostas"

---

## Problema

Quando o usuário enviava uma mensagem de **texto**, o bot enviava:
1. Texto da resposta
2. Nota de voz com a mesma resposta

Isso acontecia porque `deliverFinalAnswer` forçava `requiresAudio = true` sempre que TTS estava disponível.

---

## Solução

Corrigido em dois arquivos:

### 1. `internal/telegram/input_pipeline.go`

```go
// ANTES (bug):
func (bc *BotController) deliverFinalAnswer(...) {
    if bc.tts != nil && bc.tts.IsAvailable() {
        requiresAudio = true  // ❌ FORÇAVA audio para TODAS mensagens
    }
    return bc.deliverWithParallelTTS(...)
}

// DEPOIS (corrigido):
func (bc *BotController) deliverFinalAnswer(...) {
    // S-34: Only send TTS audio when user sent a voice message.
    // Text messages get text-only responses to avoid duplicate outputs.
    return bc.deliverWithParallelTTS(..., requiresAudio)
}
```

### 2. `internal/telegram/output.go`

```go
// ANTES:
func deliverWithParallelTTS(sender, chat, synthesizer, text string, opts ...interface{}) error {
    if synthesizer != nil && synthesizer.IsAvailable() {  // ❌ sempre sintetizava
        // synthesize...
    }
}

// DEPOIS:
func deliverWithParallelTTS(..., requiresAudio bool, opts ...interface{}) error {
    // Only synthesize audio when user sent a voice message.
    if requiresAudio && synthesizer != nil && synthesizer.IsAvailable() {
        // synthesize...
    }
}
```

---

## Comportamento Resultante

| Input do Usuário | Resposta do Bot |
|------------------|-----------------|
| Mensagem de texto | Apenas texto |
| Mensagem de voz | Texto + nota de voz |

---

## Arquivos Alterados

```
internal/telegram/input_pipeline.go  (deliverFinalAnswer)
internal/telegram/output.go         (deliverWithParallelTTS)
internal/telegram/bootstrap.go      (chamada atualizada)
```

---

## Deploy

```bash
go build -o aurelia-linux ./cmd/aurelia
docker cp aurelia-linux aurelia-api:/usr/local/bin/aurelia
docker restart aurelia-api
```
