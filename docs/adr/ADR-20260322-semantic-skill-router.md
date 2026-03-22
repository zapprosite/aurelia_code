# ADR-20260322-semantic-skill-router: Sub-5 Semantic Skill Router

## Status
Ativo

## Contexto
Atualmente, o Agent Loop (Aurelia) injeta todas as definições das skills e ferramentas disponíveis no _System Prompt_ via `ToolRegistry`. Conforme o número de Skills instaladas crescer (20+), isso resultará em _Context Bloat_, estourando os limites de token, encarecendo a inferência e gerando alucinações onde a IA confunde qual ferramenta usar.

## Decisão
Implementar o `internal/skill/semantic_router.go`.
O Semantic Router irá:
1. Usar o Ollama `embed` local para vetorizar a intenção do usuário (`query`).
2. Consultar o Qdrant (coleção `aurelia_skills`) para retornar os IDs apenas do **Top-K** das Skills mais relevantes.
3. Fornecer ao `agent.Loop` apenas esse subset afunilado de ferramentas.
*Nota: Ferramentas nativas do SO (run_command, read_file, etc) mantêm-se como "Core Tools" e escapam desta restrição.*

## Consequências
- Economia massiva de Tokens no Zero-Shot.
- Redução de alucinações (LLM focará apenas em 3~5 skills fortemente relacionadas ao problema visualizado).
- Adiciona uma latência imperceptível (< 50ms) no pre-flight do Loop.

## Testes e Rollout
1. Implementar `SemanticRouter` exportando `MatchForTask(query string, limit int) []string`.
2. O catálogo de ferramentas atual do loop já foi parcialmente preparado no Sub-1 (`ToolCatalog`/`catalog.MatchForTask`). Vamos conectá-lo ao Qdrant.
3. Testar a injeção em modo dry-run / log debug com um pedido real.
