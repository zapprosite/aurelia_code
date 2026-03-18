---
description: Define os limites de permissão automática baseados em risco.
id: 03-tiers-autonomy
---

# 🛡️ Regra 03: Tiers de Autonomia (Risco)

A execução é governada por níveis de risco definidos no `PRD.md`.

<directives>
1. **Tier A (Read-only)**: Auto-approve em 100% dos casos para leitura e análise.
2. **Tier B (Local Edit)**: Auto-approve permitido apenas em **Worktrees**. Edições na `main` são proibidas sem revisão.
3. **Tier C (High-risk)**: Aprovação Humana OBRIGATÓRIA para:
   - Modificações em rede/firewall.
   - Gestão de segredos e chaves API.
   - Operações de deploy ou exclusão em massa.
</directives>
