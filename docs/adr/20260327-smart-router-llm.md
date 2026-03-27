# ADR 20260327: Implementação do Smart Router (Soberano 2026) 🤖🚦

## Status
Proposto (Design Técnico)

## Contexto
O ecossistema Aurélia necessita de uma camada de roteamento inteligente para otimizar o uso de recursos. 
- **Gemma 3 12b (Local)**: Excelente para tarefas de infraestrutura, triagem e respostas rápidas com custo zero.
- **OpenRouter (Kimi, Qwen, MiniMax)**: Modelos de maior capacidade para raciocínio complexo, codificação avançada e contextos longos.

## Decisão
Implementar um **Roteador Automático** no LiteLLM via `model_group` nomeado `aurelia-smart`. 

### Arquitetura de Roteamento:
1. **Tier 0 (Local - Prioridade 1)**: `ollama/gemma3:12b`.
2. **Tier 1 (External - Fallback/Scale)**:
   - `openrouter/qwen/qwen-2.5-coder-32b-instruct` (Lógica/Codificação).
   - `openrouter/minimax/minimax-01` (Criatividade/Contexto).
   - `openrouter/moonshotai/kimi-k2.5` (Eficiência).

### Lógica de Seleção:
- O roteador será exposto como um endpoint unificado.
- A seleção inicial prioriza o cluster local. Em caso de sobrecarga ou necessidade explícita de "High Reasoning" (via tags no prompt), o LiteLLM fará o escalonamento para o OpenRouter.

## Consequências
- **Positivas**: Redução drástica de custos (Sovereign First), resiliência em falhas locais, acesso a modelos SOTA globais.
- **Negativas**: Pequena latência adicional na camada de gateway do LiteLLM (ms).
- **Segurança**: Auditoria contínua de chaves via `audit-secrets.sh`.
