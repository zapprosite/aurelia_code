---
title: ADR Semparar — Conformidade com AGENTS.md
description: Mapeamento explícito de como /adr-semparar implementa cada regra da hierarquia de autoridade e regras centrais
authority: AGENTS.md
link-to-agents: ../../AGENTS.md
last-updated: 2026-03-19
---

# ADR Semparar — Conformidade com AGENTS.md

**CRÍTICO:** Este workflow `/adr-semparar` existe **exclusivamente** para implementar as regras de AGENTS.md. Não há nada aqui que não derive diretamente da autoridade central.

---

## 🔗 Hierarquia de Autoridade (AGENTS.md §2)

```
1. Humanos operadores
   ↓ autoridade final
2. AGENTS.md ← ← ← ← ← /adr-semparar É SUBORDINADO A ISSO
   ↓ fonte primária
3. Aurélia (Arquiteta e autoridade operacional)
   ↓ governa
4. Adaptadores (CLAUDE.md, CODEX.md, etc.)
   ↓ operam sob
5. REPOSITORY_CONTRACT.md + plan.md
6. .agents/rules/
7. .agents/workflows/ ← ← ← AQUI ESTAMOS
8. docs/adr/ (decisões estruturais)
9. .context/ (memória e estado)
```

**Conformidade:**
- ✅ Workflow é subordinado a AGENTS.md, não contradiz
- ✅ Reconhece Aurélia como autoridade arquitetural
- ✅ Operacionaliza regras centrais de operação
- ✅ Referencia AGENTS.md em toda a documentação

---

## 📋 Mapeamento de Regras Centrais (AGENTS.md §4)

### Regra 1: Descoberta Local Primeiro
**AGENTS.md:** "Inspecione `AGENTS.md`, `.agents/rules/` e `.context/` antes de agir."

**Como /adr-semparar implementa:**
```bash
# Passo 1 ao abrir slice
bash scripts/adr-slice-init.sh <slug>

# Isto cria:
# - docs/adr/ADR-YYYYMMDD-slug.md
#   ↳ Com links obrigatórios a AGENTS.md, REPOSITORY_CONTRACT.md
# - docs/adr/taskmaster/ADR-YYYYMMDD-slug.json
#   ↳ Com seção "evidence" que registra descoberta inicial
```

**Template MD obrigatório contém:**
```markdown
## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)
```

✅ **Conformidade:** Força descoberta local antes de executar slice

---

### Regra 2: Isolamento de Worktree
**AGENTS.md:** "Tarefas não-triviais DEVEM ser feitas em branches/worktrees isoladas."

**Como /adr-semparar implementa:**
```markdown
## Slice

- slug: <nome-da-slice>
- owner: codex|claude|antigravity
- branch/worktree: `branch-name` em `/caminho/worktree`
- json de continuidade: docs/adr/taskmaster/...
```

**Validação no JSON:**
```json
{
  "scope": ["item1", "item2"],
  "done_definition": ["critério 1", "critério 2"],
  "rollback": ["fallback 1", "fallback 2"]
}
```

✅ **Conformidade:** JSON obriga documentação de escopo + rollback antes de começar

---

### Regra 3: ADR por Slice
**AGENTS.md:** "Toda mudança estrutural ou slice não-trivial DEVE nascer com ADR ou estar registrada no backlog oficial."

**Como /adr-semparar implementa:**
```bash
# Fluxo obrigatório
1. bash scripts/adr-slice-init.sh <slug> --title "Título"
   → Cria ADR-YYYYMMDD-slug.md
   → Cria JSON taskmaster

2. Preenche MD com:
   - Contexto
   - Decisão
   - Escopo
   - Arquivos afetados
   - Smoke tests
   - Rollout + Rollback

3. Preenche JSON com:
   - goal
   - scope
   - done_definition
   - next_actions
   - handoff (CRÍTICO para continuidade)

4. Registra em docs/adr/README.md (automático via .context/)
```

✅ **Conformidade:** Workflow obriga ADR ANTES de executar qualquer mudança estrutural

