---
description: //super-git - O Combo Definitivo de Industrialização e Entrega SOTA 2026
---

# //super-git

> 🚀 O workflow soberano que sincroniza, compila, pole e entrega tudo em uma única sequência ininterrupta.

## Quando usar
- Finalização de grandes features ou marcos de industrialização.
- Quando o sistema está estável e pronto para `main` com tag de release.
- **Uso Obrigatório** para consolidar sessões de codificação complexas.
- **Pré-requisito**: Resolver todos os blockers pendentes (`//pr-review`).

---

## ⚠️ Gate de Segurança — Secret Scanning

Antes de qualquer `git push`, **OBRIGATÓRIO** verificar:

```bash
# Verificar se há secrets bloqueantes
gh api repos/zapprosite/aurelia_code/secret-scanning/alerts --jq '.[].secret_type'

# Se houver alertas, resolver ANTES de prosseguir:
# Opção 1: Allowlist no GitHub
# https://github.com/zapprosite/aurelia_code/security/secret-scanning/unblock-secret/<ID>

# Opção 2: BFG (⚠️ REESCREVE HISTÓRICO)
# java -jar bfg.jar --replace-text <(echo "configs/litellm/config.yaml:34:REDACTED") .
```

---

## Passos da Sequência Soberana

### 0. Pré-flight Check (Obligatório)
```bash
# Verificar status
git status --short

# Verificar build local
go build ./... 2>&1 | tail -5

# Verificar secrets no diff
git diff --staged | grep -i "api_key\|secret\|password" && echo "⚠️ SECRETS DETECTADOS!"
```

### 1. Sincronização de Contexto (@/sincronizar-ai-context)
Garante que o `codebase-map.json` e a documentação técnica estão em paridade com o código atual.
```bash
# Verificar paridade .env
./scripts/audit/audit-env-parity.sh

# Regenerar docs via ai-context
mcp__ai-context__context({ action: "fill", target: "all" })
```

### 2. Slice ADR (@/adr-semparar) — Se Nova Feature
Toda mudança estrutural deve nascer com ADR:
```bash
# Criar slice
bash scripts/adr-slice-init.sh <slug>

# Isso cria:
# - docs/adr/ADR-YYYYMMDD-slug.md
# - docs/adr/taskmaster/ADR-YYYYMMDD-slug.json
```

### 3. Build Industrial (Linux Go)
Compilação estática (CGO_ENABLED=0) com injeção de metadados de versão.
```bash
export CGO_ENABLED=0
VERSION=$(date +%Y.%m)-SOTA
go build -v \
  -ldflags="-s -w -X main.Version=${VERSION}" \
  -trimpath \
  -o bin/aurelia ./cmd/aurelia

# Validar binário
ls -lh bin/aurelia
file bin/aurelia
```

### 4. Suite de Testes (//test-all)
```bash
# Testes Go
go test ./... -v -count=1

# Testes de integração (se existirem)
./e2e/run.sh
```

### 5. Polish & Runtime Check
Garante que o binário gerado e o ambiente de execução estão nos padrões de 2026.
```bash
# Limpeza de binários via .gitignore
git clean -fd bin/

# Verificação de permissões
chmod +x bin/aurelia

# Smoke test
./bin/aurelia --version
```

### 6. Commit Semântico
```bash
# Gerar mensagem via //commit-message
@commit-message

# Ou manualmente (Conventional Commits):
git commit -m "feat(scope): descrição curta

- Bullet point 1
- Bullet point 2

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

### 7. Análise de PR (@/pr-review)
```bash
# Review antes do push
gh pr list --state open --repo zapprosite/aurelia_code
gh pr view <NUMBER> --json title,state,mergeable
```

### 8. Push com Auditoria de Segredos
```bash
# 🚨 ANTES DO PUSH — verificar secrets
scripts/audit/audit-secrets.sh

# Se audit passar, prosseguir:
git push origin $(git branch --show-current)

