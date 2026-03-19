---
description: Slice nonstop para portar gateway + voz para a worktree de deploy e provar o runtime live.
status: in_progress
---

# ADR-20260319-deploy-gateway-voice

## Status

- Em execução

## Slice

- slug: deploy-gateway-voice
- owner: codex
- branch/worktree: `feat/24x7-system-service` em `/home/will/aurelia-24x7`
- json de continuidade: docs/adr/taskmaster/ADR-20260319-deploy-gateway-voice.json

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O repositório principal já fechou:

- gateway enforcement com budgets/breaker
- spool/processador de voz
- capture worker comandado por contrato
- capturador real com `openWakeWord` + VAD via script
- store/mirror locais em SQLite

O gap restante é o deploy: a worktree `/home/will/aurelia-24x7` está com drift próprio e o serviço live do host ainda roda o binário de `/usr/local/bin/aurelia`.

## Decisão

O rollout deve ser feito em cima da worktree de deploy existente, sem descartar o drift atual, e validado com provas reais:

- build verde na worktree de deploy
- health live com rotas de gateway/voice
- capture local sem falso positivo em idle
- `spool -> STT -> resposta` e status routes funcionando

## Escopo

- integrar mudanças de gateway/voz na worktree `/home/will/aurelia-24x7`
- build/test nessa worktree
- reinstalar binário/serviço live a partir dela
- validar `/health`, `/v1/router/status`, `/v1/voice/status`, `/v1/voice/capture/status`

## Fora de escopo

- merge em `main`
- alteração da worktree do Gemini Flash
- TTS premium

## Arquivos afetados

- `/home/will/aurelia-24x7/cmd/aurelia/*`
- `/home/will/aurelia-24x7/internal/gateway/*`
- `/home/will/aurelia-24x7/internal/voice/*`
- `/home/will/aurelia-24x7/scripts/*`
- `/home/will/aurelia-24x7/docs/*`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/health`
  - `curl -fsS http://127.0.0.1:8484/v1/router/status`
  - `curl -fsS http://127.0.0.1:8484/v1/voice/status`
  - `curl -fsS http://127.0.0.1:8484/v1/voice/capture/status`
- testes:
  - `go test ./... -count=1`
  - `./scripts/build.sh`
- scripts:
  - `./scripts/install-system-daemon.sh` ou equivalente da worktree
  - `bash ./scripts/voice-capture-smoke.sh`
- fallback:
  - voltar ao binário anterior em `/usr/local/bin/aurelia`
  - `voice_capture_enabled=false`

## Rollout

1. integrar mudanças sem perder o drift da worktree
2. validar suite completa na worktree de deploy
3. instalar/reativar o serviço live
4. validar endpoints e logs

## Rollback

- restaurar binário/serviço anterior
- desabilitar capture/voice no config

## Evidência esperada

- serviço live ativo com o binário novo
- `/health` sem falso `200 ok`
- status routes do gateway e voz vivos
- spool processando transcript real

## Pendências / bloqueios

- a worktree de deploy já tem mudanças locais em arquivos que o rollout também toca
- a unit de serviço live atual é de sistema, não `systemd --user`
