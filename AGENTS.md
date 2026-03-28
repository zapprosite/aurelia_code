# AGENTS.md

Contrato operacional canônico do repositório `aurelia`.

## Autoridade

1. Operador humano
2. `AGENTS.md`
3. [`docs/governance/REPOSITORY_CONTRACT.md`](docs/governance/REPOSITORY_CONTRACT.md)
4. [`docs/governance/SKILL-CATALOG.md`](docs/governance/SKILL-CATALOG.md)
5. [`.agent/rules/README.md`](.agent/rules/README.md)
6. [`docs/adr/README.md`](docs/adr/README.md)
7. [`.context/docs/README.md`](.context/docs/README.md)

## Leitura obrigatória

1. [`AGENTS.md`](AGENTS.md)
2. [`README.md`](README.md)
3. [`docs/governance/REPOSITORY_CONTRACT.md`](docs/governance/REPOSITORY_CONTRACT.md)
4. [`docs/governance/SKILL-CATALOG.md`](docs/governance/SKILL-CATALOG.md)
5. [`.agent/rules/README.md`](.agent/rules/README.md)
6. [`.agent/workflows/README.md`](.agent/workflows/README.md)
7. [`.agent/skills/README.md`](.agent/skills/README.md)
8. [`docs/adr/README.md`](docs/adr/README.md)
9. [`.agent/rules/16-testing-governance.md`](.agent/rules/16-testing-governance.md)

## Estrutura canônica

- Regras executáveis: [`.agent/rules/`](.agent/rules)
- Workflows do workspace: [`.agent/workflows/`](.agent/workflows)
- Catálogo canônico de skills: [`.agent/skills/`](.agent/skills)
- Contratos duráveis: [`docs/governance/`](docs/governance)
- ADRs: [`docs/adr/`](docs/adr)
- Contexto operacional: [`.context/`](.context)

`aurelia` não usa `.agents` como fonte canônica neste repositório. Qualquer referência a `.agents` deve ser tratada como drift legado e corrigida.

## Skills

- A fonte de verdade para skills versionadas é [`.agent/skills/`](.agent/skills).
- Overlays opcionais continuam existindo em `~/.aurelia/skills` e `<repo>/.aurelia/skills`, mas não são a autoridade do catálogo.
- Em colisão de nome, o catálogo canônico do repo deve vencer o overlay.
- Toda skill versionada deve ter `SKILL.md` com frontmatter válido.
- O índice semântico de skills usa a collection `aurelia_skills` no Qdrant com embeddings locais e chunking por seções do markdown.
- Auditoria do catálogo deve detectar drift estrutural, links quebrados e nomes duplicados.

## Adapters por motor

Estes arquivos são adapters finos e não devem duplicar a governança:

- [`CLAUDE.md`](CLAUDE.md)
- [`GEMINI.md`](GEMINI.md)
- [`MODEL.md`](MODEL.md)
- [`.github/copilot-instructions.md`](.github/copilot-instructions.md)

Todos devem delegar para `AGENTS.md` e para o catálogo em `docs/governance/SKILL-CATALOG.md`.

## Guardrails

- Alterações estruturais exigem ADR quando afetarem runtime, memória, modelos, storage, voz, segurança, deploy ou governança.
- Não alterar [`.agent/rules/`](.agent/rules) sem ordem direta de Will.
- Manter `.env` e `.env.example` em paridade estrutural.
- Antes de mudanças de rede, portas ou subdomínios, consultar `NETWORK_MAP.md` em `/srv/ops/ai-governance/`.
- A política de modelos é imutável sem ADR: [`docs/governance/MODEL-STACK-POLICY.md`](docs/governance/MODEL-STACK-POLICY.md)

## Referências rápidas

- Contrato do repositório: [`docs/governance/REPOSITORY_CONTRACT.md`](docs/governance/REPOSITORY_CONTRACT.md)
- Catálogo de skills: [`docs/governance/SKILL-CATALOG.md`](docs/governance/SKILL-CATALOG.md)
- Skill `/add-subdomain`: [`.agent/skills/add-subdomain/SKILL.md`](.agent/skills/add-subdomain/SKILL.md)
- Workflow `//test-all`: [`.agent/workflows/test-all.md`](.agent/workflows/test-all.md)
- Workflow nonstop: [`.agent/workflows/adr-semparar.md`](.agent/workflows/adr-semparar.md)
- Voz local: [`docs/jarvis_local_voice_blueprint_20260319.md`](docs/jarvis_local_voice_blueprint_20260319.md)

## Serviços de Sistema (SOTA 2026)
- **Aurelia System API**: Gateway de governança em Go na porta `8080`. Gerencia paridade de `.env` e saúde do cluster.
- **Smart Router**: Roteamento inteligente na porta `4000` (LiteLLM).
- **QvC Monitor**: Auditoria horária de Qualidade vs Custo (Cron).

## AI Context References
- Documentation index: `.context/docs/README.md`
- Agent playbooks: `.context/agents/README.md`

