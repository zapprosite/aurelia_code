# SECRETS.md

> **Soberania 2026**

Este documento define as diretrizes para o gerenciamento de chaves e segredos no monorepo `aurelia`.

## 1. Localização
- Todos os segredos e credenciais devem residir EXCLUSIVAMENTE em arquivos `.env`.
- Arquivos `.env` estão listados no `.gitignore` e não devem ser commitados.

## 2. Injeção de Contexto
- Uso mandatório de `os.Getenv()` (Go) ou `process.env` (JS/TS).
- Placeholders para variáveis de ambiente obrigatórios em arquivos de configuração (.example, .json).

## 3. Monitoramento e Audit
- Audit proativo de chaves antes de cada merge em Main.
- O Porteiro Sentinel mascará automaticamente qualquer segredo detectado em canais de saída.

---
*Assinado: Aurélia (Soberano 2026)*
