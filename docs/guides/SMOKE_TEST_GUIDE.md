# 🧪 Smoke Test — Integração Telegram + Homelab

Guide completo para testar integração end-to-end com respostas de nível senior/arquiteto.

## Prerequisitos

```bash
# Variáveis de ambiente
export TELEGRAM_BOT_TOKEN="..."          # Bot token
export TELEGRAM_CHAT_ID="..."            # Chat ID para teste
export CLAUDE_API_KEY="..."              # Para respostas via Claude
export ANTHROPIC_API_KEY="..."           # Alternativa
```

## Opção 1: Script Shell (Rápido)

```bash
# Rodar smoke test local — sem Telegram
./scripts/smoke-test-homelab.sh

# Output esperado:
# ✅ Health check completo
# ✅ Container diagnostics
# ✅ Architecture analysis
# ✅ Safe automation verified
# ✅ Orchestration ready
# ✅ DR readiness validated
```

**O que valida:**
- ✅ 30+ containers rodando
- ✅ GPU VRAM > 5GB livre
- ✅ ZFS pool saudável
- ✅ Voice stack up e respondendo
- ✅ Network tunnel ativo
- ✅ Backups e snapshots OK
- ✅ Disaster recovery prontidão

**Saída esperada:**
```
🔬 SMOKE TEST — Aurelia Homelab Integration
=============================================

1️⃣  HEALTH CHECK FULL
   📦 Containers: 30 ativos ✅
   💾 GPU VRAM: 18GB livre ✅
   📦 ZFS: state: ONLINE ✅
   🌐 Tunnel: Cloudflared ✅

2️⃣  CONTAINER DIAGNOSTICS
   📦 Container 'n8n': Up ✅
   💾 RAM: 512MiB
   ✅ Logs: OK

3️⃣  ARCHITECTURE DECISIONS
   🎵 Voice component: speaches ✅
   💾 Database strategy: PostgreSQL local ✅
   📊 Monitoring: Prometheus ✅

4️⃣  SAFE AUTOMATION
   📸 ZFS Snapshot capability verified ✅
   💾 Backups: 219M ✅

5️⃣  MULTI-STEP ORCHESTRATION
   Step 1: VRAM Pre-check ✅
   Step 2: Check backends ✅
   Step 3: Health verification ✅

6️⃣  DISASTER RECOVERY
   💾 Últimos backups ✅
   📸 Snapshots ZFS: 63 ✅

=============================================
✅ SMOKE TEST COMPLETO
```

---

## Opção 2: Testes Unitários (CI/CD)

```bash
# Rodar testes de skill
go test ./internal/skill/ -v -run TestLoad

# Output esperado:
# PASS: TestLoader_SingleDir
# PASS: TestLoader_DuplicateName
# ok  github.com/kocar/aurelia/internal/skill	0.004s
```

---

## Opção 3: Teste Telegram (Interativo)

### Preparação

1. **Setup bot local**
   ```bash
   # Iniciar daemon Aurelia
   systemctl start aurelia --user

   # Ou manualmente:
   go run ./cmd/aurelia/main.go
   ```

2. **Configurar variáveis**
   ```bash
   export TELEGRAM_BOT_TOKEN="BOT:TOKEN"
   export TELEGRAM_CHAT_ID="123456"
   ```

3. **Validar conexão**
   ```bash
   curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getMe" | jq '.ok'
   # Deve retornar: true
   ```

### Teste Manual via Telegram

Enviar mensagens para o bot:

#### 1️⃣ Health Check Completo
```
Você: saúde completa do homelab
Bot: [Resposta senior - métricas + recomendações]
```

**Esperado:**
- ✅ Status de 30+ containers
- ✅ GPU VRAM utilizado vs budget
- ✅ ZFS pool health
- ✅ Network latency
- ✅ Recomendações específicas (e.g., "VRAM OK para voice stack")

#### 2️⃣ Diagnóstico de Container
```
Você: diagnóstico do container n8n
Bot: [Análise + recomendações]
```

**Esperado:**
- Status: Up/Down/Restarting
- RAM/CPU usage
- Logs recentes (se houver errors)
- Recomendação: "Restart recomendado se houver OOM"

#### 3️⃣ Decisão Arquitetural
```
Você: deveria mover DB para fora?
Bot: [Análise de trade-offs]
```

**Esperado:**
- Option 1: Local (latency <1ms, single point of failure)
- Option 2: Remote (managed, but 50ms latency)
- Option 3: Hybrid (critical local, analytics remote)
- Recomendação com justificativa

#### 4️⃣ Automação Segura
```
Você: cria snapshot do tank/data
Bot: ⚠️ Confirmar? snapshot será tank@<timestamp>
Você: sim, confirmo
Bot: ✅ Snapshot criado
```

