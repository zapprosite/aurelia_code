# Testing Strategy

O Aurelia prioriza testes de unidade e integridade para componentes críticos (lock, memory, agent).

## Tipos de Teste
- **Unit Tests**: Testam a lógica isolada de pacotes (`go test ./internal/...`). Foco em transformações de dados e regras de negócio.
- **Integration Tests**: Validação de interações de subsistemas (ex: lock de instância, persistência SQLite).
- **Static Analysis**: Uso rigoroso de `go vet`, `gofmt` e linters recomendados para Go.

## Requisitos de Validação
- Nenhuma alteração em `internal/runtime` ou `internal/observability` deve ser aceita sem testes automatizados correspondentes.
- Verificações de `git diff --check` para evitar whitespace e conflitos básicos.
- Scripts bash devem passar por `bash -n` para verificação de sintaxe.

## Ambiente de Teste
- Testes locais devem ser executáveis sem dependências externas complexas (usando stubs/mocks para APIs de LLM se necessário).
- O uso de `XDG_CONFIG_HOME` ou mocks do sistema de arquivos é encorajado para testes de runtime.
