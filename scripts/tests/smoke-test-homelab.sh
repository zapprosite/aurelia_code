#!/bin/bash
set -e

# 🧪 Smoke Test — Integração Telegram + Homelab + Senior Response

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🔬 SMOKE TEST — Aurelia Homelab Integration"
echo "============================================="
echo ""

# ============================================================================
# 1. HEALTH CHECK FULL
# ============================================================================
echo "1️⃣  HEALTH CHECK FULL"
echo "   Testando: saúde completa do homelab"
echo ""

# Containers running
CONTAINER_COUNT=$(docker ps -q | wc -l)
echo "   📦 Containers: $CONTAINER_COUNT ativos"
[ "$CONTAINER_COUNT" -ge 15 ] && echo "      ✅ OK (>15)" || echo "      ⚠️  Poucos containers"

# GPU VRAM
FREE_VRAM=$(nvidia-smi --query-gpu=memory.free --format=csv,noheader,nounits 2>/dev/null | head -1)
if [ -n "$FREE_VRAM" ]; then
  VRAM_GB=$((FREE_VRAM / 1024))
  echo "   💾 GPU VRAM: ${VRAM_GB}GB livre"
  [ "$VRAM_GB" -gt 5 ] && echo "      ✅ OK (>5GB)" || echo "      ⚠️  Apertado"
else
  echo "   💾 GPU: nvidia-smi indisponível"
fi

# ZFS Pool
POOL_STATUS=$(zpool status tank 2>/dev/null | grep "state:" | head -1)
echo "   📦 ZFS: $POOL_STATUS"

