---
name: Code Review
description: Revisão de código sênior focada em padrões Go, tRPC, Zod e Soberania 2026.
phases: [R, V]
---

# 🔍 Code Review: Sovereign Standards 2026

Habilita o Antigravity a realizar revisões de código de nível industrial, garantindo que o PR ou a mudança respeite a governança da Aurélia.

## 🏛️ Checklist de Revisão (Industrial)

### 1. Arquitetura e Contratos
- **Zod-First**: A lógica de validação está no `packages/zod-schemas/`?
- **tRPC**: As rotas estão bem estruturadas e tipadas?
- **Dependency Check**: Existem novas dependências desnecessárias?

### 2. Padrões Go (Sênior)
- **Concorrência**: Goroutines usam `context` corretamente? Há risco de leak?
- **Tratamento de Erro**: Os erros são tratados e logados com `slog` de forma estruturada?
- **Performance**: Existem alocações desnecessárias em loops críticos?

### 3. Soberania e Segurança
- **Secrets**: Existe algum vazamento de chave no diff?
- **Bash over PowerShell**: Se houver scripts, eles são compatíveis com Ubuntu 24.04?
- **Sudo=1 Awareness**: O código que executa comandos de sistema faz validações de segurança?

## 🚀 Workflow de Review
1. **Analise o Diff**: Use `git diff` ou o MCP `github` para ler as mudanças.
2. **Execute Testes**: Antes de aprovar, garanta que `go test ./...` ou `npm test` passe.
3. **Comente com Contexto**: Não diga apenas "mude isso". Diga "mude isso para seguir o ADR-XXXX".

## 📍 Quando usar
- Durante a fase `R` (Review) do ciclo PREVC.
- Para validar contribuições de outros agentes.
- Antes de realizar o merge final na branch `main`.