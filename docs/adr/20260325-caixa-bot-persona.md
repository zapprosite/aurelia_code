# ADR 20260325: Implementação do Bot CAIXA_PF_PJ

## Status
Proposto

## Contexto
O usuário Master (Will) necessita de um bot especializado na gestão de suas contas Caixa Econômica Federal (Pessoa Física e Jurídica). O requisito principal é uma persona que atue como uma secretária executiva "insistente", utilizando o sistema de agendamento (cron jobs) para garantir que boletos, lembretes de conta e prazos bancários não sejam esquecidos.

## Decisão
Implementar um novo bot no pool de bots da Aurélia com as seguintes características:
- **ID**: `caixa-pf-pj`
- **Persona**: `secretaria-caixa`
- **Comportamento**: A persona deve monitorar menções a obrigações financeiras e proativamente sugerir ou criar agendamentos via `create_schedule`.
- **Identidade Vocal**: Utilizar a voz feminina premium `pt-br` do Kokoro (padrão soberano).

## Consequências
- Aumento da proatividade do sistema em tarefas administrativas.
- Necessidade de garantir que o bot tenha permissão para criar schedules vinculados ao ID do Master.
- O pool de bots agora gerencia mais de um bot simultaneamente, testando a escalabilidade do `BotPool`.


---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
