---
name: instance-manager
description: Gerencia silos isolados de dados (Postgres, Qdrant, Obsidian) por aplicação no Homelab.
---

# Instance Manager: SOTA 2026

Esta skill garante a soberania e organização do Homelab através do isolamento estrito de dados entre diferentes aplicações e domínios.

## Comandos

- `/instance create <slug>`: Provisiona uma nova instância completa (Postgres Schema, Qdrant Collection, Folder no Obsidian).
- `/instance list`: Exibe todas as instâncias ativas registradas no `INSTANCE_REGISTRY.json`.
- `/instance delete <slug>`: Arquiva (não deleta fisicamente) uma instância.

## Padrão de Naming

| Componente | Padrão | Exemplo |
|------------|--------|---------|
| Postgres   | `app_<slug>` | `app_aurelia` |
| Qdrant     | `app_<slug>_memory` | `app_aurelia_memory` |
| Obsidian   | `20-apps/<slug>/` | `20-apps/aurelia/` |

## Governança

- Toda nova instância é registrada automaticamente em `docs/governance/INSTANCE_REGISTRY.json`.
- O isolamento deve ser reforçado via RLS no Postgres sempre que possível.

---
**Criado via /master-skill**
