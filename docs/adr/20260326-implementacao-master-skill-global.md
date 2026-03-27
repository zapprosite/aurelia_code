# ADR-20260326-Implementacao-Master-Skill-Global

**Status:** ✅ Proposta
**Data:** 2026-03-26
**Autor:** Antigravity AI
**Slices Context:** Cross-Cutting / Tooling

---

## 🟢 Contexto

Atualmente, o monorepo `aurelia` e seus sub-projetos sofrem com a duplicação de *skills* (arquivos `SKILL.md`) em múltiplas pastas `.agents`, `.claude` e `.context`. Esta redundância gera:
1.  **Dificuldade de Manutenção**: Atualizações em uma *skill* comum não se propagam automaticamente.
2.  **Alto Custo de Tokens**: Agentes de IA indexam repetidamente as mesmas instruções em diferentes contextos, inflando o consumo de tokens.
3.  **Falta de Orquestração**: Não há um ponto central para inicializar *frameworks* de agentes (como Antigravity Kit ou BMED) de forma padronizada.

O problema de "Skill Fatigue" e dispersão de ferramentas exige uma camada de orquestração global.

## 🔵 Decisão

Adotamos a **Master Skill Global** como a ferramenta central de governança e orquestração de ambientes de desenvolvimento para agentes de IA (Antigravity, Claude Code e Codex).

### 1. Instalação e Invocação Global
A *skill* não reside nos projetos individuais, mas no nível global do agente:
-   **Antigravity**: `~/.gemini/antigravity/skills/master-skill/`
-   **Claude**: `~/.claude/skills/master-skill/`
-   **Invocação**: Comando prioritário `/master skill` (ou `/master-skill`), configurado via definições em `.agent/commands/`.

### 2. Gestão de Frameworks
A Master Skill possui lógica de bootstrap para instalar automaticamente os seguintes *frameworks* de projeto:
-   **Antigravity Kit**: Ferramental de design e UI.
-   **BMED (BMad Method)**: Metodologia de desenvolvimento estruturado.
-   **Specit (Spec-Kit)**: Motor de especificações técnicas.

### 3. Importação Seletiva de Skills Externas
Para evitar a replicação por projeto, a Master Skill implementa um "Skill Store" local:
-   Lê um diretório central de *skills* do sistema (ex: `~/dev/skills/`).
-   Extrai metadados (YAML frontmatter) para permitir busca por intenção.
-   **Importação sob Demanda**: O agente importa (via symlink ou cópia) apenas a *skill* estritamente necessária para a tarefa atual, minimizando o contexto enviado ao LLM.

### 4. Configuração Persistente (`/master skill init`)
O comando de inicialização inicializa o ambiente e salva os parâmetros em `config/settings.json`:
-   **Inputs**: Identificação do agente em uso e o caminho absoluto da pasta global de *skills*.
-   **Persistência**: Evita a necessidade de repetir configurações em sessões futuras.

## 🔴 Consequências

-   **Eficiência de Tokens**: Redução estimada de 30-50% na indexação inicial de projetos complexos ao evitar o carregamento de *skills* irrelevantes.
-   **Soberania de Ferramental**: Padronização sênior de como agentes interagem com o sistema de arquivos e *skills*.
-   **Ambiente Unificado**: Experiência consistente, independentemente de estarmos usando Claude CLI, Antigravity IDE ou Codex.
-   **Manutenção Centralizada**: Uma única fonte de verdade para as lógica de orquestração global.

---
**Referências:**
- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../../docs/governance/REPOSITORY_CONTRACT.md)
- [Sandeco Master-Skill Prompts](https://github.com/sandeco/prompts/tree/main/master-skill)