# Se BLOQUEADO por secret-scanning:
# → Verificar alertas: gh api repos/zapprosite/aurelia_code/secret-scanning/alerts
# → Allowlist ou BFG (ver seção ⚠️ Gate de Segurança acima)
```

### 9. PR Merge (//git-turbo) — Se PR Aprovado
```bash
# Checkout main e merge
git checkout main
git pull origin main
git merge --no-ff feature/<slug>

# Tag semântica
TAG="v$(date +%Y.%m.%d)-$(git rev-parse --short HEAD | cut -c1-4)"
git tag -a "$TAG" -m "Release: $(git log -1 --format=%s)"

# Push tags
git push origin main --tags

# Cleanup
git branch -d feature/<slug>
git push origin --delete feature/<slug> 2>/dev/null || true
```

---

## Fluxograma

```
┌─────────────────────────────────────┐
│  //super-git START                  │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  ⚠️ Secret Scanning Check            │ ← BLOQUEANTE
│  gh api secret-scanning/alerts       │
└─────────────────┬───────────────────┘
                  │ OK
                  ▼
┌─────────────────────────────────────┐
│  1. @/sincronizar-ai-context       │
│  2. //adr-semparar (se new feature)│
│  3. Build Go (CGO_ENABLED=0)       │
│  4. //test-all                      │
│  5. Polish & Runtime Check          │
│  6. Commit Semântico                │
│  7. @/pr-review                     │
│  8. Push (audit-secrets.sh)         │ ← BLOQUEANTE
│  9. //git-turbo (se PR approved)    │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  ✅ RELEASE COMPLETO                 │
│  - Binário em /bin/                 │
│  - PR merged em main                │
│  - Tag criada e pushada             │
│  - Nova branch pronta               │
└─────────────────────────────────────┘
```

---

## Output Esperado

| Artefato | Status |
|---|---|
| Contexto 100% sincronizado | ✅ |
| Binário compilado em `/bin/` | ✅ |
| ADRs criados (se new feature) | ✅ |
| PR criado/atualizado no GitHub | ✅ |
| Suite de testes passando | ✅ |
| Secret scanning audit OK | ✅ |
| Pull Request mesclado em `main` | ✅ |
| Tag de release criada e enviada | ✅ |
| Nova branch de trabalho pronta | ✅ |

---

## Troubleshooting

### Push Bloqueado por Secret Scanning
```bash
# Verificar alertas
gh api repos/zapprosite/aurelia_code/secret-scanning/alerts \
  --jq '.[].secret_type, .[].created_at'

# Allowlist (recomendado)
# Ir até: https://github.com/zapprosite/aurelia_code/security/secret-scanning

# BFG (⚠️ Rewrites history)
# Available in: /srv/ops/tools/bfg-1.14.0.jar
java -jar /srv/ops/tools/bfg-1.14.0.jar \
  --replace-text <(echo "configs/litellm/config.yaml:34:REDACTED_API_KEY") \
  --no-blob-protection .
git reflog expire --expire=now --all && git gc --prune=now --aggressive
```

### Build Falhou
```bash
# Verificar dependências
go mod tidy
go mod verify

# Verificar CGO
echo $CGO_ENABLED  # deve ser 0
```

### Tests Falhando
```bash
# Ver logs detalhados
go test ./... -v -count=1 2>&1 | tail -50

# Verificar se há fixtures pendentes
ls -la tests/fixtures/
```

### Git Zumbi (@/git-unblocker)
```bash
# Se "index.lock exists"
rm -f .git/index.lock

# Se "another git process running"
ps aux | grep git
kill -9 <PID>
```

---

## Referências

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](../../docs/adr/README.md)
- [Secret Scanning GitHub](https://github.com/zapprosite/aurelia_code/security/secret-scanning)
- [Skill: sincronizar-ai-context](../../.agent/skills/sync-ai-context/SKILL.md)
- [Skill: pr-review](../../.agent/skills/pr-review/SKILL.md)
- [Skill: git-unblocker](../../.agent/skills/git-unblocker/SKILL.md)
