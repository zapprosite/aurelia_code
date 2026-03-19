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
