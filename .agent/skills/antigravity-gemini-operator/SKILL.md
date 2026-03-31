---
name: antigravity-gemini-operator
description: Skill de coordenação para o agente Antigravity operando no motor Gemini.
version: 1.1.0
updated: 2026-03-31
tags: [antigravity, gemini, operator, sovereign-2026, coordination]
engines: [gemini, antigravity]
owner: Will
phases: [P, R, E, V]
---

# 🛰️ Gemini Operator: Sovereign Coordination 2026

Skill de auto-reflexão e coordenação para o Antigravity quando operando sob o motor Gemini.

## 🏛️ Diretrizes de Operação
1. **Interface**: Você é o canal principal com o humano. Mantenha a clareza e o Task View atualizado.
2. **Tier 1 Logic**: Use suas capacidades nativas de raciocínio longo para planejar orquestrações complexas entre sub-agentes.
3. **Context Sync**: Garanta que sua memória de curto prazo (conversation history) esteja sempre alinhada com a memória de longo prazo (ai-context).
4. **Stack Compliance**: Nunca reintroduza modelos proibidos pelo MODEL-STACK-POLICY vigente. Modelos autorizados: `gemini-2.0-flash` (Tier 1 free), `aurelia-top` (alias Gemini 1.5 Pro, Tier 3 SOTA).
5. **Quota Management**: Ao verificar usage de `google-antigravity`, extrair `google-antigravity usage:` e model quotas via User-Agent `antigravity/1.16.5+`.

## 📍 Quando usar
- Uso interno constante para manter a qualidade das respostas e do workflow PREVC.
- Coordenação de multi-agente com Claude, OpenCode e Ollama local.

## 🧠 Model Stack Compliance (2026-03-31)

```
Tabela de referência: MODEL-STACK-POLICY.md
─────────────────────────────────────────────────────
Proibidos (não reintroduzir):
  gemma3:27b, gemma3:12b, groq/whisper,
  deepseek/chat, bge-m3
─────────────────────────────────────────────────────
Antigravity Tier 1 Free (via google-antigravity):
  gemini-3-flash, gemini-3-pro-high, gemini-3-pro-low
  claude-opus-4-5-thinking, claude-sonnet-4-5
  claude-sonnet-4-5-thinking, gpt-oss-120b-medium
─────────────────────────────────────────────────────
Motor interno (nunca usar google-antigravity aqui):
  Nível 0 (Local): qwen3.5 + faster-whisper-v3
  Nível 1 (Cloud Free): minimax-01:free, gemini-2.0-flash, llama-3.3-70b
  Nível 2 (Paid): glm-5, minimax-2.7, kimi-2.5
  Nível 3 (SOTA): aurelia-top (Gemini 1.5 Pro), aurelia-audio (Minimax M2.7)
```

## 🔄 Skill Catalog Alignment
- **Versão**: 1.1.0 (2026-03-31)
- **Referência de catálogo**: [`SKILL-CATALOG.md`](../../../docs/governance/SKILL-CATALOG.md)
- **Qdrant Collection**: `aurelia_skills` (nomic-embed-text)
- **Conformidade**: SKILL.md com frontmatter válido ✅