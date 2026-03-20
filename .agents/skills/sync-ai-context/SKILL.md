---
name: sync-ai-context
description: Sincroniza o ai-context do repositório via MCP tools — padrão 2026.
---

# Skill: Sync AI Context

## Quando usar

- Após mudanças em `cmd/`, `internal/`, `pkg/`, `scripts/` ou `docs/`
- Antes de handoff entre agentes ou merge final
- Quando `.context/docs/codebase-map.json` estiver desatualizado

## Quando dispensar

- Typo ou comentário sem impacto estrutural
- Rename local sem drift semântico

---

## Execução via MCP (padrão)

O agente executa as tools do MCP `ai-context` diretamente — sem shell:

```
1. context({ action: "check", repoPath: "<repo>" })
   → confirmar initialized: true

2. context({ action: "listToFill", target: "docs" })
   → ver quais .md precisam atualização

3. context({ action: "fill", target: "docs" })
   → regenerar codebase-map.json e docs curados
   (ou fillSingle para um arquivo específico)

4. context({ action: "getMap", section: "stats" })
   → validar mapPath e lastUpdated
```

Se houver drift real, commitar:
```bash
git add .context/ && git commit --no-verify -m "chore(context): sync ai-context pos-<slug>"
```

## Fallback (MCP indisponível)

```bash
./scripts/sync-ai-context.sh
```

## Output esperado

- `.context/docs/codebase-map.json` atualizado (data recente)
- `.context/docs/*.md` coerentes com o checkout
- Commit registrado se houve drift real
