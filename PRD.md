---
description: Adaptador de Autotransformação PRD SOTA - Arquiteto & DevOps Sênior.
last-updated: 2026-03-17
---

# 🪄 Magic PRD: SOTA Requirements Engine

<system_prompt>
VOCÊ É UM ARQUITETO DE SISTEMAS E ESPECIALISTA EM DEVOPS SÊNIOR.
SUA MISSÃO É transformar o `DRAFT_INPUT` abaixo em um PRD técnico, denso em lógica e pronto para execução multi-agente.

### DIRETRIZES DE EXECUÇÃO:
1. **Lógica de Camadas**: Divida a implementação entre infraestrutura (IaC), backend (contracts) e frontend (UX).
2. **Distribuição (@agents)**: Atribua tarefas com base nas especialidades em `.claude/agents/`.
3. **Sem Alucinação**: Se a ideia for vaga, aplique as melhores práticas de mercado (SOTA) para preencher a lacuna tecnicamente (ex: usar Redis para cache, Zod para validação).
4. **Formatação**: Retorne APENAS o conteúdo do PRD formatado em Markdown Profissional, preservando o cabeçalho YAML.
5. **Data**: Utilize como data de atualização: 2026-03-17.

### ESTRUTURA OBRIGATÓRIA:
- **Phase 1 (Planning)**: `@researcher` + `@planner` (Análise de impacto, Tech Spec, ADR).
- **Phase 2 (Build)**: `@implementer` (Worktree, Feature, Tests).
- **Phase 3 (Audit)**: `@reviewer` (Security, Performance, Diffs).
- **Phase 4 (Deploy/Sync)**: `@mcp-operator` (Context Sync, CI/CD, Documentation).
</system_prompt>

---

**DRAFT_INPUT**: [SUA IDEIA BRUTA AQUI]

---

# 🎯 PRD: [O AGENTE GERARÁ O TÍTULO TÉCNICO AQUI]

> **Visão Geral**: Descreva brevemente o objetivo final e o valor de negócio.

---

## 🏛️ 1. Objetivo e Requisitos
- **Meta Principal**: [Descreva o que deve ser alcançado]
- **Requisitos Chave**:
  - [ ] Requisito 1
  - [ ] Requisito 2

---

## 🏗️ 2. Distribuição por Fases (Phases & Tasks)

### 🟢 Phase 1: Research & Planning
**Primary Agents**: `@researcher`, `@planner`

- [ ] **Task 1.1**: Analisar o estado atual e identificar lacunas. ➡️ `@researcher`
- [ ] **Task 1.2**: Criar Tech Spec e Plano de Implementação. ➡️ `@planner`
- [ ] **Task 1.3**: Registrar ADR se houver mudança arquitetural. ➡️ `@planner`

### 🔵 Phase 2: Implementation (Worktree)
**Primary Agents**: `@implementer`

- [ ] **Task 2.1**: Configurar ambiente/worktree isolado. ➡️ `@implementer`
- [ ] **Task 2.2**: Codificar a funcionalidade seguindo os padrões do repo. ➡️ `@implementer`
- [ ] **Task 2.3**: Refatorar código existente se necessário. ➡️ `@implementer`

### 🟡 Phase 3: Verification & Security
**Primary Agents**: `@reviewer`

- [ ] **Task 3.1**: Auditoria de código e diff. ➡️ `@reviewer`
- [ ] **Task 3.2**: Executar testes unitários e de integração. ➡️ `@reviewer`
- [ ] **Task 3.3**: Validar conformidade com as 10 Regras Elite. ➡️ `@reviewer`

### 🔴 Phase 4: Context Hygiene & Delivery
**Primary Agents**: `@mcp-operator`

- [ ] **Task 4.1**: Sincronizar Contexto via MCP `ai-context`. ➡️ `@mcp-operator`
- [ ] **Task 4.2**: Gerar Walkthrough de entrega. ➡️ `@mcp-operator`
- [ ] **Task 4.3**: Merge seguro na branch principal. ➡️ `@mcp-operator`

---

## 🛡️ Definição de Pronto (DoR/DoD)
- [ ] Código sem segredos expostos.
- [ ] Documentação atualizada em Português (BR).
- [ ] Contexto sincronizado e denso.

---
*Este PRD serve como o "Contrato de Voo" para a frota de agentes.*
