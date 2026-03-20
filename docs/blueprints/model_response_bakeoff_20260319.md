---
title: Model Response Bakeoff
status: active
created: 2026-03-19
owner: codex
scope: quality-latency-routing-eval
---

# Model Response Bakeoff

## Objetivo

Testar qualidade e latencia real de modelos locais e remotos com prompts curtos e uteis para a Aurelia:

- resposta operacional SRE
- saida estruturada de roteamento
- curadoria compacta para RAG

## Modelos testados

- local: `qwen3.5:9b`
- remoto premium: `minimax/minimax-m2.7`
- remoto tool/output: `deepseek/deepseek-v3.2`
- remoto barato/contexto: `qwen/qwen3.5-flash-02-23`

## Casos

### Caso 1. SRE curto

Prompt:

- cenario de reboot com `nvidia-gpu` down e restante do lab saudavel
- resposta em 4 linhas: diagnostico, primeiro comando, rollback, risco residual

Leitura:

- `minimax-m2.7` foi o mais polido
- `deepseek-v3.2` foi bom, mais seco
- `qwen3.5-flash` foi aceitavel, mas mais superficial
- o comportamento do `qwen3.5:9b` local nao foi confiavel com o budget curto usado neste bakeoff

### Caso 2. JSON de roteamento

Prompt:

- classificar lane para uma tarefa de screenshot + resumo de erro + comando bash seguro

Leitura:

- `deepseek-v3.2` entregou JSON valido e boa justificativa
- `qwen3.5-flash` entregou JSON valido e boa disciplina
- `minimax-m2.7` consumiu todo o budget em `reasoning` e terminou com `content=null`
- `qwen3.5:9b` local nao entregou conteudo util nesse teste

### Caso 3. Curadoria de RAG

Prompt:

- transformar uma nota operacional do Groq em `3 facts` e `3 tags`

Leitura:

- `deepseek-v3.2` foi o melhor equilibrio entre precisao e utilidade
- `qwen3.5-flash` foi bom e compacto
- `minimax-m2.7` novamente consumiu o budget em reasoning e nao entregou conteudo final
- `qwen3.5:9b` local, no endpoint atual, tambem consumiu o budget em reasoning e terminou vazio

## Latencias observadas

Valores aproximados da rodada:

- `qwen3.5:9b` local:
  - `~4.25s` no caso SRE
  - `~1.63s` nos prompts curtos seguintes
- `minimax-m2.7`:
  - `~4.5s-6.6s`
- `deepseek-v3.2`:
  - `~6.1s-6.4s`
- `qwen3.5-flash`:
  - `~27.8s-32.0s` nesta rodada

## Achado critico

No stack atual, `minimax-m2.7` e `qwen3.5:9b` mostraram um comportamento importante:

- em prompts curtos e estruturados, gastaram o budget em `reasoning`
- o resultado foi `finish_reason=length` com `content` vazio ou `null`

Isso significa:

- eles nao sao bons defaults para JSON curto/curadoria compacta sem governar `reasoning`
- para lanes de estrutura curta, hoje `deepseek-v3.2` e `qwen3.5-flash` foram mais previsiveis

## Decisao recomendada

### Melhor por lane

- SRE/office/polimento: `minimax-m2.7`
- JSON curto/roteamento: `deepseek-v3.2`
- curadoria curta para RAG: `deepseek-v3.2`
- remoto barato com disciplina boa: `qwen3.5-flash`

### Cuidado

- `minimax-m2.7` nao deve ser usado cegamente em prompts curtos com `max_tokens` baixo
- `qwen3.5:9b` local precisa de ajuste real de `thinking/reasoning` antes de virar o default de saida estruturada

## Consequencia para a Aurelia

A matriz de roteamento deve assumir:

- `qwen3.5:9b` continua forte como cerebro local, mas precisa governanca melhor para respostas estruturadas
- `deepseek-v3.2` e o melhor lane remoto para JSON/curadoria/tool-use compacto
- `minimax-m2.7` fica melhor como lane premium de workflow e resposta humana mais polida
