# ADR 20260328-fallback-gateway-soberano.md

## Contexto
O Gateway de LLM da Aurélia (LiteLLM + Smart Router Go) estava apresentando erro **401 Unauthorized** nas rotas `aurelia-smart` e `gemma3:27b`. O diagnóstico revelou que o LiteLLM estava tentando resolver modelos locais via OpenRouter devido a mapeamentos incorretos e placeholders no `config.yaml`.

## Decisão
Implementar a **Soberania do Tier 0** através das seguintes ações:
1. **Mapeamento Explícito**: Adicionado `aurelia-smart` e `gemma3:27b` diretamente ao `model_list` do LiteLLM apontando para `http://localhost:11434` (Ollama nativo).
2. **Correção de Provedores**: Removidos placeholders (`*_PLACEHOLDER`) no `config.yaml` em favor da sintaxe `os.environ/KEY`, permitindo que o LiteLLM use as chaves reais do `.env`.
3. **Bypass do Judge**: Recomendada a alteração do `OLLAMA_URL` para a porta 11434 no `.env` para garantir que o classificador de tarefas (Judge) não dependa do gateway LiteLLM.

## Consequências
- **Positivas**: O sistema volta a ser funcional mesmo sem conectividade com a nuvem (Modo Soberano). Redução de latência para tarefas categorizadas como `simple_short` ou `general`.
- **Negativas**: O log de custos do LiteLLM não capturará as chamadas diretas ao Ollama feitas pelo Judge, mas a estabilidade do sistema é prioritária.

---
**Status**: Implementado
**Autor**: Antigravity (Gemini 2.0)
**Data**: 2026-03-28
