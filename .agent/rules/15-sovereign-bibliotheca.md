# Governança Sovereign-Bibliotheca SOTA 2026

> **Status**: Consolidado (Março 2026)
> **Objetivo**: Manter ativos de skills Node/OpenClaw orquestrados pelo motor Soberano Go (Aurélia).

## 1. Princípio da Soberania Go (Go-First)

- **Core (Go)**: Única autoridade de estado, memória (SQLite/Qdrant/Supabase) e orquestração de ferramentas.
- **Skills (Node/OpenClaw)**: Mantidas como ativos especializados em `homelab-bibliotheca/skills/`.
- **Interoperabilidade**: A Aurélia (Go) invoca skills via `markdownbrain` ou execuções CLI diretas, eliminando a dependência de orquestradores Bash externos.

## 2. Higiene de Diretórios

- **Proibido** carregar índices JSON gigantes (`skills-registry.json`). A descoberta de skills deve ser feita pelo `markdownbrain` do core Go.
- **Proibido** duplicar lógica de sincronização em Bash. O `app.go` é o responsável térmico pela integridade dos dados.

---
*Assinado: Aurélia (Arquiteta Líder) & Antigravity (Operador SOTA 2026)*
