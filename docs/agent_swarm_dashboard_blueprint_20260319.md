---
title: Aurelia Hierarchical Agent Swarm Blueprint
status: proposed
created: 2026-03-19
owner: codex
---

# Objetivo

Desenhar um enxame hierárquico de agentes para a Aurélia, com dashboard operacional, fila de tarefas compartilhada, ajuda entre agentes quando ociosos e memória durável em `PostgreSQL + Qdrant`, sem quebrar o contrato atual do repositório.

# Decisão de alto nível

Não adotar um framework Python como runtime principal.

Em vez disso:

- `Aurélia` continua sendo a autoridade arquitetural e operacional
- o `control plane` do enxame nasce em `Go`
- `PostgreSQL` vira a fonte de verdade operacional do swarm
- `Qdrant` vira a memória semântica derivada
- `Grafana` continua para operação
- um dashboard leve em `Go` mostra canais, tarefas, agentes, leases e handoffs

Frameworks externos entram como referência de desenho, não como coração do runtime.

# O que copiar do estado da arte em 2026-03-19

## LangGraph

Vale copiar:

- o conceito de `assistants`, `threads` e `crons`
- a separação entre estado persistente e execução
- o modelo de autorização por recurso

Fonte:

- https://docs.langchain.com/langgraph-platform/auth

## AutoGen Studio

Vale copiar:

- a ideia de uma interface de times/agentes
- visualização de runs, mensagens e experimentos
- uso de estúdio/dash para observabilidade humana

Fonte:

- https://www.microsoft.com/en-us/research/wp-content/uploads/2024/08/AutoGen_Studio-12.pdf

## CrewAI

Vale copiar:

- observabilidade
- noção de crews/flows
- disciplina de execução por papel

Fonte:

- https://docs.crewai.com/

# Regra de arquitetura

O swarm segue esta autoridade:

1. humanos
2. `AGENTS.md`
3. `Aurélia`
4. agentes fixos
5. agentes efêmeros

Nenhum agente efêmero pode ultrapassar a Aurélia, abrir deploy sozinho ou reescrever política global.

# Como eu faria

## 1. Control plane em PostgreSQL

`PostgreSQL` deve guardar tudo que é estado operacional, transacional e auditável:

- agentes
- canais
- threads
- mensagens
- tarefas
- subtarefas
- leases
- heartbeats
- handoffs
- decisões
- incidentes
- resumos
- permissões

Padrão de execução:

- fila por `SELECT ... FOR UPDATE SKIP LOCKED`
- leases curtos por agente
- heartbeats renovando posse
- reaper para devolver tarefa órfã à fila

Isso é melhor que usar só fila externa porque:

- mantém consistência com o dashboard
- facilita auditoria
- simplifica retomada após queda
- reduz moving parts no homelab

## 2. Qdrant como memória semântica derivada

`Qdrant` não deve ser a verdade do sistema.

Ele deve indexar:

- resumos de threads
- ADRs
- runbooks
- postmortems
- observações operacionais
- contexto de repositório
- memória de conversa já consolidada

Contrato:

- embedding único: `bge-m3`
- distância: `cosine`
- payload sempre com:
  - `memory_type`
  - `source_table`
  - `source_id`
  - `agent_id`
  - `channel_id`
  - `thread_id`
  - `scope`
  - `created_at`
  - `expires_at`
  - `authority_level`

“Memória infinita” aqui não é literal. O desenho correto é:

- bruto em `PostgreSQL`
- resumo semântico em `Qdrant`
- arquivamento/compactação periódica

## 3. Agentes fixos

Eu começaria com seis agentes fixos:

- `aurelia-chief`
  - roteia, arbitra, decide
- `ops-sre`
  - saúde, incidentes, recuperação
- `librarian`
  - ADR, docs, curadoria, RAG
- `browser-operator`
  - Antigravity, browser, páginas
- `voice-operator`
  - STT, TTS, fila de áudio
- `local-executor`
  - CLI, scripts, automações seguras

Eles ficam sempre registrados, com capacidade, custo e limite de concorrência explícitos.

## 4. Agentes efêmeros

Agentes efêmeros nascem de uma tarefa ou thread e morrem por:

