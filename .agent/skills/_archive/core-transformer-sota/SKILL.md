---
name: core-transformer-sota
description: Orquestrador de transformação profunda (Polish + Hardening) baseado no padrão Pinned Data Center 2026.
---

# Core Transformer SOTA 2026 ⚙️

Esta skill realiza a leitura analítica de cada "core" (serviço/módulo) e aplica endurecimento industrial e polish de governança.

## 🏛️ Padrão Pinned Data Center (SOTA 2026.2)

1. **Version Pining**: Proibido o uso de `:latest`. Todas as imagens em `docker-compose.yml` devem ter tags fixas.
2. **Secret Sovereignty**: Transição de `.env` para `EnvironmentFile` em serviços systemd (isolamento por serviço).
3. **Network Mapping**: Todos os serviços devem respeitar o `NETWORK_MAP.md` global.
4. **Static Binary**: Garantir compilação estática (`CGO_ENABLED=0`) para portabilidade total.

## 🔄 Fluxo de Operação

### Passo 1: Auditoria Analítica (LLM Top)
- Invocar `repo-health-audit`.
- Analisar logs de build e teste.
- Identificar débitos técnicos de governança.

### Passo 2: Hardening (Endurecimento)
- Aplicar permissões restritas (600) em segredos.
- Configurar healthchecks universais.
- Isolar volumes por core.

### Passo 3: Polish de Interface
- Refinar ADRs associados.
- Atualizar documentação técnica conforme `CONSTITUTION.md`.

## 🛠️ Comandos Sugeridos
- `/core-transformer-sota run` — Executa o ciclo completo em um core específico.
- `/core-transformer-sota audit` — Auditoria profunda pré-transformação.
