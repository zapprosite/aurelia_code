---
title: JARVIS Local Voice Blueprint
status: active
created: 2026-03-19
owner: codex
scope: voice-background-browser-use-antigravity-terminal
---

# JARVIS Local Voice Blueprint

## Objetivo

Transformar a Aurelia em um assistente local estilo JARVIS, com:

- escuta em background por microfone
- wake word local
- transcricao via Groq
- modelo local forte para seguir instrucoes e codar
- browser-use / Antigravity / terminal sob o mesmo orquestrador
- memoria semantica consistente em Qdrant
- persistencia autoritativa em Supabase

## Estado Real Medido no Homelab

Fonte: Prometheus por tras do Grafana local em `2026-03-19`.

- Grafana: `12.4.1`
- targets `up`: `node`, `cadvisor`, `nvidia-gpu`, `prometheus`
- CPU logica: `24`
- `avg_over_time(node_load1[15m])`: `4.37`
- CPU media 15m: `21.84%`
- memoria disponivel agora: `31.16%`
- GPU media 15m: `11.03%`
- VRAM usada agora: `18.61%`
- `nvidia-smi`:
  - GPU: `RTX 4090`
  - VRAM total: `24564 MiB`
  - VRAM livre no momento da leitura: `19791 MiB`

Leitura operacional:

- ha folga de CPU para wake word, VAD, fila de audio e automacao
- ha folga de GPU para um modelo local serio
- a folga de VRAM nao justifica concorrencia agressiva
- a regra correta e `1` inferencia pesada por vez

## Arquitetura Recomendada

```text
MIC
  -> wake word local
  -> VAD + ring buffer
  -> Groq STT
  -> intent router
      -> reply to Telegram/chat
      -> Antigravity IDE prompt
      -> browser-use / agent-browser
      -> CLI / terminal tools
  -> local LLM
  -> memory write
      -> Supabase
      -> Qdrant
  -> optional PT-BR TTS
```

## Escolhas de Stack

### 1. Wake Word

Use `openWakeWord`.

Motivo:

- roda localmente
- inclui `hey jarvis`
- custo de CPU baixo
- permite VAD embutido e supressao de ruido

### 2. VAD

Use `Silero VAD`.

Motivo:

- reduz falso positivo
- evita enviar silencio para Groq
- reduz custo e latencia

### 3. STT

Use `Groq whisper-large-v3-turbo`.

Parametros padrao:

- `model=whisper-large-v3-turbo`
- `language=pt`
- `temperature=0`

Fallback:

- Whisper local apenas quando Groq falhar

### 4. Modelo Local Principal

Escolha principal:

- `gemma3:27b-it-q4_K_M`

Alternativas:

- `qwen3.5:27b-q4_K_M`
- `qwen3-coder:30b` apenas para escalonamento manual

Escolha de operacao:

- concorrencia: `1`
- contexto operacional inicial: `8K`
- `OLLAMA_NUM_PARALLEL=1`
- `OLLAMA_FLASH_ATTENTION=1`
- `OLLAMA_KV_CACHE_TYPE=q4_0`

Motivo:

- melhor encaixe para o papel de orquestrador local
- function calling e structured output pesam a favor no comando do agente
- mais prudente para `RTX 4090 24 GiB` deixar o contexto inicial curto e estavel

Observacao:

- `gemma3:27b-it-q4_K_M` vira o default local
- `qwen3.5:27b-q4_K_M` fica como alternativa tecnica, nao como padrao
- `qwen3-coder:30b` continua forte, mas apertado demais para ficar residente no host

### 5. Embedding Unico para Qdrant

Regra 1:

- usar `bge-m3` como contrato unico

Contrato:

- dimensao: `1024`
- distancia: `cosine`
- chunk alvo: `350-500 tokens`
- overlap: `50-80 tokens`
- sem misturar embedding de outro modelo na mesma collection

Colecoes:

- `conversation_memory`
- `operator_notes`
- `runbook_memory`
- `code_context`

Payload minimo:

- `source`
- `kind`
- `role`
- `project`
- `conversation_id`
- `timestamp`
- `tags`
- `version`

### 6. Browser Layer

Camada primaria:

- `agent-browser` / Playwright controlado

Camada secundaria:

- `browser-use` para sessao persistente, tarefas de navegador mais longas e fluxo mais agentico

Regra:

- `agent-browser` primeiro para tarefas estruturadas e previsiveis
- `browser-use` quando a tarefa realmente exigir raciocinio navegador-centrico de varios passos

### 7. Antigravity IDE

Modo correto:

- a Aurelia nao domina a IDE por GUI cega
- ela abre a aba certa
- cola prompt estruturado
- coleta resposta
- transforma a saida em handoff executavel

### 8. Terminal

Modo correto:

- terminal via tools/CLI nativas
- nao via digitacao visual como caminho principal

## Rate Limits Alinhados

## Groq oficial

Para `whisper-large-v3-turbo`:

