---
description: Impede improviso estrutural em Supabase, Qdrant, SQLite e Obsidian sem contrato e ADR.
id: 14-data-stack-governance
---

# 🗃️ Regra 14: Governança do Data Stack

`Supabase`, `Qdrant`, `SQLite` e `Obsidian CLI` operam sob contrato rígido.
Nenhum agente pode reorganizar essas camadas por conveniência local.

<directives>
1. **Fonte de verdade obrigatória**:
   - `Supabase` = registro canônico
   - `SQLite` = runtime local e estado curto
   - `Qdrant` = índice derivado
   - `Obsidian` = superfície editorial controlada

2. **Leitura obrigatória** antes de qualquer mudança estrutural:
   - `docs/governance/DATA_STACK_STANDARD.md`
   - `docs/governance/SCHEMA_REGISTRY.md`
   - `docs/governance/OBSIDIAN_VAULT_STANDARD.md`
   - `docs/governance/DATA_GOVERNANCE.md`

3. **É proibido sem ADR**:
   - criar schema novo fora de `core`, `ops`, `memory`, `app_<slug>`
   - criar collection por bot, repo, chat ou experimento
   - escrever dado original diretamente no `Qdrant`
   - adicionar arquivo SQLite persistente fora do registry
   - criar pasta canônica nova no vault fora do padrão
   - mudar os campos-base `app_id`, `repo_id`, `environment`, `canonical_bot_id`, `source_system`, `source_id`

4. **Lineage obrigatório**:
   Todo registro canônico ou derivado precisa carregar provenance mínima conforme `DATA_STACK_STANDARD.md`.

5. **Templates oficiais**:
   Novos apps devem começar pelos templates em `docs/governance/templates/`.
   Se os templates não servirem, isso é sinal de ADR, não de improviso.

6. **Owner operacional**:
   `controle-db` é o guardião da higiene e governança operacional da camada de dados.
</directives>

## Referências

- [`docs/adr/20260325-data-stack-contract-and-templates.md`](../../docs/adr/20260325-data-stack-contract-and-templates.md)
- [`docs/governance/DATA_STACK_STANDARD.md`](../../docs/governance/DATA_STACK_STANDARD.md)
- [`docs/governance/SCHEMA_REGISTRY.md`](../../docs/governance/SCHEMA_REGISTRY.md)
- [`docs/governance/OBSIDIAN_VAULT_STANDARD.md`](../../docs/governance/OBSIDIAN_VAULT_STANDARD.md)
- [`docs/governance/DATA_GOVERNANCE.md`](../../docs/governance/DATA_GOVERNANCE.md)
