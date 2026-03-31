---
description: Prioriza o uso de inteligência encapsulada (skills) sobre improvisação.
id: 09-skills-usage
---

# ⚡ Regra 09: Priorização de Skills

Antes de criar uma solução do zero, use as ferramentas prontas.

<directives>
1. **Descoberta**: Verifique `.agents/skills/` para habilidades de arquitetura, segurança ou pesquisa profunda.
2. **Padronização**: Skills garantem que todos os agentes produzam código no mesmo estilo e rigor técnico.
3. **Contribuição**: Se criar um padrão novo e útil, proponha-o como uma nova Skill.
</directives>
---
description: Mantém a limpeza e utilidade dos documentos gerados.
id: 10-artifact-discipline
---

# 📁 Regra 10: Disciplina de Artefatos

Documentos devem ser úteis, enxutos e sem placeholders.

<directives>
1. **Anti-Placeholder**: É proibido criar ou manter arquivos com "TODO" ou conteúdo genérico. Se não tem informação, não crie o arquivo.
2. **Links Reais**: Todos os links de arquivos em artefatos devem ser funcionais e absolutos.
3. **Poda de Contexto**: Remova artefatos temporários ao concluir tarefas para evitar ruído em sessões futuras.
</directives>
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
