# ADR 20260329-consolidacao-seguranca-sentinel

**Data**: 2026-03-29
**Status**: Aceito
**Autor**: Antigravity (Antigravity Code)

## Contexto
O ecossistema Aurélia possuía dois sistemas de segurança de entrada redundantes: o `InputGuard` legado (baseado em Qwen 3.5 sem cache persistente) e o novo `PorteiroMiddleware` (baseado em Qwen 0.5b com cache Redis). Além disso, a introdução da persona **Junior Developer** exigia travas de segurança específicas para evitar ações destrutivas acidentais.

## Decisão
1. **Consolidação**: Remover o componente `InputGuard` e centralizar toda a lógica de segurança no `PorteiroMiddleware`.
2. **Endurecimento (Hardening)**: 
    - Implementar "Zero Destruição" no prompt de sistema da persona Junior.
    - Unificar banners de bloqueio (`[🛑 BLOQUEIO DE SEGURANÇA]`).
    - Permitir bypass automático para o proprietário (Master) para evitar atrito operacional.
3. **Soberania**: Utilizar o modelo Qwen 0.5b local para garantir que nenhum dado de telemetria ou prompt saia do homelab para fins de segurança básica.

## Consequências
- **Positivas**: Redução da dívida técnica, latência unificada via Redis, maior segurança contra injeções semânticas.
- **Negativas**: Dependência total do Redis para performance ótima (fallback fail-open implementado).
