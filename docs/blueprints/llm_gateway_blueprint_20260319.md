---
title: LLM Gateway Blueprint
status: active
created: 2026-03-19
owner: codex
scope: llm-gateway-proxy-go-rust
---

# LLM Gateway Blueprint

## Objetivo

Extrair do repositório atual a ideia de um gateway/proxy inteligente para vários modelos e redesenhar isso como um serviço de primeira classe, com foco em:

- roteamento entre vários provedores e modelos
- tool-calling consistente
- fallback explícito
- latência previsível
- observabilidade real
- contrato único para o resto da Aurelia

## O que já existe hoje no repositório

O projeto já tem partes de gateway, mas espalhadas.

### 1. Camada comum de transporte

Arquivo:

- [openai_compatible.go](/home/will/aurelia/pkg/llm/openai_compatible.go)

O que já faz bem:

- normaliza request para `chat/completions`
- injeta headers e auth por provedor
- reaproveita o mesmo transport para OpenRouter, Kilo, Z.ai, Alibaba e Ollama

Leitura:

- isso já é um mini data-plane
- falta política, score, health, budget e decisão de rota

### 2. Adapters por provedor

Arquivos:

- [openrouter.go](/home/will/aurelia/pkg/llm/openrouter.go)
- [kilo.go](/home/will/aurelia/pkg/llm/kilo.go)
- [kimi.go](/home/will/aurelia/pkg/llm/kimi.go)
- [ollama.go](/home/will/aurelia/pkg/llm/ollama.go)

O que já faz bem:

- cada provedor configura endpoint, auth e descoberta
- Kimi já trata um caso real de incompatibilidade de resposta e corrige `tool_calls`

Leitura:

- isso é a parte mais inteligente do desenho atual
- o projeto já sabe que compatibilidade “OpenAI-like” não basta

### 3. Catálogo de modelos

Arquivo:

- [catalog.go](/home/will/aurelia/pkg/llm/catalog.go)

O que já faz bem:

- mistura descoberta remota e fallback curado
- ordena modelos locais/remotos com heurística simples

Leitura:

- isso já é um embrião de registry de capacidades
- hoje ainda falta registrar latência, custo, ferramentas, contexto, visão e confiabilidade

### 4. Router local de skill

Arquivo:

- [router.go](/home/will/aurelia/internal/skill/router.go)

O que já faz bem:

- separa classificação de execução
- já trabalha com degradação graciosa

Leitura:

- a ideia do roteamento já existe
- ainda está aplicada a skill, não a provedores/modelos

## O que falta para virar gateway de verdade

Hoje faltam 8 peças centrais:

1. Registry de capacidades por modelo
2. Policy engine de roteamento
3. Health scoring por provedor/modelo
4. Circuit breaker por rota
5. Budget/custo/latência por classe de tarefa
6. Normalização séria de resposta e tool-calling
7. Session affinity e sticky routing
8. Métricas e traces por decisão

Sem isso, o código atual é uma coleção boa de adapters, mas não um gateway inteligente.

## Blueprint alvo

```text
Client / Aurelia
  -> Gateway API
      -> Request normalizer
      -> Model registry
      -> Policy engine
      -> Route scorer
      -> Circuit breaker
      -> Provider adapter
      -> Response normalizer
      -> Telemetry
```

## Contrato do gateway

O gateway deve expor um contrato único para o resto da Aurelia:

- `POST /v1/chat/completions`
- `GET /v1/models`
- `GET /health`
- `GET /metrics`
- `POST /v1/router/dry-run`

Estado atual:

- `POST /v1/router/dry-run` ja tem um primeiro corte implementado no server interno da Aurelia
- a politica inicial mora em `internal/gateway/`
- `GET /v1/router/status` ja expõe o snapshot interno do gateway
- o runtime principal ja consegue selecionar lane/modelo via `gateway.Provider`
- guardas de resposta, budgets por lane e circuit breaker em memoria ja estao ativos
- telemetria Prometheus ja foi ligada no server interno em `GET /metrics`
- rollout na worktree de deploy continua faltando

O `dry-run` é importante:

- recebe prompt, tools e hints
- não executa chamada real
- só explica qual rota seria escolhida e por quê

Isso evita “mágica invisível”.

## Slice implementado agora

O gateway saiu do modo apenas documental e entrou no primeiro enforcement real.

Arquivos centrais:

- [provider.go](/home/will/aurelia/internal/gateway/provider.go)
- [app.go](/home/will/aurelia/cmd/aurelia/app.go)
- [health_checks.go](/home/will/aurelia/cmd/aurelia/health_checks.go)
- [server.go](/home/will/aurelia/internal/health/server.go)

O que ja esta entregue:

- escolha de lane/modelo no runtime principal
- guardas de `reasoning/output` por lane
- budgets por lane em memoria
- circuit breaker por `provider:model`
- rota de status interno para inspeção rapida

O que ja foi fechado depois desse corte:

- telemetria exportada para Prometheus
- `GET /v1/router/status`
- suite verde com gateway enforcement ativo

O que continua faltando:

- persistencia de budgets/breakers
- rollout validado na worktree de deploy

## Registry de capacidades

Cada modelo precisa de um registro local parecido com isto:

```json
{
  "id": "qwen3.5:9b",
  "provider": "ollama",
  "class": "local-balanced",
  "supports_tools": true,
  "supports_vision": true,
  "supports_stream": true,
  "supports_json_mode": false,
  "context_window": 262144,
  "cost_class": "local",
  "latency_class": "low",
  "quality_class": "medium-high",
  "stability_class": "high",
  "warmable": true
}
```

