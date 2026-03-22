> [!NOTE]
> Status: ✅ Arquivado / Concluído em 22/03/2026

---
title: Dashboard de Agentes — Workspace ULTRATRINK
status: active
date: 2026-03-20
decision-makers: [humano, aurelia]
tags: [ui, agents, ultratrink, dashboard]
---

# ADR-20260320: Dashboard de Agentes (ULTRATRINK Workspace)

## 📌 Contexto

À medida que o Aurelia evolui de um único agente para um enxame (swarm) autônomo (Antigravity, Claude, Codex, OpenCode, Aurelia Core), a dependência exclusiva do terminal (`stderr`/`stdout`) e do Telegram torna-se insuficiente para a observabilidade operacional e a governança visual do sistema.

Precisamos de um **QG visual sênior**, nomeado **ULTRATRINK**. O conceito não é apenas uma "tela de logs da infra", mas sim um **"Notion/Escritório Virtual"** dinâmico, onde os agentes colaboram ativamente, trocam contexto, e o desenvolvedor humano tem visão de raio-X do enxame.

## 🎯 Conceito "ULTRATRINK"

A experiência deve afastar-se do amadorismo visual. Princípios chaves:
1. **Premium & Dark Mode:** Interface "glassmorphism" ou tech-noir. Muito fluida, tipografia premium, estado da arte em dashboarding. Micro-interações.
2. **Equipe Unida:** Os agentes devem ser representados como membros de uma task force operando em conjunto. Quando o Claude empaca num erro local, a view deve mostrar a "conversa/handoff" com o *Antigravity* no terminal. É uma orquestração viva.
3. **Timeline e Task-Boards:** Sincronizado dinamicamente com os planos em `.context/plans/` e as `task.md`. Semelhante ao kanban do Notion, porém vivo e real-time.

## 🛠️ Decisões Técnicas Principais

1. **Porta Oficial:** O ULTRATRINK responderá exclusivamente em `http://localhost:3333/`.
2. **Runtime Local-First:** Nenhuma dependência de cloud dashboards (como Grafana Cloud ou Datadog). O estado será consumido das fontes locais reais:
   - Os arquivos físicos em `.context/` e `.agents/`
   - O banco `Postgres` e `Qdrant` do memory-sync local
   - A subscrição via WebSocket (ou SSE) ao daemon Go do `aurelia`.
3. **Frontend Stack:** 
   - SPA moderna (Next.js export estático, Vite + React, ou Svelte — preferência por performance pura sem SSR pesado, rodando 100% via client-side consumindo a API Go).
   - Componentes próprios ou Tailwind altamente customizado para visual corporativo sênior.
4. **Endpoint Bridge no Go:** 
   - A porta `3333` será servida pelo próprio daemon do `Aurelia` (adicionando um listener HTTP gin/fiber/stdlib em `cmd/aurelia/http.go`). O SPA ficará embutido via `go:embed` na build estática final.

## 📊 Estrutura Desejada do Dashboard

- **Main Floor (Activity Stream):** O que está acontecendo AGORA. Timeline central de tool calls, bash commands e handoffs. Estilo "feed de rede social dos robôs".
- **The Brain (.context):** Árvore de conhecimento e ADRs visualizáveis, busca embutida (Qdrant).
- **Squad Room:** Status on-line de cada profile (Gemini/Antigravity idle, Claude executando, Aurelia monitorando `systemd`). Burn-rate e VRAM consumidos por agente.
- **Mission Board:** Parse visual do `roadmap-mestre-slices.md` e pastas `.context/plans/`. 

## ⚖️ Consequências

- **A favor:** Transformará a experiência de Pair-Programming Multi-Agente em algo imersivo. Elimina a fadiga de ler 2000 linhas de console para entender o estado de um plano.
- **Custos:** Adiciona complexidade transversal (backend Go precisa cuspir SSE/WebSockets com eventos estruturados de agent-bus; Frontend precisa ser mantido com alto rigor estético).
- **Integração no Roadmap:** Entra como Slice P9 no Roadmap Mestre (`ADR-20260320-roadmap-mestre-slices.md`). O desenvolvimento frontend deve acontecer em worktree isolada (`feat/ultratrink-dashboard`).
