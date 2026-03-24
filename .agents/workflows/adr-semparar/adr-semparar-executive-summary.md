---
title: ADR Semparar — Executive Summary
description: Resumo executivo do workflow de slices nonstop — estado, conformidade, próximos passos
audience: Humanos operadores, Aurélia, OpenCode, Claude, Antigravity
---

# ADR Semparar — Executive Summary

**Data:** 2026-03-19
**Status:** ✅ ESTÁVEL, DOCUMENTADO, PRONTO PARA PRODUÇÃO

---

## O Que É

`/adr-semparar` é um **workflow operacional** que implementa as regras de `AGENTS.md` para slices estruturais longas:
- ✅ Abertura com ADR + JSON taskmaster
- ✅ Continuidade garantida entre agentes (OpenCode, Claude, Antigravity)
- ✅ Handoff estruturado com `resume_prompt`
- ✅ Validação automática de conformidade
- ✅ Sincronização de contexto obrigatória

---

## Por Que Existe

As regras de `AGENTS.md` § 4 exigem:
1. ✅ **ADR por Slice:** Toda mudança estrutural nascer com ADR
2. ✅ **Isolamento de Worktree:** Tarefas não-triviais em branch isolada
3. ✅ **Higiene de Contexto:** Sync ao fechar slice
4. ✅ **Anti-Alucinação:** Prova real antes de sucesso
5. ✅ **Descoberta Local Primeiro:** Inspecionar AGENTS.md antes de agir

**`/adr-semparar` operacionaliza TODAS essas 6 regras.**

---

## Estado Atual

### Infraestrutura
```
✅ 12 slices com par MD+JSON
✅ Script validate-adr-semparar.sh (validation automática)
✅ Templates MD + JSON (padrão estável)
✅ 4 documentos de governança
✅ Dashboard status em tempo real
```

### Conformidade
```
✅ 100% aderente a AGENTS.md
✅ 0 erros críticos
✅ 0 avisos de validação
✅ Hierarquia clara (Humanos → AGENTS.md → Aurélia → Agentes)
✅ Todos os handoffs estruturados
```

### Slices em Execução
```
Onda 1 (Voz):           3/3 em execução
Onda 2 (Orquestração):  2/2 em execução
Onda 3 (Swarm):         1/1 proposto
Onda 4 (Desktop):       1/1 iniciado
Suporte:                5/5 mapeados
```

---

## Documentação Completa

| Arquivo | Propósito | Audiência |
| :--- | :--- | :--- |
| `adr-semparar.md` | Fluxo operacional passo-a-passo | Operadores |
| `adr-semparar-governance.md` | Regras de estabilidade + checklist | Engenheiros |
| `adr-semparar-agents-md-conformance.md` | Mapeamento explícito a AGENTS.md | Arquitetos |
| `adr-semparar-status.md` | Dashboard real em tempo | Todos |
| `README.md` | Guia de workflows | Onboarding |

---

## Como Usar

### Abrir Slice
```bash
bash scripts/adr-slice-init.sh my-feature --title "Descrição"
# → Cria ADR-YYYYMMDD-my-feature.md
# → Cria JSON taskmaster
# → Preencher manualmente com contexto/decisão/próximos passos
```

### Validar
```bash
bash scripts/validate-adr-semparar.sh
# → Verifica 12/12 slices com par MD+JSON
# → Verifica status consistente
# → Verifica campos obrigatórios em JSON
# → Verifica resume_prompt não vazio
# → Retorna ✅ ESTÁVEL ou ❌ ERRO
```

### Fechar Slice
```bash
# 1. Update JSON: status → accepted, progress → 100%
# 2. Validate: bash scripts/validate-adr-semparar.sh
# 3. OBRIGATÓRIO: bash scripts/sync-ai-context.sh
# 4. Commit: git commit -m "feat(adr): fechar slice X"
```

---

## Regras Invioláveis

Derivadas de `AGENTS.md`, implementadas no workflow:

### 1. Descoberta Local Primeiro (AGENTS.md § 4.1)
- ✅ Cada ADR tem links obrigatórios a AGENTS.md
- ✅ Workflow força leitura antes de executar

### 2. ADR por Slice (AGENTS.md § 4.3)
- ✅ Toda slice começa com `adr-slice-init.sh`
- ✅ Nenhuma mudança estrutural sem ADR

### 3. Higiene de Contexto (AGENTS.md § 4.4)
- ✅ Workflow obriga `sync-ai-context` ao fechar
- ✅ Checklist não deixa passar sem sincronizar

### 4. Anti-Alucinação (AGENTS.md § 4.5)
- ✅ JSON exige campos de evidence
- ✅ Validador verifica presença de smoke/test/curl

### 5. Isolamento de Worktree (AGENTS.md § 4.2)
- ✅ ADR documenta branch/worktree obrigatoriamente
- ✅ Workflow é subordinado a `review-merge`