**Esperado:**
- Pediu confirmação (segurança)
- Executou operação
- Reportou sucesso + detalhes

#### 5️⃣ Orquestração Multi-step
```
Você: subir voice stack completo
Bot: [Step-by-step com validações]
```

**Esperado:**
```
Step 1: VRAM Pre-check → ✅ 18GB livre (precisa 8GB)
Step 2: Start Whisper → ✅ Container saudável
Step 3: Start Chatterbox → ✅ Container saudável
Step 4: Start Proxy → ✅ Respondendo
Final: Health validation → ✅ Voice stack pronto
```

#### 6️⃣ Validação de Disaster Recovery
```
Você: validar prontidão de recovery
Bot: [Status de backups + snapshots]
```

**Esperado:**
- Último backup: < 24h
- Snapshots: > 30
- Espaço de backup: > 5GB
- Recomendação: "Pronto para DR"

---

## Teste Automático (Go)

```bash
# Rodar teste smoke completo
go test ./e2e -run TestSmokeHomelabIntegration -v

# Output esperado:
=== RUN TestSmokeHomelabIntegration
=== RUN TestSmokeHomelabIntegration/HealthCheck-FullStack
  ✅ Health check completo
=== RUN TestSmokeHomelabIntegration/ContainerDiagnose-WithRecommendations
  ✅ Diagnóstico completo
=== RUN TestSmokeHomelabIntegration/ArchitectureDecision-GPUOptimization
  ✅ Análise arquitetural completa
=== RUN TestSmokeHomelabIntegration/SafeAutomation-ZFSSnapshot
  ✅ Operações seguras verificadas
=== RUN TestSmokeHomelabIntegration/MultiStepOrchestration-VoiceStackDeploy
  ✅ Orquestração verificada
=== RUN TestSmokeHomelabIntegration/IntelligentRecovery-DRValidation
  ✅ DR Status completo
--- PASS: TestSmokeHomelabIntegration (45s)
PASS
ok  github.com/kocar/aurelia/e2e	45.123s
```

---

## Validação de Respostas Senior/Arquiteto

Ao verificar respostas do bot, procurar por:

### ✅ Boas Sinais
- **Trade-off Analysis**: "Option 1 tem X, Option 2 tem Y, portanto recomendo Z porque..."
- **Metrics**: "VRAM em 18GB/24GB", "Latência 45ms", "3569GB espaço livre"
- **Recomendações Específicas**: Não "faça backup", mas "pg_dump a cada 24h e snapshot ZFS antes de mudanças"
- **Risk Awareness**: "Single point of failure", "Network latency impact", "Recovery time"
- **Alternatives**: Sempre 2-3 opções com pros/cons
- **Preconditions**: "PostgreSQL 5435 deve estar acessível"

### ❌ Red Flags
- ⚠️ Respostas genéricas ("faça backup")
- ⚠️ Sem justificativa ("faça isso porque sim")
- ⚠️ Sem números (vago vs concreto)
- ⚠️ Operação destrutiva sem confirmação
- ⚠️ Sem alternativas consideradas
- ⚠️ Ignorar constraints (VRAM, latency, cost)

---

## Troubleshooting

### Bot não responde
```bash
# Verificar daemon rodando
systemctl status aurelia --user

# Verificar logs
tail -50 ~/.cache/aurelia/daemon.log | grep -i error

# Verificar API key
echo $TELEGRAM_BOT_TOKEN | wc -c  # Deve ter ~50 chars
```

### Skills não carregam
```bash
# Verificar diretório
ls ~/.aurelia/skills/ | wc -l  # Deve ter 32

# Testar loader
go test ./internal/skill/ -v
```

### Latência alta
```bash
# Verificar latência Telegram
time curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getMe" > /dev/null

# Verificar resposta bot
time curl -s http://localhost:8080/health
```

---

## Próximas Iterações

Depois de validar smoke test:

1. **Adicionar Skills**: Incorporar mais automações (deploy, scaling, etc.)
2. **Memory Integration**: Persistir learnings sobre homelab
3. **Multi-Agent**: Colaboração entre Claude/Codex/Gemini
4. **CI/CD Gates**: Rodar smoke test em cada deploy
5. **Performance Baseline**: Registrar métricas para comparação futura

---

## Referências

- Smoke test script: `scripts/smoke-test-homelab.sh`
- Teste Go: `e2e/smoke_test.go`
- Persona: `internal/persona/senior_architect_prompt.go`
- Skills: `~/.aurelia/skills/` (32 skills)
- Homelab status: `docs/ARCHITECTURE.md`

