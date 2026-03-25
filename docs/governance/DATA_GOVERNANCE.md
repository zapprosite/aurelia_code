# Data Governance — Aurélia Homelab

> Autoridade: Will (Principal Engineer)
> Última revisão: 2026-03-24
> Escopo: Qdrant, SQLite, Supabase, multi-bot (S-32)

---

## 1. Visão Geral da Camada de Dados

```
┌─────────────────────────────────────────────────────┐
│                  Aurélia + Team Bots                │
└──────────┬──────────────────────────────────────────┘
           │
     ┌─────┴──────────────────────────────────────┐
     │              Dados em uso                  │
     │                                            │
     │  SQLite        Qdrant          Supabase    │
     │  (memória      (semântico      (relacional │
     │   curta,        + vector)       + REST)    │
     │   tasks,                                   │
     │   crons)                                   │
     └────────────────────────────────────────────┘
```

### Regra mestra
> **SQLite** para estado de runtime e memória curta.
> **Qdrant** para busca semântica e contexto histórico.
> **Supabase** para dados estruturados do negócio (obras, leads, agenda).
> Nunca misture responsabilidades entre as três camadas.

---

## 2. SQLite

### Arquivos
| Arquivo | Conteúdo | Owner |
|---|---|---|
| `~/.aurelia/data/aurelia.db` | Histórico de conversas (`MemoryManager`), cron jobs, voice mirror | Aurélia |
| `~/.aurelia/data/aurelia.db.teams` | Task store do swarm (SQLiteTaskStore), dependency graph | BotPool |

### Política
- Nenhum bot cria tabelas adicionais no SQLite sem ADR.
- `aurelia.db` é memória curta — janela de `MemoryWindowSize` mensagens (padrão 20).
- Backup diário automático antes de qualquer migração estrutural.
- **Não compartilhar** `aurelia.db` entre instâncias diferentes do processo.

### Limpeza
- Cron semanal de vacuum: `VACUUM; ANALYZE;`
- Mensagens com mais de 90 dias no histórico: truncar via cron job.

---

## 3. Qdrant

**URL:** `http://localhost:6333`
**API Key:** em `~/.aurelia/config/secrets.env` → `QDRANT_API_KEY`
**Modelo de embedding:** `bge-m3` (1024 dimensões, Cosine similarity)

### Collections — Catálogo Oficial

| Collection | Propósito | Namespace | on_disk |
|---|---|---|---|
| `conversation_memory` | Histórico semântico de todas as conversas | `bot_id` + `chat_id` | false (hot) |
| `aurelia_skills` | Embeddings de skills para semantic routing | `skill_id` + `category` | false (hot) |
| `knowledge_hvac` | Docs técnicos DAIKIN VRV, normas, especificações | `category` | true (large) |
| `knowledge_personal` | Agenda, contexto familiar, notas pessoais | `type` | true (large) |

### Schema de payload obrigatório — `conversation_memory`

Todo ponto inserido DEVE ter os seguintes campos no payload:

```json
{
  "text":       "conteúdo da mensagem ou resumo",
  "bot_id":     "aurelia",
  "persona_id": "aurelia-leader",
  "chat_id":    123456789,
  "domain":     "professional | personal | system",
  "ts":         1711234567
}
```

Pontos sem `bot_id` são considerados **legado** e serão limpos pelo Sentinel (S-34).

### Schema de payload — `aurelia_skills`

```json
{
  "text":        "descrição semântica da skill",
  "skill_id":    "commit-message",
  "category":    "git | code | planning | ...",
  "version":     "1.0"
}
```

### Schema de payload — `knowledge_hvac`

```json
{
  "text":       "conteúdo técnico",
  "category":   "especificacao | preco | norma | case",
  "source":     "daikin-catalog-2025 | obra-xyz | manual",
  "updated_at": "2026-03-24"
}
```

### Schema de payload — `knowledge_personal`

```json
{
  "text":   "nota ou evento",
  "type":   "agenda | familia | saude | financeiro",
  "pessoa": "filha | namorada | familia | will",
  "date":   "2026-03-24"
}
```

### Regras de Collections

1. **Proibido** criar collections com nome de bot (`conversation_memory_hvac_sales`).
2. Toda collection nova exige ADR aprovado antes de criar.
3. Collections com 0 pontos por mais de 30 dias serão deletadas pelo Sentinel.
4. Índices de payload DEVEM ser criados no momento de criação da collection (não depois).

### Índices de payload criados (estado atual)

```
conversation_memory: bot_id, persona_id, chat_id, domain  (keyword)
aurelia_skills:      skill_id, category                    (keyword)
knowledge_hvac:      category                              (keyword)
knowledge_personal:  category                              (keyword)
```

---

## 4. Supabase Local

**Status:** Instalado. 13 containers (Kong, PostgREST, Auth, Storage, Realtime, Studio).
**Studio:** `http://localhost:3000` (porta padrão Supabase local)
**PostgREST:** usado pelos bots via HTTP para CRUD sem SQL direto.

