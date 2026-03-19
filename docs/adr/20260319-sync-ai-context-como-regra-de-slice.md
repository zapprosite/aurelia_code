# ADR 20260319-sync-ai-context-como-regra-de-slice

**Status**: Aceito  
**Data**: 2026-03-19

## Contexto

O repositório Aurelia depende de `.context/` como memória operacional curta, evidência de execução e ponte entre código, documentação e handoff entre agentes. Antes desta decisão, o `sync-ai-context` já era usado como prática recomendada, mas ainda podia ser interpretado como ritual genérico "no fim de qualquer feature", o que gera ruído em microedições triviais e fragiliza a disciplina quando o escopo realmente importa.

Ao mesmo tempo, mudanças estruturais em `cmd/`, `internal/`, `pkg/`, `scripts/`, `docs/` e nos blueprints do runtime demonstraram que deixar o `.context/` sem sincronização aumenta o risco de drift semântico entre:

- código real
- `.context/docs/*.md`
- `codebase-map.json`
- handoffs entre Antigravity, Claude e Codex

## Decisão

O repositório passa a tratar `sync-ai-context` como **regra obrigatória por slice não trivial**, e não como ritual cego em qualquer microedição.

Obrigatório em:

- mudanças estruturais
- slices não triviais
- handoff relevante entre agentes/motores
- preparação para review/merge final

Dispensável em:

- typo isolado
- comentário sem impacto comportamental
- rename local sem impacto estrutural
- teste muito pequeno sem drift semântico relevante

Comando canônico:

- `./scripts/sync-ai-context.sh`

## Consequências

- `AGENTS.md` passa a explicitar a regra com escopo profissional
- `.agents/rules/05-context-state.md` deixa de falar em "fim de feature" de modo genérico e passa a falar em slice/handoff/merge
- a skill `.agents/skills/sync-ai-context/SKILL.md` fica alinhada a essa interpretação
- `docs/REPOSITORY_CONTRACT.md` passa a carregar a regra como parte do contrato central do repositório

## Trade-off

O custo operacional do `sync-ai-context` continua existindo, mas agora é aplicado onde ele dá retorno real: slices relevantes, revisão final e transferência de contexto. Isso reduz ruído sem relaxar a governança.

## Referências

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [.agents/rules/05-context-state.md](../../.agents/rules/05-context-state.md)
- [.agents/skills/sync-ai-context/SKILL.md](../../.agents/skills/sync-ai-context/SKILL.md)
