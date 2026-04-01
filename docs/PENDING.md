# PENDING — Slices Aguardando Implementação
Última auditoria: 01/04/2026

> POLÍTICA: Apenas slices com código no repo são marcados ✅.
> Este arquivo é a única fonte de verdade do backlog.

## P0 — Hotfix Ativo 🔥 (impede o bot de funcionar)

Nenhum. Sistema SOTA 2026 Sovereign operando de forma estável.

## P1 — Crítico 🔴 (infra sem a qual o home lab não é estável)

Nenhum limitador conhecido no momento.

## P2 — Alto 🟡 (qualidade da experiência)

| Slice | Descrição | Pré-requisito |
|---|---|---|
| S-62 | Health checks passivos para serviços do homelab | — |

## P3 — Médio 🟢 (expansão da capacidade nativa)

| Slice | Descrição | Pré-requisito |
|---|---|---|
| S-63 | Computer Use E2E (Automação Visual) | S-57 (Feature Flag) |
| S-64 | OS Native God Mode | S-63 |
| S-65 | Jarvis Voice + Computer Use em conjunção | S-64 |

---

## ✅ Slices Concluídos Recentes (Sovereign 2026)
- **S-51 a S-55**: Estabilização do Homelab, TTS via Edge, Groq STT, remoção do Supabase.
- **S-56**: Extração do `AURELIA.md` injetado no pipeline de input.
- **S-57**: Setup de Feature Flags.
- **S-58**: Construção do `DreamConsolidator` (SQLite -> LLM -> Qdrant).
- **S-59**: Refatoração do `input_pipeline.go` limitando escopo das funções (< 100 linhas).
- **S-60**: `SharedMemory` e infraestrutura de Redis PubSub para coordenação de Swarm Local.
- **S-61**: Limpeza e sincronização do PENDING.md (Este estado atual).
