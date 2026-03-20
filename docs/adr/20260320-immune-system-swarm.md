---
description: Evolução final para Arquitetura de "Sistema Imunológico" (Fluid + Holographic + Formal + AR).
status: proposed
---

# ADR 20260320-immune-system-swarm

## Status

- Proposto

## Contexto

Migração para um organismo vivo capaz de autorregulação e defesa ativa. O enxame não é mais apenas coordenado, mas **fluidodinâmico**, com memória distribuída e verificação formal constante.

## Decisão

### 1. Fluid Workflow (Pressure Router)
- Implementar balanceamento de carga baseado em pressão em tempo real.
- Uso de algoritmos de fluxo de rede para desviar tarefas de agentes saturados.

### 2. Memória Holográfica (gRPC P2P)
- Agentes trocam fragmentos de memória comprimidos diretamente entre si.
- Supabase/Qdrant passam a ser `Backbone de Auditoria` (escrita pesada, leitura leve).

### 3. Smart Collaboration Contracts
- Sistema de créditos e reputação efêmera para gerenciar ajuda mútua (RFH + Bidding).
- Verificação de qualidade via PicoLisp.

### 4. Zero-Trust Formal Verifier
- DSL de leis corporativas em PicoLisp (`laws.l`).
- Bloqueio preventivo de ações não-conformes.

### 5. Frontend AR Cognitiva
- Visualização de heatmaps de "Atenção" e balões de pensamento AR sobrepostos aos dados.

## Impacto na Arquitetura

- **Redução de Latência**: ~70% via P2P.
- **Resiliência**: Sobrevivência com 80% de eficiência mesmo com perda parcial de agentes.
- **Segurança**: Prevenção ativa de alucinações e violações.
