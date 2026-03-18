---
description: Estratégia de Interpretação e Merge (Repository Interpretation & Merge - RIM).
status: Aceito
project: Aurelia Elite
---

# Estratégia RIM (Repository Interpretation & Merge)

Esta estratégia define como o repositório Aurelia (Elite Edition) deve ser interpretado por agentes e como as atualizações devem ser integradas para evitar a degradação da governança local.

## 1. Hierarquia de Interpretação (Agent Read Order)

Para uma inicialização rápida e precisa, os agentes DEVEM ler o repositório nesta ordem:

1.  **`AGENTS.md`**: Define quem manda e quais são as regras de engajamento.
2.  **`README.md`**: Visão geral e "Landing Page" do projeto.
3.  **`.context/docs/codebase-map.json`**: Mapa semântico do código Go.
4.  **`docs/architecture.md`**: Explicação técnica da stack.
5.  **`docs/adr/`**: Histórico de decisões críticas.

## 2. Política de Merge (Dual-Remote)

O Aurelia opera com dois fluxos de atualização:
-   **Upstream (Comunidade)**: `https://github.com/Lordymine/aurelia.git`. Foco em: `/internal`, `/pkg`, `/cmd`, core Go logic.
-   **Template (Elite)**: `https://github.com/zapprosite/Templete-Master.git`. Foco em: `.agents/`, `.context/`, `docs/`, frameworks de agentes.

### Regras de Conflito:
-   Conflitos em pastas de **Governança** (`.agents`, `.context`, `.claude`): **SEMPRE** priorizar o `Template (theirs)`.
-   Conflitos em pastas de **Lógica Core** (`internal`, `cmd`): **SEMPRE** priorizar o `Upstream (ours/upstream)`.
-   Conflitos em **README.md**: Realizar merge manual preservando os diferenciais Elite no topo.

## 3. Eliminação de placeholders

Placeholder genéricos (ex: `[Nome]`, `TODO`, `TBD`) devem ser substituídos por:
-   **Ações Reais**: Se for um `TODO` técnico, deve ser convertido em uma issue ou documentado no `task.md`.
-   **Tags de Agente**: Usar `{{USER_NAME}}` ou `{{AGENT_NAME}}` para permitir interpolação dinâmica se o bot suportar.

## 4. Manutenção de Contexto

-   O comando `ai-context fill` deve ser disparado após qualquer alteração estrutural em `internal/`.
-   Nenhuma nova funcionalidade "Elite" deve ser aceita sem um ADR correspondente em `docs/adr/`.
