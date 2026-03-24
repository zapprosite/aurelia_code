# Workflows — Padrões Operacionais Reutilizáveis

Este diretório contém workflows documentados que padronizam como OpenCode, Claude e Antigravity devem operar em tarefas recorrentes e estruturais.

## Workflows Disponíveis

### `/adr-semparar` — Slices Nonstop com Continuidade

**Arquivo:** `adr-semparar.md`
**Governança:** `adr-semparar-governance.md`
**Status:** `adr-semparar-status.md`

Abre slices estruturais longas com ADR + JSON taskmaster, garantindo continuidade sem perda de contexto entre agentes.

**Quando usar:**
- Tarefa estrutural > 30 minutos de execução
- Envolve múltiplos agentes (OpenCode, Claude, Antigravity)
- Precisa de handoff explícito

**Como começar:**
```bash
bash scripts/adr-slice-init.sh my-feature --title "Meu Título"
bash scripts/validate-adr-semparar.sh  # Verificar estabilidade
```

**Autoridade:**
- Template: `docs/adr/TEMPLATE-NONSTOP-SLICE.md` + `.json`
- Governança: Leia `adr-semparar-governance.md`
- Status real: Veja `adr-semparar-status.md`

---

## Princípios dos Workflows

1. **Documentação Primeira:** Toda ação está documentada em Markdown antes de executar
2. **Validação Contínua:** Scripts automatizam verificação de estabilidade
3. **Handoff Estruturado:** Resumos (`resume_prompt`) permitem retomada sem contexto externo
4. **Rastreabilidade:** Cada workflow tem um JSON complementar com estado executável

---

## Estrutura de um Workflow

Cada workflow .md deve ter:

```markdown
# /workflow-name
Description: Uma frase clara

---

1. Passo operacional 1
2. Passo operacional 2
3. ...
```

Complementado por:
- `.json` com estado executável (se multi-agente)
- `governance.md` com regras de estabilidade
- `status.md` com dashboard real

---

## Validação de Workflows

Todos os workflows com estado executável devem passar em:

```bash
bash scripts/validate-adr-semparar.sh  # Para workflows ADR
# (Adicione validadores adicionais conforme necessário)
```

---

## Autoridade e Escalação

- **Novo workflow?** → Comece em draft aqui, promova para `.agents/rules/` após 2 semanas estável
- **Mudança em workflow ativo?** → Atualize governança, rode validação, commit com `docs:` prefix
- **Conflito de workflow?** → Escale em `AGENTS.md` — workflows não podem contradizer autoridade central

---

## Histórico

- **2026-03-19:** Workflow `/adr-semparar` criado com 12 slices estáveis
- **Future:** Workflows para planejamento, review, deploy, incident response

---

**Mantido por:** Global AI Governance
**Sincronização:** `.context/workflow/docs/`
