---
title: Aurelia Master Blueprint
status: active
created: 2026-03-19
owner: codex
scope: master-plan-architecture-rollout-tests
---

# Aurelia Master Blueprint

## Objetivo

Juntar em um unico documento:

- arquitetura alvo
- decisoes fechadas
- o que falta
- ordem de execucao
- gates
- testes de tudo

Este documento vira o mapa mestre do restante.

## Resultado final esperado

A Aurelia deve operar como:

- Jarvis local
- bot autonomo de manutencao do homelab
- agente de browser, terminal e Antigravity
- bibliotecaria de docs, runbooks e memoria
- runtime local-first, barato e previsivel

## Estado ja fechado

- `qwen3.5:9b` como cerebro local residente
- `qwen3.5:4b` sob demanda
- `Groq` para STT
- `bge-m3` como embedding unico
- `SQLite` como verdade operacional
- `Qdrant` como indice semantico derivado
- `Supabase` como estado compartilhado
- `LiteLLM` como borda
- `Gemini web` apenas como pesquisa curada
- `internal/gateway/` com policy engine inicial
- `POST /v1/router/dry-run` implementado
- `GET /v1/router/status` e telemetria Prometheus do gateway implementados
- spool/processador de voz com fallback STT, heartbeat e mirrors opcionais implementados

## Arquitetura alvo

```text
Human / Telegram / Antigravity / Browser
  -> Aurelia Control Plane (Go)
      -> gateway
      -> runbook engine
      -> maintenance loop
      -> governor
      -> memory coordinator
      -> watchdog + health

  -> Voice Plane
      -> wake word
      -> VAD
      -> ring buffer
      -> Groq STT
      -> local fallback STT
      -> spool

  -> Inference Plane
      -> Ollama qwen3.5:9b
      -> Ollama qwen3.5:4b sob demanda
      -> OpenRouter DeepSeek / MiniMax / Qwen Flash
      -> LiteLLM edge

  -> Knowledge Plane
      -> SQLite
      -> Qdrant
      -> Supabase

  -> Execution Plane
      -> CLI tools
      -> agent-browser
      -> desktop fallback
      -> systemd / docker / zfs / network
```

## Regras de ouro

- `1` modelo residente only
- `Groq` so no lane de audio
- `OpenRouter` so por capacidade explicita
- `Gemini web` nunca entra direto no runtime
- `browser` antes de `desktop`
- `SQLite` manda no runtime
- `Qdrant` nunca vira fonte primaria
- nenhum `/health` pode mentir

## O que falta

### 1. Gateway rollout

- validar lane/modelo no runtime de deploy
- parar de decidir provedor de forma espalhada no live

### 2. Response guards

- controlar reasoning
- impedir `content=""` em JSON/curadoria

### 3. Budgets e breaker

- budget por lane
- breaker por `provider:model`
- fallback coerente

### 4. Voice plane real

- wake word
- VAD
- ring buffer
- captura continua de microfone
- servico dedicado opcional

### 5. Memory plane completo

- ligar `SQLite + Qdrant + Supabase` no fluxo real
- ingestao curada do Gemini web

### 6. Execution plane seguro

- click seguro
- digitacao segura
- kill-switch
- limite de passos

### 7. Rollout final

- portar tudo para `/home/will/aurelia-24x7`
- validar runtime live

## Ordem de execucao

1. voice capture live
2. memory plane real
3. gateway + voice rollout em deploy
4. orchestration/browser E2E
5. desktop fallback seguro
6. extensoes opcionais e rollback

## Gates

### Gate 1. Unitario

Tem que passar:

- `go test ./internal/gateway ./internal/health ./internal/telegram ./internal/skill ./pkg/llm ./pkg/stt -count=1`

### Gate 2. Suite geral

Tem que passar:

- `go test ./... -count=1`

### Gate 3. Runtime local

Tem que provar:

- `POST /v1/router/dry-run`
- `/health`
- browser smoke
- Ollama smoke
- STT smoke

### Gate 4. Deploy worktree

Tem que provar em `/home/will/aurelia-24x7`:

- build
- suite
- bot online
- health real
- endpoint de dry-run

