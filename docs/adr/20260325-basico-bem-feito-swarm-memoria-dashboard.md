# ADR 20260325: Basico Bem Feito para Swarm, Memoria Infinita e Dashboard Operacional

**Status:** Substituída por `20260325-basico-bem-feito-v2-implementation.md`
**Data:** 2026-03-25
**Autoridade:** Will (Principal Engineer) + Aurelia (Arquiteta Principal)
**Escopo:** runtime multi-bot, swarm/team, memoria, dashboard e contratos de verdade

---

## 1. Contexto e problema real

Este repositorio ja possui varias pecas uteis, mas ainda nao opera como um "basico bem feito" estavel.
Hoje existe capacidade real de:

- multi-bot no pool Telegram
- orquestracao `master -> task -> worker -> mailbox`
- cron jobs persistidos em SQLite
- dashboard com SSE e APIs locais
- memoria local em SQLite
- indexacao/consulta parcial em Qdrant

Mas o estado atual ainda tem divergencias importantes entre discurso e implementacao:

1. **O swarm nao e um swarm cooperativo pleno.**
   O que existe hoje e um runtime de equipe baseado em `TeamManager`, `MasterTeamService`, tasks, dependencies e mailbox, principalmente em:
   - `internal/agent/master_team_service.go`
   - `internal/agent/worker.go`
   - `internal/agent/task_store_schema.go`

   Isso e melhor que um prototipo, mas ainda nao e "todos se ajudam como equipe" de forma nativa e auditavel. O nome "swarm" esta a frente da realidade.

2. **O dashboard ainda e mais observacional do que operacional.**
   Hoje ele publica eventos em memoria via SSE e registra rotas dinamicas em:
   - `internal/dashboard/dashboard.go`

   Isso atende a UX de monitoramento, mas ainda nao atende o criterio "realmente funciona" como superficie oficial de operacao. Falta snapshot consolidado, status unificado e degradacao explicita.

3. **A memoria ainda nao tem contrato canonico executado de ponta a ponta.**
   O repositorio diz que SQLite, Qdrant e Supabase possuem papeis diferentes, mas a implementacao ainda esta parcial:
   - `internal/memory/context_assembler.go` ainda nao filtra memoria por `bot_id`
   - `cmd/aurelia/brain_handlers.go` ainda usa `scroll` + filtro textual no payload
   - `docs/governance/DATA_GOVERNANCE.md` define Supabase local como camada de negocio, mas isso ainda nao esta integrado no runtime Go

4. **Existe risco de drift entre identidade dos bots, memoria e observabilidade.**
   O print mostra um time real de trabalho:
   - `Aurelia_Code`
   - `Caixa CPF e CNPJ`
   - `AC VENDAS`
   - `ORGANIZACAO DE OBRAS`
   - `AGENDA CPF`

   Mas o backend ainda nao trata `aurelia_code` como identidade canonica compatibilizada com o legado `aurelia`.

5. **O cron atual e simples demais para virar espinha dorsal de operacao.**
   O scheduler em `internal/cron/scheduler.go` faz polling e executa jobs em serie. Isso e suficiente para rotinas administrativas leves, mas nao para um sistema que pretende coordenar equipe, memoria, watchdogs e dashboard sem bloqueio em cascata.

Esta ADR existe para parar a inflacao de nome e consolidar um padrao simples, coerente e operacionalmente defensavel.

---

## 2. Decisao

### 2.1. Fonte de verdade e papeis por camada

O ecossistema da Aurelia passa a operar com quatro camadas bem definidas:

1. **Supabase local = sistema de registro canonico**
   - fonte principal dos dados estruturados de negocio
   - fonte principal da memoria canonica curada
   - fonte principal do cadastro de bots, identidade, papeis, permissoes e snapshots operacionais
   - toda entidade que precise sobreviver a reprocessamento, sincronizacao ou auditoria deve existir aqui

2. **SQLite local = runtime operacional e estado curto**
   - queue de cron
   - leases de worker
   - mailbox de equipe
   - cache local
   - working memory
   - replay apos restart

   SQLite deixa de ser tratado como deposito de longo prazo ou pseudo-sistema de registro de tudo.

3. **Qdrant = indice semantico derivado**
   - indexa somente dados ja canonizados
   - nunca e fonte primaria
   - qualquer ponto precisa referenciar um `source_id` canonico
   - toda leitura precisa respeitar namespace minimo por bot e dominio