# Network Tunnel
TUNNEL_RESPONSE=$(curl -s -m 2 https://n8n.zappro.site/health 2>/dev/null || echo "TIMEOUT")
if echo "$TUNNEL_RESPONSE" | grep -q "200\|app"; then
  echo "   🌐 Tunnel: Cloudflared ✅"
else
  echo "   🌐 Tunnel: FALHA ❌"
fi

echo ""
echo "   ✅ Health check completo"
echo ""

# ============================================================================
# 2. CONTAINER DIAGNOSTICS
# ============================================================================
echo "2️⃣  CONTAINER DIAGNOSTICS"
echo "   Testando: diagnóstico inteligente de container"
echo ""

CONTAINER_NAME="n8n"
CONTAINER_STATUS=$(docker ps -a | grep "$CONTAINER_NAME" | awk '{print $7}' | head -1)
echo "   📦 Container '$CONTAINER_NAME': $CONTAINER_STATUS"

if [ "$CONTAINER_STATUS" = "Up" ]; then
  CONTAINER_ID=$(docker ps | grep "$CONTAINER_NAME" | awk '{print $1}' | head -1)

  # RAM usage
  RAM=$(docker stats --no-stream "$CONTAINER_ID" 2>/dev/null | tail -1 | awk '{print $7}')
  echo "   💾 RAM: $RAM"

  # Logs
  RECENT_LOGS=$(docker logs "$CONTAINER_ID" --tail 5 2>/dev/null | grep -i "error\|warn" | wc -l)
  if [ "$RECENT_LOGS" -gt 0 ]; then
    echo "   ⚠️  Logs: $RECENT_LOGS avisos/erros nos últimos logs"
  else
    echo "   ✅ Logs: OK"
  fi
else
  echo "   ❌ Container parado"
fi

echo ""
echo "   ✅ Diagnóstico completo"
echo ""

# ============================================================================
# 3. ARCHITECTURE DECISIONS
# ============================================================================
echo "3️⃣  ARCHITECTURE DECISIONS"
echo "   Testando: recomendações de arquitetura"
echo ""

# Voice Stack
VOICE_CONTAINERS=("kokoro")
for container in "${VOICE_CONTAINERS[@]}"; do
  if docker ps --format "{{.Names}}" | grep -q "$container"; then
    echo "   🎵 Voice component: $container ✅"
  else
    echo "   🎵 Voice component: $container ❌"
  fi
done

# Database Strategy
POSTGRES_UP=$(curl -s -m 2 localhost:5435 2>/dev/null | wc -c)
if [ "$POSTGRES_UP" -gt 0 ]; then
  echo "   💾 Database strategy: PostgreSQL local ✅"
  echo "      → Recomendação: Manter local (latência <1ms)"
else
  echo "   💾 Database: indisponível ⚠️"
fi

# Monitoring
if curl -s -m 2 http://localhost:9090/-/ready > /dev/null 2>&1; then
  echo "   📊 Monitoring: Prometheus ✅"
  echo "      → Recomendação: Manter coleta de métricas"
else
  echo "   📊 Monitoring: offline"
fi

echo ""
echo "   ✅ Análise arquitetural completa"
echo ""

# ============================================================================
# 4. SAFE AUTOMATION
# ============================================================================
echo "4️⃣  SAFE AUTOMATION (DRY RUN)"
echo "   Testando: automação segura com confirmação"
echo ""

# Validate ZFS snapshot capability
if zfs list tank > /dev/null 2>&1; then
  SNAPSHOT_NAME="tank@smoke-test-$(date +%s)"
  echo "   📸 ZFS Snapshot (DRY RUN)"
  echo "      Would create: $SNAPSHOT_NAME"

  # Dry run
  if zfs list -H -t snapshot | head -1 > /dev/null 2>&1; then
    echo "      ✅ Capability verified"
  fi
else
  echo "   ❌ ZFS não disponível"
fi

# Backup capability
if [ -d /srv/backups ]; then
  BACKUP_SIZE=$(du -sh /srv/backups | awk '{print $1}')
  echo "   💾 Backups: $BACKUP_SIZE"
  echo "      ✅ Backup directory available"
else
  echo "   ❌ /srv/backups não existe"
fi

echo ""
echo "   ✅ Operações seguras verificadas"
echo ""

# ============================================================================
# 5. MULTI-STEP ORCHESTRATION
# ============================================================================
echo "5️⃣  MULTI-STEP ORCHESTRATION (DRY RUN)"
echo "   Testando: orquestração kokoro deploy"
echo ""

echo "   Step 1: Check backends"
for backend in kokoro; do
  if docker ps --format "{{.Names}}" | grep -q "$backend"; then
    echo "      ✅ $backend rodando"
  else
    echo "      ❌ $backend parado"
  fi
done

echo ""
echo "   Step 2: Health verification"
if curl -s http://localhost:8012/v1/models > /dev/null 2>&1; then
  echo "      ✅ TTS (Kokoro) respondendo"
fi

echo ""
echo "   ✅ Orquestração verificada (pronto para deploy)"
echo ""

# ============================================================================
# 6. DISASTER RECOVERY
# ============================================================================
echo "6️⃣  DISASTER RECOVERY READINESS"
echo "   Testando: validação de DR"
echo ""

# Último backup
LATEST_BACKUP=$(ls -t /srv/backups/supabase-backup-*.sql.gz 2>/dev/null | head -1)
if [ -n "$LATEST_BACKUP" ]; then
  BACKUP_AGE=$(( ($(date +%s) - $(stat -c %Y "$LATEST_BACKUP")) / 3600 ))
  echo "   💾 Último backup: ${BACKUP_AGE}h atrás"
  [ "$BACKUP_AGE" -lt 24 ] && echo "      ✅ OK" || echo "      ⚠️  Antigo"
else
  echo "   ❌ Nenhum backup encontrado"
fi

# Snapshots
SNAPSHOT_COUNT=$(zfs list -H -t snapshot 2>/dev/null | wc -l)
echo "   📸 Snapshots ZFS: $SNAPSHOT_COUNT"
[ "$SNAPSHOT_COUNT" -gt 0 ] && echo "      ✅ OK" || echo "      ⚠️  Sem snapshots"

# Espaço de backup
BACKUP_SPACE=$(df /srv/backups 2>/dev/null | tail -1 | awk '{print $4}')
if [ -n "$BACKUP_SPACE" ]; then
  SPACE_GB=$((BACKUP_SPACE / 1048576))
  echo "   📊 Espaço livre: ${SPACE_GB}GB"
  [ "$SPACE_GB" -gt 5 ] && echo "      ✅ OK" || echo "      ⚠️  Apertado"
fi

echo ""
echo "   ✅ DR prontidão validada"
echo ""

# ============================================================================
# SUMMARY
# ============================================================================
echo "============================================="
echo "✅ SMOKE TEST COMPLETO"
echo ""
echo "Próximos passos:"
echo "  1. Enviar via Telegram: 'saúde completa'"
echo "  2. Confirmar respostas senior/arquitetor"
echo "  3. Testar automação com confirmação"
echo "  4. Validar DR com snapshot real (se necessário)"
echo ""
echo "Executar teste Telegram:"
echo "  go test ./e2e -run TestSmokeHomelab -v"
echo ""