## Matriz de testes

### A. Gateway

#### Unit

- route `maintenance` -> `local-balanced`
- route `routing` -> `remote-tool-long-output`
- route `audio` -> `groq`
- route `vision` -> lane remoto de visao
- `local_only=true` bloqueia remoto

#### Integration

- `POST /v1/router/dry-run` devolve `lane/provider/model/guards`
- request invalido retorna `400`
- rota custom aparece no server interno

#### Acceptance

- logs mostram lane escolhida
- resposta estruturada nao sai vazia

### B. LLM local/remoto

#### Unit

- guardas por lane:
  - `structured_json`
  - `curation`
  - `premium_text`
  - `maintenance_local`

#### Bakeoff

- SRE curto
- JSON de roteamento
- curadoria de RAG

#### Acceptance

- `minimax` fica lane premium
- `deepseek` fica lane curta estruturada
- `qwen3.5:9b` nao fica responsavel por JSON curto sem guardas

### C. Voice

#### Unit

- wake word gating
- VAD corta silencio
- spool aceita item valido
- clip acima do maximo e recusado ou truncado

#### Integration

- audio curto PT-BR -> Groq -> transcript
- 429 -> cooldown
- hard cap -> fallback local

#### Acceptance

- sem wake word nao chama STT
- sem falso positivo explosivo
- backlog nao cresce indefinidamente

### D. Memory

#### Unit

- facts vao para `SQLite`
- embeddings vao para `Qdrant`
- payload semantico tem campos minimos

#### Integration

- doc curado -> markdown -> embedding -> Qdrant
- evento operacional -> fact -> SQLite
- nota compartilhada -> Supabase

#### Acceptance

- facts operacionais existem sem depender do vetor
- busca semantica encontra material curado

### E. Execution

#### Browser

- abrir pagina
- capturar screenshot
- preencher campo
- click previsivel

#### Desktop fallback

- localizar janela
- foco seguro
- click seguro
- digitacao segura
- kill-switch

#### Acceptance

- browser resolve primeiro
- desktop nao roda cego

### F. Maintenance autonomy

#### Unit

- health loop classifica `ok/warning/error`
- runbook minimo e selecionado
- daily summary monta saida curta

#### Integration

- `nvidia-gpu` down -> diagnostico -> acao proposta
- service missing -> runbook correspondente
- drift de backup/firewall -> incidente operacional

#### Acceptance

- incidente recorrente vira runbook
- diarios e cron refletem o estado real

### G. Deploy

#### Pre-deploy

- `go test ./... -count=1`
- smoke local
- docs alinhados

#### Deploy worktree

- portar slice
- build
- suite
- reinicio do daemon
- `/health`
- bot online

#### Post-deploy

- dry-run responde
- logs sem erro de provider
- Telegram sem regressao
- governor sem falso degradado

### H. Rollback

Tem que existir teste de:

- voltar `app.json`
- voltar provider/modelo
- reiniciar daemon
- confirmar `/health` e Telegram

## Evidencias obrigatorias

Nao declarar sucesso sem:

- saida de teste
- log
- endpoint respondendo
- medida real de host quando a decisao depender de recurso

## Documentos ligados

- [plan.md](/home/will/aurelia/plan.md)
- [aurelia_general_blueprint_20260319.md](/home/will/aurelia/docs/aurelia_general_blueprint_20260319.md)
- [gateway_rollout_blueprint_20260319.md](/home/will/aurelia/docs/gateway_rollout_blueprint_20260319.md)
- [homelab_jarvis_operating_blueprint_20260319.md](/home/will/aurelia/docs/homelab_jarvis_operating_blueprint_20260319.md)
- [jarvis_local_voice_blueprint_20260319.md](/home/will/aurelia/docs/jarvis_local_voice_blueprint_20260319.md)
- [model_routing_matrix_20260319.md](/home/will/aurelia/docs/model_routing_matrix_20260319.md)
- [model_response_bakeoff_20260319.md](/home/will/aurelia/docs/model_response_bakeoff_20260319.md)

## Fonte de verdade

Se houver conflito entre blueprints, este documento vira a consolidacao operacional do restante.
