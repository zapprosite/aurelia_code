---
type: doc
name: testing-strategy
description: Test frameworks, patterns, coverage requirements, and quality gates
category: testing
generated: 2026-03-30
status: filled
scaffoldVersion: "2.0.0"
---

# Estratégia de Testes — Aurélia Sovereign 2026.2

## Stress Test via API

O teste principal é via rota HTTP `/v1/telegram/impersonate` — envia mensagem como usuário mock e valida resposta.

```bash
BASE="http://localhost:8585/v1/telegram/impersonate"
ID=7220607041  # Telegram ID do owner

# Teste unitário rápido
curl -X POST "$BASE" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":$ID,\"message\":\"qual a diferença entre goroutine e thread?\"}"

# Saída esperada: {"status":"ok", "message":"Pipeline processing started"}
```

## Validação de Saúde

```bash
# Health do bot
curl http://localhost:8585/health

# Health de cada serviço
curl http://localhost:8585/health

# LiteLLM cascade
curl http://localhost:4000/health

# Whisper STT
curl http://localhost:8020/health

# Kokoro TTS
curl http://localhost:8012/health

# Redis
docker exec aurelia-redis-1 redis-cli ping
# → PONG

# Qdrant
curl http://localhost:6333/healthz
```

## Build Verification

```bash
CGO_ENABLED=0 go build -trimpath -o bin/aurelia ./cmd/aurelia
# Saída vazia = build passou
```

## Critérios de Gate

- `go build` sem erros
- Todos os serviços Docker healthy (`docker ps`)
- Health check do bot retorna `{"status":"ok"}`
- Stress test de 6 mensagens sem falso positivo do Porteiro
- Commit passa em `scripts/audit/audit-secrets.sh`
