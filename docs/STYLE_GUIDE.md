# STYLE_GUIDE.md

> **Padrão SOTA 2026**

Este guia define os padrões de codificação e estilo para o monorepo `aurelia`.

## 1. Idioma
- **Documentos (.md)**: Português (Brasil) por padrão.
- **Código (Go, TS, Python)**: Inglês para nomenclatura técnica e comentários.

## 2. Padrões por Linguagem
- **Go**: `gofmt` + `golangci-lint`.
- **TypeScript**: `eslint` (strict) + `zod` para validação de contratos.
- **CSS**: Vanilla CSS com design Antigravity (Glassmorphism).

## 3. Versionamento Semântico
- Cada alteração estrutural deve ser acompanhada de uma ADR em `docs/adr/`.
- Uso de **Conventional Commits** para mensagens de versionamento.

---
*Assinado: Aurélia (Soberano 2026)*
*Atualizado: 2026-03-31*
