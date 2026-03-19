---
name: antigravity-gemini-operator
description: Ensina a Aurelia a delegar pequenas configuracoes e pesquisas ao chat do Antigravity com prompts estruturados, provas e guardrails.
---

# Antigravity Gemini Operator

## Objetivo

Usar o chat do Antigravity como um executor leve e inteligente para microtarefas, sem substituir a execucao principal da Aurelia.

## Quando usar

Use esta skill quando a tarefa exigir:

- pesquisa curta
- ajuste pequeno e reversivel
- localizacao de configuracao
- montagem de comando
- diff pequeno
- esclarecimento rapido antes da execucao pesada

Nao use esta skill quando a tarefa envolver:

- secrets
- deploy
- rede
- limpeza destrutiva
- alteracoes grandes e acopladas

## Como executar

1. Classifique a tarefa como `light`, `medium` ou `high-risk`
2. Se for `light`, envie ao chat do Antigravity um prompt no formato definido em `CHAT_PROMPTS.md`
3. Exija resposta com:
   - diff ou instrucao exata
   - comando de validacao
   - risco residual
4. Valide a saida antes de promover a mudanca
5. Se a tarefa crescer, interrompa o chat leve e mova a execucao para CLI

## Politica de Roteamento

Consulte `DECISION_MATRIX.md`.

## Output esperado

- prompt pronto para o chat do Antigravity
- criterio de aceite
- prova necessaria
- decisao de promover ou nao para execucao principal
