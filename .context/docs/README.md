# Documentation Index

Bem-vindo à base de conhecimento da Aurélia. Comece com a visão geral do projeto e explore os guias específicos conforme necessário.

## Core Guides
- [Visão Geral do Projeto](./project-overview.md)
- [Workflow de Desenvolvimento](./development-workflow.md)
- [Estratégia de Testes](./testing-strategy.md)
- [Guia de Ferramentas & Produtividade](./tooling.md)

## Snapshot do Repositório (Contexto Semântico)

### Arquitetura
- **Modelos**: `packages/zod-schemas`
- **Utils**: `frontend/src/lib`
- **Componentes**: `frontend/src/components/ui`, `frontend/src/sidebar`, `frontend/src/dashboard`

### API Pública e Símbolos Principais
- `SentinelEvent` (Zod Schema) @ `packages/zod-schemas/index.ts`
- `useSystemMetrics` (React Hook) @ `frontend/src/hooks/useSystemMetrics.ts`
- `HomelabTab` (Dashboard) @ `frontend/src/components/dashboard/HomelabTab.tsx`
- `wakeword` (Voice Gateway) @ `cmd/voice-gateway/wakeword.py`
- `skill-indexer` (Core) @ `cmd/skill-indexer/main.py`

## Mapa de Documentação
| Guia | Arquivo | Entradas Primárias |
| --- | --- | --- |
| Visão Geral do Projeto | `project-overview.md` | Roadmap, README, notas de stakeholders |
| Workflow de Desenvolvimento | `development-workflow.md` | Regras de branching, config CI, guia de contribuição |
| Estratégia de Testes | `testing-strategy.md` | Configs de teste, gates de CI, suites instáveis |
| Guia de Ferramentas | `tooling.md` | Scripts CLI, configs de IDE, workflows de automação |

---
*Última sincronização completa do codebase-map: 2026-04-01*
