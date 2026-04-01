# ADR — Authority Documents

**Diretório:** `docs/adr/`
**Última atualização:** 01/04/2026
**Próximo slice:** S-50

---

## Como criar um novo slice

```bash
./scripts/adr-slice-init.sh <slug> --title "Título do Slice"
# Gera: docs/adr/YYYYMMDD-<slug>.md
```

Ou manualmente usando TEMPLATE-NONSTOP-SLICE.md.

## Índice

| Arquivo | Slice | Status |
|---|---|---|
| 0001-HISTORY.md | S-33 a S-49 | ✅ Histórico |
| PENDING.md | — | 🔵 Vazio |
| 20260331-architecture-mcp-a2a.md | S-35 | Ativo |
| 20260331-duplicate-response-fix.md | S-34 | ✅ |
| 20260331-remover-redis-qwen-porteiro.md | S-36 | ✅ |
| 20260331-sota-industrializacao.md | S-37 | ✅ |
| 20260331-telegram-fix-s34.md | S-33 | ✅ |
| 20260331-polda-2-industrializacao.md | — | 🗄️ Arquivado |
| 20260331-telegram-duplicate-fix.md | — | 🗄️ Obsoleto |
| 20260401-smart-router-homelab-cron.md | S-39 | ✅ |
| 20260401-industrializacao-opencode.md | S-38 | ✅ |

## Model Stack Policy (Ativo — 01/04/2026)

Cascade oficial (referência: docs/governance/MODEL-STACK-POLICY.md):

  Nível 0-local: gemma3:27b (ops/cron/fiscal — dados não saem do host)
  Nível 1-free:  nemotron-3-super-120b:free → qwen3.6-plus:free
  Nível 2-pago:  minimax-m2.7 → glm-5.1 → kimi-k2.5 (long-context)
  Embed:         nomic-embed-text (local)

## Cadeia de Autoridade

AGENTS.md → docs/governance/REPOSITORY_CONTRACT.md → docs/adr/

## Referências

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../governance/REPOSITORY_CONTRACT.md)
- [MODEL-STACK-POLICY.md](../governance/MODEL-STACK-POLICY.md)
