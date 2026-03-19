# Homelab Tutor v2 - 2026-03-19

## Goal

Criar uma base profissional para a Aurelia operar como tutor do homelab com:

- governança
- incident response
- runbooks executáveis
- memória operacional

## External Skills Installed

- `incident-response`
- `system-architect`
- `c-level-advisor`

Install target:

- `~/.agents/skills/`
- sincronizado para Antigravity, Codex, Claude Code e Gemini CLI

## Local Tutor Structure Created

Nova skill principal:

- `~/.aurelia/skills/homelab-tutor-v2/`

Arquivos criados:

- `SKILL.md`
- `INDEX.md`

## Runbooks Created

Monitoring:

- `gpu-metrics-recover.md`
- `prometheus-target-down.md`
- `grafana-no-data-triage.md`

Docker / Compose:

- `compose-service-missing.md`
- `n8n-health-recover.md`

Data / Platform:

- `qdrant-health-recover.md`
- `supabase-health-triage.md`

Tunnel / Network:

- `cloudflare-tunnel-recover.md`
- `firewall-drift-audit.md`

AI Runtime:

- `ollama-health-recover.md`
- `voice-stack-recover.md`

DR / Backup:

- `backup-age-enforcer.md`
- `zfs-scrub-review.md`

Additional infrastructure runbooks:

- `caprover-health-recover.md`
- `litellm-health-recover.md`
- `postgres-direct-check.md`
- `tailscale-access-check.md`
- `gpu-contention-triage.md`

## Domain Catalog Added

Tutor catalog created:

- `~/.aurelia/skills/homelab-tutor-v2/DOMAIN_RESPONSES.md`

This catalog standardizes how Aurelia should answer and route by domain:

- monitoring
- docker / compose
- data / platform
- AI runtime
- network / exposure
- DR / backup / storage

## Operating Model

O tutor agora deve operar assim:

1. classificar o domínio
2. escolher skill/runbook
3. coletar prova
4. aplicar correção mínima
5. validar
6. registrar prevenção

## Self-Healing Rule

Todo incidente relevante deve resultar em um dos seguintes:

- novo runbook
- atualização de runbook existente
- registro em `.context/workflow/docs/`
- refinamento de guardrail

## Audio Architecture Extension

Foi criado um blueprint específico para áudio PT-BR com Groq:

- `groq_ptbr_audio_blueprint.md`

Direção arquitetural registrada:

- `Groq` como camada de `speech-to-text`
- `Supabase` como fonte de verdade de sessões, mensagens e jobs
- `Qdrant` como memória semântica
- `LLM local` como cérebro e executor de instruções
- `TTS PT-BR` como camada separada de saída

## Audio Execution Slice

Primeiro corte executável do blueprint de áudio aplicado no código:

- `internal/telegram/input.go`
- `internal/telegram/input_audio_test.go`
- `scripts/groq-stt-curl-smoke.sh`
- `groq_stt_simulation.md`

O que mudou:

- transcrição de voz agora é persistida na memória local logo após o STT
- o arquivo arquivado usa `message_type=audio_transcript`
- o conteúdo salvo registra `provider`, `basename` do arquivo e texto transcrito
- foi adicionado smoke script de `curl` para validar a Groq com ou sem token

Provas registradas:

- `go test ./internal/telegram ./pkg/stt` passou
- smoke local sem token imprime o `curl` exato em `dry-run`
- chamada real com token falso retorna `401 invalid_api_key`, confirmando contrato HTTP

## Audio Cost/Quality Optimization

Slice adicional aplicado após ativação real da Groq:

- `pkg/stt/groq.go`
- `scripts/groq-stt-curl-smoke.sh`
- `groq_stt_simulation.md`

Decisões executadas:

- troca do modelo padrão para `whisper-large-v3-turbo`
- envio explícito de `language=pt`
- envio explícito de `temperature=0`

Motivo:

- reduzir custo por hora de transcrição
- reduzir latência
- melhorar previsibilidade para PT-BR

Provas:

- `GET /openai/v1/models` retornou `HTTP 200`
- `POST /openai/v1/audio/transcriptions` com WAV mínimo retornou JSON válido
- `go test ./pkg/stt` passou

## Antigravity Gemini Operator Slice

Foi criada a base para a Aurelia usar o chat do Antigravity como copiloto leve, sem misturar isso com execucao principal:

- `implementation_plan.md`
- `task.md`
- `docs/antigravity_gemini_operator_blueprint.md`
- `docs/PROJECT_PLAYBOOK.md`
- `.aurelia/skills/antigravity-gemini-operator/SKILL.md`
- `.aurelia/skills/antigravity-gemini-operator/CHAT_PROMPTS.md`
- `.aurelia/skills/antigravity-gemini-operator/DECISION_MATRIX.md`
- `.agents/skills/antigravity-gemini-operator/SKILL.md`
- `.agents/skills/antigravity-gemini-operator/CHAT_PROMPTS.md`
- `.agents/skills/antigravity-gemini-operator/DECISION_MATRIX.md`

Direcao registrada:

- o runtime do projeto passa a ter um playbook carregavel por memoria canonica
- a skill do projeto fica no path real carregado pela Aurelia
- a skill tambem fica espelhada em `.agents/skills/` para governanca versionada
- Gemini Flash / Minimax 2.7 no Antigravity vira executor de microtarefas
- Codex / CLI continua como executor principal, validador e agente de commit

## Telegram Light Delegation Slice

Foi ligado ao pipeline real do Telegram um atalho para tarefas `light`:

- `internal/telegram/antigravity_prompt.go`
- `internal/telegram/antigravity_prompt_test.go`
- `internal/telegram/input_pipeline.go`

Comportamento:

- a entrada do usuario e classificada por heuristica leve
- se a tarefa for `light`, a Aurelia gera automaticamente um prompt estruturado para o chat do Antigravity
- a resposta evita passar pelo loop completo quando a melhor acao e delegar pesquisa curta, config pequena ou revisao de diff
- tarefas com sinais de alto risco continuam fora desse atalho

Prova:

- `go test ./internal/telegram` passou

## JARVIS Local Voice Blueprint

Foi registrado um blueprint direto ao ponto para a Aurelia virar um assistente local de voz, browser-use, Antigravity e terminal:

- `docs/jarvis_local_voice_blueprint_20260319.md`

Baseado em:

- metricas reais de `Grafana/Prometheus`
- VRAM real da `RTX 4090`
- limites oficiais atuais da Groq
- stack local ja presente no host

Direcao registrada:

- `openWakeWord` + `Silero VAD` no CPU
- `Groq whisper-large-v3-turbo` para STT
- `qwen3-coder:30b` como cerebro local
- `bge-m3` como contrato unico de embedding
- `Supabase + Qdrant` como memoria/persistencia
- `agent-browser` primeiro e `browser-use` como camada avancada
- rate limits conservadores e governor por recurso

## Gemini Fallback Runtime

Foi consolidado o papel da Gemini API no runtime da Aurelia:

- `scripts/gemini-smoke.sh`
- `docs/gemini_fallback_runtime_20260319.md`

Decisao registrada:

- Gemini entra como fallback LLM e pesquisa curta
- `gemini-2.5-flash` e o default remoto
- `gemini-2.5-pro` fica para escalonamento manual
- Groq continua no STT
- `bge-m3` continua como embedding unico do Qdrant

Tambem foi endurecido o health:

- checks auxiliares agora podem emitir `warning`
- `warning` nao gera falso degradado
- `error` continua derrubando o `/health`

Observacao adicional validada em `2026-03-19`:

- a oferta estudantil do ecossistema Gemini nao deve ser confundida com aumento automatico da cota da Gemini API
- o beneficio de estudante agrega mais no `Gemini app`, `NotebookLM`, integracoes Google e armazenamento do plano `Google AI Pro`
- a Gemini Developer API continua dependente de `project`, `usage tier` e `Cloud Billing`
- conclusao arquitetural mantida: Gemini segue auxiliar no runtime da Aurelia, nao caminho critico
