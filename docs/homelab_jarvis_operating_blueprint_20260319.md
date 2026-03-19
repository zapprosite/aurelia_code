---
title: Homelab Jarvis Operating Blueprint
status: active
created: 2026-03-19
owner: codex
scope: homelab-jarvis-aurelia-litellm-rag-maintenance
---

# Homelab Jarvis Operating Blueprint

## Objetivo

Fechar um desenho único para a Aurelia operar como:

- Jarvis local
- bot de manutenção do homelab
- agente de código
- bibliotecária de runbooks, docs e memória
- operador de browser, Telegram e terminal

Sem perder:

- estabilidade
- previsibilidade
- bancos organizados
- documentação viva

## Resposta curta sobre o LiteLLM

O LiteLLM ajuda, mas **não** deve virar o cérebro da Aurelia.

## Onde o LiteLLM ajuda de verdade

- dar um endpoint OpenAI-compatible único para clientes externos
- expor aliases estáveis para modelos locais e remotos
- concentrar auth, rate limit, budgets e logs para apps humanos
- servir OpenWebUI, IDEs, scripts e experimentos sem cada cliente conhecer Ollama/OpenRouter/etc.
- fazer retry/fallback e custo em uma borda de consumo

## Onde o LiteLLM não deve mandar

- wake word
- VAD
- spool de áudio
- roteamento principal da Aurelia
- health crítico do bot
- manutenção autônoma do homelab
- RAG core

Motivo:

- ele é um gateway útil de borda
- a Aurelia precisa de um plano de controle próprio, em Go, que conheça skills, runbooks, prioridades, maintenance mode e bancos locais

## Papel correto do LiteLLM no seu stack

Papel recomendado:

- **gateway de acesso para consumidores**

Consumidores:

- OpenWebUI
- IDEs
- Antigravity quando precisar de endpoint compatível
- testes manuais
- clientes externos internos ao homelab

Não é o plano de controle da Aurelia.

Fonte de verdade do roteamento de modelos:

- [model_routing_matrix_20260319.md](/home/will/aurelia/docs/model_routing_matrix_20260319.md)

## Arquitetura geral

```text
Human / Telegram / Antigravity / Browser
  -> Aurelia Control Plane (Go)
      -> Task router
      -> Runbook engine
      -> Maintenance engine
      -> Agent orchestration
      -> Memory coordinator
      -> Health / governor / watchdog

  -> Inference Plane
      -> Ollama local
      -> OpenRouter remoto
      -> Groq STT
      -> LiteLLM edge gateway

  -> Knowledge Plane
      -> SQLite core state
      -> Qdrant semantic index
      -> Supabase shared/app state

  -> Execution Plane
      -> CLI / tools nativas
      -> Browser automation
      -> Desktop fallback
      -> systemd / docker / zfs / network
```

## Planos corretos

## 1. Control Plane

Implementar em Go dentro da Aurelia.

Responsável por:

- classificar tarefas
- escolher rota de inferência
- aplicar runbooks
- tomar decisão de manutenção
- atualizar memória
- provar health real

Esse é o núcleo.

## 2. Inference Plane

Separar em 4 classes:

- `local-fast`
- `local-balanced`
- `remote-premium`
- `audio`

Sugestão:

- `qwen3.5:4b` = `local-fast`
- `qwen3.5:9b` = `local-balanced`
- `openrouter/minimax-m2.7` = `remote-premium`
- `groq whisper-large-v3-turbo` = `audio`

Regra operacional:

- apenas `1` modelo local residente por vez no caminho ativo do bot
- o residente padrao deve ser `qwen3.5:9b`
- `qwen3.5:4b` entra frio ou aquecido sob demanda, nao residente junto por padrao
- `gemma3:27b` pode continuar como modelo manual/offline de laboratório, não como runtime ativo do bot

Justificativa numerica do host:

- VRAM total medida: `24564 MiB`
- uso base do desktop/lab: ~`4810 MiB`
- VRAM livre em idle: ~`19238 MiB`
- `qwen3.5:9b` carregado no Ollama usa ~`9.2 GiB` de VRAM e deixa ~`10.5 GiB` livres
- `qwen3.5:9b + qwen3.5:4b` juntos deixam so ~`3.8 GiB` livres
- um `27B` de ~`17 GiB`, somado ao uso base do host, deixa folga perto de `2 GiB`

## 3. Voice Plane

O caminho correto continua sendo:

- wake word local
- VAD local
- buffer local
- Groq STT
- roteamento
- TTS separado

LiteLLM pode até proxyar `/audio`, mas não deve ser o coordenador dessa pipeline.

## Regra de VRAM

A conta correta da Aurelia/Jarvis e:

