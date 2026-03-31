# Slice 5: Go Build Industrial Fix — Permission & .gitignore

**ADR Pai:** [20260330-enterprise-skills-governance.md](../20260330-enterprise-skills-governance.md)
**Status:** ✅ Concluída
**Data:** 2026-03-30

## Problema
`go build ./...` travava indefinidamente porque o diretório `data/redis/appendonlydir` (propriedade `dnsmasq`, permissão `drwx------`) bloqueava o walk do toolchain Go.

## Causa Raiz
Docker Redis criou `appendonlydir` com UID do container (dnsmasq) e permissão restritiva. O Go scanner tentava ler e ficava esperando permissão do OS.

## Correções
1. `sudo chmod -R a+rX data/` — corrigiu permissões
2. `data/` adicionado ao `.gitignore` como entrada raiz (não apenas subentradas)
3. Redis no `docker-compose.yml` agora usa `--dir /data --appendonly yes` para gravar no volume correto

## Aprendizado
- Volumes Docker com UID diferente do host causam travamento silencioso no Go build
- Sempre adicionar `data/` como entrada raiz no `.gitignore` de projetos com Docker
