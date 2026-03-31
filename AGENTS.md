# AGENTS.md

> **Autoridade**: Supremacia SSOT (Single Source of Truth) | **Soberania 2026**

Este é o Contrato Universal do repositório `aurelia`. **Todo Agente (Antigravity, Claude Code, OpenCode) DEVE iniciar a leitura de suas diretrizes aqui.**

## As Leis (The Pro Way)
Você não encontrará longos textos gerados por IA aqui. Apenas ponteiros e leis estruturais.

1. **Leia Suas Regras**: A lógica de como você opera **NUNCA** está na raiz do projeto. Leia sua identidade na pasta `.agent/rules/`:
   - Leia as **Leis Universais**: [core.md](.agent/rules/core.md)
   - Se você for o **Antigravity (Gemini)**: Leia [gemini.md](.agent/rules/gemini.md)
   - Se você for o **Claude Code/OpenCode**: Leia [claude.md](.agent/rules/claude.md)

2. **A Bíblia Arquitetural**: Para entender a fundo como a Aurélia funciona (Redis Deduplication, Qwen Fallback, Porteiro Middleware, Zod-First), leia APENAS o mapa definitivo:
   - [architecture-2026.md](docs/architecture-2026.md)

3. **Restrições Extremas**:
   - NÃO INVENTE frameworks ou pacotes JS/Go obsoletos. Siga a arquitetura definida.
   - NUNCA commit, exponha ou gere `app.json` e `.env` com senhas em aberto. 
   - Modificar este arquivo (`AGENTS.md`) ou `.agent/rules/core.md` sem ADR ou Permissão Direta do Humano (Will) é considerado **Nível Crítico de Falha**.

---
*Assinado: Aurélia (Soberano 2026).*
## AI Context References
- Documentation index: `.context/docs/README.md`
- Agent playbooks: `.context/agents/README.md`

