# Skills

Catálogo canônico de skills versionadas do repositório `aurelia`.

## Fonte de verdade

- Skills oficiais moram em subdiretórios de [`.agent/skills/`](.).
- Cada skill deve ter `SKILL.md` com frontmatter válido.
- O nome no frontmatter é a identidade lógica da skill.
- O diretório do repo vence overlays locais em caso de colisão de nome.

## Overlays aceitos

- Global opcional: `~/.aurelia/skills`
- Project overlay opcional: `<repo>/.aurelia/skills`

Esses overlays existem para extensões locais e instalações temporárias. Eles não substituem o catálogo canônico do repo.

## Indexação semântica

- A collection canônica de skills no Qdrant é `aurelia_skills`.
- O sync usa embeddings locais e chunking por seções do `SKILL.md`.
- O payload indexado deve preservar `name`, `description`, `section`, `path`, `chunk_id` e `checksum`.

## Auditoria

O catálogo deve permanecer limpo:

- sem links quebrados
- sem frontmatter inválido
- sem drift `.agents/` nos entrypoints de governança
- sem duplicação silenciosa de nomes entre fonte canônica e overlays

## Referências

- Governança principal: [`../../AGENTS.md`](../../AGENTS.md)
