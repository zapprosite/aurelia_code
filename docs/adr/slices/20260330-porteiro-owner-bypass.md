# ADR Slice: Porteiro Sentinel — Owner Bypass

## Contexto
O `ProcessExternalInput` (rota `/v1/telegram/impersonate`) não tinha bypass para o dono. Mensagens do owner passavam pelo Porteiro (qwen2.5:0.5b) e eram bloqueadas por falso positivo (ex: "qual a diferença entre goroutine e thread?" → cache Redis retornava `UNSAFE`). O bot Telegram normal tinha bypass, mas a rota API não.

## Decisão

### Antes (`internal/telegram/input_pipeline.go`)
```go
if bc.porteiro != nil && !requiresAudio {
    safe, err := bc.porteiro.IsSafe(ctx, text)
    // ... bloqueia owner indiscriminadamente
}
```

### Depois
```go
ownerBypass := false
for _, id := range bc.allowedUserIDs {
    if id == userID {
        ownerBypass = true
        break
    }
}
if !ownerBypass {
    safe, err := bc.porteiro.IsSafe(ctx, text)
    // ... só não-owner passa pelo guard
}
```

### Cache Redis
- Cache do Porteiro em Redis com TTL 30 dias
- Mensagens falsas positivas ficam cacheadas como `UNSAFE`
- Flush: `docker exec aurelia-redis-1 redis-cli KEYS "porteiro:cache:*" | xargs docker exec aurelia-redis-1 redis-cli DEL`

## Consequências
- **Positivo**: Owner nunca bloqueado, mesmo com cache envenenado
- **Negativo**: Não aplicável — owner bypass é estritamente seguro (mesmos IDs do Telegram jávalidados no auth)
