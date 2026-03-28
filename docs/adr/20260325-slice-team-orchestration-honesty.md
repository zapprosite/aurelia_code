# ADR 2026-03-25: Slice 3 — Team Orchestration Honesty

## Status
Implementada

## Objetivo

Parar de vender internamente “swarm” onde o comportamento real é `team orchestration`, e consolidar contratos de delegação, handoff e assistência lateral.

## Contexto

O runtime atual já é útil e persistente, mas sua linguagem ainda mistura branding aspiracional com execução real. Isso induz erro humano, erro de agente e documentação mentirosa.

O alvo desta slice não é enfraquecer o sistema. É tornar o naming compatível com o comportamento real.

## Escopo

- renomear superfícies internas de baixo risco de `swarm` para `team` quando o comportamento não for swarm cooperativo real
- formalizar os modos `delegation`, `handoff` e `assist`
- alinhar snapshot, dashboard e docs à mesma semântica
- reforçar testes de recuperação e notificação do team runtime

## Fora de escopo

- reescrever o sistema de tasks do zero
- introduzir microserviços
- criar colaboração multiagente complexa nova nesta slice

## Mudanças esperadas

1. inventário de nomes mentirosos ou ambíguos
2. renome de baixo risco em runtime/docs/handlers
3. contrato explícito de `delegation`, `handoff`, `assist`
4. alinhamento de dashboard/status/notifications
5. ajustes de teste para recovery e rehydration

## Smoke obrigatório

```bash
go test ./internal/agent ./cmd/aurelia ./e2e
curl -sS http://127.0.0.1:3334/api/status
curl -sS http://127.0.0.1:3334/api/bots
```

## Critério de saída

- naming interno reflete o comportamento real
- notificações e snapshot param de misturar “swarm” com “team”
- recovery continua verde após renome/contrato

## Dependência

Deve começar depois da Slice 2 para não renomear antes de endurecer contratos.


---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
