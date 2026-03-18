# Política de Segurança — Aurelia Elite Edition

## 1. Redação de Logs (Redaction)
Para mitigar riscos de vazamento de dados sensíveis (PII) e segredos em logs:
- **Core Loop**: Logs de execução de ferramentas em `internal/agent/loop.go` omitem argumentos brutos.
- **Roteador**: Logs de erro do classificador de intenções em `internal/skill/router.go` não imprimem o JSON bruto da resposta.
- **Recomendação**: No desenvolvimento de novas ferramentas, evite `log.Printf("%+v")` em structs que contenham tokens ou chaves.

## 2. Gerenciamento de Segredos
- **Configuração Local**: Todos os segredos residem em `~/.aurelia/config/app.json`.
- **Git Hygiene**: O `.gitignore` está configurado para nunca versionar `app.json`, `*.db` ou qualquer arquivo em `.env`.
- **Auditoria**: Auditorias periódicas devem ser feitas para garantir que nenhum placeholder como `[CHAVE_AQUI]` ou `TODO` sobreviva em arquivos `.md`.

## 3. Conformidade Tier C
- **Observabilidade**: O stdout/stderr não deve ser exposto a logs persistentes de terceiros se o `verbose_log` estiver ativado para depuração avançada local.
- **Acesso**: O acesso ao bot de Telegram é restrito exclusivamente aos IDs listados em `telegram_allowed_user_ids`.

---
*Atualizado em 18 de Março de 2026 por Antigravity.*
