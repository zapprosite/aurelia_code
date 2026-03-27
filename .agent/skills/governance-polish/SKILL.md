---
name: Governance Polish — Industrial Homelab Governance
description: Orquestra o polish de governança do homelab em fases incrementais (Sovereign 2026).
phases: [P, E, V]
---

# 🛡️ Governance Polish: Sovereign Guard 2026

Esta skill governa a evolução da infraestrutura e segurança da Aurélia, garantindo que o ambiente permaneça robusto, observável e em conformidade com o **Plano Mestre 2026**.

## 🏗️ Fases de Governança (Triple-Tier)

### 1. Hardening de Infra (Sudo=1)
- **Secrets**: Migração total para `EnvironmentFile` do systemd e auditoria semanal via `scripts/secret-audit.sh`.
- **Vault**: Consolidação do KeePassXC (deadline: 2026-03-27).
- **Observabilidade**: Dashboards de saúde via tRPC refletindo o status real da CPU/GPU.

### 2. Soberania de Modelos
- **Tiering**: Fiscalização do uso correto de MiniMax (Arquitetura) vs DeepSeek (Roteamento) vs Gemma3 (Local).
- **VRAM**: Manutenção do teto de memória para evitar pânico de kernel ou OOM na RTX 4090.

### 3. Documentação e Evolução
- **Slices**: Garantir que todo novo "slice" tenha ADR, Plano e Walkthrough.
- **Contexto**: Uso obrigatório da skill `/sync-ai-context` após mudanças estruturais.

## 📊 Status de Implementação (Roadmap)

| Marco | Status | Obs |
|-------|--------|-----|
| Secrets Env Overrides | ✅ | Implementado em `internal/config` |
| Roadmap Mestre Sincronizado | ✅ | Ver docs/adr/README.md |
| Secret Audit (Crontab) | ✅ | Ativo |
| Vault KeePassXC | ⏳ | Blocked: Humano (Março 27) |
| Triple-Tier Router | ✅ | Estabilizado Março 24 |

## 📍 Quando usar
- Para auditar a segurança do Homelab.
- Para verificar o progresso do roadmap mestre de desenvolvimento.
- Antes de grandes atualizações de sistema operacional ou drivers.

## 🛡️ Referências Obrigatórias
- [AGENTS.md](../../../AGENTS.md)
- [AURELIA-AUTHORITY-DECLARATION.md](../../../docs/governance/AURELIA-AUTHORITY-DECLARATION.md)
- [ADR Index](../../../docs/adr/README.md)