Isso é o coração do gateway.

## Policy engine

O roteamento não deve ser por `if provider == x`.

Deve ser por política.

Fonte de verdade da politica inicial:

- [model_routing_matrix_20260319.md](/home/will/aurelia/docs/model_routing_matrix_20260319.md)

Exemplo:

- `task=router` -> `qwen3.5:4b`
- `task=tool_orchestrator` -> `qwen3.5:9b`
- `task=deep_reasoning` -> `openrouter/minimax-m2.7`
- `task=manual_offline_chat` -> `gemma3:27b-it-q4_K_M`
- `task=stt` -> `groq`

Hints úteis:

- `requires_tools=true`
- `latency_budget_ms=2000`
- `local_only=true`
- `cost_sensitive=true`
- `vision_required=true`
- `json_strict=true`

## Score de rota

Cada rota recebe score por:

- capacidade
- saúde
- latência recente
- custo
- erro recente
- aquecimento
- afinidade de sessão

Exemplo:

```text
final_score =
  capability_score *
  health_score *
  latency_score *
  stability_score *
  session_affinity_score
```

Se `supports_tools=false`, a rota nem entra na disputa para requests com tools.

## Circuit breaker

Cada `provider:model` precisa de estado próprio:

- `closed`
- `open`
- `half-open`

Abrir circuito quando:

- `5xx` acima do threshold
- timeout acima do threshold
- tool-call inválido recorrente
- resposta não parseável recorrente

## Normalização

O gateway precisa devolver um contrato uniforme mesmo quando o provedor for inconsistente.

Exemplos reais do repositório:

- Kimi exige correção de `tool_calls` em [kimi.go](/home/will/aurelia/pkg/llm/kimi.go)
- Ollama pode ter modelo instalado sem suporte a tools

Então o normalizer deve:

- validar `tool_calls`
- corrigir formatos conhecidos
- marcar `degraded=true` quando precisou reparar resposta
- registrar `provider_behavior` para observabilidade

## Observabilidade mínima

Métricas por rota:

- `gateway_requests_total`
- `gateway_errors_total`
- `gateway_route_latency_seconds`
- `gateway_route_selected_total`
- `gateway_circuit_state`
- `gateway_tool_parse_failures_total`
- `gateway_fallback_total`
- `gateway_degraded_responses_total`

Logs por decisão:

- task class
- route candidates
- route escolhida
- motivo da escolha
- fallback ocorrido ou não

## Go ou Rust?

## Minha recomendação

Para este projeto: **Go primeiro**.

Motivo:

- o repositório já é Go
- adapters e contratos já existem em Go
- observabilidade, HTTP, workers e integração com a Aurelia já estão no mesmo ecossistema
- o ganho de tempo e coesão vence a tentação de reescrever

## Quando Rust faria sentido

Rust só vira melhor se o gateway for separado como produto próprio com foco em:

- throughput muito alto
- multi-tenant real
- isolamento mais rígido
- streaming pesado
- plugin ABI/control plane mais sério

Se o objetivo é “deixar a Aurelia inteligente agora”, Rust atrasa.

## Arquitetura recomendada em Go

Pacotes:

```text
internal/gateway/
  api/
  registry/
  policy/
  scoring/
  health/
  breaker/
  adapters/
  normalize/
  telemetry/
```

Interfaces:

```go
type Adapter interface {
    Name() string
    Capabilities(ctx context.Context, model string) (ModelCaps, error)
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
    Health(ctx context.Context, model string) HealthStatus
}

type Policy interface {
    Select(ctx context.Context, req RouteRequest, candidates []RouteCandidate) (RouteDecision, error)
}
```

## Arquitetura recomendada em Rust

Se insistir em Rust, eu faria assim:

- `axum` para API
- `reqwest` para adapters HTTP
- `tower` para middleware, retry e circuit breaker
- `serde` para normalização
- `tokio` para execução assíncrona

Mas eu não começaria por aqui neste repo.

## Slices de implementação em Go

Blueprint de rollout restante:

- [gateway_rollout_blueprint_20260319.md](/home/will/aurelia/docs/gateway_rollout_blueprint_20260319.md)

### Slice 1

- extrair `OpenAICompatibleProvider` para `internal/gateway/adapters/openai_compatible`
- criar `ModelCaps`
- criar `GET /health` por rota

### Slice 2

- criar registry local
- registrar `supports_tools`, `supports_vision`, `latency_class`
- ligar Ollama via `/api/show`

### Slice 3

- criar policy engine
- classificar requests por `task class`
- introduzir `dry-run router`

### Slice 4

- circuit breaker
- fallback explícito
- telemetry no Prometheus

### Slice 5

- transformar Aurelia em cliente do gateway
- remover decisão de provedor espalhada do app principal

## Decisão sênior

Não copiar o gateway atual “como está”.

O caminho certo é:

- reaproveitar adapters e contratos já bons
- centralizar decisão
- formalizar capacidades
- tratar ferramentas como requisito de primeira classe
- implementar em Go primeiro

Resumo:

- o repositório já tem bons pedaços de gateway
- o núcleo inteligente ainda não existe
- a melhor ideia para agora é um gateway em Go, com adapters reaproveitados, registry de capacidades e policy engine explícito
- Rust só entra se isso virar um produto separado de infra
