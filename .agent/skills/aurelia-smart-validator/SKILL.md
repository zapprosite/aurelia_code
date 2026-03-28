---
name: aurelia-smart-validator
description: Skill de auditoria e validação da Soberania Híbrida (T0 -> T1 -> T2). SOTA 2026.1.
tags: [ai, infrastructure, audit, litellm, sovereign]
---

# Aurelia Smart Validator (v2026.1)

Esta skill audita a saúde e a resiliência do orquestrador de inferência `aurelia-smart`. Ela garante que a hierarquia de tiers esteja operando conforme o contrato de Soberania Híbrida.

## Capacidades
- **Auditoria de Fallback**: Verifica se o Tier 0 (Ollama) é o ponto de entrada e se o Tier 1 assume em 10s.
- **Validação de Embeddings**: Garante que o `nomic-embed-text` local está gerando vetores.
- **Coleta de Métricas**: Extrai tempo de resposta e modelo ativo para o log de governança.

## Uso
Execute o script de auditoria para validar a infraestrutura após mudanças no `config.yaml`.

```bash
./scripts/audit-llm.sh
```

## Contrato de Tiers
1. **Tier 0**: Gemma 3 (Local).
2. **Tier 1**: Cloud Free (OpenRouter/Gemini/Groq).
3. **Tier 2**: Paid Expert (Minimax/Qwen).

## Manutenção
Sempre que uma nova API Key for adicionada ao `.env`, rode esta skill para validar se o roteador está reconhecendo o novo provedor na cascata.
