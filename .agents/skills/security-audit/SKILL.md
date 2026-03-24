---
type: skill
name: security-audit
description: Realiza auditorias de segurança proativas e varredura de vulnerabilidades no Homelab.
skillSlug: security-audit
phases: [O, V]
generated: 2026-03-22
updated: 2026-03-24
status: active
scaffoldVersion: "2.0.0"
---

# 🚔 Security Audit: Sovereign Surveillance 2026

Habilita a detecção ativa de falhas no sistema e a conformidade industrial com os padrões 2026.

## 📋 Checklist de Auditoria
- **Docker Leak**: Verifique se algum container tem privilégios excessivos.
- **Secrets Audit**: Varredura por strings padrão de tokens (GitHub, OpenAI, etc).
- **Update Checks**: Monitoramento de CVEs críticas no Ubuntu 24.04.

## 🛠️ Comandos de Auditoria
- `scripts/security-audit.sh`: Script mestre de verificação.
- `docker exec -it <container> lynis audit system`: Auditoria interna de containers críticos.
- `netstat -plnt`: Inspeção de sockets abertos.

## 📍 Quando usar
- Mensalmente para manutenção de rotina.
- Após qualquer incidente suspeito reportado pelo Bot.
- Após a instalação de novos softwares de sistema.
