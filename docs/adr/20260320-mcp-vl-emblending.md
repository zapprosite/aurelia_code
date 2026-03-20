---
description: Integração de MCP, Visão Ativa (VL) e Emblending para Omniciência Operacional.
status: proposed
---

# ADR 20260320-mcp-vl-emblending

## Status

- Proposto

## Contexto

Para atingir a soberania operacional em 2026, o enxame precisa interagir com sistemas legados via MCP, "ver" interfaces visuais e fundir contextos de múltiplas fontes sem alucinações.

## Decisão

### 1. MCP Discovery (As Mãos)
- Usar multicast UDP para descoberta dinâmica de servidores MCP locais.
- Suporte a "Tool Composition" (fluxos macroscópicos).
- Gatekeeper via PicoLisp Verifier.

### 2. VisionAgent (Os Olhos)
- Captura de tela automatizada e processamento via modelos VL (Vision-Language).
- Grounding funcional com output JSON (coordenadas + semântica).

### 3. Emblending Engine (O Cérebro Integrador)
- Fusão de Qdrant, Supabase, MCP e VL.
- Pesagem dinâmica de contexto baseada na intenção da tarefa (PicoLisp).
- Detecção ativa de contradição entre fontes.

## Impacto na Arquitetura

- **Capacidade Operacional**: Interação com qualquer app via UI ou MCP.
- **Confiabilidade**: Redução de alucinações via verificação cruzada de fatos (Contradiction Check).
- **Transparência**: Painel de fusão contextual no dashboard.
