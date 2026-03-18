<div align="center">

# Aurelia (Elite Edition)

<img src="assets/aurelia_cover.png" alt="Aurelia cover" width="720" />

**A local-first autonomous coding agent in Go.**

---

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Runtime](https://img.shields.io/badge/Runtime-Local--First-0F172A)](#runtime-model)
[![Architecture](https://img.shields.io/badge/Architecture-Elite_Workspace-1F2937)](docs/ARCHITECTURE.md)
[![Memory](https://img.shields.io/badge/Memory-Layered-0E7490)](#why-aurelia)
[![Interfaces](https://img.shields.io/badge/Interfaces-Multi--Agent-6D28D9)](#🤖-2-orquestração-e-comandos)

</div>

## 🏛️ Multi-Agent Workspace (Elite Edition)
> O template definitivo para **Google Antigravity** + **Claude Code** + **Codex**. 
> Orquestração disciplinada, segurança total e alta performance.

---

## 🚀 1. Diferenciais Competitivos
Este template supera referências mundiais ao integrar o melhor de cada ecossistema:
- **BMAD Orchestration**: Fluxo rigoroso PRD ➡️ Architect ➡️ Dev ➡️ QA.
- **Multi-Runtime**: Pronto para inter-operar com Claude, Codex e Gemini.
- **Global Intelligence**: Subagentes e skills injetados via `ag-init`.
- **Anti-Hallucination**: Regras de autoridade única e relatórios baseados em diff.

## 🤖 2. Orquestração e Comandos
Invoque especialistas diretamente no seu terminal ou IDE:

<commands>
- `/pm` ➡️ Requisitos e Critérios de Aceite.
- `/architect` ➡️ Tech Specs e Design de Sistemas.
- `/dev` ➡️ Implementação e Correções Ágeis.
- `/qa` ➡️ Testes, Validação e Auditoria.
- `/sincronizar-tudo` ➡️ Commit semântico e Push Sênior.
- `/pesquisa-profunda` ➡️ Pesquisa profunda via Gemini Web.
</commands>

## 📁 3. Estratégia de Pastas
- **`docs/`**: Verdade arquitetural, ADRs e benchmarks.
- **`.agents/`**: Governança local (Rules, Workflows).
- **`.context/`**: Memória de trabalho e estado (Gerido via MCP).
- **`internal/` / `pkg/`**: Implementação Core em Go.

---

## 🏗️ Technical Overview (Aurelia Core)

`Aurelia` é um agente de codificação autônomo projetado para rodar localmente com uma pegada operacional mínima e comportamento de runtime explícito.

### Lightweight Baseline
- **Binary size**: `23.22 MB`
- **Idle working set**: `25.66 MB`
- **Startup average**: `15.75 ms`

### Setup
Requisitos: Go `1.25+`, Telegram Token, Provedor de LLM.

```bash
# Onboarding interativo
go run ./cmd/aurelia onboard

# Iniciar agente
go run ./cmd/aurelia
```

A configuração principal reside em `~/.aurelia/config/app.json`.

---
*Este repositório foi construído para ser o #1 do GitHub em orquestração multi-agente.* 🚀