---

### Regra 4: Higiene de Contexto
**AGENTS.md:** "Ao concluir mudança estrutural, slice não trivial, handoff relevante ou preparação para merge, atualize `.context/` via `sync-ai-context`."

**Como /adr-semparar implementa:**
```markdown
## Checklist para Fechar Slice

1. Update JSON:
   - Mude status para `accepted` ou `blocked`
   - Atualize `progress`
   - Adicione `evidence`

2. Rode validação:
   bash scripts/validate-adr-semparar.sh

3. OBRIGATÓRIO:
   bash scripts/sync-ai-context.sh

4. Commit com referência a sync:
   git commit -m "feat(adr): fechar slice X — sync-ai-context executado"
```

✅ **Conformidade:** Workflow força `sync-ai-context` como pré-requisito de fechamento

---

### Regra 5: Anti-Alucinação
**AGENTS.md:** "Nunca declare sucesso sem prova real (logs, testes, capturas)."

**Como /adr-semparar implementa:**
```json
{
  "evidence": [
    "ADR criada e indexada",
    "Smoke tests passam",
    "JSON validado"
  ],
  "test_commands": [
    "go test ./... -count=1",
    "bash ./scripts/smoke-*.sh"
  ],
  "curl_checks": [
    "curl -fsS http://endpoint/health"
  ]
}
```

**Validação:**
```bash
bash scripts/validate-adr-semparar.sh
# Verifica:
# - JSON válido
# - Campos obrigatórios presentes
# - Resume prompt não vazio
# - 12/12 slices com par MD+JSON
```

✅ **Conformidade:** Workflow obriga prova (testes + curl + evidence) antes de marcar aceito

---

### Regra 6: Sem Commits Diretos
**AGENTS.md:** "Use o workflow de `review-merge` para a branch principal."

**Como /adr-semparar implementa:**
```markdown
## Rollout

- Abrir ADR em branch de trabalho isolada
- Validar com `validate-adr-semparar.sh`
- Sync com `sync-ai-context.sh`
- Submeter para review via `review-merge`
- Nunca fazer push direto em main
```

✅ **Conformidade:** Workflow documenta isolamento em branch + necessidade de review

---

## 🎭 Papéis dos Agentes (AGENTS.md §3)

### Aurélia — Arquiteta Soberana
**Responsabilidade:** "Definir direção técnica, governar roteamento, manter coerência, arbitrar conflitos, preservar estabilidade."

**Implementação em /adr-semparar:**
```markdown
## Status

Cada slice tem status que reflete decisão de Aurélia:
- Proposto: Candidato a execução
- Em execução: Aurélia aprovou, agentes executam
- Aceito: Aurélia validou entrega
- Bloqueado: Bloqueio orquestrado por Aurélia
- Cancelado: Decisão arquitetural final
```

✅ **Conformidade:** Workflow trata status como decisão de Aurélia, não de agente individual

---

### Antigravity — Interface e Cockpit
**Responsabilidade:** "Orquestração de tarefas, handoff, interação com humano."

**Implementação em /adr-semparar:**
```json
{
  "handoff": {
    "owner_engine": "codex|claude|antigravity",
    "resume_prompt": "Estruturado para retomada sem contexto",
    "last_updated": "ISO 8601 timestamp"
  }
}
```

✅ **Conformidade:** Handoff explícito permite Antigravity coordenar sem perda

---

### Claude — Motor de Execução Principal
**Responsabilidade:** "Implementação técnica e revisões complexas."

**Implementação em /adr-semparar:**
```markdown
- owner: claude → Tarefas estruturais complexas
- next_actions: Sempre 3+ itens estruturais
- test_commands: Cobertura de implementação
```

✅ **Conformidade:** Workflow rastreia owner e exige testes de implementação

---

### Codex — Executor Rápido
**Responsabilidade:** "Refatorações e correções rápidas."

**Implementação em /adr-semparar:**
```markdown
- owner: codex → Slices de suporte, voice, deploy, testes rápidos
- progress: Pode ir de 0 a 100% rapidamente
- fallback_commands: Sempre presentes para contingência
```

