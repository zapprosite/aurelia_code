---
name: Refactoring
description: Refatoração segura e incremental seguindo os padrões de arquitetura da Aurélia.
phases: [E]
---

# 🛠️ Refactoring: Sovereign Evolution 2026

Habilita o Antigravity a realizar mudanças estruturais no código sem quebrar a funcionalidade existente, priorizando a manutenibilidade e a remoção de dívida técnica.

## 🏛️ Estratégia de Refatoração

### 1. Ciclo de Segurança
- **Vermelho**: Verifique se os testes atuais passam (ou crie-os).
- **Verde**: Realize a mudança mínima necessária.
- **Refatorar**: Limpe o código seguindo os padrões `slog`, `Zod` e `tRPC`.

### 2. Purga de Legado
- Identifique e remova menções a Anthropic/Claude ou motores de nuvem obsoletos.
- Converta lógica ad-hoc para os pacotes compartilhados em `packages/`.
- Substitua scripts PowerShell por scripts Bash puros compatíveis com Ubuntu 24.04.

### 3. Alinhamento Triple-Tier
- Garanta que a refatoração facilite a escolha de modelos (ex: Injeção de dependência para provedores de LLM).

## 📍 Quando usar
- Para reduzir a complexidade de funções "gigantes".
- Para desacoplar módulos (ex: extrair lógica de Gateway de `cmd/`).
- Para padronizar a nomenclatura seguindo a Governança 2026.

## 🛡️ Guardrails
- **No Breaking Changes**: Se a mudança afetar a API pública, crie um ADR primeiro.
- **Atomic Commits**: Refatore pequenas partes por vez.