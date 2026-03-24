---
name: memory-sync-vector-db
description: Sincroniza a memória do repositório (markdown) -> Qdrant (embeddings) + Postgres (metadata) para acesso semântico local.
---

# 🧠 Memory Sync: Sovereign Vector Architecture

Esta skill governa o fluxo de persistência e recuperação de memória semântica da Aurélia, permitindo que LLMs locais (Gemma 3) ou remotas (MiniMax 2.7) acessem o histórico de decisões e código sem depender de busca textual pura.

## 🏗️ Arquitetura de Memória (2026)

### 1. Camada de Embeddings (Local Only)
- **Modelo**: `bge-multilingual-gemma` ou `bge-m3` via Ollama/Local Server.
- **Dimensão**: 1024 (ou conforme configurado no Qdrant).
- **Soberania**: O processamento vetorial nunca sai do Home Lab.

### 2. Camada Vetorial (Qdrant)
- **Engine**: Qdrant operando em container Docker.
- **Coleção**: `aurelia_semantic_index`.
- **Ação**: Inserção de chunks de arquivos `.md`, `docs/`, `.context/` e `logs/`.

### 3. Camada de Metadados (Postgres)
- **Engine**: PostgreSQL local.
- **Função**: Armazenar timestamps, referências de arquivos, `conversation_id` e tags de governança.
- **Join**: Consultas semânticas via Qdrant retornam IDs que são resolvidos no Postgres para contexto completo.

## 🔄 Fluxo de Sincronização (Automático)
- **Trigger**: Script `scripts/memory-sync-vector-db.sh` ou tool `sync-ai-context`.
- **Frequência**: A cada commit relevante ou via crontab diário.
- **Higiene**: Deduplicação automática de vetores baseada no hash do conteúdo original.

## 🛠️ Comandos de Troubleshooting
- `curl http://localhost:6333/collections`: Verificar coleções no Qdrant.
- `psql -U aurelia -c "SELECT count(*) FROM memory_metadata;"`: Contagem de registros.
- `docker logs memory-sync-worker`: Debugar falhas de sincronização.

## 📍 Quando usar
- Quando a Aurélia precisa "lembrar" de uma decisão de ADR de meses atrás.
- Antes de iniciar um novo slice complexo para resgatar padrões similares.
- Para gerar relatórios de evolução do projeto baseados em semântica.

## 🛡️ Guardrails
- **Não vaze segredos**: O sync deve ignorar `.env` e arquivos ignorados pelo `.gitignore`.
- **Custo Computacional**: Evite re-indexar todo o repo sem necessidade (use diffs).