- TTL
- conclusão
- cancelamento
- inatividade

Exemplos:

- `incident-investigator-20260319-01`
- `rag-curator-20260319-01`
- `browser-session-20260319-01`
- `deploy-review-20260319-01`

Cada efêmero tem:

- `owner_agent_id`
- `scope`
- `allowed_tools`
- `write_scope`
- `ttl_seconds`
- `max_steps`
- `authority_ceiling`

## 5. Ajuda entre agentes ociosos

Esse é o ponto central.

Eu faria um `assistance queue` separado do `primary queue`.

Quando um agente fixa ou efêmero fica ocioso:

1. ele atualiza `presence=idle`
2. consulta tarefas elegíveis para ajuda
3. tenta claim com `SKIP LOCKED`
4. se ganhar o lease, executa como `supporting_agent`
5. devolve output estruturado para o `owner_agent`

Ocioso não significa livre para qualquer coisa.

O match precisa respeitar:

- hierarquia
- capacidade
- custo
- afinidade de domínio
- política de segurança
- dependências abertas

## 6. Dashboard

Eu faria dois dashboards:

### Dashboard operacional em Grafana

Para:

- fila
- p95 por lane
- falhas
- saturação de CPU/RAM/GPU
- backlog de voz
- breaker state
- agentes sem heartbeat

### Dashboard de controle em Go

Para:

- lista de agentes fixos e efêmeros
- canais
- threads
- tarefas em progresso
- tarefas bloqueadas
- leases
- handoffs
- decisões pendentes
- assistência entre agentes

Transporte:

- `LISTEN/NOTIFY` no Postgres
- ou SSE/WebSocket

## 7. Canais estilo Slack

Eu usaria canais explícitos:

- `ops-live`
- `incidents`
- `board-room`
- `browser`
- `voice`
- `rag-curation`
- `deploy-review`

Cada mensagem precisa ter:

- `channel_id`
- `thread_id`
- `sender_agent_id`
- `target_agent_id`
- `message_kind`
- `body`
- `refs`
- `priority`
- `visibility`

## 8. Brainstorm frequente sem caos

Eu não faria brainstorm aberto.

Usaria este ciclo:

1. agente líder propõe
2. `challenger` critica
3. `constraint guardian` valida custo, segurança e performance
4. `user advocate` valida clareza e usabilidade
5. Aurélia arbitra

Isso combina com a skill instalada:

- `brainstorming`
- `multi-agent-brainstorming`

## 9. Como guardar tudo

### PostgreSQL

Tabelas mínimas:

- `agents`
- `agent_capabilities`
- `agent_presence`
- `channels`
- `threads`
- `agent_messages`
- `tasks`
- `task_dependencies`
- `task_claims`
- `agent_handoffs`
- `decisions`
- `incidents`
- `memory_events`
- `thread_summaries`

### Qdrant

Coleções mínimas:

- `thread_memory`
- `decision_memory`
- `incident_memory`
- `runbook_memory`
- `repo_memory`

## 10. O que eu não faria

- não usaria `Qdrant` como fila
- não deixaria agente ocioso assumir tarefa fora de sua autoridade
- não faria memória “infinita” sem compactação
- não colocaria framework Python inteiro como core do runtime
- não misturaria dashboard operacional com observabilidade de logs

# Rollout sugerido

## Slice 11A

`agent_bus` e esquema SQL básico

## Slice 11B

dashboard leve em `Go`

## Slice 11C

assistência entre agentes ociosos

## Slice 11D

resumo de thread + indexação em `Qdrant`

## Slice 11E

board-room e brainstorm sequencial sob autoridade da Aurélia

# Veredito

Se fosse meu homelab, eu faria um swarm nativo da Aurélia em `Go`, com:

- `PostgreSQL` como bus operacional e verdade do sistema
- `Qdrant` como memória semântica derivada
- dashboard leve próprio
- `Grafana` para operação
- agentes fixos e efêmeros
- ajuda entre agentes ociosos via `assistance queue`
- brainstorm sequencial com papéis, não enxame caótico

Esse desenho aguenta crescer para múltiplas GPUs depois sem mudar o contrato central.
