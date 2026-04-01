# Report: gemma3:27b-it-qat Pull - 20260401

## Status
- **Modelo**: gemma3:27b-it-qat
- **Fonte**: Ollama library (QAT optimized)
- **Status**: ✅ PRONTO

## Timeline
- 2026-04-01: Pull executado com sucesso via `ollama pull gemma3:27b-it-qat`

## Verificação
```bash
$ ollama list
NAME                    ID          SIZE      MODIFIED
gemma3:27b-it-qat       ccc0cddac561  18 GB    3 minutes ago
nomic-embed-text:latest e4c4a6c     274 MB   3 months ago
```

## Teste PT-BR
```bash
$ ollama run gemma3:27b-it-qat "Responda em PT-BR: qual seu nome?"
Meu nome é Gemma, um modelo de linguagem grande treinado pelo Google DeepMind...
```

## .env atualizado
```bash
OLLAMA_MODEL=gemma3:27b-it-qat
```

## Notas
- Modelo QAT (Quantization-Aware Training) - melhor qualidade que Q4_K_M
- 18GB total, ~100% VRAM na RTX 4090
- Suporta visão (gemma3:27b-it-qat)
