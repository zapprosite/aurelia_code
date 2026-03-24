---
name: pr-review
description: Analisa Pull Requests contra padrões de equipe e melhores práticas de segurança e performance.
---

# 🔎 PR Review: Sovereign Quality 2026

Habilita o Antigravity a realizar auditorias exaustivas em Pull Requests, garantindo que o código que chega à `main` seja de nível industrial.

## 🏛️ Checklist de Auditoria (Sênior)

### 1. Conformidade Arquitetural
- O código segue os ADRs relevantes?
- Existe duplicidade de lógica? (Check `packages/`)
- O contrato de API foi afetado?

### 2. Segurança & Segredos
- **Secret Audit**: Varredura por tokens ou chaves no diff.
- **OWASP**: Verificação de injeções e falhas de lógica.

### 3. Testes e Cobertura
- Foram adicionados novos testes?
- Os testes existentes continuam passando?
- A cobertura de código é aceitável (> 80%)?

### 4. Performance & Hardware
- O novo código impacta significativamente a CPU/VRAM?
- Existem goroutines sem `context` ou fechamento?

## 📍 Quando usar
- Antes de aprovar ou mergear qualquer PR.
- Para fornecer feedback técnico a outros agentes ou colaboradores.
- Durante a fase final de integração de um "slice" complexo.

## 🛡️ Guardrails
- **Sem aprovação cega**: Nunca aprove um PR sem entender 100% da mudança.
- **Higiene**: Foque no código, não no autor. Seja objetivo e técnico.
