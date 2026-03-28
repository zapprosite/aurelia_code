# ADR 20260328: Reorganização do .claude — Padrão Claude Code SOTA 2026

## Status
✅ Aceito (Slice Completa)

---

## Contexto

O diretório `/home/will/aurelia/.claude/` está com uma estrutura não padronizada:

### Problemas Identificados

| Problema | Severidade | Descrição |
|---|---|---|
| Symlinks redundantes | 🟡 MÉDIA | `.claude/agents` e `.claude/skills` apontam para `.agent/` |
| Estrutura não padronizada | 🟡 MÉDIA | Não segue o padrão Claude Code `.claude/{commands,agents,skills}` |
| Commands em lugar errado | 🟠 ALTA | Workflows estão em `.agent/workflows/`, não em `.claude/commands/` |
| Settings分散 | 🟡 MÉDIA | Configurações em locais diferentes |

### Estrutura Atual (Não Padrão)

```
aurelia/
├── .claude/              # Symlinks pointing to .agent/
│   ├── agents -> ../.agent/workflows
│   └── skills -> ../.agent/skills
├── .agent/               # Padrão Aurelia (não padrão Claude Code)
│   ├── agents/
│   ├── skills/
│   └── workflows/
└── .aurelia/            # Config local do Aurelia
```

### Estrutura Alvo (Padrão Claude Code)

```
aurelia/
├── .claude/              # Padrão Claude Code SOTA 2026
│   ├── CLAUDE.md        # Instruções do projeto
│   ├── settings.json    # Configurações do projeto
│   ├── commands/        # Slash commands (/review, /deploy, etc)
│   ├── agents/         # Agentes especializados
│   ├── skills/         # Skills de agentes
│   └── hooks/          # Event hooks (opcional)
├── .agent/              # Mantido para compatibilidade (legacy)
│   ├── skills/
│   └── workflows/
└── .aurelia/           # Config local do Aurelia (não versionado)
```

---

## Decisões

### 1. Manter `.claude/` com Estrutura Padrão

O `.claude/` deve seguir o padrão Claude Code:

```text
.claude/
├── CLAUDE.md           # Instruções do projeto (obrigatório)
├── settings.json      # Configurações do projeto
├── commands/          # Slash commands
│   ├── review.md      # /review
│   ├── pr-review.md  # /pr-review
│   └── super-git.md  # /super-git
├── agents/            # Agentes especializados
└── skills/          # Skills de agentes
    └── *.md
```

### 2. Migrar Commands para `.claude/commands/`

Comandos de workflow migrar de `.agent/workflows/` para `.claude/commands/`:

| De | Para |
|---|---|
| `.agent/workflows/super-git.md` | `.claude/commands/super-git.md` |
| `.agent/workflows/pr-review.md` | `.claude/commands/pr-review.md` |
| `.agent/workflows/adr-semparar.md` | `.claude/commands/adr-semparar.md` |

### 3. Atualizar Symlinks

Remover symlinks antigos e criar estrutura limpa:
- Remover `.claude/agents` (symlink)
- Remover `.claude/skills` (symlink)
- Criar estrutura padrão

---

## Plano de Execução

### Fase 1: Backup ✅
- [x] Documentar estrutura atual

### Fase 2: Criar Estrutura Padrão
- [ ] Criar `.claude/commands/`
- [ ] Criar `.claude/agents/`
- [ ] Criar `.claude/skills/`
- [ ] Copiar/migrar commands

### Fase 3: Migrar Commands
- [ ] Copiar `super-git.md` → `.claude/commands/`
- [ ] Copiar `pr-review.md` → `.claude/commands/`
- [ ] Copiar `adr-semparar.md` → `.claude/commands/`
- [ ] Criar symlinks ou copiar skills

### Fase 4: Cleanup
- [ ] Remover symlinks antigos
- [ ] Atualizar `.gitignore` se necessário
- [ ] Testar comandos

### Fase 5: Validação
- [ ] Verificar estrutura com `tree .claude/`
- [ ] Testar slash commands

---

## Consequências

### Positivas
- Estrutura padronizada seguindo Claude Code SOTA 2026
- Comandos disponíveis via `/command-name`
- Melhor organização e descobribilidade
- Alinhamento com documentação oficial

### Negativas
- Necessidade de migrar comandos existentes
- Potencial quebra de symlinks/integrações
- Ajuste de习惯了 para novos caminhos

---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)

---

**Data**: 2026-03-28
**Autor**: Code Review Agent (Sovereign 2026.1)
**Slice**: `feature/neon-sentinel`