4. **Obsidian CLI = superficie humana bidirecional controlada**
   - recebe ADRs, notas, memoria curada e conhecimento exportavel
   - pode originar conhecimento humano, mas somente via pipeline sincronizado e auditavel
   - nao pode virar banco primario implicito

### 2.2. Lideranca e time oficial

O time oficial do print passa a ser reconhecido como estrutura-alvo do runtime:

- `aurelia_code` = lider e orquestradora principal
- `caixa` = dominio financeiro/caixa
- `ac_vendas` = dominio comercial
- `organizacao_de_obras` = dominio operacional de obras
- `agenda_cpf` = dominio agenda/pessoal

Regras:

1. `aurelia_code` vira o nome principal do lider.
2. O alias legado `aurelia` continua aceito durante transicao por compatibilidade.
3. O lider concentra:
   - despacho
   - consolidacao
   - arbitragem de conflito
   - visao global do dashboard
4. Especialistas podem pedir ajuda lateral, mas o ownership de execucao nao muda por acidente.

### 2.3. Modelo oficial de colaboracao entre bots

O runtime oficial deixa de vender "swarm pleno" sem contrato. Passa a haver tres modos formais:

1. **Delegacao**
   - `aurelia_code` cria task para especialista
   - ownership inicial vai para o especialista

2. **Handoff**
   - ownership muda explicitamente de um agente para outro
   - requer trilha de motivo e resolucao

3. **Assistencia lateral**
   - o agente dono continua owner
   - outro agente ajuda via thread, mailbox ou task auxiliar
   - ajuda lateral nao encerra nem toma posse automaticamente da tarefa principal

Nao e permitido chamar isso de swarm cooperativo completo enquanto o estado real ainda for apenas task queue central com mailbox.

### 2.4. Dashboard oficial

`https://aurelia.zappro.site/` passa a ser definido como **superficie operacional oficial** do sistema, e nao apenas tela de eventos.

Para isso, o dashboard deve refletir:

- identidade canonica dos bots
- status do time
- saude de Supabase, SQLite, Qdrant, cron e runtime
- backlog operacional
- execucoes recentes
- degradacao explicita quando um subsistema falhar

Regra obrigatoria:

- SSE deixa de ser tratado como verdade do sistema
- SSE passa a ser somente mecanismo de entrega em tempo real
- a verdade operacional vem de snapshots consolidados e auditaveis

### 2.5. Memoria infinita

"Memoria infinita" neste repositorio nao significa contexto magico ilimitado em prompt. Significa:

- armazenamento continuo com origem rastreavel
- projecao semantica reutilizavel
- sincronizacao com camada humana
- recuperacao por dominio, bot, tempo e relevancia

Contrato minimo da memoria canonica:

- `canonical_bot_id`
- `domain`
- `kind`
- `source_system`
- `source_id`
- `content`
- `metadata`
- `created_at`
- `updated_at`
- `version`

Contrato minimo da projecao vetorial no Qdrant:

- `canonical_bot_id`
- `persona_id`
- `chat_id`
- `domain`
- `source_system`
- `source_id`
- `text`
- `ts`
- `version`

Qualquer escrita vetorial sem namespace ou sem referencia ao registro canonico passa a ser considerada invalida.

---

## 3. Consequencias praticas

### Fora de escopo desta ADR

Esta ADR **nao** autoriza, por si so:

- trocar o model stack canonico definido em `docs/governance/MODEL-STACK-POLICY.md`
- mudar portas, subdominios ou Cloudflare Access fora da governanca existente
- rebaixar Supabase, Qdrant ou Obsidian a atalhos ad hoc sem contrato
- mascarar lacunas de implementacao com rename cosmetico de feature

Ela define contrato arquitetural e direcao operacional. Alteracoes de infraestrutura ou de stack continuam sujeitas as governancas proprias.

### Positivas

- o repositorio ganha um mapa de autoridade claro
- o dashboard pode deixar de mentir por omissao
- a equipe de bots ganha papel real, e nao apenas nomes e cards
- a memoria passa a ter rastreabilidade entre runtime, indice semantico e superficie humana
- `aurelia_code` vira identidade lider sem quebrar o legado de uma vez

### Negativas

- esta ADR reduz a liberdade de "ligar tudo em tudo"
- Supabase local vira dependencia arquitetural central
- o trabalho de integracao aumenta no curto prazo
- varios atalhos existentes passam a ser tecnicamente proibidos

### Trade-off aceito

Aceita-se mais estrutura e menos improviso para ganhar estabilidade operacional e evitar um sistema "bonito no nome e fraco na realidade".

