---
title: Local Model Kit Blueprint
status: active
created: 2026-03-19
owner: codex
scope: local-llm-kit-ollama-hf-qdrant
---

# Local Model Kit Blueprint

## Objetivo

Fechar o kit local de modelos da Aurelia para:

- agent bot local em Go
- instruido por repositorio
- memoria semantica em Qdrant
- browser, Antigravity e terminal ao mesmo tempo
- sem depender de Gemini API no runtime ativo

## Kit Novo

### Modelo principal

- `qwen3.5:9b`

Papel:

- instrucao local
- orquestracao
- tool use
- repositorio
- raciocinio principal do agente

### Modelo alternativo

- `qwen3.5:4b`

Papel:

- roteamento curto
- fallback de latencia
- aquecimento sob demanda

### Modelo de escalonamento

- `qwen3-coder:30b`

Papel:

- uso manual
- tarefas de codigo mais agressivas
- nao residente por padrao

### Modelo leve

- `gemma3:27b-it-q4_K_M`

Papel:

- laboratorio manual
- comparacao
- teste offline fora do caminho ativo

### Embedding

- `bge-m3`

Papel:

- contrato unico de embedding no Qdrant

## Decisao Final

Escolher `qwen3.5:9b` como default local.

Motivo:

- suporte real a `tools`
- latencia melhor para o papel de controle local
- margem de VRAM profissional para browser, automacao e cache
- `gemma3:27b-it-q4_K_M` continua disponivel para laboratorio sem entrar no runtime ativo

## Leitura do Host

- GPU: `RTX 4090`
- VRAM total: `24564 MiB`
- VRAM livre medida: ~`19 GiB`
- CPU: `Ryzen 9 7900X`
- RAM disponivel medida: ~`8.9 GiB`

Leitura operacional:

- da para rodar um modelo serio
- nao da para brincar de concorrencia agressiva
- spill pesado para RAM nao compensa
- o caminho certo e `1` modelo residente por vez

Conta pratica:

- uso base do host: ~`4810 MiB`
- VRAM livre em idle: ~`19238 MiB`
- `qwen3.5:9b` carregado: ~`9.2 GiB` de VRAM
- folga restante com `9b`: ~`10.5 GiB`
- folga restante com `9b + 4b`: ~`3.8 GiB`
- folga restante com `27B`: ~`2 GiB`

## Politica de Ollama

Usar:

- `OLLAMA_NUM_PARALLEL=1`
- `OLLAMA_FLASH_ATTENTION=1`
- `OLLAMA_KV_CACHE_TYPE=q4_0`
- `OLLAMA_CONTEXT_LENGTH=8192`

Subir para `12288` ou `16384` so depois de prova real de estabilidade.

## Kit de Pull

```bash
ollama pull qwen3.5:9b
ollama pull qwen3.5:4b
ollama pull bge-m3:latest
```

Kit opcional:

```bash
ollama pull gemma3:27b-it-q4_K_M
ollama pull qwen3-coder:30b
```

## Scripts Operacionais

Padrao do kit:

- `scripts/update-ollama.sh`
- `scripts/ollama-local-kit-smoke.sh`

Uso:

```bash
bash ./scripts/update-ollama.sh
bash ./scripts/ollama-local-kit-smoke.sh
```

## Papel do HF Token

Segredo salvo localmente em:

- `~/.aurelia/config/secrets.env`

Campo:

- `HF_TOKEN=...`

Uso:

- downloads autenticados da Hugging Face
- modelos nao disponiveis diretamente pelo fluxo normal
- contingencia para artefatos/model cards privados ou limitados

Regra:

- nao entra no repo
- nao entra em prompt
- nao entra em `app.json`

## Contrato de Runtime

- LLM remoto principal quando necessario: `openrouter/minimax/minimax-m2.7`
- STT: `Groq`
- LLM local principal: `qwen3.5:9b`
- embedding: `bge-m3`
- browser: `agent-browser`
- desktop fallback: `xdotool` / `wmctrl`

Estado do codigo:

- o app agora aceita `llm_provider=ollama`
- o catalogo local lista modelos via `http://127.0.0.1:11434/v1/models`
- o onboarding pula a etapa de API key para `ollama`
- o `/health` verifica endpoint local + modelo configurado quando `provider=ollama`

Regra operacional:

- nao virar a config ativa para `ollama` enquanto o daemon live estiver em um binario que ainda nao recebeu esse slice

## Gate de Aceite

Para chamar o kit de pronto:

1. `ollama list` mostra o kit base instalado
2. `bash ./scripts/ollama-local-kit-smoke.sh` passa com `num_ctx=8192`
3. `OLLAMA_NUM_PARALLEL=1` fica enforced
4. browser-use + `qwen3.5:9b` coexistem sem degradacao visivel
5. `bge-m3` segue como unico embedding do Qdrant