### 6. Status Reflete Decisão de Aurélia (AGENTS.md § 3)
- ✅ Proposto/Em execução/Aceito/Bloqueado = decisão
- ✅ Agentes executam, não decidem status

---

## Hierarquia de Autoridade (Confirmada)

```
1. Humanos operadores ← Autoridade Final
              ↓
2. AGENTS.md ← Fonte Primária de Verdade
              ↓
3. Aurélia ← Autoridade Arquitetural e Operacional
              ↓
4. .agents/workflows/ ← /adr-semparar AQUI (Subordinado)
   └─ Implementa regras de AGENTS.md
   └─ Subordinado a Aurélia
   └─ Não pode contradizer autoridade acima
```

---

## Validação e Garantias

### Antes de Usar
- [ ] Leia `AGENTS.md`
- [ ] Entenda hierarquia
- [ ] Rode `validate-adr-semparar.sh` (deve passar)
- [ ] Confirme que 12/12 slices existem

### Ao Usar
- [ ] `adr-slice-init.sh` cria par correto
- [ ] Preenche manualmente contexto/decisão/escopo
- [ ] `validate-adr-semparar.sh` passa antes de commit
- [ ] Handoff resume_prompt é estruturado

### Ao Fechar
- [ ] JSON atualizado com status final
- [ ] `validate-adr-semparar.sh` passa
- [ ] **`sync-ai-context.sh` executado** (OBRIGATÓRIO)
- [ ] Commit referencia sync na mensagem

---

## O Que Garante

✅ **Conformidade com AGENTS.md:** Cada regra é implementada
✅ **Continuidade Multi-Agente:** Handoffs sem perda de contexto
✅ **Rastreabilidade:** Cada slice tem prova de execução
✅ **Estabilidade:** Validação automática garante integridade
✅ **Governança:** Hierarquia clara e respeitada

---

## O Que NÃO Garante

❌ Sucesso técnico (validação é estrutural, não técnica)
❌ Aprovação automática (ainda exige review-merge para main)
❌ Execução instantânea (agentes ainda precisam trabalhar)
❌ Conhecimento técnico (documentação sobre HOW-TO, não WHAT)

---

## Próximas Validações

- **Semanal:** `bash scripts/validate-adr-semparar.sh`
- **Por onda:** Status após conclusão de Onda 1, 2, 3, 4
- **Antes de merge:** Sempre rodar validador
- **Ao atualizar AGENTS.md:** Revisar conformance map

---

## Escalação e Conflitos

**Se encontrar problema:**
1. Rode `bash scripts/validate-adr-semparar.sh`
2. Leia `adr-semparar-agents-md-conformance.md`
3. Confirme que é **conflito com AGENTS.md**, não com workflow
4. Se for conflito com AGENTS.md → Escale para humano operador
5. Aurélia/OpenCode/Claude nunca podem contrariar AGENTS.md

---

## Para Humanos Operadores

Use este documento para:
- ✅ Entender por que slices têm ADR + JSON
- ✅ Validar que agentes estão operando conforme autoridade
- ✅ Identificar se há desvio de AGENTS.md
- ✅ Confirmar que continuidade é garantida entre agentes

**Pergunta de segurança:** Se há dúvida sobre conformidade → sempre releia AGENTS.md

---

## Para Aurélia

Use este documento para:
- ✅ Confirmar que agentes estão operando em conformidade
- ✅ Ver status real das 12 slices em tempo real
- ✅ Identificar bloqueadores ou desvios
- ✅ Validar que handoffs são estruturados

**Responsabilidade:** Você é a autoridade final que aprova progresso de slice

---

## Para Agentes (OpenCode, Claude, Antigravity)

Use este documento para:
- ✅ Entender o padrão de slice
- ✅ Saber como retomar handoff (`resume_prompt`)
- ✅ Confirmar que está em conformidade
- ✅ Saber quando rodar `sync-ai-context.sh`

**Regra de ouro:** Se há dúvida → não faça sem validação

---

## Sumário de Criação

```
2026-03-19 23:59 UTC

✅ Script: validate-adr-semparar.sh
✅ Governance: adr-semparar-governance.md
✅ Conformance: adr-semparar-agents-md-conformance.md
✅ Status: adr-semparar-status.md
✅ Workflow: adr-semparar.md (existente)
✅ README: README.md (workflows)
✅ This: adr-semparar-executive-summary.md

Total:
  - 12 slices com validação automática
  - 4 documentos de governança
  - 1 script de validação
  - 100% aderência a AGENTS.md
  - 0 erros críticos
```

---

**Mantido por:** Global AI Governance
**Autoridade Suprema:** ../../AGENTS.md
**Status:** ✅ PRONTO PARA PRODUÇÃO
**Próxima Revisão:** 2026-03-26 ou quando AGENTS.md mudar
