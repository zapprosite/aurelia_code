---
description: Define o papel do diretório .context como memória efêmera.
id: 05-context-state
---

# 🧠 Regra 05: .context como Estado, Não Política

O diretório `.context/` é memória operacional, não fonte de "leis".

<directives>
1. **Evidência**: Use para armazenar evidências de testes, logs de auditoria e planos de tarefa.
2. **Subordinação**: O conteúdo do `.context/` nunca substitui as regras de `.agents/rules/`.
3. **Higiene**: Ao final de uma feature, execute o `mcp ai-context sync` para manter a integridade operacional.
</directives>
