---
name: template-mult-clouds
description: Bootstrap universal do repositório multi-agente com Antigravity + Claude + Codex + OpenCode.
---

# 🏗️ Skill: template-mult-clouds

Automatiza a implantação do padrão de autoridade única e governança Elite para repositórios multi-agente.

<directives>

## Checklist de Implantação

### 1. Contratos Soberanos (Raiz)
- [ ] `AGENTS.md` — contrato de autoridade única com hierarquia, papéis e tiers
- [ ] `CLAUDE.md` — adaptador fino para Claude Code, subordinado a AGENTS.md
- [ ] `CODEX.md` — adaptador fino para Codex CLI, subordinado a AGENTS.md
- [ ] `GEMINI.md` — adaptador fino para Antigravity, subordinado a AGENTS.md
- [ ] `MODEL.md` — política de modelos, voz e roteamento
- [ ] `README.md` — porta de entrada para agentes e humanos
- [ ] `CONTRIBUTING.md` — regras de contribuição
- [ ] `SECURITY.md` — política de segurança

### 2. Regras Operacionais (`.agents/rules/`)
- [ ] `01-authority.md` — fonte de verdade e hierarquia
- [ ] `02-local-first.md` — descoberta local antes de agir
- [ ] `03-tiers-autonomy.md` — Tier A/B/C com sudo=1 se habilitado
- [ ] `04-worktree-isolation.md` — branches/worktrees isoladas
- [ ] `05-context-state.md` — higiene do .context/
- [ ] `06-planning-first.md` — ADR antes de execução
- [ ] `07-shared-mcp.md` — MCP compartilhado entre agentes
- [ ] `08-diff-reporting.md` — relatórios baseados em diff
- [ ] `09-skills-usage.md` — skills como extensão de capacidade
- [ ] `10-artifact-discipline.md` — disciplina de artefatos
- [ ] `11-adr-slice-contract.md` — ADR por slice estrutural

### 3. Skills de Elite (`.agents/skills/`)
- [ ] `architect-planner` — planejamento e tech spec
- [ ] `security-first` — segurança proativa
- [ ] `deep-researcher` — pesquisa profunda
- [ ] `homelab-control` — controle de home lab
- [ ] `go-telegram-expert` — bots Telegram em Go
- [ ] `self-healing` — auto-recuperação e watchdog
- [ ] `sync-ai-context` — sincronização do contexto AI

### 4. Workflows (`.agents/workflows/`)
- [ ] `/pm`, `/architect`, `/dev`, `/qa` — ciclo completo BMAD
- [ ] `/git-feature`, `/git-ship` — Git workflow senior
- [ ] `/adr-semparar` — slices longas com taskmaster JSON
- [ ] `/sincronizar-ai-context` — higiene de contexto
- [ ] `/review-merge` — merge na main com revisão

### 5. Contexto e Estado (`.context/`)
- [ ] Inicializado via `ai-context init`
- [ ] `docs/README.md` — índice de documentação
- [ ] `agents/README.md` — índice de agentes
- [ ] `codebase-map.json` — mapa do repositório

### 6. Governança (`docs/`)
- [ ] `REPOSITORY_CONTRACT.md` — índice de governança
- [ ] `adr/` — decisões arquiteturais por slice
- [ ] `adr/INDEX.md` — índice de ADRs

</directives>

## Verificação de Sucesso

O repositório é Elite quando:

1. Um agente estranho compreende toda a governança lendo apenas `README.md` e `AGENTS.md`
2. `go test ./... -count=1` retorna 100% verde
3. `go build ./...` compila sem erros
4. Nenhum placeholder "TODO" nos contratos soberanos
5. Nenhum segredo exposto no repositório (verificar pré-push)
