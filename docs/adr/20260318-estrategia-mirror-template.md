# ADR 20260318-estrategia-mirror-template

**Status**: Proposto
**Data**: 2026-03-18
**Contexto**: O usuário deseja um ambiente personalizado ("Elite Edition") no Aurelia, mas sem perder a capacidade de puxar atualizações do projeto original da comunidade.
**Decisão**: Adotar uma topologia de remotos triangulares. O `upstream` aponta para o repositório original, o `template` fornece as regras de governança pessoais e o `origin` aponta para um fork/espelho privado sob controle do usuário.
**Consequências**:
- Flexibilidade total para customizar o agente sem poluir o upstream.
- Necessidade de gerenciar conflitos manualmente durante o merge de histórias não relacionadas.
- Isolamento total das pastas `.claude`, `.antigravity` e `.codex`.
