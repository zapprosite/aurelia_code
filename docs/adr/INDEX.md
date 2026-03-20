---
description: Índice Soberano de Registros de Decisão Arquitetural (ADR) e Mapa de Slices.
---

# 🏛️ Architectural Decision Records (ADR) - INDEX

Este documento é a **fonte única de verdade** para decisões técnicas e o roadmap operacional de fatias (slices) do projeto Aurelia.

## 📜 Contratos de Base
- [AGENTS.md](../../AGENTS.md) — Autoridade e Governança Multi-Agente.
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md) — Governança Industrial.

---

## 🌊 Mapa de Execução (Roadmap de Slices)

| Onda | Slice | Prioridade | Status | ADR |
| :---: | :--- | :---: | :---: | :--- |
| **1** | Voz e Experiência (E2E) | 🔴 Crítico | 🔵 Em execução | [ADR-20260319-voice-plane-unified](./ADR-20260319-voice-plane-unified.md) |
| **2** | Orquestração (Handoff) | 🔴 Crítico | 🔵 Em execução | [ADR-20260319-antigravity-handoff-e2e](./ADR-20260319-antigravity-handoff-e2e.md) |
| **2** | Login Seguro Browser | 🔴 Crítico | 🔵 Em execução | [ADR-20260319-browser-safe-login.md](./ADR-20260319-browser-safe-login.md) |
| **3** | Swarm Hierárquico (Bus) | 🟡 Médio | 🟡 Proposto | [ADR-20260319-hierarchical-agent-bus.md](./ADR-20260319-hierarchical-agent-bus.md) |
| **4** | Desktop Fallback Seguro | 🟡 Médio | 🔵 Em execução | [ADR-20260319-desktop-safe-fallback.md](./ADR-20260319-desktop-safe-fallback.md) |

---

## 📂 Índice Completo de ADRs

| Data | Título | Status |
| :--- | :--- | :--- |
| 2026-03-20 | [Auditoria, Limpeza e Restauração](./ADR-20260320-repository-clean-audit-restore.md) | 🔵 Em execução |
| 2026-03-20 | [Restauração de Visão no Telegram](./ADR-20260320-restore-telegram-vision.md) | 🟢 Concluído |
| 2026-03-19 | **[Voice Plane Unificado (STT/TTS/Capture)](./ADR-20260319-voice-plane-unified.md)** | ✅ Aceito |
| 2026-03-19 | [Governança Industrial (Secrets/Data/Net)](./ADR-20260319-Polish-Governance-All.md) | 🟡 Proposto |
| 2026-03-19 | [Higiene Documental da Raiz](./ADR-20260319-root-document-hygiene.md) | ✅ Aceito |
| 2026-03-19 | [Aurélia como Autoridade Arquitetural](./ADR-20260319-aurelia-authority-model.md) | ✅ Aceito |
| 2026-03-19 | [Governança de Extensões](./ADR-20260319-extensions-governance.md) | ✅ Aceito |
| 2026-03-19 | [Manual Offline Homelab (Qdrant)](./ADR-20260319-offline-homelab-manual-qdrant.md) | 🟡 Proposto |
| 2026-03-18 | [Implementando ai-context](./ADR-20260318-implementando-ai-context.md) | ✅ Aceito |
| 2026-03-18 | [Integração MCP Antigravity](./ADR-20260318-integracao-mcp-antigravity.md) | ✅ Aceito |
| 2026-03-18 | [Estratégia RIM](./ADR-20260318-estrategia-rim.md) | ✅ Aceito |
| 2026-03-17 | [Rebalanceamento de Contexto Elite](./ADR-20260317-rebalanceamento-elite.md) | ✅ Aceito |

---

## 🛠️ Manutenção
- Para novas slices: `./scripts/adr-slice-init.sh <slug>`
- Padrão de nome: `ADR-YYYYMMDD-slug.md`
- Localização: `docs/adr/`