- `Groq` cuida do `STT`, entao o audio nao entra no budget de VRAM local
- a VRAM local deve ser guardada para:
  - um unico LLM residente
  - browser e desktop
  - embeddings assincronos
  - margem de estabilidade

Conclusao pratica:

- `Groq` e uma boa escolha porque tira o Whisper/STT da GPU local
- isso nao autoriza dois LLMs locais residentes ao mesmo tempo
- o desenho seguro continua sendo `1` modelo residente only

## 4. Knowledge Plane

Aqui está a parte mais importante para “bancos organizados”.

### SQLite

Fonte de verdade local do runtime.

Guardar:

- sessões
- jobs
- heartbeats
- incidentes
- facts
- notes
- cron
- estado do bot
- evidência resumida

SQLite é o banco do **controle operacional**.

### Qdrant

Apenas índice semântico derivado.

Guardar:

- chunks de docs
- runbooks
- conversas relevantes
- notas operacionais
- contexto de código

Qdrant é o banco da **recuperação semântica**.

Nunca deve ser a única fonte de verdade.

### Supabase

Use como plano compartilhado/app-level, não como coração do bot local.

Guardar:

- sessões compartilhadas
- auditoria mais rica
- jobs multi-interface
- relatórios
- espelhamento de eventos

Supabase é o banco de **integração e colaboração**.

## Regra de ouro dos bancos

- SQLite manda no runtime local
- Qdrant deriva de SQLite/docs/runbooks
- Supabase complementa, não substitui

## RAG e semântico

O desenho correto para RAG:

1. documento ou evento nasce em SQLite/docs
2. worker de indexação gera embedding
3. Qdrant recebe vetor e payload
4. recuperação semântica devolve candidatos
5. resposta final ainda cruza com facts/estado do SQLite

Assim você evita:

- alucinação semântica
- drift entre vetor e estado real
- “Qdrant manda mais que o runtime”

## Bibliotecária da Aurelia

Aurelia como bibliotecária deve ter 3 funções:

### 1. Catalogadora

- indexa docs
- indexa runbooks
- classifica incidentes
- mantém taxonomia por domínio

### 2. Curadora

- detecta duplicação
- marca docs obsoletos
- sobe prioridade de docs críticos
- reforça links entre runbook e evidência

### 3. Arquivista

- move histórico frio
- mantém retenção
- deixa bancos pequenos e legíveis

## Home lab autônomo

Para a Aurelia tocar a manutenção sozinha, faltam só 5 laços bem definidos:

### 1. Health loop

- varrer serviços
- validar health real
- distinguir warning de erro
- abrir incidente quando necessário

### 2. Repair loop

- escolher runbook mínimo
- aplicar correção pequena
- revalidar
- abortar se risco subir

### 3. Memory loop

- registrar causa raiz
- atualizar facts
- gerar note
- empurrar para Qdrant se merecer indexação

### 4. Hygiene loop

- revisar tamanho de logs
- revisar snapshots/backups
- revisar drift de firewall
- revisar drift de docker/systemd

### 5. Documentation loop

- atualizar `.context`
- manter índice de runbooks
- registrar mudanças arquiteturais

## Onde o LiteLLM entra nesse desenho

Entra aqui:

```text
Human-facing model access
  -> LiteLLM
      -> Ollama
      -> OpenRouter
      -> outros provedores
```

Não aqui:

```text
Aurelia maintenance brain
  -> runbook decision
  -> homelab self-healing
  -> core memory truth
```

## Recomendação final de stack

### Control Plane

- Go

### Gateway interno da Aurelia

- Go

### LiteLLM

- manter como gateway de borda para consumidores

### Modelos

- `qwen3.5:9b` principal local equilibrado e unico residente por padrao
- `qwen3.5:4b` roteador rápido, nao residente junto por padrao
- `openrouter/minimax-m2.7` premium remoto
- `groq whisper-large-v3-turbo` STT
- `bge-m3` embedding único

### Bancos

- SQLite = runtime truth
- Qdrant = semantic index
- Supabase = app/shared state

## Decisão sênior

Se a pergunta é “o LiteLLM ajuda em quê?”, a resposta é:

- ajuda muito como **porta única de modelos**
- ajuda pouco como **cérebro operacional do homelab**

Se a pergunta é “como manter foco no home lab estável, Jarvis local, RAG, agentes e bancos organizados?”, a resposta é:

- Aurelia em Go como plano de controle
- LiteLLM como borda de consumo
- SQLite como verdade local
- Qdrant como índice derivado
- Supabase como camada compartilhada
- Groq só no ouvido
- `qwen3.5:9b` como cerebro local residente
- runbooks + health + watchdog + `.context` como memória operacional viva
