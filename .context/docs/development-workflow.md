# Development Workflow

O desenvolvimento do Aurelia segue regras rígidas de governança para manter a integridade do monorepo.

## Fluxo de Trabalho (BMAD)
1. **Planning (P)**: Requisitos definidos em PRD e Tech Specs.
2. **Review (R)**: Revisão de design e arquitetura antes da implementação.
3. **Execution (E)**: Implementação em branches ou worktrees isoladas.
4. **Verification (V)**: Testes automatizados e validação manual.
5. **Completion (C)**: Merge na main via processo de review.

## Padrões de Código
- **Idioma**: Comentários e documentação interna em Inglês. Documentação de usuário e relatórios finais em Português (BR).
- **Go**: Seguir `gofmt`. Evitar `log.Fatalf` fora de inicializações críticas; preferir retorno de erro.
- **Logging**: Sempre usar `observability.Logger(component)` e evitar logar valores sensíveis.

## Comandos Úteis
- `scripts/build.sh`: Compila o binário otimizado.
- `go test ./...`: Executa a suíte de testes.
- `scripts/install-user-daemon.sh`: Instala/Reinicia o serviço de usuário.

## Regra Operacional para MCP
- MCP opcional não deve ser tratado como dependência fatal de bootstrap.
- Em configs MCP do repositório, `enabled=false` deve ser respeitado no loader. Evite lógica que reative servidores apenas por estarem listados no arquivo.
- Antes de reativar um servidor MCP, valide:
  - comando/binário acessível pelo processo real;
  - permissões de execução;
  - variáveis de ambiente obrigatórias;
  - handshake/conexão básica.
