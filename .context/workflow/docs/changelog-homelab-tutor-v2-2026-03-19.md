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
