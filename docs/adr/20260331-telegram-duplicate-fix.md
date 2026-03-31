# ADR-20260331: Telegram Duplicate Response Fix

**Status:** ✅ Implementado
**Data:** 2026-03-31
**Autor:** Claude Code
**Slice:** S-33

---

## Problema

O Telegram envia retry requests quando não recebe ack dentro do timeout, causando **duas respostas iguais** ao usuário.

## Solução

**Deduplicação por Message ID** no handler de texto.

### Implementação

**Arquivos alterados:**
- `internal/telegram/bot.go` — adiciona `seenMessageIDs map[int]time.Time`
- `internal/telegram/input.go` — verifica duplicação antes de processar

```go
// bot.go
type BotController struct {
    seenMu       sync.Mutex
    seenMessageIDs map[int]time.Time // S-33: deduplication
}

func (bc *BotController) isDuplicateMessage(msgID int) bool {
    // Clean old entries (>60s)
    // Return true if seen, else mark and return false
}

// input.go
func (bc *BotController) handleText(c telebot.Context) error {
    if bc.isDuplicateMessage(c.Message().ID) {
        return nil // Skip duplicate
    }
    return bc.processInput(c, c.Text(), false)
}
```

### Janela de Deduplicação
- **60 segundos** — limpa automaticamente entries velhos
- **Thread-safe** — usa mutex para evitar race conditions

---

## Verificação

```bash
cd ~/aurelia && go build -o aurelia-linux ./cmd/aurelia/
docker restart aurelia-api
# Teste: envie mensagem ao bot — não deve duplicar
```

---

## Resultado

- Mensagens duplicadas são ignoradas silenciosamente
- Usuário recebe apenas 1 resposta
- Sem impacto na latência (verificação em memória)
