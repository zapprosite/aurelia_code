# ADR 20260325: Bot CONTROLE_DB para Governança de Supabase, Qdrant, Obsidian e SQLite

## Status
Proposto

## Contexto
O ecossistema Aurélia já opera com múltiplas camadas de dados e memória:

- `SQLite` para runtime local
- `Qdrant` para busca semântica
- `Supabase local` para registro estruturado
- `Obsidian CLI` para camada editorial e reconciliação humana

O risco operacional atual não é apenas falta de feature. É desorganização cumulativa:

- coleções e artefatos de teste esquecidos
- drift entre runtime e documentação
- payloads legados sem namespace
- instâncias auxiliares ligadas à camada de dados sem dono explícito
- ausência de um guardião operacional transversal para saneamento e governança

## Decisão
Instituir um novo bot oficial no pool da Aurélia:

- **ID:** `controle-db`
- **Nome:** `CONTROLE DB`
- **Persona:** `data-governance`
- **Missão:** organizar, auditar e higienizar Supabase, Qdrant, Obsidian CLI e SQLite com trilha auditável

## Regras de Operação

1. O bot atua como guardião de governança de dados e instâncias ligadas à camada de dados.
2. Toda limpeza relevante deve seguir a sequência:
   - inventário
   - classificação
   - proposta
   - backup ou snapshot quando aplicável
   - execução
   - relatório
3. O bot pode remover artefatos de teste e experimentos esquecidos quando houver evidência suficiente.
4. O bot não pode tratar ausência de documentação como licença para apagar dados potencialmente canônicos.
5. Em caso de dúvida entre teste e produção, deve preservar, isolar e reportar.

## Consequências

### Positivas
- A camada de dados passa a ter um dono operacional explícito
- Supabase, Qdrant, Obsidian e SQLite entram em disciplina comum de higiene
- artefatos de teste deixam de se acumular sem revisão
- a governança de dados deixa de depender apenas da memória do operador

### Negativas
- aumenta a responsabilidade do pool multi-bot
- o bot precisará operar com alto rigor para não virar agente destrutivo
- a utilidade dele depende de trilha auditável e de critérios claros de limpeza

## Mitigação
- prompt com política de limpeza segura
- foco em evidência e relatório
- governança registrada em `docs/governance/DATA_GOVERNANCE.md`
