---
title: Model Routing Matrix
status: active
created: 2026-03-19
owner: codex
scope: local-remote-audio-research-cost-quality
---

# Model Routing Matrix

## Objetivo

Fechar uma politica unica de roteamento para a Aurelia que reduza custo mensal sem perder qualidade nem estabilidade do homelab.

Esta matriz governa:

- modelos locais no Ollama
- modelos remotos via OpenRouter
- audio via Groq
- pesquisa profunda via Gemini web
- ingestao curada para RAG e semantico

## Principios

- o caminho critico do homelab deve ser majoritariamente local
- o remoto entra por escalonamento, nao por padrao
- audio e um lane proprio
- pesquisa profunda e um lane proprio
- `SQLite` manda no runtime
- `Qdrant` indexa material curado
- `Supabase` complementa o estado compartilhado

## Lane por classe

| Classe | Primario | Papel | Regra |
|---|---|---|---|
| `local-fast` | `qwen3.5:4b` | triagem, classificacao, resumo curto | nao fica residente por padrao |
| `local-balanced` | `qwen3.5:9b` | tool use local, repo, manutencao, respostas do bot | unico modelo residente do caminho ativo |
| `remote-cheap-long-context` | `openrouter/qwen/qwen3.5-flash` | contexto grande, triagem remota barata, resumo longo | usar so quando local nao for suficiente |
| `remote-cheap-vision` | `openrouter/qwen/qwen3.5-9b` | visao barata, multimodal, leitura de tela/imagem | nao usar como default do bot |
| `remote-tool-long-output` | `openrouter/deepseek/deepseek-v3.2` | tool use remoto e saidas longas | usar por escalonamento |
| `remote-premium-workflow` | `openrouter/minimax/minimax-m2.7` | workflow digital premium, browser, redacao melhor | lane premium, nao default |
| `audio-stt` | `groq/whisper-large-v3-turbo` | transcricao PT-BR | lane isolado, com governor proprio |
| `deep-research` | `Gemini web` | pesquisa profunda e relatorio com fontes | nunca entra direto no runtime |

## Matriz por tipo de tarefa

| Tarefa | Primario | Fallback | Evitar |
|---|---|---|---|
| classificar tarefa | `qwen3.5:4b` | `qwen3.5:9b` | `minimax` |
| tool use local | `qwen3.5:9b` | `deepseek-v3.2` | `gemma3:27b` |
| manutencao do homelab | `qwen3.5:9b` | `minimax-m2.7` so se explicito | qualquer remoto como default |
| browser/agent task premium | `minimax-m2.7` | `deepseek-v3.2` | `qwen3.5:4b` |
| leitura de imagem/tela barata | `qwen3.5-9b` | `qwen3.5-flash` | `minimax` por custo |
| resumo de log/doc grande | `qwen3.5-flash` | `qwen3.5:9b` | `minimax` |
| saida longa remota | `deepseek-v3.2` | `minimax-m2.7` | `qwen3.5-9b` |
| STT de audio | `groq whisper-large-v3-turbo` | whisper local | qualquer lane LLM |
| pesquisa profunda | `Gemini web` | pesquisa manual web | uso automatico no runtime |
| ingestao para RAG | markdown curado + `bge-m3` | fila manual de curadoria | indexacao bruta de chat remoto |

## Politica de custo

### Default economico

- `qwen3.5:9b` deve resolver a maior parte das tarefas do bot
- `qwen3.5:4b` entra para triagem e respostas curtas
- `OpenRouter` entra so em:
  - visao
  - contexto muito longo
  - workflow premium
  - escalonamento explicito
- `Groq` entra so no ouvido
- `Gemini web` alimenta conhecimento, nao runtime

### Distribuicao alvo

- `70-85%` local
- `10-20%` audio STT
- `5-10%` remoto premium e visao

### Regras de economia

- cron e maintenance loop nunca devem usar remoto por padrao
- watchdog, health e repair loop devem ser locais
- `OpenRouter` so entra com motivo explicito de capacidade
- `Gemini web` nunca indexa direto no `Qdrant` sem curadoria

## Regra de VRAM

- `qwen3.5:9b` e o unico residente do caminho ativo
- `qwen3.5:4b` so aquece sob demanda
- `gemma3:27b-it-q4_K_M` fica fora do runtime ativo
- `qwen3-coder:30b` e manual e nao residente

Motivo:

- o host usa ~`4.8 GiB` de VRAM em base
- `qwen3.5:9b` carregado deixa ~`10.5 GiB` livres
- `qwen3.5:9b + qwen3.5:4b` juntos deixam ~`3.8 GiB` livres
- um `27B` deixa folga perto de `2 GiB`

## Politica para Gemini web

Uso correto:

- pesquisar temas novos
- gerar relatorio profundo com fontes
- exportar para markdown, doc ou nota operacional
- submeter a curadoria da Aurelia
- so depois indexar no `Qdrant`

Uso incorreto:

- alimentar runtime automaticamente
- substituir runbook
- substituir facts do `SQLite`
- virar dependencia do maintenance loop

## Fluxo de ingestao curada

```text
Gemini web / Deep Research
  -> export/share/copy
  -> markdown curado
  -> facts extraidos
  -> armazenamento local
      -> docs/
      -> SQLite notes/facts
  -> embeddings `bge-m3`
  -> Qdrant
```

## Regras de falha e fallback

- `Groq` falhou ou bateu budget: usar STT local ou modo texto
- `OpenRouter` falhou: voltar para `qwen3.5:9b`
- `minimax` indisponivel: usar `deepseek-v3.2`
- `vision` indisponivel remoto: cair para fluxo manual ou screenshot + descricao

## Enforcement atual

Primeiro corte implementado:

- endpoint interno `POST /v1/router/dry-run`
- policy engine inicial em `internal/gateway/`
- guardas iniciais de `reasoning_mode`, `max_output_tokens` e `soft_timeout_ms`

Escopo do corte:

- ainda nao altera o runtime principal automaticamente
- ja permite explicar qual lane/modelo seria escolhido e com quais guardas

## Decisao final

Para economia mensal sem perder qualidade:

- `qwen3.5:9b` fica como cerebro local do bot
- `Groq` fica isolado no audio
- `LiteLLM` serve clientes e aliases, nao decide manutencao
- `Gemini web` vira motor de pesquisa e inteligencia curada
- `Qdrant` indexa so material aprovado
- `SQLite` continua como verdade operacional do homelab
