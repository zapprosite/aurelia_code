---
description: Modo turbo — commit, push, merge em main, tag e nova feature branch, tudo de uma vez.
---

# Fluxo Turbo (quando você está com pressa)

> ⚡ Faz tudo em sequência sem perguntas. Use quando precisar avançar rápido.

## Passo 1 — Salvar o que tem
1. `git add -A` — adiciona rastreados + novos não-ignorados pelo `.gitignore`
   - O `.gitignore` é a proteção real contra secrets e build artifacts.
2. Verificar se há algo staged: `git diff --cached --stat`
   - Se não houver nada, pular para o Passo 5 diretamente.

## Passo 2 — Commit com nome aleatório criativo
3. Gerar mensagem no formato `chore([escopo]): [verbo] [substantivo]` com vocabulário técnico:
   - Exemplos: `chore(core): patch signal-router`, `chore(api): wire async-conduit`,
     `chore(infra): align void-matrix`, `chore(db): sync iron-ledger`
   - Detectar escopo pelo path dos arquivos alterados (api, web, ui, core)
4. Commitar:
   ```bash
   git commit -m "chore([escopo]): [descrição-aleatória]

   Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
   ```

## Passo 3 — Push da branch atual
5. Obter branch atual: `git branch --show-current`
6. Push com proteção:
   ```bash
   git push --force-with-lease origin HEAD
   ```

## Passo 4 — Merge em main
7. Ir para main e atualizar:
   ```bash
   git checkout main
   git pull origin main
   ```
8. Merge da branch de origem:
   ```bash
   git merge --no-ff [branch-de-origem] -m "chore: merge [branch-de-origem] → main"
   ```
9. Push main:
   ```bash
   git push origin main
   ```

## Passo 5 — Tag aleatória
10. Gerar nome de tag no formato `v0.[número].[número]-[codinome]`:
    - Exemplos: `v0.9.1-phantom`, `v0.7.3-nebula`, `v0.4.2-forge`, `v0.6.0-eclipse`
11. Criar e enviar tag:
    ```bash
    git tag -a [nome-da-tag] -m "release: [nome-da-tag]"
    git push origin [nome-da-tag]
    ```

## Passo 6 — Nova feature branch para o próximo trabalho
12. Gerar nome criativo no formato `[adjetivo]-[substantivo]`:
    - Exemplos: `dark-runtime`, `swift-conduit`, `nano-kernel`, `flux-engine`, `zero-payload`
13. Criar branch e configurar upstream:
    ```bash
    git checkout -b feature/[nome-gerado]
    git push -u origin feature/[nome-gerado]
    ```

## Resumo final
14. Exibir tabela rápida:
    | Ação | Resultado |
    |------|-----------|
    | Commit | `chore([escopo]): [descrição]` |
    | Merge | `[branch] → main` |
    | Tag | `[nome-da-tag]` |
    | Nova branch | `feature/[nome-gerado]` |
