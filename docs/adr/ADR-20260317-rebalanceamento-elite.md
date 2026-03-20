# ADR: 20260317-rebalanceamento-elite

**Status**: Aceito
**Data**: 2026-03-17
**Contexto**: O repositório estava acumulando placeholders vazios e carecia de uma estrutura de metadados clara para agentes de IA (YAML/XML).

## Decisão
Implementamos o padrão **Lean Documentation** em todo o workspace:
1.  **YAML Frontmatter**: Obrigatório em todos os arquivos Markdown para otimizar a descoberta de contexto por LLMs.
2.  **XML Tags**: Utilizadas para delimitar seções críticas (regras, autoridade, comandos).
3.  **Tríade de Pastas**: Consolidação em `.agents` (Governança), `docs` (Arquitetura) e `.context` (Execução).

## Consequências
- **Positivas**: Redução drástica no custo de tokens (agentes lêem apenas o necessário). Maior assertividade na execução de comandos slash.
- **Negativas**: Requer disciplina manual para manter o frontmatter atualizado em novos documentos.
