---
title: ADR-20260319 Governance Documentation Index
description: Inventário de toda a documentação de governança (dados, operação, observabilidade, segredos e rede)
owner: codex
updated: 2026-03-20
---

# ADR-20260319 — Índice de Documentação de Governança

**ADR mestre:** [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)

**Propósito:** Navegar pela documentação de governança organizada por áreas críticas e centralizar o status de cada artefato.

**Status geral:** Governança de dados, operacional e de observabilidade estão concluídas; governança de segredos foi adiada conforme solicitação, e governança de rede segue na fase de planejamento.

## Roteiro de navegação

### 1️⃣ GOVERNANÇA_DE_DADOS — ✅ CONCLUÍDA
**Objetivo:** definir onde cada dado vive, quem é responsável e como ele é protegido permanentemente.

| Documento | Propósito | Público |
|-----------|-----------|---------|
| [Schema Registry — SQLite](./schema-registry-sqlite.md) | Especifica cinco tabelas essenciais em `aurelia.db` (gateway, voice, cron, memory, health) com índices e retenção. | Desenvolvedores, DBA |
| [Schema Registry — PostgreSQL](./schema-registry-postgres.md) | Mapeia quatro instâncias Postgres (n8n, supabase, litellm, dev) e tabelas futuras para cada workload. | Desenvolvedores, DBA |
| [Domain Ownership Table](./data-governance-domain-ownership.md) | Registra owner, SLA e criticidade de cada tabela. | Operações, Compliance |
| [Store Selection Matrix](./data-governance-store-selection.md) | Árvore de decisão entre SQLite, Postgres, Qdrant e Supabase. | Arquitetos, Desenvolvedores |
| [Data Lifecycle Policy](./data-governance-lifecycle.md) | Rotinas: 30 dias hot, 90 dias warm e arquivamento gz. | Operações, Compliance |
| [Qdrant Collection Contract](./qdrant-collection-contract.md) | Define bge-m3 384-dim, payload e fluxo de sincronização. | Engenheiros de ML, Desenvolvedores |
| [Compliance Matrix](./data-governance-compliance.md) | Alinha CONTRACT.md e GUARDRAILS.md com trilhas de auditoria. | Compliance, Operações |

**Resumo:** Os sete documentos estão completos e aprovados, garantindo contratos de dados claros e auditáveis.

### 2️⃣ GOVERNANÇA_OPERACIONAL — ✅ CONCLUÍDA
**Objetivo:** manter a infraestrutura operacional saudável e pronta para incidentes.

| Documento | Propósito | Público |
|-----------|-----------|---------|
| [Health Checks](./operational-governance-health-checks.md) | Crons a cada 5/15/60/1440 minutos para containers, disco, VRAM, Qdrant e backfills. | Operações, SRE |
| [Backup Verification](./operational-governance-backup-verification.md) | Verificações diárias de frescor e simulação mensal de restore. | Operações, DBA |
| [Incident Response](./operational-governance-incident-response.md) | Runbooks P1-P4 para containers, OOM, disco, banco, tunelamento e breach. | Operações, On-call |

**Resumo:** Os três núcleos operacionais estão documentados e integrados aos cron jobs e runbooks.

### 3️⃣ OBSERVABILIDADE_GOVERNANCE — ✅ CONCLUÍDA
**Objetivo:** monitorar métricas e alertas para detectar problemas antes que escalem.

| Documento | Propósito | Público |
|-----------|-----------|---------|
| [Metrics Contract](./observability-governance-metrics.md) | Métricas obrigatórias por serviço, Prometheus, dashboards e regras de alerta. | Desenvolvedores, SRE |

**Resumo:** Contrato de métricas implementado com dashboards para saúde, banco, memória e LLM, além de alertas críticos configurados.

### 4️⃣ GOVERNANÇA_DE_SEGREDOS — ⏳ ADIADA (por solicitação)
**Objetivo:** proteger credenciais sensíveis e rotacionar chaves trimestralmente.

**Documentos planejados:** KeePassXC Vault & migração, playbook de rotação, checklist de limpeza de credenciais em texto puro.

**Status:** deferida em respeito à instrução do usuário; permanece um ponto crítico não coberto.

