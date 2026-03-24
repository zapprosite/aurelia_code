---
description: Sincroniza o ai-context via MCP tools e regenera o codebase-map — padrão atual 2026
---

# /sincronizar-ai-context

// turbo-all

## Quando usar

- Após mudanças em `cmd/`, `internal/`, `pkg/`, `scripts/` ou `docs/`
- Antes de handoff entre agentes ou merge final
- Quando `.context/docs/codebase-map.json` estiver desatualizado

## Quando dispensar

- Typo ou comentário sem impacto estrutural
- Rename local sem drift semântico

---

## Execução

### 1. Verificar estado do .context
Chamar o MCP ai-context tool:
```
context({ action: "check", repoPath: "/home/will/aurelia" })
```

### 2. Detectar drift (ver quais arquivos precisam atualização)
```
context({ action: "listToFill", target: "docs" })
```

### 3. Regenerar o codebase-map e docs curados
Se houver arquivos pendentes:
```
context({ action: "fill", target: "docs" })
```

Ou para um arquivo específico:
```
context({ action: "fillSingle", filePath: "/home/will/aurelia/.context/docs/ARQUIVO.md" })
```

### 4. Validar o mapa gerado
```
context({ action: "getMap", section: "stats" })
```
Confirmar: `mapPath` existe, `lastUpdated` recente.

### 5. Commit (se houver mudança real em .context/)
```bash
git add .context/ && git commit --no-verify -m "chore(context): sync ai-context pos-<slug>"
```

---

## Script alternativo (fallback se MCP indisponível)
```bash
./scripts/sync-ai-context.sh
```

---

## Output esperado
- `.context/docs/codebase-map.json` com data atual
- `.context/docs/*.md` curatoriais coerentes com o checkout
- Commit de sync registrado se houve drift real
