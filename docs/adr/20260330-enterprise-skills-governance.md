# ADR-20260330: Enterprise Skills Governance — Sovereign 2026 Q2

**Status:** 🟡 Em Execução (Nonstop Slice)
**Data:** 2026-03-30
**Autor:** Antigravity AI (Will Aprovado)
**Tipo:** Cross-Cutting / Tooling / Governança

---

## 🟢 Contexto

Com a consolidação do modelo "Sovereign Enterprise 2026", o repositório `aurelia` passou a operar como um sistema multi-agente de produção. Identificamos 3 lacunas críticas no catálogo de skills que impedem a formalização do padrão de qualidade:

1. **Organização de Repositório**: Não há uma skill que force o padrão de scaffolding enterprise (estrutura de pastas, linting, ARCHITECTURE.md, .editorconfig).
2. **Polimento de Código**: O catálogo carece de uma skill automatizada para detecção de segredos hardcoded e CVEs em containers.
3. **Governança Agêntica**: Faltam regras formalizadas para codificação assistida por IA (AGENTS.md, .cursorrules, guardrails de output).

## 🔵 Decisão

Adotar as seguintes skills como padrão oficial do catálogo (`.agent/skills/`):

| Slice | Skill | Prioridade |
|-------|-------|-----------|
| 1 | `system-architect-enterprise` | 🔴 Alta — estrutura base |
| 2 | `security-guardian-enterprise` | 🔴 Alta — CI/CD gate |
| 3 | `ai-coding-toolkit` | 🟡 Média — guardrails e regras |

## 🔴 Consequências

- O repositório passa a ter um **Arquiteto de Sistema** ativo que gera scaffolding e padrões de código.
- A **Security Guardian** torna-se gate obrigatório antes de qualquer `git push`.
- O **AI Coding Toolkit** formaliza as regras de codificação assistida para toda a equipe (humano + agente).

---

## Slices de Execução

- [x] Slice 0: Instalação das 3 skills no catálogo canônico.
- [ ] Slice 1: `system-architect-enterprise` — Gerar ARCHITECTURE.md e .editorconfig.
- [ ] Slice 2: `security-guardian-enterprise` — Executar scan de segredos e containers.
- [ ] Slice 3: `ai-coding-toolkit` — Criar/atualizar AGENTS.md enterprise com guardrails.
