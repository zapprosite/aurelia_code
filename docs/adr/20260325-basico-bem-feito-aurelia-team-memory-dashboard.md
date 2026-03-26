# ADR 20260325: Basico Bem Feito para Aurélia, Team Bots, Memória Infinita e Dashboard Operacional

**Status:** Proposto
**Data:** 25 de Março de 2026
**Autoridade:** Will (Principal Engineer) + Aurélia (Arquiteta Principal)
**Escopo:** Runtime Go, team bots Telegram, memória infinita, dashboard, governança de dados

---

## 1. Contexto

Em 25/03/2026, o repositório já possui peças relevantes para uma operação séria:

- runtime Go com `team -> task -> worker -> mailbox`
- cron jobs persistidos em SQLite
- dashboard web em `https://aurelia.zappro.site/`
- memória local em SQLite
- indexação e busca parcial em Qdrant
- Supabase local disponível no homelab
- uso operacional de Obsidian/Obsidian CLI no ecossistema
- múltiplos bots Telegram especializados, conforme a operação real mostrada no print:
  - `Aurelia_Code`
  - `Caixa CPF e CNPJ`
  - `AC VENDAS`
  - `ORGANIZACAO DE OBRAS`
  - `AGENDA CPF`

O problema não é falta de tecnologia. O problema é falta de um contrato simples, honesto e estável entre essas peças.

Hoje o discurso do sistema está adiantado em relação ao que o runtime realmente entrega. O repositório fala em swarm, brain semântico, dashboard operacional e memória unificada, mas a implementação atual ainda está mais perto de:

- orquestração master-worker com mailbox, não swarm cooperativo pleno
- cron síncrono e simples, não orquestração resiliente por classes de jobs
- dashboard observacional leve, não painel operacional com fonte de verdade
- Qdrant parcialmente integrado, com contratos de payload divergentes
- Supabase documentado como centro de negócio, mas ainda não integrado ao runtime principal
- Obsidian presente no ecossistema, mas sem governança explícita como parte da arquitetura oficial

Esta ADR existe para forçar o repositório a voltar ao "básico bem feito": poucas camadas, responsabilidades claras, liderança explícita, memória infinita auditável e dashboard que reflita o estado real.

---

## 2. Diagnóstico Crítico do Estado Atual

### 2.1 Swarm / Team

O núcleo atual de equipe é bom o suficiente para ser base, mas não deve ser vendido como algo que ainda não é.

- o estado persistido real está em `teams`, `tasks`, `task_dependencies`, `mail_messages` e `task_events`
- há reidratação e recuperação básica
- o fluxo dominante é centralizado no master
- `swarm_channels`, `swarm_threads` e `assistance_tasks` ainda não sustentam colaboração lateral real no dia a dia

Decisão implícita desta ADR: **o sistema atual sera tratado oficialmente como team orchestration com mailbox, e nao como swarm cooperativo pleno**.

### 2.2 Dashboard

O dashboard atual serve para visualização, mas ainda não pode ser tratado como superfície operacional canônica.

- SSE é fan-out em memória
- não há replay, persistência curta de eventos nem backpressure handling robusto
- `/status`, scripts shell e dashboard não compartilham a mesma verdade consolidada
- o dashboard aparenta mais maturidade do que realmente existe

Decisão implícita desta ADR: **dashboard sem snapshot consolidado nao sera considerado "operacional"**.

### 2.3 Memória e Dados

Hoje há drift entre documentação, código e destino dos dados.

- SQLite já é a verdade do runtime local
- Qdrant existe, mas com uso e schema inconsistentes
- Supabase local existe, mas ainda não é o sistema de registro do runtime
- Obsidian não tem contrato claro na arquitetura
- há risco real de vazamento de contexto entre bots sem namespace obrigatório

Decisão implícita desta ADR: **memória infinita sem contrato de namespace, origem e versionamento nao sera aceita como arquitetura valida**.

---

## 3. Decisão

Implementar a arquitetura alvo "Basico Bem Feito Aurélia 2026" com as seguintes decisões normativas.

### 3.1 Liderança do time

`aurelia_code` passa a ser a identidade principal visível e operacional do bot líder do sistema.

Regras:

- `aurelia_code` é a líder do time e a única autoridade de coordenação global
- o identificador legado `aurelia` continua aceito como alias compatível durante a transição
- nenhum bot especialista fala em nome do sistema inteiro sem passar por `aurelia_code`
- especialistas podem pedir ajuda lateral, mas ownership de coordenação e fechamento do trabalho volta para `aurelia_code`

### 3.2 Modelo de colaboração

O time passa a operar oficialmente no modelo:

- **líder + especialistas**
- colaboração lateral controlada
- ownership explícito
- trilha auditável

Bots oficiais desta camada:

