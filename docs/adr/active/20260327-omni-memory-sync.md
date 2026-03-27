# ADR 20260327: Omni-Memory Sync (SOTA 2026) 🧠🔄

## Status
Proposto

## Contexto
A memória da Aurélia está fragmentada:
- **Qdrant (L2)**: Busca semântica rápida, mas volátil/difícil de auditar em massa.
- **Supabase/Postgres (L3)**: Persistência relacional robusta, mas sem busca vetorial nativa otimizada na camada de aplicação.
- **Obsidian (L1)**: Interface humana da memória.

## Decisão
Implementar o **Omni-Memory Sync**, um worker Go que:
1.  **Bidirectional Pulse**: Sincroniza metadados do Postgres -> Qdrant (garantindo que todo fato tenha um embedding).
2.  **Conflict Resolution**: Usa timestamps para resolver divergências entre o Obsidian Vault e o banco de dados.
3.  **Zod-First Enforcement**: Valida todos os payloads de memória através dos esquemas em `packages/zod-schemas/`.

## Arquitetura
- **Engine**: Goroutines na System API monitorando canais de eventos ou varredura periódica (cron).
- **Integridade**: Checksum de conteúdo para evitar re-embedding desnecessário.

## Consequências
- **Positivas**: Memória imutável, auditável e extremamente rápida. Fim da "amnésia" entre sessões.
- **Negativas**: Aumento leve no uso de CPU durante a sincronização inicial de grandes vaults.
