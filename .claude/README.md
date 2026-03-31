# .claude/ — Claude Code Template Pro (Sovereign 2026)

> Estrutura mínima. Skills e commands derivam de `.agent/`.

## Estrutura

```
.claude/
├── settings.json        # Permissões de ferramentas
├── settings.local.json  # Override local (não commitado)
└── README.md            # Este arquivo
```

## Fontes Canônicas

| Recurso | Localização |
|---------|-------------|
| Skills | `.agent/skills/` (SSOT) |
| Workflows/Commands | `.agent/workflows/` (SSOT) |
| Rules | `.agent/rules/` + `AGENTS.md` |
| Instruções Claude | `CLAUDE.md` (raiz do repo) |

## Referências

- [AGENTS.md](../AGENTS.md) — Autoridade central
- [CLAUDE.md](../CLAUDE.md) — Instruções do projeto
- [.agent/](../.agent/) — Skills, workflows e rules
