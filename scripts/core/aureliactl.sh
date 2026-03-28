#!/bin/bash
# Aurelia Runtime Controller (aureliactl) - SOTA 2026
# Sênior Direto ao Ponto. Sovereign Governance.

# Encontrar o local real do script (mesmo via symlink)
REAL_PATH="$(readlink -f "${BASH_SOURCE[0]}")"
REPO_ROOT="$(cd "$(dirname "$REAL_PATH")/../.." && pwd)"
BIN_DIR="$REPO_ROOT/bin"
LOG_DIR="$REPO_ROOT/logs"
SKILLS_DIR="$REPO_ROOT/.agent/skills"

mkdir -p "$LOG_DIR"

function show_status() {
    echo "🛰️  Aurelia Runtime Status (SOTA 2026)"
    echo "------------------------------------"
    pgrep -af "system-api" > /dev/null && echo "✅ System API: running" || echo "🔴 System API: stopped"
    pgrep -af "aurelia" > /dev/null && echo "✅ Aurelia Core: running" || echo "🔴 Aurelia Core: stopped"
    echo ""
    curl -s --max-time 2 http://localhost:8081/health | jq . || echo "⚠️  System API not responding on :8081"
}

function start_api() {
    pgrep -af "system-api" > /dev/null && { echo "⚠️  System API já está rodando."; return; }
    echo "🚀 Ligando Aurelia System API..."
    nohup "$BIN_DIR/system-api" > "$LOG_DIR/system-api.log" 2>&1 &
    sleep 2
    show_status
}

function show_logs() {
    echo "📄 Exibindo logs do System API (Ctrl+C para sair)..."
    tail -f "$LOG_DIR/system-api.log"
}

function list_skills() {
    echo "🛠️  Catálogo de Skills Soberanas (SOTA 2026)"
    echo "-------------------------------------------"
    if [ -d "$SKILLS_DIR" ]; then
        ls -1 "$SKILLS_DIR" | grep -v "README"
    else
        echo "🔴 Erro: Diretório de skills não encontrado em $SKILLS_DIR"
    fi
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
    restart) stop_all; sleep 1; start_api ;;
    status) show_status ;;
    logs) show_logs ;;
    skills) list_skills ;;
    *) echo "Uso: $0 {start|stop|restart|status|logs|skills}" ;;
esac
