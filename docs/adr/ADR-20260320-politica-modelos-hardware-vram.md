---
title: Política de Modelos, Hardware e Budget de VRAM
status: accepted
date: 2026-03-20
decision-makers: [humano, aurelia]
supersedes: MODEL.md (raiz)
---

# ADR-20260320: Política de Modelos, Hardware e Budget de VRAM

## Contexto

O projeto operava com um `MODEL.md` na raiz que definia a política de modelos de forma ad-hoc. Com a padronização do repositório como template multi-agente, a política de modelos precisa ser registrada como ADR para manter rastreabilidade e permitir evolução governada.

## Decisão

### Hardware

- CPU: `AMD Ryzen 9 7900X` (Zen 4, iGPU RDNA 2 integrada)
- GPU: `AMD RX 7900 XTX` (24 GiB VRAM, ROCm)
- RAM: `32 GiB DDR5 Gen5`
- Storage: `4 TB NVMe Gen5 (ZFS) + 1 TB NVMe Gen3 (Ubuntu Desktop)`
- Mobo: `ASUS X670E-Plus`
- Regra operacional: desktop na iGPU (cabo no onboard), GPU 100% para ML

### Stack de Modelos Locais

| Função | Modelo | VRAM | Notas |
|--------|--------|------|-------|
| Residente principal (agêntico) | `gemma3:12b` | ~8.1 GiB | Melhor p/ agent instructions (ref. comunidade) |
| Fallback contexto longo | `qwen3.5:9b` | ~6.6 GiB | 262K ctx, vision, tools, thinking |
| Escalonamento pesado | `gemma3:27b-it-q4_K_M` | ~17 GiB | Sob demanda, contexto curto |
| Embedding | `bge-m3` (384-dim) | ~1.2 GiB | Sempre carregado, busca semântica Qdrant |
| STT | `Groq` (remoto) | 0 GiB | Grátis, PT-BR excelente |
| TTS | `Kokoro` (local) | 0 GiB | CPU-only |

### Budget de VRAM (contabilizando KV cache + overhead ROCm)

- `gemma3:12b` + bge-m3 + KV + overhead = ~13 GiB → **~11 GiB de folga** ✅
- `gemma3:27b` + bge-m3 + KV + overhead = ~22 GiB → **~2 GiB de folga** ⚠️ apenas contexto curto

### Regras Fechadas

1. `gemma3:12b` é o residente principal para orquestração do enxame agêntico.
2. `qwen3.5:9b` é o fallback quando a task exige contexto longo (262K vs 128K).
3. `gemma3:27b-it-q4_K_M` é o escalonamento sob demanda para raciocínio profundo.
4. `Groq` fica isolado no lane de áudio/STT.
5. `Kokoro` é o TTS local (CPU-only, zero VRAM).
6. `Gemini TTS / Sulafat` é a voz pronta imediata da Aurelia (fallback remoto).
7. `MiniMax Audio` é o lane premium de clonagem autorizada da voz oficial da Aurelia.
8. `OpenRouter` só entra por capacidade explícita.
9. `Gemini web` não entra no runtime automático.
10. `bge-m3` (384-dim) é o embedding único para code history + memory sync (Qdrant).
11. Embedding é SEMPRE local (não remoto).
12. Toda mudança de política de embedding, LLM ou voz exige ADR.
13. Modelos de **code** não são residentes — foco é **instruction** para chatbot agêntico em enxame.
14. Nenhum motor externo pode contornar a política decidida pela Aurélia para o runtime.

### Limites e Degradação

- LLM pesado concorrente: `1`
- Fila máxima do LLM pesado: `1`
- Embeddings concorrentes: `1`
- Browser-use ativo em paralelo: `1`
- Degradar quando:
  - CPU média 15m > `70%`
  - Memória disponível < `20%`
  - GPU média > `70%`
  - VRAM usada > `85%`

### Operação

- Não abrir novo lane de modelo sem ADR.
- Não mudar default local sem medir VRAM e rodar smoke.
- Não introduzir segundo residente local sem justificativa formal.
- Não mudar voz/STT/TTS sem explicitar custo, fallback e impacto no homelab.

## Consequências

- `MODEL.md` da raiz removido — este ADR é a fonte de verdade.
- Qualquer alteração no stack de modelos exige novo ADR ou revisão deste.
- Agentes externos consultam este ADR via `docs/adr/README.md`.
