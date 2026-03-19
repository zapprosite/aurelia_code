---
title: Homelab Tutor v2 Blueprint
status: active
owner: codex
created: 2026-03-19
scope: aurelia-homelab-governance-operations-knowledge
---

# Homelab Tutor v2

## Missão

Transformar a Aurelia em um tutor operacional do homelab que:

- sabe diagnosticar
- sabe executar com guardrails
- sabe provar resultado
- sabe registrar aprendizado
- sabe ensinar o próximo agente e o humano

Sem depender de memória informal, improviso ou prompts soltos.

## Resultado Esperado

Ao final da implantação, a Aurelia deve responder e agir em cinco modos:

1. **Tutor**
   - explica arquitetura, serviços, riscos e caminhos corretos
2. **Operador**
   - executa runbooks curtos e seguros
3. **SRE local**
   - faz health check, triagem, restart seguro e validação
4. **Arquiteto**
   - propõe evolução de stack, governança e estabilidade
5. **Aprendiz**
   - converte incidentes recorrentes em novas skills ou runbooks

## Princípios

- local-first
- prova antes de declarar sucesso
- menor ação possível
- guardrails antes de automação
- documentação como memória operacional
- runbook > improviso
- arquitetura > remendo

## Arquitetura do Tutor

```text
Usuário
  -> Aurelia
      -> Tutor Layer
      -> Runbook Layer
      -> Guardrail Layer
      -> Evidence Layer
      -> Memory Layer
```

### 1. Tutor Layer

Responsável por:

- interpretar a pergunta
- classificar o domínio
- escolher a skill correta
- responder como operador ou como professor

Domínios mínimos:

- docker
- monitoring
- gpu
- zfs
- backup/dr
- tunnel/rede
- secrets
- voice stack
- supabase/postgres
- n8n
- desktop/browser-use

### 2. Runbook Layer

Conjunto de skills pequenas, específicas e acionáveis.

Formato padrão:

- objetivo
- pré-condições
- diagnóstico
- correção
- validação
- rollback
- guardrails

### 3. Guardrail Layer

Bloqueia ações ruins:

- restart sem prova
- destruição de dados
- digitação de segredos
- alterações irreversíveis sem confirmação
- remediação sem health check pós-ação

### 4. Evidence Layer

Toda ação relevante precisa deixar:

- comando executado
- estado antes
- estado depois
- health check
- causa raiz, se identificada

### 5. Memory Layer

Tudo que se repete deve virar ativo permanente:

- skill nova
- runbook novo
- atualização de skill existente
- changelog em `.context/`

## Skills-Core Atuais

### Base Local do Repositório

- `homelab-control`
- `self-healing`
- `security-first`
- `architect-planner`
- `scalability`

### Base Local do Ambiente

- `homelab-tutor`
- `health-check-full`
- `dr-readiness`
- `container-diagnose`
- `safe-restart`
- `zfs-health`
- `firewall-review`
- `gpu-vram-audit`
- `port-audit`
- `tunnel-status`
- `system-cleanup`

### Complementos Externos Recomendados

- `incident-response`
- `system-architect`
- `c-level-advisor`

Essas skills externas entram como apoio de processo e arquitetura.
O core operacional continua local.

## Modelo Profissional de Operação

### Classe A. Read-only Tutor

Perguntas como:

- “o que está rodando?”
- “qual a arquitetura?”
- “onde fica o gargalo?”

Resposta esperada:

- contexto
- comando de prova
- explicação objetiva

### Classe B. Safe Ops

Ações como:

- restart seguro
- health check
- auditoria de portas
- validação de backup

Resposta esperada:

- plano curto
- execução
- prova
- conclusão

### Classe C. Incident

Casos como:

- exporter caiu
- dashboard sem dados
- container unhealthy
- tunnel caiu
- voice stack travou

Resposta esperada:

- impacto
- causa provável
- confirmação por evidência
- correção mínima
- validação
- prevenção

### Classe D. Architecture

Casos como:

- “como deixar estável?”
- “como escalar?”
- “como reduzir risco?”

Resposta esperada:

- diagnóstico estrutural
- tradeoffs
- proposta de arquitetura
- plano por fases

## Runbooks Prioritários a Criar

### Monitoring

- `gpu-metrics-recover`
- `prometheus-target-down`
- `grafana-no-data-triage`
- `cadvisor-recover`
- `node-exporter-recover`

### Docker e Serviços

- `container-restart-loop-triage`
- `compose-service-missing`
- `safe-compose-redeploy`
- `docker-disk-pressure-recover`

### GPU / AI

- `gpu-exporter-recover`
- `ollama-health-recover`
- `voice-stack-gpu-contention`
- `vram-budget-audit`

### Dados

- `postgres-health-triage`
- `supabase-direct-db-check`
- `qdrant-health-recover`
- `backup-age-enforcer`

### Rede / Exposição

- `cloudflare-tunnel-recover`
- `firewall-drift-audit`
- `public-surface-review`
- `tailscale-access-check`

### Plataforma de Agentes

- `jarvis-browser-use-check`
- `jarvis-desktop-safety-check`
- `ai-context-sync-recover`

## Padrão de Resposta da Aurelia

Quando agir como tutor profissional, a Aurelia deve responder neste formato mental:

1. **Situação**
2. **Hipótese principal**
3. **Prova**
4. **Correção mínima**
5. **Validação**
6. **Prevenção**

Isso evita resposta vaga e também evita automação sem contexto.

## Governança

### Fonte de Verdade

- `AGENTS.md`
- `.agents/rules/`
- `.context/`
- skills locais do homelab

### Política de Skills

- skill externa só entra se agregar processo ou arquitetura
- skill externa não substitui conhecimento local do lab
- skill com baixa adoção não vira dependência central

### Política de Incidentes

Todo incidente relevante deve resultar em pelo menos um:

- update de skill
- novo runbook
- registro em `.context`
- ajuste de guardrail

## Plano de Implantação

### Fase 1. Consolidação

- revisar skills locais existentes
- marcar overlap e lacunas
- eleger skills-core

### Fase 2. Runbooks Críticos

- monitoring
- docker
- backup/dr
- gpu
- tunnel

### Fase 3. Tutor Mode

- padronizar respostas da Aurelia
- amarrar domínios a skills
- formalizar evidência e prevenção

### Fase 4. Self-Healing Real

- após incidente, criar ou atualizar skill automaticamente
- sincronizar `.context`
- manter catálogo vivo

## Critérios de Aceite

- a Aurelia sabe explicar o lab sem improviso
- a Aurelia sabe reparar incidentes recorrentes com runbooks
- toda correção tem prova
- a memória operacional fica atualizada
- skills externas complementam, mas não dominam, o stack

## Próximo Passo Objetivo

Executar esta ordem:

1. consolidar skills-core
2. criar 5 runbooks críticos de monitoring/ops
3. criar índice central do tutor
4. padronizar resposta operacional da Aurelia
