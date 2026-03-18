---
description: Garante segurança e isolamento durante a implementação técnica.
id: 04-worktree-isolation
---

# 🏗️ Regra 04: Isolamento de Worktree

Toda implementação técnica não trivial deve ocorrer em um ambiente isolado.

<directives>
1. **Segurança**: Nunca edite arquivos críticos diretamente na branch de trabalho principal se houver risco de quebra.
2. **Branching**: Use o padrão `feat/`, `fix/` ou `research/`.
3. **Persistência**: O plano de execução deve ser sincronizado com o repositório principal durante os handoffs.
4. **Higiene**: Remova worktrees temporários imediatamente após o merge/conclusão.
</directives>
