# ADR-2026-ROADMAP-FUTURO: Estratégia de Evolução e Slices Pendentes 🛰️

**Status:** 🔄 Ativo (Em Evolução)
**Autoridade:** Aurélia (Arquiteta Principal)
**Data de Revisão:** 2026-03-22
**Nível de Governança:** Industrial Homelab

---

## 1. Resumo Executivo
Este roadmap define a trajetória da subida da Aurélia do estado de assistente para "Jarvis Local-First". O foco é a **Autonomia Total (Level 5)**, eliminando latências externas e garantindo que a cognição ocorra prioritariamente no hardware local (RTX 4090 + 7900X), utilizando a nuvem apenas como fallback de alta inteligência ou ferramenta estruturada.

## 2. Status dos Slices (Backlog Oficial)

| ID | Slice Name | Descrição | Status | Prioridade |
|:---|:---|:---|:---:|:---:|
| **S-15** | Tool Introspection | Filtro semântico dinâmico de ferramentas via Qdrant. | 📅 Pendente | ALTA |
| **S-16** | Execution DNA | Templates de workflow nativos (Go) por tipo de tarefa. | 📅 Pendente | MÉDIA |
| **S-17** | Planning Loop | Loop PREV (Plan-Review-Exec-Verify) nativo no daemon. | 📅 Pendente | CRÍTICA |
| **S-18** | Codebase Symbol Map | Parseamento AST (.ast) para localização de símbolos. | 📅 Pendente | ALTA |
| **S-19** | Semantic Skill Router| Roteamento de skills via embeddings (Vector-First). | 📅 Pendente | ALTA |

## 3. Detalhamento Técnico dos Slices

### S-17: Planning Loop (PREV Phase) — CRITICAL PATH
O objetivo é mover a lógica de orquestração do Antigravity/Claude para dentro do binário `aurelia` em Go.
- **Garantia**: Bloqueio de segurança (`gate`) para planos que afetam infraestrutura (sudo) sem intervenção humana.
- **Componentes**: `internal/orchestrator/`, `internal/gatekeeper/`.

### S-18: Codebase Symbol Map
Substituição do `grep` puro por uma busca de símbolos consciente do tipo.
- **Garantia**: Redução de alucinações em renames e refatorações complexas.
- **Tecnologia**: Go Parser + Qdrant Metadata Sync.

## 4. Visão de Longo Prazo (S20+)
- **Autonomous HW Management**: Auto-escalonamento de VRAM e gestão de energia via `nvidia-smi` monitorado pelo `homelab-control`.
- **Global Auth Proxy**: Unificação de sessões entre Dashboard (Next.js), Telegram e CLI.
- **Cognitive Self-Healing**: Capacidade da Aurélia de detectar falhas em logs de containers e aplicar correções via PR interno.

## 5. Critérios de Sucesso Industrial
- **Autonomia**: Intervenção humana < 5% em fluxos de desenvolvimento.
- **Performance**: Latência de "Pensamento" < 500ms usando modelos locais otimizados (Gemma 3).
- **Higiene**: Zero discrepância entre o estado do código e o `.context/`.

---
