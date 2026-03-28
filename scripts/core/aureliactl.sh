#!/bin/bash
# Aurelia Runtime Controller (aureliactl) - SOTA 2026
# Sênior Direto ao Ponto. Sovereign Governance.

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BIN_DIR="$REPO_ROOT/bin"
LOG_DIR="$REPO_ROOT/logs"

mkdir -p "$LOG_DIR"

function show_status() {
    echo "🛰️  Aurelia Runtime Status (SOTA 2026)"
    echo "------------------------------------"
    pgrep -af "system-api" > /dev/null && echo "✅ System API: running" || echo "🔴 System API: stopped"
    pgrep -af "aurelia" > /dev/null && echo "✅ Aurelia Core: running" || echo "🔴 Aurelia Core: stopped"
    echo ""
    curl -s http://localhost:8081/health | jq . || echo "⚠️  System API not responding on :8081"
}

function start_api() {
    echo "🚀 Ligando Aurelia System API..."
    nohup "$BIN_DIR/system-api" > "$LOG_DIR/system-api.log" 2>&1 &
    sleep 2
    show_status
}

function stop_all() {
    echo "🛑 Desligando serviços..."
    pkill -f "system-api"
    pkill -f "aurelia"
    echo "Sessão encerrada."
}

case "$1" in
    start) start_api ;;
    stop) stop_all ;;
    status) show_status ;;
    *) echo "Uso: $0 {start|stop|status}" ;;
esac
