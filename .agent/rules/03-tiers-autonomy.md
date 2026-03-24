---
description: Define os limites de permissão automática baseados em risco.
id: 03-tiers-autonomy
---

# 🛡️ Regra 03: Tiers de Autonomia (Risco)

A execução é governada por níveis de risco definidos em `AGENTS.md` § 5.

> **Diretiva do Humano (2026-03-20):** Autonomia total habilitada (`sudo=1`).
> O operador mantém backup completo e aceita os riscos operacionais.

<directives>
1. **Tier A (Read-only)**: Auto-approve 100%.
2. **Tier B (Local Edit)**: Auto-approve 100%. Preferência por Worktrees para isolamento.
3. **Tier C (High-risk)**: Auto-approve com **log obrigatório** para:
   - Modificações em rede/firewall.
   - Gestão de segredos e chaves API.
   - Operações de deploy ou exclusão em massa.
   - Comandos `sudo`.
4. **Segurança Compensatória**:
   - Todo comando `sudo` deve ser registrado em log estruturado.
   - Dry-run sempre que possível para `docker-compose` e scripts bash.
   - Auditoria de segredos antes de `git push` continua ativa.
</directives>
