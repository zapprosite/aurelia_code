# ADR 20260325: Data Stack Contract, Templates e Persistência de Governança

## Status
Aprovado

## Contexto
O repositório já tinha boas intenções sobre `Supabase`, `Qdrant`, `SQLite` e `Obsidian CLI`, mas ainda operava com drift entre ADRs, governança e runtime. O risco real não era falta de tecnologia. Era a repetição do padrão mais comum em stacks agentic: o mesmo dado aparecendo em três lugares, com nomes diferentes, sem fonte de verdade, sem lineage e sem regra de descarte.

Isso produz a "salada":

- `Supabase` usado como dump genérico de tudo
- `Qdrant` tratado como segundo banco de dados
- `SQLite` crescendo além do papel de runtime local
- `Obsidian` misturando rascunho, runbook, ADR, backlog e memória canônica
- LLMs criando schema, collection, pasta ou payload por conveniência

O operador pediu um padrão persistente, contratual e difícil de ser degradado por improviso futuro.

## Decisão
O repositório passa a ter um contrato normativo único para o data stack, composto por:

1. `docs/governance/DATA_STACK_STANDARD.md`
2. `docs/governance/SCHEMA_REGISTRY.md`
3. `docs/governance/OBSIDIAN_VAULT_STANDARD.md`
4. `docs/governance/templates/*`
5. `.agent/rules/14-data-stack-governance.md`

Esses documentos deixam de ser recomendação e passam a ser regra operacional.

## Regras obrigatórias

### 1. Papel de cada camada
- `Supabase` é o sistema de registro canônico.
- `SQLite` é runtime local, fila, cron, mailbox, cache e estado curto.
- `Qdrant` é índice derivado para busca semântica.
- `Obsidian CLI` é superfície editorial humana controlada.

### 2. Proibição de ambiguidade
Nenhum agente pode:

- criar collection por bot, chat, repo ou experimento
- criar tabela sem migration e sem schema previsto
- escrever dado original diretamente no `Qdrant`
- tratar nota do `Obsidian` como verdade canônica automática
- usar `SQLite` como banco principal de negócio

### 3. Contrato de identidade e lineage
Toda entidade relevante deve carregar, no mínimo:

- `app_id`
- `repo_id`
- `environment`
- `canonical_bot_id`
- `source_system`
- `source_id`
- `created_at`
- `updated_at`
- `version`

### 4. Organização por app/repo
- O isolamento principal é por `app_id` e `environment`.
- `repo_id` organiza conhecimento, notas, governança e provenance.
- O mesmo `Supabase` e a mesma instância `Qdrant` podem atender múltiplos apps do mesmo ambiente, desde que o isolamento seja feito por schema/payload e não por improviso nominal.
- Para múltiplas instâncias locais, cada app/ambiente deve usar seu próprio `AURELIA_HOME`.

### 5. Templates oficiais
O repositório passa a fornecer templates oficiais para:

- `Supabase`
- `Qdrant`
- `SQLite`
- `Obsidian`

Nenhum novo app deve começar sem escolher um dos perfis oficiais documentados.

## Perfis padrão suportados

### Perfil A — `app-lite`
Para apps simples sem workflow intensivo.

- `Supabase`: `core`, `ops`, `app_<slug>`
- `Qdrant`: opcional, no máximo `knowledge_items`
- `SQLite`: `runtime.sqlite`
- `Obsidian`: `20-apps/<app_id>/`

### Perfil B — `app-business`
Para apps transacionais e dashboards.

- `Supabase`: `core`, `ops`, `memory`, `app_<slug>`
- `Qdrant`: `memory_items`, `knowledge_items`
- `SQLite`: `runtime.sqlite`, `teams.sqlite`
- `Obsidian`: `20-apps/<app_id>/`, `40-runbooks/`

### Perfil C — `app-knowledge`
Para apps focados em documentos, notas e recuperação semântica.

- `Supabase`: `core`, `ops`, `memory`
- `Qdrant`: `knowledge_items`, `skills_index`
- `SQLite`: `runtime.sqlite`
- `Obsidian`: uso intenso com curadoria obrigatória

### Perfil D — `app-agentic`
Para sistemas multi-bot, cron, mailbox e governança operacional.

- `Supabase`: `core`, `ops`, `memory`, `app_<slug>`
- `Qdrant`: `memory_items`, `knowledge_items`, `skills_index`
- `SQLite`: `runtime.sqlite`, `teams.sqlite`, `cache.sqlite`
- `Obsidian`: `10-governance/`, `20-apps/<app_id>/`, `30-repos/<repo_id>/`, `40-runbooks/`

## Consequências

### Positivas
- reduz a chance de duplicação semântica entre stores
- cria uma fronteira clara entre verdade, projeção e rascunho humano
- melhora portabilidade entre apps e repositórios
- dificulta regressões causadas por LLMs que "inventam organização"

### Negativas
- reduz liberdade de prototipagem caótica
- obriga naming consistente e discipline de migrations/manifests
- exige ADR para desvios legítimos

## Enforcement
Qualquer mudança estrutural em `Supabase`, `Qdrant`, `SQLite` ou `Obsidian` fora deste contrato exige:

1. nova ADR em `docs/adr/`
2. atualização do `SCHEMA_REGISTRY`
3. atualização da rule `14-data-stack-governance`
4. atualização dos templates, se o novo padrão passar a ser reutilizável

## Artefatos
- [DATA_STACK_STANDARD.md](/home/will/aurelia/docs/governance/DATA_STACK_STANDARD.md)
- [SCHEMA_REGISTRY.md](/home/will/aurelia/docs/governance/SCHEMA_REGISTRY.md)
- [OBSIDIAN_VAULT_STANDARD.md](/home/will/aurelia/docs/governance/OBSIDIAN_VAULT_STANDARD.md)
- [14-data-stack-governance.md](/home/will/aurelia/.agent/rules/14-data-stack-governance.md)
- [templates/README.md](/home/will/aurelia/docs/governance/templates/README.md)
