---
type: skill
name: Commit Message
description: Geração de mensagens de commit seguindo o padrão Conventional Commits com detecção de escopo.
skillSlug: commit-message
phases: [E, C]
generated: 2026-03-18
updated: 2026-03-24
status: active
scaffoldVersion: "2.0.0"
---

# 📝 Commit Message: Sovereign Standards 2026

Garante que o histórico do Git seja limpo, legível e semanticamente correto para humanos e agentes.

## 🛠️ Padrão Conventional Commits

Formato: `<type>(<scope>): <subject>`

### Tipos (Types)
- `feat`: Nova funcionalidade (ex: novo Tier no roteador).
- `fix`: Correção de bug (ex: pânico no handler).
- `docs`: Mudanças na documentação.
- `chore`: Manutenção de rotina, sync de contexto ou atualização de deps.
- `refactor`: Mudança de código que não altera comportamento externo.

### Escopos (Scopes)
- `gateway`, `agent`, `voice`, `infra`, `ui`, `context`.

## 🛡️ Regras de Ouro
1. **Subject**: Use imperativo, primeira letra minúscula, sem ponto final (ex: `add minimax tier to policy`).
2. **Body**: Explique o "Motivo" da mudança se não for óbvia.
3. **Footer**: Adicione `Closes #XXX` ou `Ref: ADR-XXXX` se aplicável.

## 📍 Quando usar
- Antes de cada commit via `git commit`.
- Durante os workflows `/git-ship` ou `/git-turbo`.
