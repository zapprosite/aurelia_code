# ADR-20260321-codebase-symbol-map: Sub-4 Codebase Symbol Map

## Status
Ativo

## Contexto
O agente da Aurelia necessita entender estruturalmente a base de código antes de iniciar mutações (mimetizando a leitura de projeto feita por engenheiros). A sincronização via `ai-context` MCP agora garante a consistência e geração do arquivo `.context/docs/codebase-map.json`. 

Precisamos implementar a capacidade da engine nativa do Go de absorver esse JSON e expor metadados chave do domínio, acoplando com a fase `PLANNING` introduzida pelo Sub-3.

## Decisão
Implementar um mecanismo nativo (`internal/agent/codebase_map.go`) que faça _parsing_ e compactação inteligente do `codebase-map.json` durante a injeção do System Prompt. Desta forma, a estrutura macro do projeto e símbolos críticos já estarão no Context Window primário antes da primeira emissão.

## Consequências
- Aceleração monstruosa no Planejamento: redução drástica das chamadas reativas a `list_dir` e `read_file` em prol da consciência _in-prompt_.
- Aumento da taxa de acerto no handoff de contexto.

## Testes e Rollout
1. Validação unitária de parse de struct.
2. Build via `go build ./cmd/aurelia`.
3. Análise visual em Log de Debugger (`slog`) de como o prompt de sistema se deforma com injeção do _Codebase Symbol Map_.
