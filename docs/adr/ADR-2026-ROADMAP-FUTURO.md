# ADR-2026-ROADMAP-FUTURO: Roadmap Estratégico e Slices Pendentes

**Status:** 🔄 Ativo
**Autoridade:** Aurélia (Arquiteta Principal)
**Foco:** Autonomia Total e Cognição Avançada

## 1. Objetivo Estratégico
Elevar a Aurelia ao estado de "Jarvis Local-First", com capacidade de auto-resolução de problemas complexos, introspecção de ferramentas e gerenciamento de contexto semântico.

## 2. Slices Prioritários (Próximas Fases)

### S-15: Tool Introspection System
- **Objetivo**: Permitir que a Aurelia entenda e filtre suas próprias ferramentas dinamicamente baseada na tarefa.
- **Componentes**: `internal/agent/tool_catalog.go`, integração com Qdrant para matching semântico.

### S-16: Execution DNA
- **Objetivo**: Templates de workflow por tipo de tarefa (debug, feature, refactor).
- **Componentes**: `internal/persona/execution_dna.go`.

### S-17: Planning Loop (PREV Phase)
- **Objetivo**: Implementar o loop de Plano -> Revisão -> Execução -> Verificação nativo no binário Go.
- **Garantia**: Bloqueio de segurança para planos de alto risco sem aprovação humana.

### S-18: Codebase Symbol Map
- **Objetivo**: Parseamento de símbolos Go (.ast) para que a Aurelia localize funções e tipos sem busca por texto puro.
- **Componentes**: `internal/agent/codebase_map.go`.

### S-19: Semantic Skill Router
- **Objetivo**: Roteamento inteligente de habilidades via embeddings, substituindo o roteador atual baseado em strings.

## 3. Planos Futuros e Inovação
- **Memory Context Assembler**: Assemblagem dinâmica de contexto (Qdrant + Git + Code) para alimentar a LLM com precisão máxima.
- **Autonomous HW Management**: Auto-gerenciamento de VRAM e containers via `homelab-control`.
- **Global Auth Proxy**: Unificação de autenticação entre Dashboard, Telegram e API.

## 4. Critérios de Sucesso
- Autonomia Level 5 (Human-in-the-loop apenas para aprovação estratégica).
- Latência de resposta < 2s para ferramentas locais.
- Zero drift semântico entre código e documentação.
