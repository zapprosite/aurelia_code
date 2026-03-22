# ADR-2026-HISTORICO-S0-S14: Consolidação Arquitetural Aurelia

**Status:** ✅ Concluído (Snapshot Março 2026)
**Autoridade:** Aurélia (Arquiteta Principal) / Antigravity (Coordenação)

## 1. Contexto e Evolução
Este documento consolida todas as Decisões Arquiteturais (ADRs) tomadas desde a concepção do projeto até o estágio atual de Autonomia Nível 5 (Pines Core). O objetivo é fornecer uma visão única e sênior da base tecnológica.

## 2. Decisões Consolidadas por Era

### Era 0-10: Fundação e Resiliência
- **Monorepo & Go-Native**: Adoção de Go como linguagem principal para performance e tipagem forte em sistemas distribuídos.
- **Telegram como Interface**: Uso do Telegram para controle remoto seguro e interativo.
- **Local-First AI (Ollama)**: Governança de hardware para rodar LLMs (Gemma, Qwen) localmente, garantindo privacidade e custo zero.

### Era 11: Observabilidade e Controle (Cockpit)
- **Dashboard React (ULTRATRINK)**: Criação de uma interface web para visualização em tempo real de logs, ferramentas e planos de ação via SSE.
- **Arquitetura de Eventos**: Implementação de um `EventHub` centralizado para desacoplar a lógica de execução da UI.

### Era 12: Restauração da Visão (Multimodalidade)
- **Gemma 3 Integration**: Restauração da capacidade de visão e OCR para processamento de fotos e documentos no Telegram.
- **Unified Vision Interface**: Criação de handlers genéricos para processamento de mídia.

### Era 13: Governança Multi-Agente
- **AGENTS.md & REPOSITORY_CONTRACT.md**: Definição clara de papéis e autoridades (Aurélia > Humanos > Adaptadores).
- **Zod-First Schema**: Centralização de contratos de dados em `packages/zod-schemas/`.

### Era 14: Pines Core (Autonomia Nível 5)
- **Voz Canônica**: Adoção da voz `aurelia-ptbr-formal-doce-v1` como identidade soberana.
- **Autonomous Engineering**: Introdução de loops de reflexão e verificações de segurança pré-execução.

## 4. Checklist de Fases Concluídas (Proof of Work)

- [x] **Phase 0: S0 (Infra)** — Go Monorepo, Docker, CapRover & CI/CD.
- [x] **Phase 1: S1 (Bot)** — Telegram Core, handlers básicos e roteamento.
- [x] **Phase 2: S2 (Voice)** — Integração STT/TTS (Gemini/OpenAI/Groq).
- [x] **Phase 3: S3 (Memory)** — Qdrant Vector DB & Persistent Context.
- [x] **Phase 4: S4 (Governance)** — ADR Process, AGENTS.md, Rules.
- [x] **Phase 5: S5 (Skills)** — Extensibilidade via scripts e MCP servers.
- [x] **Phase 6: S6 (Vision)** — OCR e Image Analysis (Gemma 3).
- [x] **Phase 7: S7 (Observability)** — Slog structured logging & tRPC tracers.
- [x] **Phase 8: S8 (Homelab)** — `homelab-control` (Ollama/Docker control).
- [x] **Phase 9: S9 (Security)** — Secret auditing, Pre-push hooks, Sudo=1.
- [x] **Phase 10: S10 (Cockpit)** — Dashboard React (Ultratrink), SSE Hub.
- [x] **Phase 11: S11 (Symbol Map)** — AST parsing p/ navegação inteligente.
- [x] **Phase 12: S12 (Multi-Agent)** — Coordenação Claude/Codex/Antigravity.
- [x] **Phase 13: S13 (Autonomous)** — Planning Loop & Human-in-the-loop gate.
- [x] **Phase 14: S14 (Pines Core)** — Voz Canônica & Autonomia Nível 5.

## 5. Consequências
- **Redução de Ruído**: Menos arquivos ADR fragmentados.
- **Clareza Executiva**: Um único ponto de entrada para entender "O Porquê" de cada componente.
- **Persistência**: Memória técnica preservada via Qdrant e documentação consolidada.
