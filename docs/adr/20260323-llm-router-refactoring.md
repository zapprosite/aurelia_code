# ADR-20260323: Refatoração do Roteador LLM para Eficiência de Custo

## Contexto
O sistema atual de roteamento LLM utiliza uma abordagem baseada em palavras-chave para classificar tarefas e selecionar modelos. Embora funcional, essa abordagem é limitada na precisão da classificação e na flexibilidade de políticas de custo/performance. Precisamos de um sistema mais inteligente que utilize um modelo local potente (Gemma 3 12B) para agir como "Juiz" e rotear tarefas entre MiniMax (principal execução), OpenRouter (modelos baratos e fallback) e Kimi (contexto longo e multimodal).

## Decisão
Implementaremos um sistema de roteamento baseado em um Juiz local estruturado e políticas de fallback dinâmicas.

### Componentes Principais:
1.  **Juiz Local (Gemma 3 12B)**: Executado via Ollama, classifica a tarefa em uma de quatro categorias: `simple_short`, `coding_main`, `long_context_or_multimodal`, `critical`.
2.  **Provedor Direto MiniMax**: Implementação de um adaptador direto para a API da MiniMax, evitando dependência exclusiva do OpenRouter para o modelo principal.
3.  **Políticas de Roteamento**:
    - `simple_short`: DeepSeek V3.1 (Primary) -> Qwen 3 (Fallback).
    - `coding_main`: MiniMax M2.7 (Primary) -> Qwen 3 (Fallback).
    - `long_context_or_multimodal`: Kimi K2.5 (Primary) -> MiniMax M2.7 (Fallback).
    - `critical`: MiniMax M2.7 (Primary) -> Kimi K2.5 (Fallback).

### Regras de Override:
- **Multimodal (Imagem/PDF)**: Força `long_context_or_multimodal`.
- **Contexto Largo**: Força `long_context_or_multimodal`.
- **Confiança Baixa (< 0.6)**: Roteia para `coding_main`.
- **Retentativas (Retry >= 2)**: Escala para MiniMax ou Kimi conforme a natureza da tarefa.

### Controle de Resposta:
- Injeção de instrução de sistema: "Be concise. Avoid unnecessary explanations. Output minimal sufficient answer."
- Limites opcionais de tokens por rota.

## Consequências
- **Latência de Decisão**: Pequeno acréscimo de latência (~300-600ms) devido à chamada ao Gemma 3 local.
- **Eficiência de Custo**: Redução significativa ao evitar modelos caros para tarefas simples.
- **Robustez**: Fallbacks explícitos e observabilidade detalhada.

## Status
- [ ] Proposto (2026-03-23)
