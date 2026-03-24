---
name: security-first
description: Garante a integridade, privacidade e conformidade de segurança do Home Lab e do Monorepo.
phases: [P, E]
---

# 🛡️ Security First: Sovereign Defense 2026

Esta skill é o pilar de defesa da Aurélia, garantindo que a autonomia (`sudo=1`) não comprometa a segurança do sistema.

## 🧱 Pilares de Segurança

### 1. Gestão de Segredos (Aurelia Vault)
- **Princípio**: Zero Secret in Code. Todo commit deve ser auditado.
- **Workflow**: Se detectar segredos expostos, aborte o push imediatamente e execute o script de revogação.

### 2. Autonomia Supervisionada
- **Auditoria de Comandos**: Monitoramento de comandos `sudo` que alteram configurações globais ou redes.
- **Local First**: Priorize execução local para dados sensíveis ou códigos privados.

### 3. Firewall & Rede
- Verifique periodicamente o `docker-compose.yml` em busca de portas expostas indevidamente (ex: Redis ou DBs abertos para o mundo).

## 📍 Quando usar
- Antes de qualquer `git push`.
- Ao configurar novos serviços Docker.
- Ao realizar o onboarding de novos tokens de API de terceiros.

## 🚫 Anti-Padrões
- Ignorar avisos de linter de segurança.
- Commitar arquivos `.env` ou chaves `.pem`.
- Rodar scripts desconhecidos com `sudo` sem inspeção prévia.