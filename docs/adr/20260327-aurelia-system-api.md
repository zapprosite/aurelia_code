# ADR 20260327: Aurelia System API 🔐⚙️⚓

## Status
Proposto (Aprovado pelo Usuário)

## Contexto
O ecossistema Aurélia cresceu para múltiplos bots, modelos e ferramentas que dependem de um arquivo `.env` complexo ("Completona"). A leitura direta do disco por múltiplos agentes é ineficiente e dificulta a manutenção do `.env.example` sincronizado.

## Decisão
Implementar uma API centralizada em Go (`services/aurelia-api`) para:
1.  **Governança de Segredos**: Atuar como um proxy seguro para metadados do `.env`.
2.  **Paridade de Repo**: Sincronizar automaticamente as chaves do `.env` para o `.env.example`.
3.  **Discovery**: Informar aos agentes o estado real de saúde dos modelos configurados no Smart Router.

## Consequências
- **Positivas**: Centralização de logs, auditabilidade de chaves e facilidade de setup para novos desenvolvedores/agentes.
- **Negativas**: Introduz uma dependência de runtime (a API precisa estar UP para os agentes funcionarem 100%).
