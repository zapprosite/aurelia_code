# Plans

Este diretório guarda artefatos de execução por slice que não devem morar na raiz do repositório.

## Regra

- `implementation_plan.md` e `task.md` vivem em `.context/plans/<slice>/`
- ADR continua em `docs/adr/`
- blueprint e runbook canônicos continuam em `docs/`

## Estrutura recomendada

```text
.context/plans/
  <slice>/
    implementation_plan.md
    task.md
```

## Observação

Se a slice virar regra permanente, a decisão precisa ser promovida para `docs/adr/`.
