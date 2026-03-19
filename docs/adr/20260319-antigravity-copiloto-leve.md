# ADR 20260319-antigravity-copiloto-leve

**Status**: Aceito  
**Data**: 2026-03-19

## Contexto

A Aurelia precisava usar o chat do Antigravity para microtarefas sem diluir a autoridade do executor principal nem abrir brechas de governança.

## Decisão

O chat do Antigravity passa a operar apenas como **copiloto leve**.

Entra no escopo:

- pesquisa curta
- pequenas configurações locais
- rascunho de patch
- revisão de diff curta
- preparação de handoff

Fica fora do escopo:

- secrets
- deploy
- rede
- merges
- ações destrutivas
- mudanças estruturais não aprovadas

Os artefatos da slice ficam em `.context/plans/20260319-antigravity-gemini-operator/`, e o comportamento canônico fica documentado em `docs/PROJECT_PLAYBOOK.md` e `docs/antigravity_gemini_operator_blueprint.md`.

## Consequências

Positivas:

- acelera microexecução sem perder governança
- reduz carga do executor principal em tarefas leves
- formaliza handoff e prova esperada

Trade-offs:

- exige roteamento correto de tarefas
- aumenta a necessidade de prompts estruturados

## Referências

- [PROJECT_PLAYBOOK.md](../PROJECT_PLAYBOOK.md)
- [antigravity_gemini_operator_blueprint.md](../antigravity_gemini_operator_blueprint.md)
- [implementation_plan.md](../../.context/plans/20260319-antigravity-gemini-operator/implementation_plan.md)
- [task.md](../../.context/plans/20260319-antigravity-gemini-operator/task.md)
