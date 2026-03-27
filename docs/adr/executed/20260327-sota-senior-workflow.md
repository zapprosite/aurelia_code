# ADR 20260327: SOTA Senior Workflow

## Contexto
O usuário Will solicitou uma persona de desenvolvedor sênior, "direto ao ponto", com integração profunda entre Obsidian, Qdrant, Supabase e SQLite, utilizando comandos `/` como interface principal.

## Decisão
Implementar o padrão **SOTA 2026 Senior Workflow**:
1. **Persona**: Tom técnico, minimalista, focado em código e execução. Evitar explicações longas a menos que solicitado.
2. **Interface**: Priorizar comandos `/` para operações de sistema e governança.
3. **Memória de 4 Camadas**:
   - **L0 (Ephemeral)**: Context window do LLM.
   - **L1 (Local)**: SQLite para transações e estados imediatos.
   - **L2 (Semantic)**: Qdrant para busca vetorial de conhecimento e histórico.
   - **L3 (Galactic)**: Supabase para persistência de longo prazo e acesso multi-agente.
4. **Conhecimento**: Obsidian como fonte fria de documentação e notas de pesquisa, sincronizado via `internal/obsidian`.

## Consequências
- Maior velocidade de resposta (menos tokens desperdiçados em chat).
- Melhor rastreabilidade de decisões históricas via Qdrant/Supabase.
- Curva de aprendizado inicial para novos comandos, compensada pela produtividade.
