# Data Stack Standard

> Status: ativo
> Autoridade: Will + Aurélia
> Enforcement: `.agent/rules/14-data-stack-governance.md`
> ADR fonte: `docs/adr/20260325-data-stack-contract-and-templates.md`

## 1. Princípios não negociáveis

1. `Supabase` é verdade canônica.
2. `SQLite` é estado operacional local.
3. `Qdrant` é índice derivado.
4. `Obsidian CLI` é interface editorial controlada.
5. Todo dado importante precisa de `app_id`, `repo_id`, `environment` e `canonical_bot_id`.
6. Se um dado existe em mais de um lugar, um deles deve ser explicitamente marcado como canônico e os outros como derivados.

## 2. Organização por app, repo e instância

### Identificadores obrigatórios
- `app_id`: fronteira do produto ou runtime.
- `repo_id`: fronteira do repositório Git.
- `environment`: `local`, `staging`, `prod`.
- `canonical_bot_id`: dono lógico da ação ou do registro.

### Regra prática
- `1 app_id` pode ter vários `repo_id`.
- `1 environment` pode ter vários `app_id`.
- `1 instância local` = `1 AURELIA_HOME` por `app_id + environment`.

### Layout recomendado de instâncias
```text
/srv/aurelia/instances/
  aurelia-local/
  aurelia-staging/
  hvac-local/
  hvac-prod/
```

Cada instância aponta seu próprio `AURELIA_HOME`.

## 3. Papel de cada tecnologia

### Supabase
Use para:
- entidades estruturadas
- memória canônica
- knowledge canonizado
- cadastro de apps, repos e bots
- auditoria operacional
- sync jobs e snapshots

Não use para:
- fila efêmera
- cache temporário
- busca vetorial

### SQLite
Use para:
- cron
- mailbox
- fila local
- leases
- cache
- replay curto
- estado transitório do runtime

Não use para:
- verdade principal de negócio
- storage semântico de longo prazo

### Qdrant
Use para:
- recuperação semântica
- ranking
- índice vetorial de memória/knowledge/skills

Não use para:
- criação canônica de registros
- metadados sem `source_id`
- isolamento por collection improvisada

### Obsidian CLI
Use para:
- ADR
- runbooks
- notas curadas
- material humano editável
- revisão editorial

Não use para:
- store canônico implícito
- backlog operacional sem status/owner
- notas sem frontmatter e sem lineage

## 4. Perfis oficiais de aplicação

### `app-lite`
Para micro-apps, utilitários e ferramentas simples.

```text
Supabase schemas: core, ops, app_<slug>
Qdrant: opcional
SQLite: runtime.sqlite
Obsidian: 20-apps/<app_id>/
```

### `app-business`
Para CRM, obras, agenda, dashboards operacionais.

```text
Supabase schemas: core, ops, memory, app_<slug>
Qdrant: memory_items, knowledge_items
SQLite: runtime.sqlite, teams.sqlite
Obsidian: 20-apps/<app_id>/, 40-runbooks/
```

### `app-knowledge`
Para bibliotecas, pesquisa, acervos e memória longa.

```text
Supabase schemas: core, ops, memory
Qdrant: knowledge_items, skills_index
SQLite: runtime.sqlite
Obsidian: forte uso editorial
```

### `app-agentic`
Para multi-bot, cron, mailbox, operações de equipe e dashboard.

```text
Supabase schemas: core, ops, memory, app_<slug>
Qdrant: memory_items, knowledge_items, skills_index
SQLite: runtime.sqlite, teams.sqlite, cache.sqlite
Obsidian: 10-governance/, 20-apps/, 30-repos/, 40-runbooks/
```

## 5. Contrato de naming

### Campos mandatórios em registros canônicos
```text
id
app_id
repo_id
environment
canonical_bot_id
source_system
source_id
created_at
updated_at
version
```

### Campos mandatórios em payload vetorial
```text
app_id
repo_id
environment
canonical_bot_id
domain
source_system
source_id
text
ts
version
```

## 6. Política de mudança

É proibido, sem ADR:
- criar schema novo fora do padrão
- criar collection nova fora do catálogo
- adicionar novo arquivo SQLite persistente sem manifesto
- criar pasta canônica nova no vault
- mudar nomenclatura de campos base

## 7. Donos do contrato

- Humano operador: aprova mudança estrutural
- `aurelia_code`: guarda coerência arquitetural
- `controle-db`: executa higiene, inventário, auditoria e enforcement operacional

## 8. Documentos subordinados
- [SCHEMA_REGISTRY.md](/home/will/aurelia/docs/governance/SCHEMA_REGISTRY.md)
- [OBSIDIAN_VAULT_STANDARD.md](/home/will/aurelia/docs/governance/OBSIDIAN_VAULT_STANDARD.md)
- [templates/README.md](/home/will/aurelia/docs/governance/templates/README.md)
