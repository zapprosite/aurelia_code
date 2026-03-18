---
name: architect-planner
description: Elite skill for architect-planner.
---



# 🏛️ Skill: Architect-Planner

Garante que toda mudança estrutural seja precedida por um design técnico sólido.

<directives>
1. **Geração de Artefatos**: Sempre crie `implementation_plan.md` e `task.md` para mudanças não triviais.
2. **Integração ai-context**: Utilize o `mcp ai-context buildSemantic` para entender o impacto da mudança no grafo de dependências antes de propor o plano.
3. **Plano de Verificação**: Todo plano deve conter passos claros de teste (Unit, Integration, E2E) e comandos de execução.
4. **Higiene Final**: Lembre o desenvolvedor de executar a sincronização de contexto ao concluir a tarefa.
</directives>

## Padrões de Documentação
- Use YAML frontmatter em todos os planos.
- Use XML tags para delimitar seções de revisão obrigatória por humanos.
