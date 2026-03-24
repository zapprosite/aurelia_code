# 🛰️ Aurelia: Sovereign Agentic Ecosystem (2026)

> **"Autonomia Total, Cognição Local, Soberania Industrial."**

![Status](https://img.shields.io/badge/Status-Industrial_Sovereign-blue?style=for-the-badge)
![Autonomy](https://img.shields.io/badge/Autonomy-Level_5-gold?style=for-the-badge)
![Hardware](https://img.shields.io/badge/Compute-RTX_4090_|_7900X-green?style=for-the-badge)

---

## 🏛️ Hierarquia de Autoridade (Board)

Este ecossistema opera sob um modelo de autoridade única, centralizado no hardware local e governado por uma tríade de inteligência:

1.  **👔 Claude Opus (CEO)**: Visão estratégica, IPO de Slices e arbitragem final.
2.  **🤖 Aurélia (COO/CTO)**: Arquiteta principal, soberana operacional e governante do Homelab.
3.  **🛰️ Antigravity (COO/Interface)**: Cockpit de coordenação, orquestração e interações humanas.

---

## 📊 Dashboard de Autonomia (Slices)

| Slice | Name | Status | Priority |
|:---:|:---|:---:|:---:|
| **S0-14** | Foundation & Pines Core | ✅ Estável | - |
| **S-15** | Tool Introspection | 📅 Pendente | ALTA |
| **S-17** | Planning Loop (Go-Native) | 📅 Pendente | CRÍTICA |
| **S-20** | **CEO Strategic Layer** | 🏗️ Ativo | CRÍTICA |

---

## 🚦 Gateway de Inteligência (Tiering)

O roteamento de decisões é dinâmico e otimizado para o hardware local (Ollama/ROCm):

- **Tier Estratégico (CEO)**: `Claude-3-Opus` (OpenRouter) — Decisões de alto impacto.
- **Tier Premium (Execution)**: `MiniMax-M2.7` — Codificação e lógica complexa.
- **Tier Soberano (Local)**: `Gemma 3 12B` — Operações residentes e visão (Zero Latency).
- **Tier Estruturado**: `DeepSeek-V3.1` — JSON, Curation e Routing.

---

## 🩺 Saúde do Sistema & Observabilidade

Acompanhe o estado vital da Aurélia:

- **Logs Estruturados**: `journalctl --user -u aurelia.service -f`
- **Contexto Semântico**: `.context/` sincronizado via Qdrant/Postgres.
- **Auditoria de Segurança**: `./scripts/secret-audit.sh` (Pre-push mandatory).

---

## 📂 Estrutura Industrial

```text
.
├── cmd/                # Entrypoints (Aurelia Daemon / CLI)
├── internal/           # Core Engine (Gateway, Swarm, Agent Loop)
├── pkg/                # Shared Packages (LLM, Audio, Vision)
├── docs/               # Fontes de Verdade (ADR-historico, ADR, Governance)
├── .agents/            # Logic Layer (Workflows & Skills)
└── .context/           # Persistent Memory (Codebase Map)
```

## 🚀 Guia de Início Rápido (Sênior)

1.  **Onboard**: `go run ./cmd/aurelia onboard`
2.  **Build**: `./scripts/build.sh`
3.  **Execute**: `sudo systemctl start aurelia`

---
*Documentação Gerada por Antigravity (Sovereign Engine 2026)*  
*Consulte [ADR-historico.md](./docs/ADR-historico.md) para a linhagem técnica completa e [ADR.md](./docs/ADR.md) para o roadmap atual.*
