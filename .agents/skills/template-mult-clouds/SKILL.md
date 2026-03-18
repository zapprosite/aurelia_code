---
description: Automatiza o setup inicial de workspaces multi-agente de elite.
id: template-mult-clouds
---

# 🏗️ Skill: template-mult-clouds

Automatiza a implantação do padrão de autoridade única e governança de 10 regras.

<directives>
1. **Checklist de Implantação**:
   - [ ] Criar `AGENTS.md` como fonte de verdade.
   - [ ] Implantar as **10 regras fundamentais** em `.agents/rules/`.
   - [ ] Configurar o `.context/` via `ai-context init`.
   - [ ] Instalar o set de skills de elite (Security, Research, Architect).
2. **Adaptadores**: Garante que arquivos adaptadores para Claude e Codex apontem para a autoridade do `AGENTS.md`.
3. **Consistência**: Verifica se não há placeholders de "TODO" nos arquivos gerados.
</directives>

## Verificação de Sucesso
O repositório é considerado "Elite" quando um agente estranho consegue entender toda a governança e arquitetura lendo apenas o `README.md` e o `AGENTS.md`.
