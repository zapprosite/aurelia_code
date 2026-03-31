# ADR-20260331: Telegram Response Fix (S-33 + S-34)

**Status:** ✅ Implementado
**Data:** 2026-03-31
**Slices:** S-33 (dedup), S-34 (TTS always-on)

---

## S-33: Deduplicação de Mensagens

### Problema
Telegram retry requests causavam duas respostas iguais.

### Solução
```go
// bot.go
seenMessageIDs map[int]time.Time
isDuplicateMessage(msgID int) bool // 60s window
```

---

## S-34: TTS Sempre Ativo

### Problema
TTS não estava enviando áudio quando `requiresAudio=false`.

### Solução
```go
// input_pipeline.go
func deliverFinalAnswer(...) {
    if bc.tts != nil && bc.tts.IsAvailable() {
        requiresAudio = true  // FORÇAR audio
    }
    return deliverWithParallelTTS(...)
}
```

### Config TTS (.env)
```
TTS_PROVIDER=openai_compatible
TTS_BASE_URL=http://127.0.0.1:8012
TTS_MODEL=kokoro
TTS_VOICE=pt-br
TTS_FORMAT=opus
```

---

## Arquivos Alterados
```
internal/telegram/bot.go         (+ deduplicação)
internal/telegram/input.go        (+ dedup check)
internal/telegram/input_pipeline.go (+ TTS always-on)
.env                           (+ TTS config)
```

---

## Deploy
```bash
docker cp aurelia-linux aurelia-api:/usr/local/bin/aurelia
docker restart aurelia-api
```