✅ **Conformidade:** Workflow permite execução rápida com fallback documentado

---

## ✅ Governança e Tiers (AGENTS.md §5)

### Tier A: Read-only (Auto-approve 100%)
- Pesquisa em `docs/adr/`
- Leitura de `.context/`
- Leitura de `.agents/rules/`

**Implementação:**
```bash
# Qualquer agente pode:
cat docs/adr/ADR-20260319-*.md
cat .agents/workflows/adr-semparar-status.md
bash scripts/validate-adr-semparar.sh  # read-only validation
```

✅ **Conformidade:** Documentação é read-only, não requer aprovação

---

### Tier B: Local Edit (Auto-approve Condicional)
- Criar slice com `adr-slice-init.sh`
- Editar ADR/JSON em branch isolada
- Executar testes locais

**Implementação:**
```bash
# Em branch isolada:
bash scripts/adr-slice-init.sh feature-x --title "Título"
# ↳ Cria par MD+JSON
# Preenche manualmente
bash scripts/validate-adr-semparar.sh  # Valida antes de commit
git commit -m "feat(adr): abrir slice feature-x"
# ↳ Commit local, não é push direto
```

✅ **Conformidade:** Workflow exige isolamento + validação local antes de push

---

### Tier C: High-Risk (Aprovação Humana OBRIGATÓRIA)
- Merge para main
- Deploy em `.aurelia-24x7`
- Mudança em AGENTS.md

**Implementação:**
```bash
# Workflow obriga:
1. git commit em branch
2. sync-ai-context
3. validate-adr-semparar.sh PASSA
4. review-merge (exige aprovação humana)
5. Merge em main (nunca direto)
```

✅ **Conformidade:** Workflow documenta que merge é Tier C

---

## 📚 Referências Cruzadas

**Documentação Oficial:**
- `.context/docs/README.md` — Índice de documentação
- `.context/agents/README.md` — Playbooks de agentes
- `docs/REPOSITORY_CONTRACT.md` — Governança oficial
- `docs/adr/README.md` — ADR index

**Workflows Relacionados:**
- `.agents/workflows/adr-semparar.md` — Fluxo operacional
- `.agents/workflows/adr-semparar-governance.md` — Regras de estabilidade
- `.agents/workflows/adr-semparar-status.md` — Dashboard real

**Regras Operacionais:**
- `.agents/rules/` — Regras do Antigravity
- `AGENTS.md` (você está aqui) — Autoridade central
- `CLAUDE.md` — Adaptador Claude
- `CODEX.md` — Adaptador Codex

---

## 🎯 Checklist de Conformidade

Antes de usar `/adr-semparar`, confirme:

- ✅ Leu `AGENTS.md` (§ 2, 3, 4, 5)
- ✅ Entendeu hierarquia (Humanos → AGENTS.md → Aurélia → Agentes)
- ✅ Sabe que workflow é **subordinado**, não soberano
- ✅ Conhece as 6 regras centrais
- ✅ Entende os 3 tiers de governança
- ✅ Vai usar `bash scripts/validate-adr-semparar.sh` ANTES de commit
- ✅ Vai rodar `bash scripts/sync-ai-context.sh` AO FECHAR slice
- ✅ Entende que Aurélia decide status final, não o agente

---

## 📌 Princípio Overarching

> **"O workflow `/adr-semparar` não é uma ferramenta autossuficiente. É uma implementação operacional das regras de AGENTS.md. Se você está usando este workflow sem entender AGENTS.md, está fazendo errado."**

**Consequência:** Se há conflito entre este documento e AGENTS.md, AGENTS.md ganha sempre.

---

**Mantido por:** Global AI Governance
**Autoridade Suprema:** ../../AGENTS.md
**Revisão:** Qualquer mudança em AGENTS.md invalida este documento e requer update imediato
**Próxima Validação:** Quando AGENTS.md mudar, ou semanalmente (o que vier primeiro)