- `aurelia_code` -> líder técnica, coordenadora e curadora de contexto
- `caixa_cpf_cnpj` -> financeiro e insistência operacional
- `ac_vendas` -> comercial
- `organizacao_de_obras` -> execução e acompanhamento de obras
- `agenda_cpf` -> agenda pessoal e compromissos

Regras:

- delegação: a líder cria e distribui trabalho
- handoff: ownership muda explicitamente
- assistência lateral: ajuda entre especialistas não muda ownership por padrão
- toda cooperação relevante precisa deixar evento auditável

### 3.3 Sistema de registro e responsabilidades por camada

O sistema de registro oficial será:

- **Supabase local** -> sistema de registro canônico
- **SQLite** -> runtime local e working set operacional
- **Qdrant** -> índice semântico derivado
- **Obsidian CLI** -> interface editorial bidirecional controlada

#### Supabase local

Supabase passa a ser o registro principal de:

- bots e perfis
- relações entre bots e políticas de acesso
- memória curada
- knowledge items
- status consolidado de componentes
- auditoria de alto nível
- entidades estruturadas do negócio e agenda

Supabase deixa de ser "planejado" e passa a ser "obrigatório" para o alvo final.

#### SQLite

SQLite fica restrito a:

- filas locais
- leases
- mailbox
- cron queue
- cache operacional
- memória curta e replay de sessão
- buffers locais tolerantes a restart

SQLite não será expandido como banco híbrido central de negócio.

#### Qdrant

Qdrant será apenas projeção vetorial derivada.

Regras:

- nenhum dado nasce canonicamente no Qdrant
- toda escrita vetorial deve apontar para um `source_system` e `source_id`
- toda leitura vetorial deve aplicar namespace mínimo por bot e domínio
- falha em busca vetorial deve expor estado degradado, não lista vazia silenciosa

#### Obsidian CLI

Obsidian entra oficialmente na arquitetura como camada bidirecional controlada.

Regras:

- Obsidian pode originar ou receber conhecimento
- nada escrito no vault se torna verdade canônica sem sync auditável
- toda sincronização precisa registrar origem, autor, versão, timestamp e destino
- conflitos são resolvidos por pipeline de reconciliação, não por sobrescrita silenciosa

---

## 4. Contratos Obrigatórios

### 4.1 Identidade de bot

Todo evento, task, cron, item de memória, snapshot e documento sincronizado deve carregar:

- `canonical_bot_id`
- `display_name`
- `persona_id`
- `domain`

Durante a transição:

- `aurelia` e `aurelia_code` devem resolver para a mesma líder
- o sistema deve sempre persistir o identificador canônico escolhido

### 4.2 Item canônico de memória

Todo item canônico de memória/conhecimento deve possuir, no mínimo:

```json
{
  "id": "uuid",
  "canonical_bot_id": "aurelia_code",
  "kind": "message|note|decision|knowledge|task_summary|event",
  "domain": "system|personal|business|finance|operations",
  "source_system": "telegram|sqlite|supabase|obsidian|cron|dashboard",
  "source_id": "external-or-local-id",
  "content": "texto canônico",
  "metadata": {},
  "version": 1,
  "created_at": "2026-03-25T00:00:00Z",
  "updated_at": "2026-03-25T00:00:00Z"
}
```

### 4.3 Projeção vetorial

Toda projeção para Qdrant deve carregar, no mínimo:

```json
{
  "canonical_bot_id": "aurelia_code",
  "persona_id": "aurelia-leader",
  "domain": "system",
  "source_system": "supabase",
  "source_id": "uuid",
  "text": "conteúdo indexável",
  "ts": 1774396800,
  "version": 1
}
```

É proibido manter payload legado sem namespace como padrão de futuro.

### 4.4 Contrato de colaboração entre bots

Toda task do time deve explicitar:

- `owner_bot_id`
- `requested_by_bot_id`
- `helper_bot_id` quando houver ajuda lateral
- `team_run_id`
- `status`
- `depends_on`
- `handoff_reason`
- `resolution_summary`

Sem esses campos, o time volta a parecer inteligente no chat, mas continua opaco operacionalmente.

### 4.5 Contrato de status operacional

Todo componente relevante deve produzir snapshot consolidado com:

- `component`
- `status` -> `healthy|degraded|offline|stale`
- `latency_ms`
- `queue_depth`
- `last_ok_at`
- `last_error`
- `source`

Este snapshot será a única base para:

- dashboard web
- `/status` no Telegram
- watchdogs
- alertas
- health scripts

---

## 5. Dashboard Operacional de Verdade

O dashboard em `https://aurelia.zappro.site/` passa a ter o seguinte contrato:

- não é vitrine; é superfície operacional oficial
- não é dono da verdade; consome a verdade consolidada
- SSE deixa de ser o centro do desenho e vira apenas meio de entrega
- o dashboard precisa refletir liderança, time, memória, cron e saúde real dos serviços

