---
title: Runtime Without Gemini Blueprint
status: active
created: 2026-03-19
owner: codex
scope: deploy-runtime-without-gemini
---

# Runtime Without Gemini Blueprint

## Objetivo

Voltar o runtime de deploy para um caminho mais simples e mais previsivel:

- `OpenRouter/Minimax` como LLM remoto principal
- `Groq` como STT
- sem `Gemini API` no caminho ativo
- health com foco no provedor primario real

## Decisao

Nao usar Gemini neste slice de deploy.

Diretriz operacional:

- manter `llm_provider=openrouter`
- manter `llm_model=minimax/minimax-m2.7`
- manter `stt_provider=groq`
- remover `google_api_key` do runtime local
- health deve provar apenas o que esta realmente em uso

## Config Esperada

Representacao polida do estado esperado da config:

```json
{
  "llm_provider": "openrouter",
  "llm_model": "minimax/minimax-m2.7",
  "stt_provider": "groq",
  "groq_api_key": "SET",
  "google_api_key": "ABSENT"
}
```

## Health Esperado

`GET /health` deve responder com:

- `status: ok`
- `checks.primary_llm.status: ok`
- sem `checks.gemini_api`

## Gate de Promocao

Para autorizar promocao:

1. `go test ./... -count=1`
2. `go test ./internal/agent -count=5`
3. `./scripts/build.sh`
4. daemon live em `/home/will/aurelia-24x7`
5. `/health` sem falso `200 ok`