### Schemas — Separação por Domínio

```
public          ← dados do negócio (obras, leads)
personal        ← dados pessoais (agenda, familia)  [RLS obrigatório]
system          ← metadados internos da Aurélia
```

### Tabelas Core (a criar via migração S-33)

#### `public.obras`
```sql
id          uuid primary key default gen_random_uuid()
nome        text not null
cliente     text
status      text check (status in ('orcamento','em_andamento','entregue','cancelado'))
orcamento   numeric(12,2)
inicio      date
entrega     date
bot_id      text not null default 'project-manager'
created_at  timestamptz default now()
updated_at  timestamptz default now()
```

#### `public.leads`
```sql
id          uuid primary key default gen_random_uuid()
nome        text
telefone    text
produto     text check (produto in ('VRV','Split','HVAC-R','Outro'))
status      text check (status in ('novo','contato','proposta','fechado','perdido'))
valor_est   numeric(12,2)
bot_id      text not null default 'hvac-sales'
created_at  timestamptz default now()
updated_at  timestamptz default now()
```

#### `personal.agenda`
```sql
id          uuid primary key default gen_random_uuid()
titulo      text not null
data        timestamptz not null
tipo        text check (tipo in ('academia','igreja','familia','obra','reuniao','pessoal'))
recorrente  boolean default false
bot_id      text not null default 'life-organizer'
created_at  timestamptz default now()
```

### Política de Acesso
- Schema `personal` → RLS habilitado. Apenas `life-organizer` + `aurelia` podem ler/escrever.
- Schema `public` → bots de negócio têm acesso total ao seu `bot_id`.
- Nenhum bot lê dados de outro bot sem passar pela Aurélia (líder).

---

## 5. Multi-Bot — Namespace e Responsabilidades

### Mapa de bots e suas fontes de dados

| Bot ID | Telegram | Qdrant reads | Supabase writes |
|---|---|---|---|
| `aurelia` | chat pessoal do Will | `conversation_memory` (todos) + `aurelia_skills` | tudo (leitura) |
| `hvac-sales` | grupo Comercial | `conversation_memory` (filtro `bot_id=hvac-sales`) + `knowledge_hvac` | `public.leads` |
| `project-manager` | grupo Obras | `conversation_memory` (filtro `bot_id=project-manager`) + `knowledge_hvac` | `public.obras` |
| `life-organizer` | chat pessoal / grupo família | `conversation_memory` (filtro `bot_id=life-organizer`) + `knowledge_personal` | `personal.agenda` |

### Regra de escrita no Qdrant
Todo insert em `conversation_memory` DEVE incluir `bot_id` no payload.
O `ContextAssembler` filtra por `bot_id` automaticamente (a implementar em S-35).

---

## 6. CLI de Operação

Localização: `scripts/bot-cli.sh`

```bash
# Listar bots ativos
./scripts/bot-cli.sh list

# Listar personas disponíveis
./scripts/bot-cli.sh personas

# Criar novo bot
./scripts/bot-cli.sh add <id> <nome> <token> <persona_id> [focus_area] [user_ids_csv]
# Exemplo:
./scripts/bot-cli.sh add hvac-sales "Bot de Vendas" "123:AAA..." hvac-sales \
  "Funil DAIKIN VRV SP" "7220607041"

# Remover bot
./scripts/bot-cli.sh remove hvac-sales

# Ping Aurélia via impersonação (sem abrir Telegram)
./scripts/bot-cli.sh ping "Qual o status das obras?"

# Ver estado do Qdrant
./scripts/bot-cli.sh qdrant
```

---

## 7. Sentinel de Dados (S-34 — a implementar)

Cron job semanal (segunda-feira 09:00) que executa:

1. **Auditoria Qdrant:** conta pontos por `bot_id`, detecta pontos sem namespace.
2. **Auditoria SQLite:** tamanho, últimas mensagens, vacuum se > 500MB.
3. **Auditoria Supabase:** obras/leads sem update há > 7 dias, agenda vencida.
4. **Relatório:** enviado via Telegram para o chat do Will.

---

## 8. Próximas Slices

| Slice | O que faz | Deps |
|---|---|---|
| **S-33** | Supabase tools (insert/query obras, leads, agenda) | Supabase up |
| **S-34** | Sentinel de dados (cron + auditoria + relatório) | S-32 |
| **S-35** | Namespace `bot_id` no ContextAssembler (Qdrant filter) | S-32 |

---

## 9. Histórico de Mudanças

| Data | Ação | Por quê |
|---|---|---|
| 2026-03-24 | Deletadas 6 collections de teste (`rag_docs`, `catalog_embeddings`, `main`, `app_voice_v1`, `app_journal_v1`, `rag_governance_v1`) | Dados de setup/teste sem valor operacional |
| 2026-03-24 | Criadas 4 collections de governança com índices de payload | S-32 multi-bot, namespace `bot_id` |
| 2026-03-24 | `qdrant_url` corrigido de `""` para `http://localhost:6333` | Config estava vazio, app usava fallback silencioso |
