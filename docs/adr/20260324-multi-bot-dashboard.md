# ADR 20260324: Multi-Bot Team Dashboard (S-32)

## Status
Aprovado

## Contexto
Will opera CNPJ HVAC-R especializado em DAIKIN VRV para alto padrão em SP. Para escalar
o negócio e a organização pessoal, precisa orquestrar múltiplos bots Telegram como um time:
líder (Aurélia) + vendas HVAC + gestão de obras + organização pessoal. O sistema atual suporta
apenas 1 bot (1 token, 1 BotController). O dashboard React existe como monitoring read-only.

## Decisão

### Backend (Go)
1. **BotConfig** — novo struct em `internal/config/config.go` com `ID`, `Name`, `Token`,
   `AllowedUserIDs`, `PersonaID`, `FocusArea`, `Enabled`. Campo `Bots []BotConfig` em
   `AppConfig`. Backward compat: se `Bots` vazio e `TelegramBotToken` presente, sintetizar
   entry `"aurelia"` automaticamente.

2. **BotPool** — `internal/telegram/pool.go` gerencia `map[string]*BotController`, com métodos
   `Add/Remove/Get/All/StartAll/StopAll/Configs`. Cada bot roda em goroutine própria.

3. **Composition root** — `cmd/aurelia/app.go` substitui `bot *BotController` por
   `pool *BotPool + primaryBot`. `initFeatures()` itera `cfg.Bots`. `start()` chama
   `pool.StartAll()`. `shutdown()` chama `pool.StopAll()`.

4. **Dashboard events** — `Event.BotID string` adicionado para filtros por bot.

5. **REST API** — `GET/POST/PATCH/DELETE /api/bots` no dashboard server.

6. **Persistence** — `internal/config/persistence.go` com `SaveBots()`.

7. **Squad sync** — `SyncBotsToSquad(bots []BotConfig)` em `internal/agent/squad.go`.

### Frontend (React)
1. Nova aba **"Bots"** no sidebar (segundo item, após Timeline), ícone `Bot`.
2. `BotsTab.tsx` — grid de BotCards com fetch `GET /api/bots`.
3. `BotCard.tsx` — card individual com status dot, persona, focus area.
4. `CreateBotModal.tsx` — formulário glass-effect para criação via `POST /api/bots`.
5. `BotDetail.tsx` — editor de config + activity feed filtrado por `bot_id`.

### Persona Templates
`internal/persona/templates.go` — 4 templates:
- `aurelia-leader` — COO, orquestra time (Crown/purple)
- `hvac-sales` — Funil DAIKIN VRV SP (Thermometer/blue)
- `project-manager` — 3+1 obras, orçamentos (ClipboardCheck/yellow)
- `life-organizer` — Gym, igreja, filha, namorada (Calendar/green)

## Consequências

**Positivas:**
- Will pode criar e gerir o time de bots pelo dashboard sem editar JSON
- Cada bot tem persona, contexto e usuários permitidos independentes
- Dashboard mostra atividade por bot (filtro `bot_id`)
- Base para futura integração Instagram (S-32b, fora desta slice)

**Negativas / Riscos:**
- Cada bot adicional consome uma goroutine de polling Telegram (~1MB RAM)
- Tokens extras devem ser mantidos secretos (não entram no git)
- BotPool precisa de shutdown graceful para evitar mensagens duplicadas

**Mitigação:**
- `pool.StopAll()` no shutdown com timeout
- Tokens armazenados apenas em config local JSON (não commitado)
- Testes unitários de `BotPool` para garantir lifecycle correto

## Artefatos
- `internal/config/config.go` — BotConfig + backward compat
- `internal/telegram/pool.go` — BotPool manager
- `internal/telegram/bot.go` — campo botID
- `cmd/aurelia/app.go` — composition root refatorado + REST API
- `internal/dashboard/dashboard.go` — BotID em Event
- `internal/agent/squad.go` — SyncBotsToSquad
- `internal/persona/templates.go` — 4 persona templates
- `internal/config/persistence.go` — SaveBots
- `frontend/src/components/sidebar/Sidebar.tsx` — aba Bots
- `frontend/src/App.tsx` — roteamento da aba
- `frontend/src/components/dashboard/BotsTab.tsx` — grid de gerenciamento
- `frontend/src/components/dashboard/BotCard.tsx` — card reutilizável
- `frontend/src/components/dashboard/CreateBotModal.tsx` — formulário
- `frontend/src/components/dashboard/BotDetail.tsx` — detalhe/edição


---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
