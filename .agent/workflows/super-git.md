---
description: //super-git - O Combo Definitivo de Industrialização e Entrega SOTA 2026
---

# //super-git

> 🚀 O workflow soberano que sincroniza, compila, pole e entrega tudo em uma única sequência ininterrupta.

// turbo-all

## Quando usar
- Finalização de grandes features ou marcos de industrialização.
- Quando o sistema está estável e pronto para `main` com tag de release.
- **Uso Obrigatório** para consolidar sessões de codificação complexas.

---

## Passos da Sequência Soberana

### 1. Sincronização de Contexto (@/sincronizar-ai-context)
Garante que o `codebase-map.json` e a documentação técnica estão em paridade com o código atual.
- Verificar paridade `.env`.
- Regenerar docs via `ai-context`.

### 2. Build Industrial (Linux Go)
Compilação estática (CGO_ENABLED=0) com injeção de metadados de versão.
- `export CGO_ENABLED=0`
- `go build -v -ldflags="-s -w -X main.Version=$(date +%Y.%m)-SOTA" -trimpath -o bin/aurelia ./cmd/aurelia`

### 3. Polish & Runtime Check
Garante que o binário gerado e o ambiente de execução estão nos padrões de 2026.
- Limpeza de `/bin/` (via `.gitignore`).
- Verificação de permissões de execução.

### 4. Entrega Inteligente (@/git-ship)
Criação de branch de staging (se necessário), commit semântico e Push + PR no GitHub.
- `gh pr create` (ou atualização do existente).

### 5. Finalização Turbo (@/git-turbo)
Merge em `main`, criação de Tag aleatória criativa e setup da próxima feature branch.
- `git checkout main && git merge`
- `git tag -a v0.9.6-infinite-tts`
- `git checkout -b feature/proximo-passo`

---

## Output Esperado
- Contexto 100% sincronizado.
- Binário compilado em `/bin/`.
- Pull Request mesclado em `main`.
- Tag de release criada e enviada.
- Nova branch de trabalho pronta.
