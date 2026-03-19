# ADR Nonstop Slice Reference

## Estrutura mínima do JSON

- `adr_id`
- `title`
- `status`
- `phase`
- `progress`
- `goal`
- `scope`
- `done_definition`
- `next_actions`
- `simulation_commands`
- `test_commands`
- `curl_checks`
- `evidence`
- `blockers`
- `handoff`

## Regras de preenchimento

- `next_actions` deve sempre apontar o próximo passo executável
- `simulation_commands` deve listar smokes, dry-runs e comandos reversíveis
- `curl_checks` deve conter endpoints locais quando existirem
- `handoff.resume_prompt` deve permitir continuidade sem releitura longa

## Simulações típicas

- `curl` em `/health`, `/metrics`, `/v1/router/status`, `/v1/voice/status`
- `go test` focado no pacote do slice
- scripts `smoke` existentes em `scripts/`
- fallback command quando a dependência principal estiver indisponível
