---
type: skill
name: Governance Polish — Industrial Homelab Governance
description: Execute ADR-20260319-Polish-Governance-All in 4 phases. Phase 1 (CRITICAL) requires human action; Phases 2-4 automated by codex.
skillSlug: governance-polish
phases: [P, E, V]
generated: 2026-03-19
status: unfilled
scaffoldVersion: "2.0.0"
---

# Skill: governance-polish

## Propósito

Orquestrar a implementação de governança industrial do homelab em 4 fases:
- **Fase 1 (CRITICAL):** Humano — vault KeePassXC, migração secrets, shred plaintext, postgres password
- **Fase 2 (HIGH):** Codex — deletar backups, refatorar MCP, schema registry, rotation policy
- **Fase 3 (MEDIUM):** Codex — health checks, backup verification, incident playbook, observability
- **Fase 4 (LOW):** Codex — cleanup, UFW, compliance matrix, audit scripts

## Invocação

```bash
# Fase 1: Mostrar checklist para humano
/governance-polish --phase 1 --show-checklist

# Fase 2: Executar após Fase 1 concluída
/governance-polish --phase 2 --execute

# Verificar status geral
/governance-polish --status

# Ver próximos passos
/governance-polish --next-action
```

## Referência

- ADR: [ADR-20260319-Polish-Governance-All](../../docs/adr/ADR-20260319-Polish-Governance-All.md)
- JSON Taskmaster: [taskmaster JSON](../../docs/adr/taskmaster/ADR-20260319-Polish-Governance-All.json)
- Tutorial: [KeePassXC Tutorial](../../keepassxc-tutorial.html)

## Integrações

- CLAUDE.md: Consultar para contexto de execução
- CODEX.md: Executar Fases 2-4
- GEMINI.md: Coordenar handoff entre agentes
- AGENTS.md: Respeitar hierarquia de autoridade
- plan.md: Atualizar com progresso de cada fase

## Critério de Sucesso

- Fase 1: Vault criado, plaintext shredded, postgres password seguro, humano confirma
- Fase 2: Backups deletados, MCP refatorado, schema registry escrito
- Fase 3: Health checks rodando, backup verification ativo, playbook escrito
- Fase 4: UFW documentado, compliance matrix completa, audit scripts rodando

## Smoke Tests

```bash
# Após Fase 1
ps aux | grep -i password | grep -v grep  # vazio

# Após Fase 2
bash scripts/secret-audit.sh  # exit 0

# Após Fase 3
curl -s localhost:9090/api/v1/rules | jq '.status'  # success

# Após Fase 4
bash scripts/governance-audit.sh  # all green
```

## Fallback/Rollback

- Cada fase é independente
- `crontab -r` (com backup) se necessário
- `sudo ufw disable` para reverter UFW
- Consulte .context/plans/ para estado detalhado
