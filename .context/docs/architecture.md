# Architecture Overview

A arquitetura do workspace é baseada no modelo **BMAD** (Business, Model, Architect, Developer) com foco em separação de preocupações entre agentes.

## Componentes Chave
- **.agents/workflows/**: Define os comandos slash interativos.
- **.agents/rules/**: Define os limites éticos e técnicos (Guardrails).
- **.claude/agents/**: Define as personalidades e ferramentas dos subagentes locais.
- **docs/adr/**: Mantém o histórico de decisões estruturais.

## Fluxo de Dados de Contexto
As LLMs iniciam a leitura por \`AGENTS.md\`, descem para as \`rules/\` e usam o \`.context/\` como cache de memória efêmera para a tarefa atual.
