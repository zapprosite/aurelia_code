---
description: Formaliza ADR obrigatório por slice estrutural e backlog oficial de pendências.
id: 11-adr-slice-contract
---

# 🧾 Regra 11: ADR por Slice

Slices estruturais não podem existir só em conversa ou `TODO` solto.

<directives>
1. **Obrigatoriedade**: Toda mudança estrutural, de arquitetura, runtime, provider, storage, segurança, deploy ou governança deve ter ADR em `docs/adr/` ou entrada explícita no backlog oficial de slices pendentes.
2. **Template oficial**: O formato base é `docs/adr/TEMPLATE-SLICE.md`.
3. **Modo Nonstop**: Slices em execução contínua devem nascer com o par `ADR-YYYYMMDD-slug.md` + `docs/adr/taskmaster/ADR-YYYYMMDD-slug.json`.
4. **Backlog oficial**: Pendências abertas devem ser listadas em `docs/adr/PENDING-SLICES-20260319.md` ou sucessor equivalente.
5. **Links mínimos**: Cada ADR de slice deve linkar `AGENTS.md`, `plan.md`, blueprint relacionado, arquivos afetados, testes esperados e plano de rollout/rollback.
6. **Fechamento**: Um slice só pode ser marcado como concluído quando ADR, JSON de continuidade, testes e `.context/` estiverem sincronizados.
</directives>
