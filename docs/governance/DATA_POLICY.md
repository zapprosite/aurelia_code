# DATA_POLICY.md

> **Padrão Pinned Data Center 2026**

Este documento define a governança de dados para o ecossistema Aurélia.

## 1. Soberania Local
- **Processamento**: Todos os dados críticos (transcrições, logs, estados de banco) residem no Homelab Ubuntu.
- **Isolamento**: Uso de instâncias isoladas de Redis e Qdrant por aplicação.

## 2. Gateways Seguros
- **Rede**: Acesso externo apenas via Cloudflare Tunnel com autenticação de portões (Porteiro).
- **Secrets**: Proibido o armazenamento de tokens em arquivos `.md`. Uso mandatório de `.env` (Zero Drift).

## 3. Observabilidade
- **Audit**: Logs estruturados em `/logs/` com rotação periódica.
- **Monitoramento**: Checagem horária de integridade via `aurelia-system-api`.

*Atualizado: 2026-03-31*
