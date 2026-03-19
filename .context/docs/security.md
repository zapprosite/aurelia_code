# Security

Segurança e privacidade são fundamentais para um agente que opera em ambiente local.

## Princípios de Segurança
- **Least Privilege**: O daemon roda como usuário, nunca como root.
- **Data Redaction**: Logs estruturados expurgam valores sensíveis (API keys, argumentos confidenciais).
- **Environment Isolation**: Preferência por comandos controlados e caminhos explícitos via `AURELIA_HOME`.
- **Secret Management**: Chaves de API gerenciadas via `app.json` com permissões 0600.

## Auditoria
- Revisão proativa de `AGENTS.md` e regras do repositório antes de execuções de risco.
- Validação de comandos (`run_command`) antes de execução automática quando em níveis de alto risco.
- Auditoria de `.gitignore` para evitar vazamento de credenciais locais em commits.
