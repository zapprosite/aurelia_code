---
name: researcher
description: Habilita pesquisa profunda e síntese de informações complexas para tomada de decisão técnica.
phases: [P]
---

# 🕵️ Researcher: Sovereign Analysis 2026

Habilita o Antigravity a realizar investigações técnicas exaustivas para embasar decisões de arquitetura e resolução de problemas.

## 🏛️ Metodologia de Pesquisa (Triple-Tier)

### 1. Pesquisa Externa (Core Integration)
- **Ferramenta**: Utilize o MCP `context7` para documentação de bibliotecas.
- **Navegação**: Use o agente `playwright` para ler documentações oficiais diretamente de sites (MDN, Go Docs, Next.js Docs).
- **Provedor**: Utilize o Tier 1 (MiniMax/Claude) para sintetizar as descobertas.

### 2. Pesquisa Interna (Aurelia Search)
- **Ferramenta**: Utilize `mcp_ai-context_search` e o Vector DB local para encontrar precedentes no monorepo.
- **Referência**: Consulte sempre `docs/adr/` para entender o "porquê" de decisões passadas.

## 🚀 Workflow de Síntese
1. **Gathering**: Reúna evidências de múltiplas fontes (Interno + Externo).
2. **Analysis**: Identifique conflitos, oportunidades e riscos (SWOT).
3. **Report**: Gere um relatório conciso com recomendações acionáveis.

## 📍 Quando usar
- Antes de iniciar um novo slice de tecnologia desconhecida.
- Para comparar bibliotecas ou frameworks antes de um ADR.
- Para descobrir a causa raiz de comportamentos bizarros do sistema.