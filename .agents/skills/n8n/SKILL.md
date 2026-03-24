---
type: skill
name: n8n Specialist
description: Especialista em automação de fluxos no n8n para integração Claude <> Antigravity <> Bot.
skillSlug: n8n
phases: [P, E]
generated: 2026-03-18
updated: 2026-03-24
status: active
scaffoldVersion: "2.0.0"
---

# 🤖 n8n Specialist: Sovereign Automation 2026

Habilita o Antigravity a projetar, depurar e orquestrar fluxos de trabalho visuais no n8n que servem como o "sistema nervoso" da Aurélia.

## 🏛️ Padrões de Fluxo (Industrial)

### 1. Webhooks & Triggers
- **Aurélia Trigger**: Todo fluxo deve ter um ponto de entrada claro vindo do daemon ou bot.
- **Error Handling**: Implemente "Error Trigger" em todos os fluxos críticos para notificar falhas no canal de suporte do Telegram.

### 2. Integração com IA (Triple-Tier)
- **Node OpenAI/Custom**: Utilize o gateway da Aurélia para as requisições de IA dentro do n8n, garantindo o uso correto do MiniMax (Tier 1) para decisões de fluxo.
- **Local Fallback**: Para tarefas simples de texto, utilize o Ollama (Tier 3) via HTTP Request local.

### 3. Persistência de Dados
- **Postgres**: Utilize o banco de dados oficial em `srv/apps/postgres` para armazenar estados de longa duração.
- **Sovereign File System**: Salve arquivos temporários em `/tmp/n8n/` e mova-os para o Monorepo após processamento.

## 📍 Quando usar
- Para automatizar newsletters, backups ou triagem de mensagens.
- Para integrar a Aurélia com APIs de terceiros (Google Calendar, Slack, etc.).
- Para criar pipelines de processamento de áudio/vídeo em lote.

## 🛡️ Guardrails
- **Security First**: Nunca insira senhas em texto puro nos nós do n8n. Use as credenciais seguras do n8n ou o vault KeePassXC.
- **Rate Limit**: Evite loops infinitos que possam consumir toda a CPU do 7900x.
