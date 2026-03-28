---
description: Define os níveis de permissão e modelos baseados em criticidade.
id: 03-tiers-autonomy
---

# 🛡️ Regra 03: Tiers de Autonomia (SOTA 2026.1)

A autonomia de ação é escalonada pelo risco e capacidade de raciocínio do modelo.

<directives>
1. **Tier 0 (Soberano/Local)**: Gestão de arquivos, builds e auditoria de segredos. 
   - **Model**: Gemma 3 (Ollama/TensorRT).
2. **Tier 1 (Design/Tech Spec)**: Decisões arquiteturais e PRDs.
   - **Model**: Claude 3.7 (Architect) / Gemini 2.0 (Deep Research).
3. **Tier 2 (Rápida Execução)**: Implementação de código desacoplado e tRPC.
   - **Model**: MiniMax 2.7 / Qwen 2.5 (Fast Reasoning).
</directives>
