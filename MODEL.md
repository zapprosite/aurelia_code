---
description: Adaptador fino da política de modelos, voz e roteamento do runtime.
---

# 🧠 MODEL.md — Contrato de Modelos

> **IMPORTANTE**: Este arquivo é subordinado a [AGENTS.md](./AGENTS.md). Ele não escolhe governança; ele só resume a política de modelos em vigor.

## Papel

Alinhar qualquer motor, agente ou UI à mesma política de:

- modelos locais
- modelos remotos
- STT/TTS
- budgets e fallback
- decisão final de lane sob a arquitetura da Aurélia

## Regras fechadas

### LLMs (Geração de Texto)

1. `qwen3.5:9b` é o único modelo local residente do caminho ativo.
2. `Groq` fica isolado no lane de áudio/STT.
3. `Gemini TTS / Sulafat` é a voz pronta imediata da Aurelia.
4. `MiniMax Audio` é o lane premium de clonagem autorizada da voz oficial da Aurelia.
5. `OpenRouter` só entra por capacidade explícita.
6. `Gemini web` não entra no runtime automático.

### Embedding Models (Busca Semântica)

7. `bge-m3` (384-dim) é o modelo de embedding para code history + memory sync (Qdrant).
8. Embedding é SEMPRE local (não remoto) — garantir Aurelia trabalhe sem web.
9. Toda mudança de política de embedding, LLM ou voz exige ADR.

### Restrições Gerais

10. Nenhum motor externo pode contornar a política decidida pela Aurélia para o runtime.

## Fontes canônicas

- [AGENTS.md](./AGENTS.md)
- [REPOSITORY_CONTRACT.md](./docs/REPOSITORY_CONTRACT.md)
- [ADR Index](./docs/adr/README.md)
- [Model Routing Matrix](./docs/model_routing_matrix_20260319.md)
- [Local Model Kit Blueprint](./docs/local_model_kit_blueprint_20260319.md)
- [LLM Gateway Blueprint](./docs/llm_gateway_blueprint_20260319.md)
- [Jarvis Local Voice Blueprint](./docs/jarvis_local_voice_blueprint_20260319.md)
- [Voz Oficial da Aurelia](./docs/aurelia_voice_profile_20260319.md)
- [ADR Gemini TTS Ready Voice](./docs/adr/20260319-gemini-tts-ready-voice.md)

## Operação

- não abrir novo lane de modelo sem ADR
- não mudar default local sem medir VRAM e rodar smoke
- não introduzir segundo residente local sem justificativa formal
- não mudar voz/STT/TTS sem explicitar custo, fallback e impacto no homelab
