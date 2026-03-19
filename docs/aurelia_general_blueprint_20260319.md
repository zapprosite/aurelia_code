---
title: Aurelia General Blueprint
status: active
created: 2026-03-19
owner: codex
scope: jarvis-gateway-homelab-rag-rollout
---

# Aurelia General Blueprint

## Objetivo

Consolidar em um unico documento tudo que falta para a Aurelia operar como:

- Jarvis local
- bot autonomo de manutencao do homelab
- agente de browser, terminal e Antigravity
- bibliotecaria de docs, runbooks e memoria
- runtime estavel, explicavel e barato

## Estado fechado

Ja existe:

- stack local decidida:
  - `qwen3.5:9b` residente
  - `qwen3.5:4b` sob demanda
  - `Groq` para STT
  - `bge-m3` para embeddings
- browser baseline validado
- skill do Antigravity pronta
- prompt `light` no Telegram
- gateway dry-run em `Go`
- gateway enforcement real com budgets, breaker e telemetria
- matriz de roteamento registrada
- bakeoff real entre modelos
- health real sem falso `200 ok` no deploy slice sem Gemini
- spool/processador de voz com fallback STT e status route

## O que falta

Falta fechar 6 frentes:

1. rollout do gateway em deploy
2. voz em background live
3. memoria operacional completa
4. governor e budgets
5. desktop fallback seguro
6. rollout final na worktree de deploy

## Arquitetura final

```text
Human / Telegram / Antigravity / Browser
  -> Aurelia Control Plane (Go)
      -> gateway policy + route enforcement
      -> runbook engine
      -> maintenance loop
      -> governor
      -> memory coordinator
      -> health + watchdog

  -> Voice Plane
      -> wake word
      -> VAD
      -> ring buffer
      -> Groq STT
      -> spool

  -> Inference Plane
      -> Ollama qwen3.5:9b
      -> Ollama qwen3.5:4b sob demanda
      -> OpenRouter DeepSeek / MiniMax / Qwen Flash
      -> LiteLLM edge gateway

  -> Knowledge Plane
      -> SQLite truth
      -> Qdrant semantic index
      -> Supabase shared state

  -> Execution Plane
      -> CLI tools
      -> agent-browser
      -> desktop fallback
      -> systemd / docker / zfs / network
```

## Sequencia obrigatoria

### Fase 1. Gateway real

**Objetivo:** sair de policy documental para policy aplicada.

Entregas:

- route enforcement no runtime
- guardas reais de reasoning/output
- budgets por lane
- circuit breaker
- telemetria

Critério de aceite:

- `maintenance` usa local por default
- `routing/curation` pode escalar para `deepseek-v3.2`
- `workflow premium` so escala para `minimax` quando apropriado
- respostas estruturadas nao saem vazias por reasoning descontrolado

Status atual:

- concluida no repositório principal
- pendente apenas rollout/validacao na worktree de deploy

### Fase 2. Voice plane

**Objetivo:** tornar a Aurelia sempre ouvindo, sem desperdiçar custo.

Entregas:

- spool local
- processador de fila com heartbeat
- fallback STT local
- `openWakeWord`
- `Silero VAD`
- ring buffer
- servico dedicado de captura

Critério de aceite:

- sem wake word nao chama STT
- fala real gera spool e transcript
- ruído e silencio nao disparam Groq
- 429 cai para fallback/cooledown

Status atual:

- spool/processador/fallback ja entregues
- captura live de microfone ainda pendente

### Fase 3. Memory plane

**Objetivo:** deixar bancos organizados e coerentes.

Entregas:

- `SQLite` como verdade do runtime
- `Qdrant` como derivado semantico
- `Supabase` como estado compartilhado
- ingestao curada do Gemini web para RAG

Critério de aceite:

- facts operacionais nao dependem do Qdrant
- docs e notas so entram no Qdrant com curadoria
- memoria semantica e runtime nao entram em drift

### Fase 4. Execution plane

**Objetivo:** consolidar a capacidade de agir com seguranca.

Entregas:

- browser-use seguro
- desktop click seguro
- digitacao segura
- limites de passo
- kill-switch

Critério de aceite:

- browser continua sendo primeira opcao
- desktop entra so em fallback
- acoes de alto risco nao executam cegamente

### Fase 5. Maintenance autonomy

**Objetivo:** bot local tocar a manutencao do homelab sozinho.

Entregas:

- health loop
- repair loop
- hygiene loop
- documentation loop
- cron diarios e horarios

Critério de aceite:

- incidente comum vira runbook
- drift operacional aparece nos diarios
- o bot detecta e corrige degradacoes simples

### Fase 6. Deploy rollout

**Objetivo:** levar o desenho final para a worktree de deploy sem surpresa.

Entregas:

- portar slices para `/home/will/aurelia-24x7`
- validar testes
- validar `health`
- validar Telegram
- validar gateway dry-run
- validar voice path

Critério de aceite:

- bot segue online
- `/health` continua real
- sem regressao no deploy
- sem custo remoto indevido

## Regras de ouro

- `qwen3.5:9b` e o unico modelo residente do caminho ativo
- `Groq` so no ouvido
- `Gemini web` so como pesquisa curada
- `LiteLLM` so como borda
- `SQLite` manda no runtime
- `Qdrant` indexa material aprovado
- `Supabase` complementa, nao substitui
- nenhum health pode mentir

## Gates

### Gate A. Codigo

- `go test ./... -count=1`

### Gate B. Runtime local

- `POST /v1/router/dry-run`
- `/health`
- browser smoke
- STT smoke

### Gate C. Deploy worktree

- build
- suite
- Telegram vivo
- health real

## Ordem pratica

1. fechar voice plane real
2. fechar estado/memoria reais
3. levar gateway + voice para deploy
4. fechar browser/Antigravity E2E
5. fechar desktop fallback seguro
6. avaliar extensoes opcionais

## Fonte de verdade

Este documento consolida o restante.

Documentos de apoio:

- [plan.md](/home/will/aurelia/plan.md)
- [gateway_rollout_blueprint_20260319.md](/home/will/aurelia/docs/gateway_rollout_blueprint_20260319.md)
- [homelab_jarvis_operating_blueprint_20260319.md](/home/will/aurelia/docs/homelab_jarvis_operating_blueprint_20260319.md)
- [jarvis_local_voice_blueprint_20260319.md](/home/will/aurelia/docs/jarvis_local_voice_blueprint_20260319.md)
- [model_routing_matrix_20260319.md](/home/will/aurelia/docs/model_routing_matrix_20260319.md)
