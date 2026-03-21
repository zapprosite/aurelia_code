# Índice de ADRs — Aurelia

| Data | ADR | Status |
|------|-----|--------|
| 2026-03-20 | [Dashboard Agentes ULTRATRINK](ADR-20260320-dashboard-agentes-ultratrink.md) | ✅ Concluído |
| 2026-03-20 | [Plano Mestre JARVIS Local-First](ADR-20260320-plano-mestre-jarvis-local-first.md) | 🔄 Ativo |
| 2026-03-20 | [Política de Modelos e Hardware/VRAM](ADR-20260320-politica-modelos-hardware-vram.md) | 🔄 Ativo |
| 2026-03-20 | [Roadmap Mestre de Slices](ADR-20260320-roadmap-mestre-slices.md) | 🔄 Ativo |
| 2026-03-20 | [Unificação Cross-Model Skills](ADR-20260320-unificacao-cross-model-skills.md) | ✅ Concluído |
| 2026-03-21 | [Agent-to-Agent Go Native](ADR-20260321-agent-to-agent-go-native.md) | ✅ Concluído |
| 2026-03-21 | [Aurelia Autonomous Engineering — ULTRATRINK](ADR-20260321-aurelia-autonomous-engineering.md) | 🔄 Em progresso |

## Taskmasters (JSON de continuidade)

| ADR | Arquivo |
|-----|---------|
| Dashboard Agentes ULTRATRINK | [JSON](taskmaster/ADR-20260320-dashboard-agentes-ultratrink.json) |
| Restore Telegram Vision | [JSON](taskmaster/ADR-20260320-restore-telegram-vision.json) |
| Agent-to-Agent Go Native | [JSON](taskmaster/ADR-20260321-agent-to-agent-go-native.json) |

## Convenções

- **Formato de nome:** `ADR-YYYYMMDD-slug.md`
- **Taskmaster JSON:** `taskmaster/ADR-YYYYMMDD-slug.json`
- **Scaffold:** `./scripts/adr-slice-init.sh <slug> --title "Title"`
- **Escopo obrigatório:** arquitetura, providers, storage, runtime, áudio, deploy, segurança
