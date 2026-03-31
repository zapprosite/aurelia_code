---
description: Define os limites de permissão automática baseados em risco.
id: 03-tiers-autonomy
---

# 🛡️ Regra 03: Tiers de Autonomia (Risco)

A execução é governada por níveis de risco definidos em `AGENTS.md` § 5.

> **Diretiva do Humano (2026-03-20):** Autonomia total habilitada (`sudo=1`).
> O operador mantém backup completo e aceita os riscos operacionais.

<directives>
1. **Tier A (Read-only)**: Auto-approve 100%.
2. **Tier B (Local Edit)**: Auto-approve 100%. Preferência por Worktrees para isolamento.
3. **Tier C (High-risk)**: Auto-approve com **log obrigatório** para:
   - Modificações em rede/firewall.
   - Gestão de segredos e chaves API.
   - Operações de deploy ou exclusão em massa.
   - Comandos `sudo`.
4. **Segurança Compensatória**:
   - Todo comando `sudo` deve ser registrado em log estruturado.
   - Dry-run sempre que possível para `docker-compose` e scripts bash.
   - Auditoria de segredos antes de `git push` continua ativa.
</directives>
---
description: Define o papel do diretório .context como memória efêmera.
id: 05-context-state
---

# 🧠 Regra 05: .context como Estado, Não Política

O diretório `.context/` é memória operacional, não fonte de "leis".

<directives>
1. **Evidência**: Use para armazenar evidências de testes, logs de auditoria e planos de tarefa.
2. **Subordinação**: O conteúdo do `.context/` nunca substitui as regras de `.agents/rules/`.
3. **Higiene Obrigatória por Slice**: Execute `sync-ai-context` ao final de toda mudança estrutural, slice não trivial, handoff relevante ou preparação para merge.
4. **Dispensa de Baixo Impacto**: Em mudanças triviais sem drift semântico relevante (typo, comentário, rename local sem impacto, teste pontual), a sincronização pode ser dispensada.
5. **Comando Canônico**: A forma padrão é `./scripts/sync-ai-context.sh`, seguida de revisão manual dos `.context/docs/*.md` impactados quando houver drift curatorial.
</directives>
---
description: Garante que o design preceda a execução técnica.
id: 06-planning-first
---

# 📐 Regra 06: Planejamento Antes do Código

Nenhum código significativo deve ser escrito sem um plano aprovado.

<directives>
1. **Atitude**: Crie sempre `.context/plans/<slice>/implementation_plan.md` para tarefas complexas.
2. **Impacto**: Avalie efeitos colaterais em arquivos de configuração e infraestrutura.
3. **Aprovação**: Busque feedback do usuário ou valide contra `.context/plans/` e `AGENTS.md` antes de iniciar a "escrita pesada".
</directives>
