# 🎓 Treinamento de Agente: Ecossistema Aurelia (SOTA 2026.1)

Este manual é o contrato operacional para qualquer agente (Claude, Gemini, Antigravity) atuando no repositório `aurelia`.

## 1. O Fluxo de Autoridade (CONTRATO)

Antes de qualquer `run_command` ou `edit_file`, você deve seguir esta hierarquia:
1. **Ler `AGENTS.md`**: A constituição do repositório.
2. **Consultar o Qdrant**: Use `aureliactl skills list` ou busca semântica para encontrar inteligência pronta.
3. **Respeitar as `.rules`**: Localizadas em `.agent/rules/`. Elas são guardrails invioláveis.

## 2. Workflows Obrigatórios

A Aurelia opera em ciclos de fatias (Slices). O uso de slash commands é o padrão:

- **@[/pm]**: Definição de requisitos e PRD.
- **@[/architect]**: Design técnico e Tech Spec (Diagramas em Mermaid).
- **@[/dev]**: Implementação técnica baseada na Tech Spec.
- **@[/qa]**: Validação, testes e auditoria de qualidade.

> [!IMPORTANT]
> Nunca pule para `@[/dev]` sem uma Tech Spec aprovada gerada pelo `@[/architect]`.

## 3. Guardrails e Segurança

- **Zero Hardcode**: Auditoria proativa de segredos. Nunca exponha chaves no Git.
- **Zod-First**: Esquemas de dados devem residir exclusivamente em `packages/zod-schemas/`.
- **Alog**: Utilize o motor de logging em Go (`internal/purity/alog`) para novos serviços.

## 4. Integração Obsidian

Este guia e todas as regras estão espelhados na sua Vault do Obsidian. 
- Use o **Obsidian CLI** para buscas rápidas offline.
- O script oficial de sincronia é `scripts/ops/obsidian-sync.sh`.

---
*Assinado: Aurélia (Arquiteta Líder) & Antigravity (Operador SOTA 2026.1)*
