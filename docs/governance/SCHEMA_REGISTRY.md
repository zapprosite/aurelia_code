# Schema Registry

> Status: ativo
> Escopo: contratos reutilizáveis de `Supabase`, `Qdrant`, `SQLite` e `Obsidian`

## 1. Catálogo Supabase

### Schemas canônicos
- `core`
- `ops`
- `memory`
- `app_<slug>`

### Tabelas mínimas

#### `core.apps`
```sql
id text primary key
name text not null
environment text not null check (environment in ('local','staging','prod'))
owner_bot_id text not null
created_at timestamptz default now()
updated_at timestamptz default now()
```

#### `core.repos`
```sql
id text primary key
app_id text not null references core.apps(id)
git_remote text
default_branch text
created_at timestamptz default now()
updated_at timestamptz default now()
```

#### `core.bots`
```sql
id text primary key
app_id text not null references core.apps(id)
display_name text not null
persona_id text
role text not null
status text not null default 'active'
created_at timestamptz default now()
updated_at timestamptz default now()
```

#### `ops.audit_events`
```sql
id uuid primary key default gen_random_uuid()
app_id text not null
repo_id text
canonical_bot_id text not null
event_type text not null
source_system text not null
source_id text
payload jsonb not null default '{}'::jsonb
created_at timestamptz default now()
```

#### `ops.sync_jobs`
```sql
id uuid primary key default gen_random_uuid()
app_id text not null
repo_id text
direction text not null
source_system text not null
target_system text not null
status text not null
started_at timestamptz
finished_at timestamptz
metadata jsonb not null default '{}'::jsonb
```

#### `memory.memory_items`
```sql
id uuid primary key default gen_random_uuid()
app_id text not null
repo_id text
environment text not null
canonical_bot_id text not null
kind text not null
domain text not null
source_system text not null
source_id text not null
content text not null
metadata jsonb not null default '{}'::jsonb
version int not null default 1
created_at timestamptz default now()
updated_at timestamptz default now()
```

#### `memory.knowledge_items`
```sql
id uuid primary key default gen_random_uuid()
app_id text not null
repo_id text
environment text not null
canonical_bot_id text not null
category text not null
source_system text not null
source_id text not null
title text
content text not null
metadata jsonb not null default '{}'::jsonb
version int not null default 1
created_at timestamptz default now()
updated_at timestamptz default now()
```

## 2. Catálogo Qdrant

### Collections permitidas por padrão
- `memory_items`
- `knowledge_items`
- `skills_index`

### `memory_items` payload obrigatório
```json
{
  "app_id": "aurelia",
  "repo_id": "aurelia",
  "environment": "local",
  "canonical_bot_id": "aurelia_code",
  "domain": "operations",
  "source_system": "supabase",
  "source_id": "memory_items/<uuid>",
  "text": "conteúdo indexado",
  "ts": 1774400000,
  "version": 1
}
```

Índices obrigatórios:
- `app_id`
- `repo_id`
- `environment`
- `canonical_bot_id`
- `domain`
- `source_system`

### `knowledge_items` payload obrigatório
```json
{
  "app_id": "aurelia",
  "repo_id": "aurelia",
  "environment": "local",
  "canonical_bot_id": "controle-db",
  "category": "governance",
  "source_system": "supabase",
  "source_id": "knowledge_items/<uuid>",
  "text": "conteúdo indexado",
  "ts": 1774400000,
  "version": 1
}
```

### `skills_index` payload obrigatório
```json
{
  "app_id": "aurelia",
  "repo_id": "aurelia",
  "environment": "local",
  "skill_id": "incident-response",
  "category": "operations",
  "source_system": "filesystem",
  "source_id": "skills/incident-response",
  "text": "descrição semântica da skill",
  "ts": 1774400000,
  "version": 1
}
```

## 3. Catálogo SQLite

### Arquivos oficiais
- `runtime.sqlite`
- `teams.sqlite`
- `cache.sqlite`

### Responsabilidades

#### `runtime.sqlite`
- sessão
- cron
- mailbox leve
- estado do worker
- mirrors transitórios

#### `teams.sqlite`
- tasks
- dependencies
- team runs
- worker leases

#### `cache.sqlite`
- cache descartável
- materialização local
- replay curto de índices e sync

### Regra
Se um novo arquivo SQLite persistente for necessário, ele deve entrar neste registry antes de existir em produção.

## 4. Catálogo Obsidian

### Pastas oficiais
- `00-inbox/`
- `10-governance/`
- `20-apps/<app_id>/`
- `30-repos/<repo_id>/`
- `40-runbooks/`
- `90-archive/`

### Tipos oficiais de nota
- `adr`
- `runbook`
- `decision`
- `knowledge`
- `sync-report`
- `draft`

O contrato detalhado do vault está em [OBSIDIAN_VAULT_STANDARD.md](/home/will/aurelia/docs/governance/OBSIDIAN_VAULT_STANDARD.md).
