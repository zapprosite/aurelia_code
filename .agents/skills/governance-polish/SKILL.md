---
type: skill
name: Governance Polish — Industrial Homelab Governance
description: Orquestra o polish de governança do homelab em fases incrementais.
skillSlug: governance-polish
phases: [P, E, V]
generated: 2026-03-19
updated: 2026-03-21
status: active
scaffoldVersion: "2.0.0"
---

# Skill: governance-polish

## Propósito

Orquestrar a implementação de governança industrial do homelab:
- **Secrets:** Consolidar `secrets.env` com systemd `EnvironmentFile`, env overrides em Go
- **Vault:** KeePassXC com masterkey em hardware token (deadline: 2026-03-27)
- **Auditoria:** secret-audit.sh semanal via crontab
- **Documentação:** Roadmap, links, índice de ADRs sincronizados

## Status Atual

| Fase | Status |
|------|--------|
| Secrets.env + systemd EnvironmentFile | ✅ Concluído |
| Env Overrides em `internal/config/config.go` | ✅ Concluído |
| Roadmap mestre sincronizado | ✅ Concluído |
| Links verificados no REPOSITORY_CONTRACT | ✅ Concluído |
| Índice de ADRs (`docs/adr/README.md`) | ✅ Concluído |
| KeePassXC vault | ⏳ Aguardando humano (2026-03-27) |
| Secret-audit no crontab semanal | ✅ Concluído |

## Referência

- Roadmap Mestre: [ADR-20260320-roadmap-mestre-slices.md](../../docs/adr/ADR-20260320-roadmap-mestre-slices.md)
- Secrets: [SECRETS.md](../../docs/governance/SECRETS.md)
- Authority: [AURELIA-AUTHORITY-DECLARATION.md](../../docs/governance/AURELIA-AUTHORITY-DECLARATION.md)

## Integrações

- AGENTS.md: Respeitar hierarquia de autoridade
- GEMINI.md: Coordenar handoff entre agentes
- REPOSITORY_CONTRACT.md: Índice de governança

## Smoke Tests

```bash
# Links verificados
grep -r '\[.*\](.*\.md)' docs/governance/ | head -5

# Secret-audit
bash scripts/secret-audit.sh  # exit 0

# Binários na raiz
ls aurelia-elite 2>/dev/null && echo "FAIL" || echo "OK"
```