- `20 RPM`
- `2000 RPD`
- `7200 ASH`
- `28800 ASD`
- preco: `$0.04 / hora de audio`

## Rate limits recomendados da Aurelia

Esses limites sao mais conservadores que os da Groq e alinhados ao host atual.

### Audio ingress

- wake word aceito: max `4 ativações/min`
- clips simultaneos em fila: max `2`
- upload STT concorrente: `1`
- upload STT burst: `2`
- upload STT sustentado: `4 req/min`
- duracao maxima por clip: `20s`
- duracao alvo por clip: `8s-12s`

Justificativa:

- `4 req/min * 20s = 4800 audio-sec/h`, abaixo de `7200 ASH`
- mantem folga para erro, retry e uso manual

### Local LLM

- inferencias concorrentes: `1`
- fila maxima: `1`
- timeout duro por tarefa longa: `120s`
- downgrade para resposta curta quando houver browser ativo + STT ativo

### Browser / Antigravity

- sessoes browser-use ativas: `1`
- tabs controladas em paralelo: `1`
- retries automaticos por acao: `2`
- screenshots por etapa: `antes/depois` apenas em acoes relevantes

### Embeddings

- embeddings fora do caminho sincrono principal
- indexacao assincrona em lotes pequenos
- `1` worker de embedding por vez
- backpressure quando GPU/CPU cruzar limiar

## Guardrails de Recursos

Ative degradacao quando:

- CPU media 15m > `70%`
- memoria disponivel < `20%`
- GPU media 15m > `70%`
- VRAM usada > `85%`

Reacoes:

- pausar `browser-use`
- reduzir contexto do LLM local
- suspender embeddings nao urgentes
- responder texto sem TTS
- aumentar debounce do wake word

## O que Falta para Virar JARVIS de Verdade

1. `mic-daemon`
- processo em background sempre ativo

2. `wake-word-service`
- `openWakeWord` + threshold ajustavel

3. `audio-buffer-service`
- ring buffer pre-roll + post-roll

4. `vad-gate`
- so sobe STT quando ha fala real

5. `intent-router`
- `reply`
- `terminal`
- `antigravity`
- `browser`

6. `antigravity-adapter`
- abrir aba certa
- colar prompt certo
- ler saida

7. `resource-governor`
- aplicar rate limits
- aplicar backpressure
- aplicar degradacao

8. `storage-writer`
- gravar em `Supabase`
- indexar em `Qdrant`

9. `observability`
- latencia STT
- fila
- wake false positive
- browser task duration
- local LLM latency
- queue depth

## Slices de Implementacao

### Slice 1. Voice Foundation

- `openWakeWord`
- `Silero VAD`
- ring buffer
- servico de usuario

### Slice 2. STT Pipeline

- Groq STT
- retries
- rate limiting
- metricas

### Slice 3. Local Brain

- instalar `gemma3:27b-it-q4_K_M`
- ajustar contexto curto e seguro
- limitar concorrencia

### Slice 4. Memory Contract

- schema de `Supabase`
- collections `Qdrant`
- contrato `bge-m3`

### Slice 5. Action Router

- browser / antigravity / terminal / reply

### Slice 6. Resource Governor

- politicas baseadas em metricas
- thresholds
- backpressure

## Decisao Final

O desenho forte em `2026-03-19` para esta maquina e:

- ouvir sempre no CPU
- acordar com wake word local
- recortar fala com VAD
- transcrever na Groq
- pensar com `gemma3:27b-it-q4_K_M`
- agir com `agent-browser` / `browser-use` / CLI
- lembrar com `Supabase + Qdrant + bge-m3`
- proteger tudo com rate limit e degradacao por recurso

O maior erro seria tentar:

- modelo gigante sempre ativo ouvindo audio bruto
- varias inferencias paralelas
- browser-use em paralelo com embeddings e IDE pesada sem governor
- terminal controlado por GUI

## Fontes

- Groq Rate Limits: https://console.groq.com/docs/rate-limits
- Groq Whisper Turbo: https://console.groq.com/docs/model/whisper-large-v3-turbo
- Groq Speech-to-Text: https://console.groq.com/docs/speech-to-text
- Groq Spend Limits: https://console.groq.com/docs/spend-limits
- Google Gemma 3: https://blog.google/innovation-and-ai/technology/developers-tools/gemma-3/
- Ollama `gemma3`: https://ollama.com/library/gemma3
- Ollama `qwen3.5`: https://ollama.com/library/qwen3.5
- Ollama `qwen3-coder`: https://ollama.com/library/qwen3-coder
- BGE-M3 oficial: https://huggingface.co/BAAI/bge-m3
- Qdrant Embeddings: https://qdrant.tech/documentation/embeddings/
- Qdrant FastEmbed: https://qdrant.tech/documentation/fastembed/
- openWakeWord oficial: https://github.com/dscripka/openWakeWord
- browser-use oficial: https://github.com/browser-use/browser-use
