# ADR-20260331: Arquitetura Telegram SOTA 2026.2 + Roadmap MCP/A2A

**Status:** Ativo
**Data:** 2026-03-31
**Slices:** S-33 (dedup), S-34 (TTS), S-35 (MCP/A2A roadmap)

---

## 1. Estado Atual do Sistema

### Stack de Produção

```
Telegram Bot (Go)
├── Porteiro Middleware (Qwen 0.5B + Redis)
│   ├── IsSafe() → Verifica input
│   ├── PolishOutput() → JSON → Markdown
│   ├── SecureOutput() → Mascaramento de segredos
│   └── Deduplicate() → Anti-retry
├── LiteLLM Router
│   ├── Tier 1: Gemini / Claude (cloud)
│   ├── Tier 2: OpenRouter fallback
│   └── Tier 0: Ollama local (Qwen 2.5)
├── Agent Loop (Go)
│   ├── Tool Registry (filesystem, docker, etc)
│   └── Memory Assembler (Qdrant RAG)
├── Kokoro TTS (GPU)
└── Redis (state) + Qdrant (vectors)
```

### Bugs Corrigidos (31/03/2026)

| Slice | Problema | Solução |
|-------|----------|---------|
| S-33 | Mensagens duplicadas (retry Telegram) | Redis SetNX com lock 10s |
| S-34 | Respostas duplicadas (texto + nota de voz) | TTS só quando usuário enviou voz |
| S-35 | Loop infinito por markdown_brain_sync | Removido do default tools (cron job é suficiente) |

---

## 2. Arquitetura Telegram — Fluxo Atual

```
[Mensagem Telegram]
       │
       ▼
[Whitelist Middleware] ──→ Bloqueia usuários não autorizados
       │
       ▼
[Anti-Retry Deduplication] ──→ Redis: porteiro:v2:dupe:{user}:{msg_id}
       │
       ▼
[Porteiro IsSafe] ──→ Qwen 0.5B: input seguro?
       │               Redis cache: 30 dias
       ▼
[Memory Context] ──→ Qdrant: conversation_memory + markdown_brain
       │
       ▼
[LLM Router] ──→ LiteLLM (Gemini/Claude/Ollama)
       │
       ▼
[Tool Execution] ──→ read_file, run_command, docker, etc.
       │
       ▼
[Porteiro PolishOutput] ──→ JSON → Markdown 2026
       │
       ▼
[Porteiro SecureOutput] ──→ Regex masking (sk-, ghp_, AUR_)
       │
       ▼
[deliverFinalAnswer] ──→ Texto (se input=text) ou Texto+Voz (se input=voice)
```

---

## 3. Roadmap MCP + A2A (SOTA 2026)

### 3.1 MCP (Model Context Protocol)

**Relevância:** O daemon Go já tem um `mcp.Manager` em `internal/mcp/`. A arquitetura atual permite ferramentas MCP via skills.

**Próximos Passos:**
- [ ] Migrar ferramentas internas (filesystem, docker) para MCP SDK
- [ ] Expor ferramentas via protocolo MCP para agents externos
- [ ] Integrar com Claude Code via MCP (já suportado)

**Fontes:**
- [MCP Specification](https://modelcontextprotocol.io/specification/2025-11-25)
- [Claude Code MCP Integrations](https://www.truefoundry.com/blog/claude-code-mcp-integrations-guide)

### 3.2 A2A (Agent-to-Agent Protocol)

**Relevância:** O swarm de agents (`/master-skill`, squads) precisa de comunicação padronizada.

**Estado Atual:**
- `handoff_to_agent` tool existe (`internal/agent/loop.go:220`)
- Task store SQLite para coordenação (`internal/agent/task_store.go`)
- Mailbox para comunicação entre agents (`internal/tools/team_mailbox.go`)

**Próximos Passos:**
- [ ] Avaliar adoção do protocolo A2A do Google (anunciado Abr/2025)
- [ ] Padronizar interface de handoff entre agents
- [ ] Integrar com Anthropic Vertex AI para multi-agent

**Fontes:**
- [A2A + MCP Protocol Stack](https://pub.towardsai.net/mcp-a2a-thanks-to-skills-the-complete-protocol-stack-your-multi-agent-system-needs-2757b2028b9f)
- [Anthropic Multi-Agent Webinar](https://www.anthropic.com/webinars/deploying-multi-agent-systems-using-mcp-and-a2a-with-claude-on-vertex-ai)

---

## 4. Limpeza Documental

### Arquivos a Manter

```
docs/
├── architecture-2026.md          ← Mapa canônico do sistema
├── STYLE_GUIDE.md                ← Convenções de código
├── governance/
│   ├── REPOSITORY_CONTRACT.md   ← Contrato do repo
│   ├── SKILL-CATALOG.md         ← Índice de skills
│   ├── MODEL-STACK-POLICY.md    ← Stack de modelos
│   └── PORTEIRO_SENTINEL_2026.md
├── adr/
│   ├── README.md                ← Índice de ADRs
│   └── [YYYYMMDD]-*.md         ← ADRs cronológicos
└── reports/
    └── [date]-repo-health.md
```

### Arquivos Obsoletos (Arquivar)

```
docs/archive/                          ← Plans legados
docs/governance/_archive/               ← Contratos superseded
.agent/rules/_archive/                  ← Regras legadas
.agent/skills/_archive/                 ← Skills legadas
.agent/workflows/_archive/              ← Workflows legados
.context/runbooks/                     ← Mover para docs/runbooks/
```

---

## 5. Decisões Arquiteturais

| ADRs Ativos | Assunto |
|------------|---------|
| 20260331-duplicate-response-fix | Correção S-34 (TTS) |
| 20260331-telegram-fix-s34 | Deduplicação S-33 |
| 20260331-sota-industrializacao | Industrialização 2026 |
| **Este ADR** | Arquitetura + MCP/A2A roadmap |

---

## 6. Referências

- [Agentic Web: AGENTS.md, MCP vs A2A](https://www.nxcode.io/resources/news/agentic-web-agents-md-mcp-a2a-web-4-guide-2026)
- [AI Agent Protocols 2026 Complete Guide](https://www.ruh.ai/blogs/ai-agent-protocols-2026-complete-guide)
- [MCP + A2A: Complete Protocol Stack](https://pub.towardsai.net/mcp-a2a-thanks-to-skills-the-complete-protocol-stack-your-multi-agent-system-needs-2757b2028b9f)
