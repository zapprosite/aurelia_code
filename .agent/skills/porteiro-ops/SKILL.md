---
name: porteiro-ops
description: Gestão industrial do Sentinela (Porteiro) e cache de segurança no Redis.
---

# 🛡️ Porteiro Ops: Sentinel Control Center (SOTA 2026)

Skill especializada em gerenciar a camada de segurança de entrada/saída do ecossistema Aurélia. Permite ajustar o nível de rigor do Porteiro e manipular o cache de injeções detectadas no Redis.

## 🕹️ Comandos de Controle

### 1. `porteiro-status`
Exibe o estado atual do Sentinela.
- **Modo**: STRICT (Bloqueio), LOG_ONLY (Aprendizado), OFF (Desativado).
- **Modelo**: `qwen2.5:0.5b` (Ollama).
- **Redis**: Contagem de entradas de cache e latência média.

### 2. `porteiro-flush`
Limpa **todo** o cache de segurança do Redis. Útil quando há muitos falsos positivos após uma atualização de modelo ou mudança de contexto.
- **Ação**: `docker exec aurelia-redis-1 redis-cli KEYS "porteiro:cache:*" | xargs -r docker exec aurelia-redis-1 redis-cli DEL`

### 3. `porteiro-unblock "<prompt>"`
Remove o bloqueio de um prompt específico.
- **Ação**: Calcula o SHA-256 do prompt e deleta a chave `porteiro:cache:<hash>` no Redis.

## 🛠️ Modos de Operação (via .env)

Para alterar o rigor do Porteiro, edite o arquivo `.env` e reinicie o serviço:

```bash
# Opções: STRICT (default), LOG_ONLY, OFF
PORTEIRO_MODE=LOG_ONLY
```

- **STRICT**: Bloqueia imediatamente qualquer suspeita de injection.
- **LOG_ONLY**: Permite a passagem de tudo, mas loga `[LEARNING MODE] UNSAFE DETECTED` com o hash para auditoria.
- **OFF**: Desativa completamente a verificação para latência mínima em ambiente local isolado.

## 📊 Troubleshooting de Desenvolvimento

Se o bot Telegram estiver bloqueando comandos de desenvolvimento:
1. Verifique se você está na lista `ALLOWED_USER_IDS`.
2. Se não estiver, use `porteiro-unblock` para liberar o comando específico.
3. Se o bloqueio for generalizado, alterne para `PORTEIRO_MODE=LOG_ONLY` temporariamente.

---
*Governança Sovereign 2026 — Segurança sem Perda de Agilidade.*
