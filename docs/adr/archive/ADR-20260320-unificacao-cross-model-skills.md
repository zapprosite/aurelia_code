> [!NOTE]
> Status: ✅ Arquivado / Concluído em 22/03/2026

---
title: Unificação Cross-Model de Skills
status: active
date: 2026-03-20
decision-makers: [humano, aurelia, antigravity]
---

# ADR: Unificação Cross-Model de Skills

## Contexto
O ecossistema multi-agentes possuía bibliotecas isoladas de habilidades (`skills/`) fragmentadas entre provedores (`.claude/skills`, `.context/skills`, `.opencode/skills`). Isso criava silos onde um LLM operando via Claude não enxergava as ferramentas exclusivas escritas no `.context` ou `.agents`.

## Decisão
1. **Centralização Absoluta**: A pasta `.agents/skills/` passa a ser o único source-of-truth (SOT) físico para todas as skills de todos os modelos.
2. **Symlinks**: Todas as outras pastas relacionadas a `skills` nos executores (`.claude/`, `.opencode/`, `.codex/`, `.context/`) foram deletadas e substituídas por [Symlinks] para `../.agents/skills`.
3. **Agosticismo**: Prompts de skills devem ser model-agnostic, focando em bash, manipulação de arquivos ou ferramentas MCP universais.

## Consequências
- Fim da duplicação de instruções e ferramentas diverentes.
- Um modelo que for rodar no Cline/Roo Code via Claude terá exatamente os mesmos poderes que a Aurélia invocando o Antigravity via Gemini ou o Codex via terminal.
