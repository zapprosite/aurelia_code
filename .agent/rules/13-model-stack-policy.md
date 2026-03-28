---
description: Define os provedores e modelos aprovados para o ecossistema.
id: 13-model-stack-policy
---

# 💠 Regra 13: Política de Stack de Modelos (SOTA 2026.1)

A stack de modelos deve priorizar a soberania local sem comprometer a potência de raciocínio.

## 1. Stack Canônica
- **Raciocínio Estrutural (SOTA)**: Claude 3.7 (Sonnet/Opus) para `@[/architect]` e `@[/pm]`.
- **Soberano Local**: Gemma 3 / Llama 3.3 (via Ollama/vLLM) para automação Tier 0 e segurança.
- **Alta Performance Co-Pilot**: MiniMax 2.7 ou Qwen 2.5 para geração rápida de código.

## 2. Guardrails de Modelos
- É proibido utilizar modelos sem suporte a **Function Calling** estruturado para tarefas críticas.
- O histórico de tokens deve ser otimizado via compressão semântica de contexto.

---
*Assinado: Aurélia (Arquiteta Líder) — Março 2026*
