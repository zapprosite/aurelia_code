# ADR 20260328: ADR Semparar — Resolução e Conformidade do Slice docs/adr

## Status
✅ Aceito (Slice Completa)

---

## Contexto

O PR #8 introduziu o workflow `/adr-semparar` com governança rigorosa de slices, porém o slice existente `docs/adr/` **não está conformante** com as regras do AGENTS.md e o padrão Semparar. A auditoria revelou:

### Achados da Auditoria (2026-03-28)

| Problema | Severidade | ADRs Afetados |
|---|---|---|
| Links obrigatórios ausentes | 🔴 CRÍTICA | 24/24 ADRs |
| Campo `## Status` ausente | 🟠 ALTA | 1 ADR (`linux-god-mode`) |
| Campo `Consequências` ausente | 🟡 MÉDIA | ~8 ADRs |
| taskmaster JSON ausente | 🔴 CRÍTICA | 0/0 ADRs |

### Detalhes da Auditoria

**Links Obrigatórios (AGENTS.md + REPOSITORY_CONTRACT.md):**
- Apenas `20260326-implementacao-master-skill-global.md` possui links (2/24)
- 23 ADRs **não referenciam** a autoridade central

**Taskmaster JSON:**
- O diretório `docs/adr/taskmaster/` **não existe**
- Nenhum ADR tem par `.md + .json` conforme o contrato Semparar

---

## Decisões

### 1. Adicionar Links Obrigatórios

Todo ADR em `docs/adr/` DEVE conter:

```markdown
## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
```

### 2. Criar taskmaster JSON para ADRs Estruturais

Para ADRs que representam **mudanças arquiteturais ou de governança**, criar par JSON:

```json
{
  "adr_id": "20260328-nome-da-slice",
  "title": "Título do ADR",
  "status": "in_progress",
  "progress": 0,
  "goal": "Objetivo em uma frase",
  "next_actions": ["ação 1", "ação 2"],
  "handoff": {
    "resume_prompt": "Prompt para retomada...",
    "owner_engine": "claude",
    "last_updated": "2026-03-28T00:00:00Z"
  },
  "test_commands": ["go build ./..."],
  "curl_checks": []
}
```

### 3. Corrigir ADR linux-god-mode

O ADR `20260328-implementacao-linux-god-mode.md`:
- Falta `## Status` com estado padronizado
- Falta `## Consequências`
- Falta links obrigatórios

### 4. Atualizar README.md

O ADR Index deve refletir:
- Conformidade Semparar (par MD+JSON)
- Status honesto (Implementada/Parcial/Proposta)
- Links para taskmaster/

---

## Plano de Execução

### Fase 1: Auditoria e Diagnóstico ✅
- [x] Identificar ADRs sem links obrigatórios
- [x] Identificar ADRs sem campo Status
- [x] Identificar ADRs sem Consequências
- [x] Documentar em `/docs/pr-review/PR-0008-INDUSTRIALIZE-SOTA-2026.md`

### Fase 2: Correção de ADRs ✅
- [x] Adicionar links obrigatórios a todos os 24 ADRs
- [x] Criar taskmaster/ directory
- [x] Gerar JSON para ADRs estruturais (6 JSONs criados)
- [x] Corrigir `linux-god-mode.md` (Status + Consequências)
- [x] Atualizar README.md do index

### Fase 3: Validação ✅
- [x] Validar JSONs taskmaster (6/6 ✅)
- [x] Verificar todos os links internos (26 ADRs conformantes)
- [x] Build do binário Aurelia (`go build ./cmd/aurelia` ✅)

---

## Resultado Final

| Métrica | Valor |
|---|---|
| ADRs com links obrigatórios | 26/26 ✅ |
| JSONs taskmaster criados | 6 |
| ADRs corrigidos | 24 |
| build Aurelia | ✅ OK |

### JSONs Criados em taskmaster/
1. `ADR-20260328-adr-semparar-docs-adr-resolve.json` — Slice de resolução
2. `ADR-20260328-jarvis-voice-computer-use.json` — Jarvis Voice + Computer Use
3. `ADR-20260328-linux-god-mode.json` — God Mode Linux
4. `ADR-20260325-basico-bem-feito-v2.json` — Basic Bot v2
5. `ADR-20260326-homelab-dashboard-native.json` — Dashboard Native
6. `ADR-20260327-markdown-brain.json` — Markdown Brain

---

## Consequências

### Positivas
- ADRs conformantes com AGENTS.md e hierarquia de autoridade
- Dashboard de slices funcional com JSON taskmaster
- Retomada de trabalho estruturada via resume_prompt
- Histórico audítavel de decisões

### Negativas
- Alteração massiva em 24+ arquivos Markdown
- Necessidade de validação manual de cada ADR
- Risco de conflitos se有人在 editando ADRs simultaneamente

---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [PR #8 Review](../../pr-review/PR-0008-INDUSTRIALIZE-SOTA-2026.md)

---

**Data**: 2026-03-28
**Autor**: Code Review Agent (Sovereign 2026.1)
**Slice**: `feature/neon-sentinel`
