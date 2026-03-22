> [!NOTE]
> Status: ✅ Arquivado / Concluído em 22/03/2026

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

### Hardware Pinned (Verificado em 2026-03-22)

- CPU: `AMD Ryzen 9 7900X` (12 Cores, 24 Threads) → **Master of Voice & Orchestration**
- GPU: `NVIDIA GeForce RTX 4090` (24 GiB VRAM) → **Pure Inference Engine**
- RAM: `32 GiB DDR5 Gen5` → **System Buffer** (Objetivo: Swap < 5%)

### Lógica Matemática Pinned (Budget de VRAM)

| Recurso | Modelo | VRAM Alocada | Notas |
|--------|--------|--------------|-------|
| **Cérebro (Primary)** | `gemma3:12b` | 8.1 GiB | Ativo 100% do tempo |
| **Contexto (Fallback)** | `qwen3.5:9b` | 6.6 GiB | Residente para tarefas longas/multimodais |
| **Memória (Cognitive)** | `bge-m3` | 1.2 GiB | Residente para busca semântica |
| **Overhead (System/KV)** | - | ~4.0 GiB | Xorg, Rustdesk, context window, KV cache |
| **TOTAL PINNED** | - | **~19.9 GiB** | **Segurança: ~4.1 GiB Livres** ✅ |

### Regras da Lógica Matemática

1.  **Voz no CPU**: O processamento de áudio (Kokoro/STT local se houver) deve usar os 12 cores do 7900X, poupando VRAM para os LLMs.
2.  **Pinned Residency**: `gemma3:12b`, `qwen3.5:9b` e `bge-m3` devem estar sempre carregados para latência zero.
3.  **Anti-Swap**: Se a VRAM exceder 22GB ou a RAM 90%, o Agente deve abortar tarefas pesadas de imagem/video para proteger o kernel.
4.  **Descarte do 27B**: O modelo `gemma3:27b` foi removido por ser matematicamente incompatível com a política de residência multi-modelo.

### Regras Fechadas

1. `gemma3:12b` é o residente principal.
2. `gemma3:27b` exige comando explícito e descarregamento de outros LLMs pesados.
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
