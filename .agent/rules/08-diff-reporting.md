---
description: Exige provas factuais de mudanças via git diff e logs.
id: 08-diff-reporting
---

# 📊 Regra 08: Relatórios Baseados em Diff

Toda mudança deve ser acompanhada de prova visual e técnica.

<directives>
1. **Prova de Fato**: Use `render_diffs` e logs de terminal para demonstrar sucessos.
2. **Específicos**: Evite descrições genéricas. Cite números de linhas e nomes de funções alteradas.
3. **Segurança de Segredos**: Audite o diff para garantir que nenhuma credencial vazou no commit.
</directives>
