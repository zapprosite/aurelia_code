# 🧪 Smoke Test — Telegram + Homelab + Senior Architect

**Teste completo de integração** que valida:
- ✅ Telegram → Bot pipeline
- ✅ Acesso às 30+ skills
- ✅ Respostas em nível senior/arquiteto
- ✅ Automação segura com confirmação
- ✅ Orquestração multi-step
- ✅ Disaster recovery prontidão

---

## Quick Start

### 1. Shell Script (Sem Telegram)
```bash
./scripts/smoke-test-homelab.sh
```
**Tempo:** 10 segundos
**Valida:** Status real do homelab

### 2. Go Tests
```bash
go test ./e2e -run TestSmoke -v
go test ./internal/skill/ -v
```
**Tempo:** 30 segundos
**Valida:** Code + skills

### 3. Telegram Interativo
```bash
# Terminal 1: Start daemon
go run ./cmd/aurelia/main.go

# Terminal 2: Enviar mensagem
# Bot Telegram: "saúde completa"
```
**Tempo:** Variável (resposta do LLM)
**Valida:** Pipeline completo

---

## Teste Smoke — 6 Cenários

### 1️⃣ Health Check Completo
**Prompt:** "saúde completa do homelab"

**Bot valida:**
- 30+ containers
- GPU VRAM (18GB livre)
- ZFS pool status
- Network tunnel
- Recomendações específicas

**Esperado:** Resposta com métricas + recomendações

---

### 2️⃣ Container Diagnostics
**Prompt:** "diagnóstico do n8n" ou "por que está lento?"

**Bot fornece:**
- Status (Up/Down)
- Logs (erro/warning)
- RAM/CPU
- Recomendação de ação

**Esperado:** Análise + próximos passos

---

### 3️⃣ Architecture Decision
**Prompt:** "deveria mover DB para fora?" ou "otimizar GPU"

**Bot analisa:**
- Option 1: Local (latência <1ms, single point of failure)
- Option 2: Remote (managed, mas 50ms latency)
- Option 3: Hybrid (critical local, analytics remote)

**Esperado:** Trade-off analysis com recomendação justificada

---

### 4️⃣ Safe Automation
**Prompt:** "cria snapshot do tank/data"

**Bot executa:**
1. ⚠️ Pede confirmação (segurança)
2. ✅ Executa operação
3. 📊 Reporta sucesso + detalhes

**Esperado:** Automação com guardrails

---

### 5️⃣ Multi-Step Orchestration
**Prompt:** "subir voice stack completo"

**Bot orquestra:**
1. VRAM Pre-check → ✅ 18GB livre
2. Start Whisper (STT) → ✅ Healthy
3. Start Chatterbox (TTS) → ✅ Healthy
4. Start Proxy → ✅ Respondendo
5. Final validation → ✅ Voice stack pronto

**Esperado:** Progression clara com validações

---

### 6️⃣ Disaster Recovery
**Prompt:** "validar prontidão de recovery"

**Bot valida:**
- Último backup: < 24h ✅
- Snapshots: 63 ✅
- Espaço: 3569GB ✅
- Recomendação: "Pronto para DR"

**Esperado:** DR status completo

---

## Métricas de Sucesso

### ✅ Respostas Senior/Arquiteto
- [ ] Trade-off analysis (2+ opções)
- [ ] Métricas concretas (números)
- [ ] Risk assessment
- [ ] Recomendação justificada
- [ ] Preconditions explícitas

### ✅ Segurança
- [ ] Confirmação para operações perigosas
- [ ] ZFS snapshots antes de mudanças
- [ ] Rollback instructions
- [ ] Health checks após operações

### ✅ Automação
- [ ] Suporta 30 skills
- [ ] Executa skills de forma orquestrada
- [ ] Reporta erros e sucessos
- [ ] Integrado com Telegram

### ✅ Performance
- [ ] Health check: < 5s
- [ ] Diagnóstico: < 10s
- [ ] Arquitetura: < 15s
- [ ] Automação: confirma em < 2s

---

## Estrutura de Arquivos

```
.
├── scripts/
│   └── smoke-test-homelab.sh       # Teste shell (rápido)
├── e2e/
│   └── smoke_test.go               # Teste Go completo
├── internal/persona/
│   └── senior_architect_prompt.go  # Prompts para resposta senior
├── docs/
│   ├── SMOKE_TEST_GUIDE.md         # Guide detalhado
│   └── SMOKE_TEST_README.md        # Este arquivo
└── ~/.aurelia/skills/
    ├── health-check-full/          # 30 skills...
    ├── gpu-vram-audit/
    ├── voice-stack-up/
    └── ... (27 skills mais)
```

---

## Como Rodar

### Local (Sem Telegram)
```bash
./scripts/smoke-test-homelab.sh
# ✅ Output: 6 cenários validados, ~10s
```

### Tests
```bash
go test ./e2e -run TestSmoke -v
go test ./internal/skill/ -v
# ✅ Output: All tests PASS, ~30s
```

### Interativo (Com Telegram)
```bash
# Terminal 1:
systemctl start aurelia --user

# Terminal 2 (Bot Telegram):
"saúde completa"
→ Resposta senior com métricas + recomendações
```

---

## Integração CI/CD

Adicionar ao `.github/workflows/ci.yml`:

```yaml
- name: Smoke Test
  run: |
    ./scripts/smoke-test-homelab.sh
    go test ./e2e -run TestSmoke -v
    go test ./internal/skill/ -v
  timeout-minutes: 5
```

---

## Troubleshooting

| Problema | Solução |
|----------|---------|
| Teste shell falha | `docker ps` deve ter 30+ containers |
| Go test falha | `go mod tidy && go test ./...` |
| Bot não responde | Verificar `TELEGRAM_BOT_TOKEN` e daemon |
| Resposta genérica | Verificar `senior_architect_prompt.go` |

---

## Próximos Passos

- [ ] Executar smoke test interativo via Telegram
- [ ] Validar respostas senior/arquiteto
- [ ] Adicionar mais skills sob demanda
- [ ] Integrar memory para learnings
- [ ] Adicionar CI/CD gate

---

## Referências

- **Scripts**: `scripts/smoke-test-homelab.sh`
- **Tests**: `e2e/smoke_test.go`
- **Personas**: `internal/persona/senior_architect_prompt.go`
- **Skills**: `~/.aurelia/skills/` (32 skills)
- **Guide**: `docs/SMOKE_TEST_GUIDE.md`

**Status:** ✅ Pronto para usar