### 5️⃣ GOVERNANÇA_DE_REDE — ⏳ EM PLANEJAMENTO
**Objetivo:** isolar tráfego, documentar exposição de portas e proteger perfis de rede.

**Documentos planejados:** hardening UFW (SSH), matriz de portas públicas vs internas, segurança do Cloudflare Tunnel, redes Docker separadas e política de ACL do Tailscale.

**Status:** nenhum artefato criado ainda; representa risco de superfície exposta se não priorizada na fase 4.

### 6️⃣ MATRIZ_DE_CONFORMIDADE — ✅ INTEGRADA
Já faz parte da governança de dados via [data-governance-compliance.md](./data-governance-compliance.md): cobre `CONTRACT.md`, `GUARDRAILS.md`, trilhas de auditoria, autorizações e caminhos de escalonamento.

## Observações de estabilidade (itens suspeitos/incompletos)
- Governança de segredos permanece skippada; documentar os playbooks e a limpeza de credenciais é necessário para tornar esse risco sustentável.
- Governança de rede ainda não tem artefatos; planejar um sprint da fase 4 para criar UFW/matriz/ACLs e reduzir a superfície exposta.
- As tarefas de fase 2 listadas na execução (remoção de `app.json.bak*`, refatoração de `mcp_servers_config.json`, política de rotação) continuam pendentes e devem constar no backlog para não perder referencial.
- Garantir que todos os textos acima fiquem em português para manter a memória técnica estável e consistente com o perfil da Aurélia.

## Resumo de execução
- ✅ Fase 1 (Crítico) — o usuário pediu para suspender o vault, portanto não houve mudança aqui.
- ✅ Fase 2 (Alta) — documentação de schema registry, timers e app.json (pendências listadas acima ainda precisam fechar).
- ✅ Fase 3 (Média) — health checks, backups, incident response, Qdrant, data lifecycle e métricas estão documentados e integrados à operação.
- ⏳ Fase 4 (Baixa) — limpeza de diretórios de deploy, mensagens de desligamento sujo e scripts de compliance permanecem na fila.

## Estatísticas dos documentos

| Área | Documentos | Linhas estimadas | Status |
|------|------------|------------------|--------|
| Governança de Dados | 7 | ~2.500 | ✅ Concluído |
| Governança Operacional | 3 | ~1.200 | ✅ Concluído |
| Observabilidade | 1 | ~600 | ✅ Concluído |
| Governança de Segredos | 0 | — | ⏳ Adiada |
| Governança de Rede | 0 | — | ⏳ Planejado |
| Matriz de Conformidade | integrada | — | ✅ Integrada |
| Total | 11 | ~4.300 | 11/16 itens (5 pendentes) |

## Links rápidos por público

### DevOps / SRE
- [Health Checks](./operational-governance-health-checks.md) — cronômetros e alertas.
- [Incident Response](./operational-governance-incident-response.md) — runbooks P1 a P4.
- [Backup Verification](./operational-governance-backup-verification.md) — restauração e frescor de cópias.
- [Metrics Contract](./observability-governance-metrics.md) — dashboards, alertas e métricas de sistema.

### Desenvolvedores
- [Schema Registry — SQLite](./schema-registry-sqlite.md) — tabelas locais do gateway, voice, cron, memory e health.
- [Schema Registry — PostgreSQL](./schema-registry-postgres.md) — instâncias e tabelas futuras.
- [Store Selection Matrix](./data-governance-store-selection.md) — como escolher entre SQLite/Postgres/Qdrant/Supabase.
- [Qdrant Collection Contract](./qdrant-collection-contract.md) — contrato de embeddings e payloads.

### Segurança / Compliance
- [Data Lifecycle Policy](./data-governance-lifecycle.md) — retenção hot/warm/archive.
- [Domain Ownership Table](./data-governance-domain-ownership.md) — responsabilidades e SLAs.
- [Compliance Matrix](./data-governance-compliance.md) — alinhamento com CONTRACT.md e GUARDRAILS.md.
- [Incident Response](./operational-governance-incident-response.md) — procedimentos de breach.

### Arquitetos
- [Store Selection Matrix](./data-governance-store-selection.md) — árvore de decisão.
- [Memory Sync Architecture](./memory-sync-architecture.md) — fluxo Qdrant + embeddings.