---

## 4. Passos obrigatorios para padrao estavel

### Fase 1 - Honestidade arquitetural

1. Parar de chamar o estado atual de swarm cooperativo pleno.
2. Declarar em docs e codigo a diferenca entre:
   - runtime operacional
   - observabilidade
   - memoria
   - indices derivados
3. Introduzir `canonical_bot_id` como chave obrigatoria nas camadas que hoje dependem de nome solto.
4. Tratar `aurelia_code` como nome principal com compatibilidade temporaria para `aurelia`.

### Fase 2 - Fundacao de dados

1. Criar no Supabase local as entidades canonicas de:
   - bots
   - team runs
   - tasks consolidadas
   - memory items
   - knowledge items
   - dashboard snapshots
   - audit events
2. Rebaixar SQLite ao papel de runtime local e estado curto.
3. Ajustar Qdrant para indexar somente dados canonicos derivados.
4. Definir pipeline formal `Obsidian CLI <-> sync controlado <-> Supabase`.

### Fase 3 - Equipe que realmente se ajuda

1. Formalizar ownership de task.
2. Formalizar ajuda lateral sem roubo silencioso de ownership.
3. Formalizar handoff com motivo, timestamp e consolidacao pelo lider.
4. Tornar o dashboard capaz de mostrar:
   - quem pediu ajuda
   - quem ajudou
   - quem continua dono
   - qual tarefa travou

### Fase 4 - Dashboard operacional

1. Criar snapshot consolidado de status como fonte unica para:
   - dashboard web
   - `/status`
   - watchdogs
   - health checks
2. Tornar falha de Qdrant, Supabase ou SQLite visivel como `degraded` ou `offline`, e nao como lista vazia silenciosa.
3. Separar fluxo de eventos de fluxo de estado.

### Fase 5 - Cron util para producao

1. Tirar o cron do papel de executor ingenuo em serie para tudo.
2. Classificar jobs por tipo:
   - sistema
   - memoria
   - negocio
   - notificacao
3. Adotar timeout, retry e visibilidade operacional minima por job.
4. Nao permitir que um job longo derrube watchdog critico por bloqueio em cascata.

---

## 5. O que fica proibido por esta ADR

Ficam proibidos os seguintes atalhos:

- tratar Qdrant como banco primario
- usar Obsidian como fonte implicita sem auditoria
- escrever memoria sem `canonical_bot_id`
- chamar dashboard de "operacional" sem snapshot consolidado
- apresentar lista vazia como se fosse sucesso quando o backend falhou
- adicionar novos bots sem papel, dominio e contrato de colaboracao
- continuar expandindo schema de swarm sem ligar isso a comportamento real

---

## 6. Criterios de aceite

Esta ADR so sera considerada cumprida quando os seguintes resultados forem verdadeiros:

1. `aurelia_code` for reconhecida como lider no runtime e no dashboard, com compatibilidade controlada ao legado.
2. Os bots especialistas puderem colaborar de forma auditavel.
3. O dashboard oficial mostrar estado verdadeiro, inclusive degradacao.
4. Supabase, SQLite, Qdrant e Obsidian tiverem papeis claros e nao sobrepostos.
5. A memoria puder ser recuperada por bot, dominio e origem sem vazamento lateral.
6. O sistema ficar mais simples de explicar do que hoje.

---

## 7. Artefatos e areas impactadas

Os trabalhos decorrentes desta ADR devem partir principalmente destes pontos do repositorio:

- `cmd/aurelia/app.go`
- `internal/agent/master_team_service.go`
- `internal/agent/task_store_schema.go`
- `internal/agent/squad.go`
- `internal/cron/scheduler.go`
- `internal/dashboard/dashboard.go`
- `internal/memory/context_assembler.go`
- `cmd/aurelia/brain_handlers.go`
- `docs/governance/DATA_GOVERNANCE.md`

---

## 8. Nota critica final

O erro mais provavel no estado atual nao e falta de tecnologia. E **misturar papeis**, **superestimar maturidade** e **aceitar interfaces que parecem prontas sem serem autoridade real**.

Esta ADR escolhe um caminho menos glamouroso e mais profissional:

- menos fantasia de swarm
- menos memoria "magica"
- menos dashboard ornamental
- mais contrato
- mais ownership
- mais rastreabilidade
- mais verdade operacional

Esse e o padrao minimo para a Aurelia funcionar como sistema serio, e nao apenas como uma colecao de componentes interessantes.
