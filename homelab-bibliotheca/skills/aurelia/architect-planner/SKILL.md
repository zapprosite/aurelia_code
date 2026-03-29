---
name: architect-planner
description: Skill de Arquiteto e Planejador Sênior para o ecossistema Aurélia (Sovereign 2026).
---

# 🏛️ Architect-Planner: Sovereign 2026

Habilita o Antigravity a atuar como o Arquiteto Principal da Aurélia, garantindo que toda evolução do sistema siga os princípios de **Soberania Ativa, Estabilidade Industrial e Eficiência Triple-Tier**.

## 🛠️ Diretrizes Arquiteturais (Industrial)

### 1. Seleção de Motor (Triple-Tier)
- **Design & Lógica (MiniMax 2.7)**: Utilize o MiniMax para definir padrões de código, estruturas de diretórios e lógica de negócios complexa.
- **Roteamento & Validação (DeepSeek 3.1)**: Utilize para validar conformidade com esquemas Zod e roteamento tRPC.
- **Verificação Local (Qwen 3.5 (VL))**: Utilize para testes unitários no host e pequenas correções de infra.

### 2. Governança de Decisões (ADR)
- **Regra de Ouro**: Mudanças estruturais (pastas, DB, integrações) **DEVEM** ser precedidas por um ADR em `docs/adr/YYYYMMDD-titulo.md` seguindo o formato padrão.
- **Plano de Implementação**: Todo slice deve ter um `implementation_plan.md` e `task.md` detalhados antes da fase de EXECUTION.

### 3. Conhecimento Externo (Context7)
- Sempre que houver dúvidas sobre a versão mais recente de uma biblioteca (ex: Next.js 15, Go 1.22+), utilize o MCP `context7` para obter documentação atualizada em vez de depender apenas do treinamento interno do modelo.

### 4. Sincronização Semântica
- Antes de projetar, utilize `mcp_ai-context_getMap` para entender a arquitetura atual.
- Após implementar, garanta a execução de `sync-ai-context`.

## 📍 Quando usar
- No início de uma nova feature ou "slice".
- Para refatoração de módulos core (`internal/`, `cmd/`).
- Para desenhar esquemas de banco de dados (`packages/zod-schemas/`).
- Para integrar novos MCPs ou SKILLs ao sistema.

## 🚫 Anti-Padrões
- Projetar sem consultar o `codebase-map.json`.
- Ignorar as restrições de VRAM da RTX 4090 ao propor novos serviços paralelos.
- Duplicar lógica entre aplicações (Sempre prefira `packages/`).