Requisitos mínimos do alvo final:

1. **Status unificado**
   - dashboard, `/status` e scripts de health leem a mesma fonte consolidada

2. **Team real**
   - mostrar `aurelia_code` como líder
   - mostrar especialistas reais
   - distinguir claramente status visual de presença e estado real de tasks

3. **Brain real**
   - indicar se está em modo semântico, lexical fallback ou degradado
   - jamais esconder falha estrutural como "sem resultados"

4. **Cron real**
   - mostrar filas, próximos jobs, atrasos e falhas
   - jobs críticos não podem ficar invisíveis atrás de um status genérico "online"

5. **Audit trail**
   - decisões, handoffs, ajudas laterais e falhas relevantes devem aparecer de forma rastreável

---

## 6. Fases de Implementação

### Fase 1: Fundacao e honestidade arquitetural

- declarar formalmente o estado atual como `team orchestration com mailbox`
- introduzir `canonical_bot_id`
- promover `aurelia_code` como líder oficial com alias compatível
- corrigir contratos de documentação que vendem maturidade inexistente
- fixar schema canônico de memória, task e status

### Fase 2: Supabase como centro canônico

- migrar registro principal de bots, knowledge items, memória curada e snapshots para Supabase local
- manter SQLite apenas como runtime local
- criar auditoria básica de sincronização e origem

### Fase 3: Memória infinita auditável

- implementar pipeline canônico `Supabase -> indexador -> Qdrant`
- implementar pipeline `Obsidian CLI <-> sync controlado <-> Supabase`
- bloquear acesso vetorial sem namespace mínimo
- remover silent failure de busca e leitura semântica

### Fase 4: Team bots como equipe real

- institucionalizar os bots do print como especialistas oficiais
- separar delegação, handoff e assistência lateral
- garantir ownership, recovery e auditoria
- impedir que especialistas virem silos soltos sem coordenação da líder

### Fase 5: Dashboard operacional

- trocar a ideia de "eventos em memória = verdade"
- consolidar snapshots persistidos
- alinhar dashboard, `/status` e watchdogs
- expor estados `healthy`, `degraded`, `offline` e `stale`

---

## 7. Consequências

### Positivas

- o repositório sai da zona cinzenta entre protótipo inteligente e sistema operacional confiável
- `aurelia_code` vira líder explícita do time, sem ambiguidade de identidade
- Supabase local passa a ter função real em vez de existir só como infraestrutura pronta
- Qdrant deixa de ser pseudo-banco de verdade e vira índice derivado
- Obsidian entra no desenho sem causar anarquia documental
- o dashboard passa a refletir o sistema real

### Negativas

- haverá trabalho de migração e saneamento de contratos
- o sistema perderá parte da "mística" de swarm total até a colaboração lateral virar fato operacional
- várias partes do runtime terão de parar de usar fallback silencioso
- haverá custo inicial de reconciliar aliases, schemas e responsabilidades

### Riscos

- tentar acelerar tudo de uma vez e acabar produzindo mais camadas sem fechar as antigas
- promover Supabase a centro sem concluir o contrato de sync com SQLite e Obsidian
- manter payloads legados no Qdrant e chamar isso de compatibilidade
- continuar usando dashboard como palco bonito em vez de painel operacional

### Mitigacao

- fases curtas e disciplina de contratos
- compatibilidade temporária explícita para `aurelia` -> `aurelia_code`
- auditoria de origem e versionamento em toda sincronização
- proibir novas extensões sem namespace, ownership e source of truth definidos

---

## 8. Verificacao de Aceite

Esta ADR só poderá ser considerada cumprida quando os seguintes cenários forem verdadeiros:

- `aurelia_code` aparece como líder oficial do time, com compatibilidade controlada para legado
- um bot não consegue recuperar memória de outro sem namespace e regra explícita
- Qdrant não recebe escrita direta sem `source_system`, `source_id` e `canonical_bot_id`
- Obsidian consegue sincronizar conhecimento sem sobrescrita silenciosa
- dashboard e `/status` mostram o mesmo estado consolidado
- falha de Qdrant, Supabase ou SQLite aparece como degradação explícita
- cron crítico não fica invisível nem bloqueado sem sinalização
- delegação, handoff e assistência lateral deixam trilha auditável

---

## 9. Regra Final

Fica proibido no ecossistema Aurélia:

- chamar de swarm cooperativo algo que ainda é apenas fila com workers
- chamar de memória infinita algo sem namespace, origem e versionamento
- chamar de dashboard operacional algo que não consome snapshot canônico
- chamar de integração com Supabase algo que ainda vive só em documento
- chamar de sync com Obsidian algo que não seja auditável

O padrão oficial daqui em diante é:

**menos fantasia de arquitetura e mais básico bem feito, com liderança clara, dados com dono, memória auditável e operação visível.**
