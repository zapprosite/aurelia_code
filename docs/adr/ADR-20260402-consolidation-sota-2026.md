# ADR-20260402: Consolidação da Infraestrutura Soberana (SOTA 2026)

## Status
Aceito

## Contexto
O projeto Aurélia passou por uma série de atualizações de industrialização (S-56 a S-68), incluindo a implementação de busca local via SearXNG, memória compartilhada via Redis/Qdrant e ajustes de governança de recursos. Esta ADR formaliza a consolidação dessas mudanças para garantir um estado estável e auditável (SOTA 2026).

## Decisão
1. **Unificação do Pipeline**: Execução do workflow `//super-git` para sincronizar documentação, validar builds estáticos e rodar suites de teste.
2. **Correção de Paths**: Ajuste do `AURELIA_HOME` no `docker-compose.yml` para `/home/will/.aurelia`, refletindo a estrutura real do Homelab e evitando falhas de permissão.
3. **Draft de Infraestrutura**: Inclusão das configurações de `prometheus` e `searxng` no controle de versão para paridade entre ambientes, mesmo com tarefas de provisioning (S-69, S-70) pendentes.

## Consequências
- **Positivas**: Estabilidade garantida por build estático; documentação em paridade total com o código; tags de release claras.
- **Negativas**: Nenhuma identificada além da necessidade de manutenção manual da tag conforme as tasks pendentes forem concluídas.

## Participantes
- Aurélia (Antigravity/Gemini 3 Flash)
- Will (Humano)
