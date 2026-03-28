---
description: Governa a stack de dados, bancos e infraestrutura de vetores.
id: 14-data-stack-governance
---

# 📦 Regra 14: Governança de Data Stack (SOTA 2026.1)

Aurelia mantém soberania absoluta sobre seus dados.

## 1. Persistência Estruturada
- **Relacional**: PostgreSQL (Docker) para faturamento, usuários e metadados.
- **Lite**: SQLite para cache local e estados voláteis de agentes.

## 2. Memória Inteligente (Sovereign Hub)
- **Vetor**: Qdrant operando em porta `6333`.
- **Sync**: A sincronização semântica deve seguir o motor `memory-sync-vector-db` em conformidade com o `obsidian-sync.sh`.

## 3. Segurança
- Zero Hardcode em qualquer camada de dados.
- Auditoria periódica via `/env-audit`.

---
*Assinado: Aurélia (Arquiteta Líder) — Março 2026*
