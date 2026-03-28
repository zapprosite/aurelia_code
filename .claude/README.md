# .claude/ — Claude Code SOTA 2026

> Estrutura padrão Claude Code para o projeto Aurelia.

## Estrutura

```
.claude/
├── commands/       # Slash commands (/super-git, /pr-review, etc)
├── skills/          # Symlinks para skills do .agent/skills/
├── agents/         # Symlinks para workflows do .agent/workflows/
├── CLAUDE.md       # Instruções do projeto (herdado)
└── README.md       # Este arquivo
```

## Commands Disponíveis

| Command | Descrição |
|---|---|
| `/super-git` | Combo soberano de build + delivery |
| `/pr-review` | Review de PRs GitHub |
| `/adr-semparar` | Workflow de slices nonstop |
| `/dev` | Inicia implementação técnica |
| `/sincronizar-tudo` | Sync Git padrão sênior |
| `/sincronizar-ai-context` | Sincroniza contexto AI |
| `/git-turbo` | Merge + tag + cleanup |
| `/git-ship` | Ship para main |
| `/git-unblock` | Destrava Git |

## Skills Disponíveis (symlinks)

| Skill | Descrição |
|---|---|
| `code-review` | Revisão de código |
| `pr-review` | Review de PRs |
| `documentation` | Geração de docs |
| `frontend-design` | UI/UX design |
| `security-audit` | Auditoria de segurança |
| `homelab-control` | Controle do homelab |

## Referências

- [.agent/skills/](../.agent/skills/) — Skills completos
- [.agent/workflows/](../.agent/workflows/) — Workflows completos
- [AGENTS.md](../AGENTS.md) — Autoridade central
