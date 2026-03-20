---
title: "Padronizar o Repositório como Template Multi-Agente Senior"
status: proposed
owner: antigravity
created: 2026-03-20
scope: repositório inteiro
slice: padronização
---

# ADR-20260320: Padronizar o Repositório como Template Multi-Agente Senior

## Contexto

O repositório `aurelia` evoluiu organicamente de "bot Telegram em Go" para um sistema multi-agente completo com:
- 15 pacotes `internal/` (agent, config, cron, gateway, health, heartbeat, mcp, memory, observability, persona, runtime, skill, telegram, tools, voice)
- 3 pacotes `pkg/` (llm, stt, tts)
- 22 skills, 30+ workflows, 24 ADRs
- Governança madura (AGENTS.md, REPOSITORY_CONTRACT, regras, tiers de autonomia)

Porém acumulou dívida técnica que impede servir como **template profissional**:

### Diagnóstico (2026-03-20)

| Área | Estado | Problema |
|---|---|---|
| **Build** | 🟢 | Compila sem erros |
| **Testes** | 🟡 | 3 pacotes falhando: `pkg/llm`, `internal/tools`, `cmd/aurelia` + 1 teste em `internal/agent` |
| **Raiz** | 🔴 | 2 binários compilados (~80MB) rastreados pelo Git |
| **Hardware Spec** | 🔴 | `plan.md` referencia RTX 4090 mas hardware real é AMD 7900X + 32GB DDR5 Gen5 |
| **Governança** | 🟢 | Sólida, mas com referências cruzadas quebradas em `docs/governance/` |
| **.gitignore** | 🟡 | Não cobre binários compilados na raiz |

## Decisão

Transformar o repositório em template **Elite Multi-Agente** nível dev senior:

### 1. Poda da Raiz (sem quebrar)
- Remover binários compilados (`aurelia`, `aurelia-elite`) do tracking Git
- Adicionar pattern ao `.gitignore`

### 2. Corrigir Testes Quebrando
- `pkg/llm/openai_compatible_test.go:132` — syntax error (`<<` encontrado)
- `internal/tools/service_control_test.go` — `blocksSelfServiceMutation` undefined
- `cmd/aurelia/main_test.go` — `shouldSuppressDuplicateLaunch`, `recordDuplicateLaunch`, `exitCodeForBootstrapError` undefined
- `internal/agent` — teste de dependent subtask falhando

### 3. Atualizar Hardware Spec
- Hardware real: AMD Ryzen 9 7900X, 32GB DDR5 Gen5, 4TB Gen3 VFS, 1TB Gen3 SSD Ubuntu
- GPU: AMD Radeon RX 7900 XTX com 24GB VRAM
- Ajustar budget de modelos locais para VRAM real

### 4. Template Multi-Agente
- Skill `template-mult-clouds` com checklist verificável
- README como porta de entrada para qualquer agente estranho
- Sem placeholders "TODO" nos contratos

## Autonomia Total (sudo=1)

> **Diretiva do Humano:** Todos os agentes têm permissão sudo total.
> O operador aceita os riscos e mantém backup completo do sistema.

### Impacto nos Tiers de Governança

| Tier | Antes | Depois |
|---|---|---|
| **A (Read-only)** | Auto-approve 100% | Auto-approve 100% |
| **B (Local Edit)** | Auto-approve condicional | Auto-approve 100% |
| **C (High-risk)** | Aprovação Humana OBRIGATÓRIA | Auto-approve com log obrigatório |

### Regras de Segurança Compensatórias

1. **Log obrigatório** — todo comando `sudo` executado deve ser registrado em log estruturado
2. **Backup existe** — o humano garante backup completo e aceita os riscos
3. **Dry-run primeiro** — para `docker-compose`, scripts bash e configurações de SO, executar dry-run quando possível (regra mantida do contrato global)
4. **Auditoria de segredos** — continua ativa antes de qualquer `git push`

### Atualizar nos Contratos

- `.agents/rules/03-tiers-autonomy.md` — remover gate de aprovação humana do Tier C
- `AGENTS.md` — seção de Governança Tiers deve refletir `sudo=1`
- `GEMINI.md`, `CLAUDE.md`, `CODEX.md` — reconhecer autonomia total

## Consequências

### Positivas
- Qualquer agente (Antigravity, Claude, Codex, OpenCode) entra e entende a governança
- Suíte de testes 100% verde
- Raiz limpa e profissional
- Hardware spec refletindo a realidade operacional
- Agentes podem operar sem interrupção (sudo sem gate)

### Negativas
- Requer validação manual do hardware spec pelo humano (GPU AMD vs NVIDIA, modelos Ollama compatíveis)
- Modelos locais podem precisar de reavaliação (Ollama no ROCm vs CUDA)
- Risco maior em operações destrutivas (mitigado por backup e logs)

## Plano de Execução

1. Poda e `.gitignore` (auto-approve)
2. Corrigir testes quebrando (auto-approve)
3. Atualizar hardware no `plan.md` (auto-approve com log)
4. Atualizar tiers de autonomia com `sudo=1` (auto-approve com log)
5. Atualizar skill e README (auto-approve)
6. `sync-ai-context` final

## Referências

- `AGENTS.md` — contrato soberano
- `docs/governance/REPOSITORY_CONTRACT.md` — índice de governança
- `plan.md` — plano mestre JARVIS
- `docs/adr/INDEX.md` — índice de ADRs
