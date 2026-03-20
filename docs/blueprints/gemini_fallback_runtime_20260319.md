---
title: Gemini Fallback Runtime Notes
status: active
created: 2026-03-19
owner: codex
scope: google-gemini-api-runtime-fallback
---

# Gemini Fallback Runtime Notes

## Objetivo

Registrar como a API Gemini agrega ao runtime da Aurelia sem competir com Groq STT, Qdrant `bge-m3` e o cerebro local.

## Decisao

Usar Gemini como provedor auxiliar de LLM e pesquisa curta:

- `gemini-2.5-flash` como padrao remoto rapido
- `gemini-2.5-pro` como escalonamento manual para sintese mais pesada

Nao usar Gemini nesta fase para:

- STT
- TTS
- embeddings do Qdrant
- memoria autoritativa

## Observacao sobre AI Studio / conta estudante

O valor da conta estudante/AI Studio esta mais na experiencia de produto e nas ferramentas do ecossistema Gemini do que em transformar a API num caminho primario sem limites.

Status oficial conferido em `2026-03-19`:

- a pagina oficial de estudantes do Gemini informa que a oferta estudantil anterior expirou em `2026-03-11` na regiao renderizada e hoje mostra `1 month Google AI Pro trial`
- essa camada de estudante/Google AI Pro libera principalmente o app Gemini, Gemini nos apps Google, recursos avancados do NotebookLM e `2 TB` de armazenamento
- a propria pagina oficial informa que uma conta `Google Workspace` da universidade pode dar acesso a ferramentas como Gemini app e NotebookLM, dependendo do admin da instituicao
- nada disso muda por si so a cota da Gemini Developer API
- na Gemini API, rate limits sao por `project` e `usage tier` no AI Studio; para sair do free tier e ir ao paid tier e preciso habilitar `Cloud Billing`

Diretriz operacional:

- tratar a Gemini API como auxiliar
- nao depender dela para o caminho critico do Jarvis
- continuar com Groq para voz e `bge-m3` para embeddings

## Runtime

- segredo salvo apenas em `~/.aurelia/config/app.json`
- health mostra se `google_api_key` esta configurada
- smoke oficial: `scripts/gemini-smoke.sh`

## Provider policy

- caminho principal: local
- caminho remoto rapido: Gemini Flash
- caminho remoto caro/manual: Gemini Pro
