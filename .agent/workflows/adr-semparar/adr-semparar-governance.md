# Governança: /adr-semparar

**Status:** ✅ Estável e Conformante com AGENTS.md
**Última validação:** 2026-03-19
**Autoridade suprema:** [AGENTS.md](../../AGENTS.md)
**Mapa de conformidade:** [adr-semparar-agents-md-conformance.md](./adr-semparar-agents-md-conformance.md)
**Total de slices:** 12 (9 Em execução, 2 Propostos, 1 Aceito)

---

## Objetivo

O workflow `/adr-semparar` padroniza a abertura e execução de **slices nonstop** — trabalhos estruturais longos que envolvem múltiplos agentes (OpenCode, Claude, Antigravity) e requerem continuidade sem perda de contexto.

Cada slice nasce com:
- ✅ ADR em `docs/adr/ADR-YYYYMMDD-slug.md`
- ✅ JSON taskmaster em `docs/adr/taskmaster/ADR-YYYYMMDD-slug.json`
- ✅ Handoff estruturado com `resume_prompt`
- ✅ Smoke tests e fallback documentados

---

## Regras de Estabilidade

### 1. Pair Obrigatório: MD + JSON

Toda slice DEVE ter:
- Um arquivo `.md` em `docs/adr/`
- Um arquivo `.json` em `docs/adr/taskmaster/`
- Mesmo `adr_id` em ambos os arquivos

**Validação:**
```bash
bash ./scripts/validate-adr-semparar.sh
```

### 2. Status em Padrão Português

MDs devem usar exatamente:
```markdown
## Status

- Proposto
- Em execução
- Aceito
- Bloqueado
- Cancelado
```

JSONs devem usar:
```json
"status": "proposed" | "in_progress" | "accepted" | "blocked" | "cancelled"
```

### 3. Campos Obrigatórios em JSON

Todos os 12 JSONs DEVEM ter:
- `adr_id`: Identificador único
- `title`: Título claro
- `status`: Estado atual
- `progress`: Percentual (0-100)
- `goal`: Objetivo em uma frase
- `next_actions`: Lista de próximos passos
- `handoff.resume_prompt`: Prompt estruturado para retomada
- `handoff.owner_engine`: Motor proprietário (opencode|claude|antigravity)
- `handoff.last_updated`: ISO 8601 timestamp

### 4. Smoke Tests Obrigatórios

Cada JSON DEVE ter pelo menos:
- `test_commands`: Testes de unidade/integração (mínimo 1)
- `curl_checks`: Validações HTTP (mínimo 1)
- `fallback_commands`: Plano de fallback (mínimo 1)

### 5. Handoff Estruturado

O `resume_prompt` DEVE ser:
- **Específico:** Não genérico, contém ação concreta
- **Executável:** Começa com verbo (Crie, Implemente, Valide, etc.)
- **Autossuficiente:** Não requer leitura de outros arquivos para começar

**Exemplo bom:**
```
"Crie a árvore docs/homelab/manual/ com 7 documentos (00-overview, 10-inventory, ...). Cada um com frontmatter: título, status, owner, data, tags, criticidade."
```

**Exemplo ruim:**
```
"Continue a trabalhar na slice"
```

---

## Arquivos Críticos

- **Workflow:** `.agents/workflows/adr-semparar.md`
- **Governança:** `.agents/workflows/adr-semparar-governance.md` (este arquivo)
- **Script de criação:** `scripts/adr-slice-init.sh`
- **Validador:** `scripts/validate-adr-semparar.sh`
- **Template MD:** `docs/adr/TEMPLATE-NONSTOP-SLICE.md`
- **Template JSON:** `docs/adr/taskmaster/TEMPLATE-NONSTOP-SLICE.json`

---

## Checklist para Abrir Slice

1. **Escolha slug:** `my-feature-name` (sem datas, sem hífens duplos)
2. **Execute:** `bash scripts/adr-slice-init.sh my-feature-name --title "Meu Título"`
3. **Preencha o MD:**
   - Status inicial
   - Contexto, Decisão, Escopo
   - Smoke tests
   - Rollout + Rollback
4. **Preencha o JSON:**
   - Goal, Done Definition
   - Next Actions (mínimo 3)
   - Smoke/test/curl commands
   - Resume prompt estruturado
5. **Teste:** `bash scripts/validate-adr-semparar.sh`
6. **Commit:** Com mensagem tipo `feat(adr): abrir slice my-feature-name`

---

## Checklist para Fechar Slice

1. **Update JSON:**
   - Mude status para `accepted` ou `blocked`
   - Atualize `progress` para 100% (aceito) ou deixe como está (bloqueado)
   - Adicione evidência no array `evidence`
2. **Update MD:**
   - Mude status para `Aceito` ou `Bloqueado`
   - Documente evidência esperada
3. **Rode validação:** `bash scripts/validate-adr-semparar.sh`
4. **Rode sync-ai-context:** `bash scripts/sync-ai-context.sh`
5. **Commit:** Com mensagem tipo `feat(adr): fechar slice my-feature-name — evidência em docs/adr/taskmaster/`

---

## Monitoramento de Saúde

### Validação Regular
```bash
# Executar semanal ou antes de merge
bash scripts/validate-adr-semparar.sh
```

### Verificação de Cobertura
```bash
# Garantir que não há slices órfãs
find docs/adr -name "ADR-*.md" | wc -l  # Deve igualar
find docs/adr/taskmaster -name "ADR-*.json" | wc -l
```

### Análise de Progresso
```bash
# Ver distribuição de status
jq -r '.status' docs/adr/taskmaster/ADR-*.json | sort | uniq -c
```

---

## Casos Especiais

### Slice Proposta (Não Aprovada)
- Status: `proposed`
- Progress: 0-10%
- Next actions: Design inicial, aprovação
- Pode permanecer aqui 2-4 dias enquanto ganha tração

### Slice Bloqueada (Aguardando Dependência)
- Status: `blocked`
- Progress: Congelar no último %
- Blockers: Documentar dependência explícita
- Handoff: Deixar pronto para retomar quando desbloqueada

### Slice em Transição (Entre Agentes)
- Handoff must-have: `resume_prompt` é 100% suficiente
- Teste com: Outro agente lê apenas o JSON e consegue continuar?
- Se não consegue → resume_prompt está incompleto

---

## Padrão de Resumo Executivo

Ao fazer report (status do projeto):
1. Leia `docs/adr/TASKMASTER-INDEX.md` (4 ondas)
2. Para cada onda, mostre:
   - Slices críticas + status
   - Bloqueadores
   - Próximas ações de maior prioridade
3. Inclua progresso agregado: `jq '.progress' docs/adr/taskmaster/ADR-*.json | awk '{s+=$1} END {print s/12 "%"}'`

---

## Autoridade e Escalação

- **Validação de padrão:** Este arquivo (`.agents/workflows/adr-semparar-governance.md`)
- **Template oficial:** `docs/adr/TEMPLATE-NONSTOP-SLICE.md` + `.json`
- **Autoridade final:** `AGENTS.md` (regras de slices e ADR)
- **Aprovação de exceção:** Precisa de consenso em AGENTS.md update

---

**Mantido por:** ADR Governance Working Group
**Última revisão:** 2026-03-19
**Próxima validação:** Pós-conclusão de Onda 1 (Voz)
