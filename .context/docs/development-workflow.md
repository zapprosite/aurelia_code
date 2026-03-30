---
type: doc
name: development-workflow
description: Day-to-day engineering processes, branching, and contribution guidelines
category: workflow
generated: 2026-03-30
status: filled
scaffoldVersion: "2.0.0"
---

# Fluxo de Desenvolvimento — Aurélia Sovereign 2026.2

## Branches e Releases

Modelo trunk-based com feature branches curtas.

```
main                           ← código em produção
feature/[adjetivo]-[substantivo]  ← branches de feature
```

Nome de feature branch: `[adjetivo]-[substantivo]`
Exemplos: `feature/nano-kernel`, `feature/halo-engine`

Tags: `v0.X.Y-[codinome]`
Exemplo: `v0.10.0-phantom`, `v0.9.7-phantom`

## Fluxo Turbo (git-turbo skill)

Quando com pressa, usar o fluxo turbo para fazer tudo em sequência:
1. `git add -A` (protegido por .gitignore)
2. `git commit` com nome criativo `chore([escopo]): [verbo] [substantivo]`
3. `git push --force-with-lease`
4. Merge em main
5. Tag `v0.X.Y-codinome`
6. Nova feature branch
7. Commit de ADR slice para cada mudança significativa

## Deploy

```bash
# Build
CGO_ENABLED=0 go build -trimpath -o bin/aurelia ./cmd/aurelia

# Deploy
sudo systemctl restart aurelia

# Health check
curl http://localhost:8585/health
```

## Estrutura de Commits

```
chore([escopo]): [verbo] [substantivo]

Exemplos:
fix(telegram): onboarding direto + porteiro bypass para owner
chore(core): cascade liteelm + faster-whisper local + legacy cleanup
feat(gateway): LiteLLM router + bot Telegram pipeline
```

## ADR Slices

Toda mudança estrutural (novo serviço, mudança de pipeline, remoção de legado) exige ADR slice em `docs/adr/slices/YYYYMMDD-nome.md`.

Formato: contexto → decisão → consequências (positivo/negativo/neutro).

## Auditoria de Segredos

O hook `scripts/audit/audit-secrets.sh` roda automaticamente em cada commit. Bloqueia commits que exponham segredos.

## Health Checks

```bash
# Bot
curl http://localhost:8585/health

# LiteLLM
curl http://localhost:4000/health

# STT local
curl http://localhost:8020/health

# TTS local
curl http://localhost:8012/health

# Redis
docker exec aurelia-redis-1 redis-cli ping

# Qdrant
curl http://localhost:6333/healthz
```

## Stress Test

```bash
BASE="http://localhost:8585/v1/telegram/impersonate"
ID=7220607041
curl -X POST "$BASE" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":$ID,\"message\":\"qual a diferença entre goroutine e thread?\"}"
```
