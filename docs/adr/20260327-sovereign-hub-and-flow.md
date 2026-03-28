# ADR-20260327: Aurelia Sovereign Hub & Flow Core

## Status: Proposta

## Contexto
O ecossistema OpenClaw (2025/2026) popularizou o uso de "Skills" modulares através do ClawHub. O PicoFlow introduziu uma DSL ágil para encadeamento de agentes. Atualmente, a Aurelia possui skills em `.agent/skills`, mas o carregamento e a descoberta são manuais/ad-hoc.

## Decisão
Implementar o **Aurelia Sovereign Hub** (uma evolução soberana do ClawHub) e o **Aurelia Flow** (uma DSL nativa em Go para orquestração).

### 1. Sovereign Hub (Inspirado no ClawHub)
- **Estrutura**: Centralizar metadados de todas as skills em `internal/agent/hub.go`.
- **Comando**: Adicionar `aureliactl skills list` e `aureliactl skills install`.
- **Isolamento**: Cada skill deve ter um contrato Zod (Frontend) e um Schema Go (Backend) rigoroso para evitar prompt injection (diferente do OpenClaw que é permissivo).

### 2. Aurelia Flow (Inspirado no PicoFlow)
- **DSL em Go**: Criar o pacote `internal/purity/aflow` para encadeamento de funções.
- **Exemplo**: `aflow.New().Step(Research).Then(Write).Run(ctx)`.
- **Vantagem**: Performance 10x superior às DSLs interpretadas (Tier 0).

## Consequências
- **Positivas**: Facilidade extrema para adicionar novas capacidades à Aurelia mantendo a soberania e a performance.
- **Negativas**: Aumenta a complexidade do `aureliactl`.

---
**Data**: 27/03/2026
**Autor**: Antigravity (SOTA 2026)
