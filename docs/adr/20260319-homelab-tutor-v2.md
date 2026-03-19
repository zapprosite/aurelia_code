# ADR 20260319-homelab-tutor-v2

**Status**: Aceito  
**Data**: 2026-03-19

## Contexto

A Aurelia precisava deixar de ser apenas agente de execução e passar a atuar como tutor operacional do homelab, com memória, runbooks e governança explícita.

## Decisão

Adotar o `Homelab Tutor v2` como camada operacional local, com:

- classificação por domínio
- runbooks executáveis
- guardrails
- prova de resultado
- memória operacional e self-healing

O blueprint canônico sai da raiz e passa a viver em `docs/`, enquanto a implementação operacional continua nas skills locais.

## Consequências

Positivas:

- melhora a previsibilidade da manutenção autônoma
- transforma incidentes em aprendizado reaproveitável
- reduz dependência de prompts soltos

Trade-offs:

- aumenta o custo de curadoria de runbooks
- exige disciplina de atualização de memória e contexto

## Referências

- [homelab_tutor_v2_blueprint_20260319.md](../homelab_tutor_v2_blueprint_20260319.md)
- [homelab_jarvis_operating_blueprint_20260319.md](../homelab_jarvis_operating_blueprint_20260319.md)
- [PENDING-SLICES-20260319.md](./PENDING-SLICES-20260319.md)
