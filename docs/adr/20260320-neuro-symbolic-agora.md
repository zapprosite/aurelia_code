---
description: Evolução para Memória Neuro-Simbólica (NSM) e Interface Generativa (GenUI) - Padrão 2026.
status: proposed
---

# ADR 20260320-neuro-symbolic-agora

## Status

- Proposto

## Contexto

Em Março de 2026, a arquitetura de agentes exige mais do que busca vetorial. Precisamos de **raciocínio causal** (quem/porquê) antes da recuperação de detalhes (o quê). O uso de **PicoLisp** para grafos simbólicos em RAM é a nossa vantagem competitiva.

## Decisão

### 1. Neuro-Symbolic Memory (NSM)
- **Camada Simbólica (PicoLisp)**: Um grafo de conhecimento em tempo real que mapeia relações causais (`agente-A -> ajudou -> agente-B`).
- **Camada Neural (Qdrant/Supabase)**: Vetores para preencher os detalhes semânticos após a filtragem simbólica.
- **Roteamento**: PicoLisp define *onde* buscar; Vetores trazem o *conteúdo*.

### 2. Protocolo A2A (Agent-to-Agent)
- Implementar um adaptador para o padrão A2A 2026, permitindo que agentes externos se conectem ao enxame como "freelancers".

### 3. Generative UI (GenUI)
- O backend Go enviará `dynamicLayoutSchema` (JSON) via WebSocket.
- O frontend React renderizará componentes sob demanda (ReactFlow + D3 + Componentes Dinâmicos).

### 4. Reflective Loop
- Injeção de uma etapa de `SELF_REFLECTION` no cérebro do agente antes de qualquer ação externa.

## Consequências

- **Positivas**: Raciocínio explicável, maior autonomia, interface adaptativa.
- **Negativas**: Aumento na complexidade do bridge Go-PicoLisp.

## Testes

- Validar se a busca simbólica no PicoLisp reduz o ruído na busca vetorial no Qdrant.
- Verificar renderização dinâmica de componentes no dashboard via schema JSON.
