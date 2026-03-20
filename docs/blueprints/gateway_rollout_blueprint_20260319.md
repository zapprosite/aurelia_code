---
title: Gateway Rollout Blueprint
status: active
created: 2026-03-19
owner: codex
scope: gateway-rollout-enforcement-runtime
---

# Gateway Rollout Blueprint

## Objetivo

Fechar o restante do gateway da Aurelia em cortes pequenos, testaveis e seguros, saindo de:

- policy engine e dry-run prontos

para:

- enforcement real no runtime
- guardas de reasoning e output
- budgets por lane
- circuit breaker
- telemetria
- rollout em deploy worktree

Sem quebrar:

- o bot Telegram
- o runtime live
- o health real
- o foco local-first do homelab

## Estado Atual

Ja existe:

- `internal/gateway/` com policy engine inicial
- `POST /v1/router/dry-run`
- matriz de roteamento versionada
- bakeoff real por qualidade/latencia

Ainda nao existe:

- provider selection dinamico no loop principal
- enforcement de lane por tarefa
- guardas aplicados na chamada real dos modelos
- budgets e counters por lane
- circuit breaker por `provider:model`
- metricas do gateway

## Principio de rollout

- mudar o minimo por slice
- provar cada slice com teste e evidencia
- nao virar o runtime live para comportamento novo sem passar na worktree de deploy

## Slice 10. Route Enforcement

**Objetivo:** fazer a Aurelia usar o planner de gateway de verdade.

Entradas:

- `internal/gateway/policy.go`
- `docs/model_routing_matrix_20260319.md`

Saidas:

- um seletor de provider/modelo por tarefa real
- integracao com `buildLLMProvider` ou camada equivalente de escolha

Arquivos alvo:

- `cmd/aurelia/app.go`
- `internal/agent/`
- `internal/telegram/`
- `internal/skill/`

Regra:

- `maintenance` continua local por default
- `routing/curation` podem escalar para `deepseek-v3.2`
- `workflow premium` pode escalar para `minimax-m2.7`
- `audio` nunca entra na lane LLM

Critério de aceite:

- tarefa simulada de `maintenance` usa lane local
- tarefa simulada de `routing` usa lane `deepseek`
- tarefa simulada de `browser_workflow` usa lane premium quando explicitamente permitido

## Slice 11. Response Guards

**Objetivo:** impedir `reasoning` desperdiçado e resposta vazia em lanes estruturadas.

Problema que resolve:

- `minimax-m2.7` e `qwen3.5:9b` consumindo budget em reasoning e entregando `content` vazio

Entregas:

- perfil por lane:
  - `structured_json`
  - `curation`
  - `premium_text`
  - `maintenance_local`
- guardas por request:
  - `max_output_tokens`
  - `soft_timeout_ms`
  - hint de reasoning minimo

Arquivos alvo:

- `pkg/llm/openai_compatible.go`
- `pkg/llm/ollama.go`
- `internal/gateway/`

Critério de aceite:

- prompts curtos estruturados nao retornam `content=""`
- testes unitarios cobrindo `structured_json` e `curation`

## Slice 12. Budgets e Circuit Breaker

**Objetivo:** colocar custo e falha no mesmo plano de controle.

Entregas:

- counters por lane:
  - `local`
  - `remote_structured`
  - `remote_premium`
  - `remote_vision`
  - `audio`
- budget soft/hard por lane
- circuit breaker por `provider:model`
- cooldown basico para 429/5xx/timeouts

Arquivos alvo:

- `internal/gateway/`
- `internal/health/`
- `cmd/aurelia/health_checks.go`

Critério de aceite:

- 429 abre breaker temporario
- erro repetido em rota remota cai para fallback
- health reflete breaker aberto como `warning` ou `error` coerente

## Slice 13. Telemetria e Dry-Run Rico

**Objetivo:** tornar a decisao do gateway observavel.

Entregas:

- metricas:
  - `gateway_route_selected_total`
  - `gateway_requests_total`
  - `gateway_errors_total`
  - `gateway_fallback_total`
  - `gateway_circuit_state`
- `dry-run` mais rico:
  - candidatos
  - score simplificado
  - motivo do veto

Arquivos alvo:

- `internal/gateway/`
- `internal/observability/`
- docs operacionais

Critério de aceite:

- `dry-run` explica por que uma rota foi escolhida
- rota vetada por `supports_tools=false` aparece como vetada

## Slice 14. Runtime Integration

**Objetivo:** ligar o gateway ao resto do sistema sem espalhar decisao de modelo.

Entregas:

- Telegram usa lane correta por classe de tarefa
- cron e maintenance loop usam local by default
- browser/Antigravity usam lane premium so quando o plano mandar
- curadoria de RAG usa lane `deepseek`

Arquivos alvo:

- `internal/telegram/`
- `internal/cron/`
- `internal/persona/`
- `internal/agent/`

Critério de aceite:

- logs mostram lane real usada
- tasks recorrentes nao escalam para remoto sem motivo

## Slice 15. Deploy e Rollout

**Objetivo:** levar o gateway para a worktree de deploy sem surpresa.

Passos:

1. portar para `/home/will/aurelia-24x7`
2. rodar `go test ./... -count=1`
3. validar `POST /v1/router/dry-run` na worktree de deploy
4. validar health
5. ativar enforcement primeiro em modo conservador
6. observar logs e fallback

Critério de aceite:

- bot continua online
- `/health` continua real
- sem regressao em Telegram
- sem escalonamento remoto indevido

## Ordem recomendada

1. Slice 10
2. Slice 11
3. Slice 12
4. Slice 13
5. Slice 14
6. Slice 15

## Regra de ouro

O gateway nao existe para “usar IA melhor”.

Ele existe para:

- reduzir custo
- manter qualidade por lane
- conter reasoning desperdicado
- proteger o homelab de surpresas
- manter a Aurelia explicavel
