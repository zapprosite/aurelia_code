---
name: git-unblocker
description: Resolve travamentos fantasmas do Git, index.lock preso e timeouts no Ubuntu.
---

# Skill: Git Unblocker (Modo Sênior)

## Problema
No ambiente de desenvolvimento local (especialmente Ubuntu/WSL com concorrência alta, vscode watchers ativos ou execuções abortadas de CI/agentes), o Git frequentemente entra em estado de "zumbi" ou trava indefinidamente em operações simples (`git status`, `git add`, `git push`). O arquivo `.git/index.lock` fica preso e os processos antigos não morrem.

## Sintomas
- Comandos Git demoram mais de 10 segundos para retornar.
- Erro: `Another git process seems to be running in this repository, e.g. an editor opened by 'git commit'. Please make sure all processes are terminated then try again.`
- O agente não consegue avançar e atinge timeout na tool `run_command`.

## Solução Definitiva (Ação Sênior)

Quando identificar esse padrão, **NÃO** tente esperar ou rodar o comando repetidas vezes. Você deve matar os precessos pendurados e limpar os locks "na marra" antes de tentar novamente.

### Passo 1: Execução Agressiva de Limpeza

Execute o seguinte comando bash exato:
```bash
pkill -9 -f "git " 2>/dev/null; pkill -9 -f "gh pr" 2>/dev/null; sleep 2; rm -f .git/index.lock .git/refs/heads/feature.*.lock 2>/dev/null; git status --short
```

### O que este comando faz:
1. `pkill -9 -f "git "`: Mata todos os processos do sistema que contenham "git " na linha de comando (com o `-9` SIGKILL, sem piedade).
2. `pkill -9 -f "gh pr"`: Mata também o GitHub CLI caso esteja travado em operações de PR.
3. `sleep 2`: Aguarda o SO liberar os file descriptors.
4. `rm -f .git/index.lock ...`: Deleta ativamente arquivos de lock que ficaram órfãos.
5. `git status --short`: Executa um comando leve e inofensivo para garantir que a árvore Git voltou a responder instantaneamente.

### Passo 2: Execução Segura e Isolada das Próximas Ações

Após o unblock:
- Evite encadear múltiplos comandos pesados do Git com `&&` (ex: `git add && git commit && git push && gh pr create`).
- Separe operações de rede (push, fetch, pr) das operações locais (add, commit, tag). Opere em passos isolados usando o `run_command` com tempos de timeout adequados (`WaitMsBeforeAsync`).

## Quando NÃO Usar
- Se o Git estiver rodando um `git clone` ou `git lfs pull` massivo que genuinamente leva tempo esperado. O `git-unblocker` abortará o progresso.
