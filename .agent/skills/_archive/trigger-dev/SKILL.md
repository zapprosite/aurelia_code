---
name: trigger-dev
description: Automação de gatilhos (triggers) e webhooks para integração contínua (CI/CD) e notificações.
---

# ⚡ Trigger-Dev: Sovereign Automation 2026

Habilita a orquestração de gatilhos que disparam ações automáticas no sistema, desde notificações no Telegram até deploys automáticos em CapRover.

## 🛠️ Padrões de Automação
1. **GitHub Actions**: Automatize linting e testes de gateway no push.
2. **Local Hooks**: Use `git hooks` para rodar o `secret-audit.sh` antes de cada commit.
3. **tRPC Triggers**: Dispare atualizações no dashboard `ULTRATRINK` ao detectar mudanças de estado na infraestrutura.

## 📍 Quando usar
- Para automatizar rotinas de backup.
- Para criar integrações entre o Bot e serviços externos.
- Para gerenciar ciclos de vida de build e deploy.