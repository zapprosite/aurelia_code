---
description: Bootstrap universal do repositório multi-agente com Antigravity + Claude + Codex.
---

# Workflow: templet-mult-clouds

Passos:
1.  **Inspecionar**: Analisar a estrutura atual do repositório.
2.  **Validar**: Verificar se as ferramentas (claude, git, node, codex) estão disponíveis.
3.  **Instalar Skills**: Instalar as skills globais recomendadas via `npx skills`.
4.  **Autoridade**: Criar ou atualizar o `AGENTS.md` como contrato global.
5.  **Adaptadores**: Criar arquivos finos `CLAUDE.md`, `GEMINI.md` e `CODEX.md`.
6.  **Regras**: Estabelecer as regras em `.agent/rules/`.
7.  **Skill Local**: Criar a skill do projeto em `.agent/skills/template-mult-clouds/`.
8.  **Plan**: Consultar `plan.md` e `AGENTS.md` para direcionamento.
9.  **Docs Context**: Criar documentação técnica em `.context/docs`.
10. **Resumo**: Gerar relatório final de inicialização.
