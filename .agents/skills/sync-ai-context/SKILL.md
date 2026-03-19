---
name: sync-ai-context
description: Sincroniza o ai-context do repositório, regenera .context/docs/codebase-map.json e revisa os .md curatoriais afetados.
---

# Skill: Sync AI Context

Atualiza a camada `.context/` do repositório quando houver mudanças estruturais, drift entre código e documentação, ou quando o `ai-context` precisar ser rodado de forma repetível.

## Quando usar

- Após mudanças relevantes em `cmd/`, `internal/`, `pkg/`, `scripts/` ou `docs/`
- Quando o `ai-context` apontar que os docs em `.context/docs/` precisam de revisão
- Quando `codebase-map.json` estiver desatualizado ou inconsistente com o checkout real
- Quando for necessário deixar evidência reexecutável de sincronização de contexto

## Diretivas

<directives>
1. **Descoberta primeiro**:
   - Leia `AGENTS.md`, `.agents/rules/` e `.context/docs/` antes de sincronizar.
   - Verifique se existe drift real entre código, `codebase-map.json` e os `.md` curatoriais.
2. **Comando canônico**:
   - Execute `./scripts/sync-ai-context.sh` a partir da raiz do repositório.
   - Esse script roda `ai-context update --dry-run` para detectar impacto e regenera `./.context/docs/codebase-map.json` de forma determinística.
3. **Limite de automação**:
   - Trate `.context/docs/*.md` como documentação curatorial.
   - Se o `ai-context` apenas sinalizar impacto, revise manualmente os `.md` afetados em vez de afirmar que o MCP os preencheu sozinho.
4. **Validação obrigatória**:
   - Confirme que `./scripts/sync-ai-context.sh` terminou com sucesso.
   - Confira `./.context/docs/codebase-map.json` para data de geração, contagem de arquivos e diretórios principais.
5. **Persistência de contexto**:
   - Se a sincronização representar uma mudança operacional importante, registre o resultado em `.context/workflow/docs/`.
</directives>

## Fluxo de Trabalho

1. Rodar `./scripts/sync-ai-context.sh`.
2. Ler o resumo de impacto emitido por `ai-context update --dry-run`.
3. Revisar os `.md` afetados em `.context/docs/` se houver drift semântico.
4. Validar o `codebase-map.json` regenerado.
5. Reportar quais arquivos foram sincronizados e o que ainda depende de curadoria humana.

## Output Esperado

- `codebase-map.json` atualizado
- `.context/docs/*.md` coerentes com o checkout atual
- resumo claro de quais docs foram revisitados e por quê
