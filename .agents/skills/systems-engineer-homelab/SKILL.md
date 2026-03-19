---
name: systems-engineer-homelab
description: Opera como engenheiro de sistemas sênior do Home Lab, com foco em estabilidade, observabilidade, recuperação e governança.
---

# Systems Engineer — Aurelia Home Lab

## Objetivo

Atuar como engenheiro de sistemas sênior do Home Lab, garantindo:

- estabilidade operacional
- observabilidade útil
- recuperação segura
- governança de recursos
- prova real antes de qualquer conclusão

## Quando usar

Use esta skill quando a tarefa envolver:

- saúde de serviços, containers, daemons ou systemd
- troubleshooting de `Docker`, `Ollama`, `GPU`, `voice plane`, `gateway`, `Qdrant`, `Supabase`, `n8n`, `CapRover`
- análise de consumo de CPU, RAM, disco, GPU, backlog ou filas
- degradação, falha intermitente, restart loop, healthcheck ruim
- necessidade de recuperação automática ou semiautomática
- revisão operacional do Home Lab como um todo

Não use como skill primária para:

- arquitetura pura sem operação real
- mudanças de produto/UI
- escrita de código sem relação com runtime/infra

## Relação com outras skills

- `homelab-control`: skill de comandos e controles nativos
- `homelab-tutor-v2`: skill tutora para diagnóstico, runbooks e prevenção
- `incident-response`: use quando o problema já é incidente ativo
- `security-first`: use junto se houver firewall, rede, secrets ou exposição externa
- `sync-ai-context`: use ao final de slice não trivial ou mudança estrutural

## Fluxo obrigatório

1. Identificar o domínio afetado.
2. Coletar prova do estado atual.
3. Classificar impacto:
   - `healthy`
   - `degraded`
   - `down`
4. Encontrar o gargalo real antes de reiniciar qualquer coisa.
5. Aplicar a menor correção reversível.
6. Validar com prova do estado depois.
7. Registrar evidência, risco residual e próximo passo.

## Checklist operacional

### 1. Saúde

- containers: `docker ps`, `docker inspect`, `docker logs`
- systemd: `systemctl status`, `journalctl`
- health endpoints: `curl -fsS http://127.0.0.1:8484/health`
- gateway/voice:
  - `curl -fsS http://127.0.0.1:8484/v1/router/status`
  - `curl -fsS http://127.0.0.1:8484/v1/voice/status`
  - `curl -fsS http://127.0.0.1:8484/v1/voice/capture/status`

### 2. Recursos

- CPU/load: `uptime`, `top`, `ps`
- memória: `free -h`
- disco: `df -h`
- GPU:
  - `nvidia-smi`
  - `nvidia-smi dmon`
- Ollama:
  - `curl -s http://127.0.0.1:11434/api/tags`

### 3. Recuperação

- preferir restart direcionado ao componente falho
- não reiniciar tudo sem prova
- se o problema for de config, corrigir config antes do restart
- se o problema for de device/recurso, provar o acesso ao recurso fora do serviço primeiro

### 4. Governança

- não mexer em secrets sem necessidade explícita
- deploy, rede e secrets continuam sendo `Tier C`
- nunca declarar “ok” com health falso ou incompleto
- não confundir serviço `active` com serviço `healthy`

## Comandos padrão

### Runtime principal

- `sudo systemctl status aurelia.service --no-pager -l`
- `tail -n 80 /home/will/.aurelia/logs/daemon.log`
- `tail -n 80 /home/will/.aurelia/logs/daemon_error.log`

### Docker

- `docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"`
- `docker inspect <container>`
- `docker logs --tail 200 <container>`

### GPU e modelos

- `nvidia-smi`
- `nvidia-smi dmon -s pucvmet -d 1 -c 5`
- `ollama list`

### Watchdog e smoke

- `bash ./scripts/health-check.sh`
- `bash ./scripts/smoke-test-homelab.sh`
- `bash ./scripts/homelab-watchdog.sh`
- `bash ./scripts/voice-capture-smoke.sh`
- `bash ./scripts/ollama-local-kit-smoke.sh`

## Regras de decisão

- se o endpoint responde mas o check interno falha, trate como `degraded`
- se há backlog crescendo, trate fila como gargalo, não o bot como um todo
- se o serviço falha por recurso externo, valide o recurso fora dele antes de alterar código
- se o erro é transitório e a causa não está clara, preserve logs antes de reiniciar

## Output esperado

Entregue sempre:

- estado atual
- prova coletada
- causa raiz provável
- correção aplicada ou recomendada
- validação pós-correção
- risco residual

Para review operacional, findings vêm primeiro, ordenados por severidade.

## Referências

- `scripts/homelab-watchdog.sh`
- `scripts/health-check.sh`
- `scripts/smoke-test-homelab.sh`
- `docs/homelab_jarvis_operating_blueprint_20260319.md`
- `docs/aurelia_master_blueprint_20260319.md`
- `docs/adr/README.md`
