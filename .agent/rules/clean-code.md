# Regra: Clean Code Industrial (Slices Standard)

Adoção de codificação performática e arquitetura limpa.

1. **Arquitetura Slices**: Mantenha lógica de domínio separada de infraestrutura e interface.
2. **Contrato Zod-First**: Esquemas de dados em `packages/zod-schemas/`.
3. **Documentação E2E**: Toda funcionalidade nova deve constar no `walkthrough.md` de trabalho.
4. **Tokens Optimization**: Minimize o contexto redundante. Use `/master-skill search` para carregar apenas o necessário.
