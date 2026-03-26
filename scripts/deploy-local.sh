#!/usr/bin/env bash
# deploy-local.sh — build + deploy atômico da Aurélia no homelab
# Uso: ./scripts/deploy-local.sh [--skip-tests]
set -euo pipefail

INSTALLED=/usr/local/bin/aurelia
BACKUP=/usr/local/bin/aurelia.bak
BUILD_OUT=/tmp/aurelia-build-$$

SKIP_TESTS=false
[[ "${1:-}" == "--skip-tests" ]] && SKIP_TESTS=true

cd "$(git rev-parse --show-toplevel)"

echo "=== Aurélia deploy local $(date) ==="
echo "Commit: $(git rev-parse --short HEAD) | Branch: $(git branch --show-current)"

if [ "$SKIP_TESTS" = false ]; then
  echo "--- Rodando testes..."
  go test ./... -count=1 -timeout 120s
  echo "--- Testes OK"
fi

echo "--- Build..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath -ldflags="-s -w" -o "$BUILD_OUT" ./cmd/aurelia
echo "    Binário: $(ls -lah $BUILD_OUT | awk '{print $5, $9}')"

echo "--- Backup do binário atual..."
sudo cp "$INSTALLED" "$BACKUP"

echo "--- Stop → swap → start..."
sudo systemctl stop aurelia
sudo cp "$BUILD_OUT" "$INSTALLED"
sudo chmod +x "$INSTALLED"
sudo systemctl start aurelia
rm -f "$BUILD_OUT"

echo "--- Aguardando health check (até 30s)..."
STATUS="unreachable"
for i in $(seq 1 6); do
  sleep 5
  STATUS=$(curl -s http://127.0.0.1:8484/health \
    | python3 -c "import json,sys; d=json.load(sys.stdin); print('ok' if d['status'] in ('ok','degraded') else 'fail')" 2>/dev/null || echo "unreachable")
  echo "    tentativa $i: $STATUS"
  [ "$STATUS" = "ok" ] && break
done

if [ "$STATUS" = "ok" ]; then
  sudo rm -f "$BACKUP"
  echo "=== Deploy OK — status: $STATUS ==="
  # Notifica via Aurélia
  curl -s -X POST http://127.0.0.1:8484/v1/telegram/impersonate \
    -H "Content-Type: application/json" \
    -d "{\"text\": \"Deploy local OK: $(git rev-parse --short HEAD) em $(date '+%d/%m %H:%M').\", \"bot_id\": \"aurelia\"}" \
    >/dev/null || true
else
  echo "!!! Health check falhou ($STATUS) — rollback..."
  sudo systemctl stop aurelia
  sudo cp "$BACKUP" "$INSTALLED"
  sudo systemctl start aurelia
  sudo rm -f "$BACKUP"
  echo "!!! Rollback aplicado. Investigue os logs: journalctl -u aurelia -n 50"
  exit 1
fi
