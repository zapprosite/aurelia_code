---
title: Groq PT-BR Audio Blueprint
status: active
created: 2026-03-19
owner: codex
scope: aurelia-audio-stt-groq-ptbr-memory-local-llm
---

# Groq PT-BR Audio Blueprint

## Objetivo

Usar a API da Groq para assumir a transcrição de áudio e liberar GPU local, mantendo:

- raciocínio e instruções no modelo local
- memória e estado em `Qdrant + Supabase`
- saída de voz em solução PT-BR separada

## Ponto de Partida Real

O repositório já possui integração Groq para STT:

- `pkg/stt/groq.go`
- `cmd/aurelia/app.go`
- `internal/config/config.go`
- fluxo de onboarding para `groq_api_key`

Ou seja:

- Groq já é conhecido pela Aurelia
- o modo certo agora é operacionalizar a arquitetura PT-BR

## Decisão Arquitetural

### O que vai para Groq

- `speech-to-text`

### O que fica local

- system prompt
- tool use
- decisão
- planejamento
- orquestração
- memória operacional

### O que fica em storage local

- `Supabase`: sessão, mensagens, jobs, auditoria
- `Qdrant`: embeddings e recuperação semântica

### O que não vai para Groq como núcleo

- memória longa
- histórico autoritativo
- política de ferramentas
- governança

## Arquitetura Final

```text
Audio Input
  -> Groq STT
      -> transcript
          -> Supabase (source of truth)
          -> Qdrant (semantic memory)
          -> Local LLM (instruction follower)
              -> text response
              -> PT-BR TTS layer
```

## Papel de Cada Componente

### 1. Groq

Responsável por:

- transcrever áudio rapidamente
- devolver texto limpo para o pipeline local

Não é responsável por:

- seguir instruções da Aurelia
- guardar memória
- decidir ações

### 2. Supabase

Responsável por:

- guardar sessão
- guardar transcrição bruta
- guardar resposta final
- guardar jobs de áudio
- guardar auditoria de tool use

Tabelas recomendadas:

- `sessions`
- `messages`
- `audio_jobs`
- `tool_runs`
- `incident_log`

### 3. Qdrant

Responsável por:

- indexar mensagens relevantes
- indexar runbooks
- indexar notas operacionais
- recuperar contexto semântico

Coleções recomendadas:

- `conversation_memory`
- `runbook_memory`
- `operator_notes`

### 4. Modelo Local

Responsável por:

- seguir instruções
- planejar
- usar tools
- responder com política local

Modelo recomendado:

- manter o principal local
- usar Groq apenas como camada de áudio

### 5. TTS PT-BR

Responsável por:

- sintetizar resposta em português

Diretriz:

- manter TTS PT-BR separado da Groq
- encaixar como etapa final do pipeline

## Fluxo Operacional

### Entrada

1. receber áudio
2. normalizar formato
3. enviar à Groq STT

### Pós-transcrição

1. salvar transcrição no `Supabase`
2. gerar embedding
3. indexar no `Qdrant`
4. recuperar contexto relevante
5. montar prompt local

### Raciocínio

1. modelo local recebe:
   - system prompt
   - histórico curto
   - memória semântica
   - contexto da sessão
2. decide resposta
3. decide ferramentas

### Saída

1. resposta textual
2. opcionalmente enviar para camada `TTS PT-BR`

## Política PT-BR

### STT

- `Groq` como padrão
- manter fallback local

### TTS

- usar solução PT-BR separada
- não acoplar a Aurelia ao TTS remoto se isso ferir latência, custo ou qualidade

### LLM

- manter local como autoridade

## Fallback

### Se Groq falhar

- fallback para Whisper local

### Se TTS PT-BR falhar

- retornar texto
- ou usar TTS local existente

### Se Qdrant falhar

- responder com histórico curto do Supabase
- registrar perda de contexto semântico

### Se Supabase falhar

- modo degradado sem persistência longa
- registrar incidente

## Configuração Mínima

### Chaves

- `groq_api_key`

### Provider

- `stt_provider = groq`

### Persistência

- `Supabase` para registro
- `Qdrant` para memória

## Roadmap

### Fase 1

- validar `Groq STT` em produção
- persistir transcrição no Supabase

### Fase 2

- indexar transcrição no Qdrant
- recuperar memória semântica no pipeline

### Fase 3

- manter o LLM local como cérebro
- integrar TTS PT-BR de forma limpa

### Fase 4

- adicionar observabilidade:
  - latência STT
  - custo
  - taxa de erro
  - fallback acionado

## Critérios de Aceite

- Aurelia transcreve áudio via Groq
- transcrição fica registrada no Supabase
- memória semântica entra no Qdrant
- resposta continua sendo decidida localmente
- PT-BR funciona na saída
- fallback local existe e é testado

## Resumo Executivo

Arquitetura correta:

- `Groq = ouvido`
- `Qdrant = memória semântica`
- `Supabase = histórico e estado`
- `LLM local = cérebro`
- `TTS PT-BR = voz`
