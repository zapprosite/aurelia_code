---
type: skill
name: Git Unblocker
description: Workflow rápido para invocar a skill de destravamento de Git zumbi.
skillSlug: git-unblocker
phases: [E]
generated: 2026-03-20
updated: 2026-03-24
status: active
scaffoldVersion: "2.0.0"
---

# 🔓 Git Unblocker: Sovereign Recovery 2026

Pequena skill utilitária para resolver travamentos de índice (.git/index.lock) ou processos zumbis que impedem operações de commit/push no host Ubuntu.

## 🛠️ Procedimento de Autocura
1. **Identificação**: Se o comando `git` retornar "Another git process seems to be running", invoque esta skill.
2. **Ação**:
   - `rm -f .git/index.lock`: Remoção do lock de índice.
   - `ps aux | grep git`: Matar processos de diff ou status travados se necessário.
3. **Verificação**: `git status` para confirmar a liberação.

## 📍 Quando usar
- Quando o Antigravity ou outros agentes ficarem travados em operações de Git.
- Após falhas de rede durante um `push` ou `pull`.
