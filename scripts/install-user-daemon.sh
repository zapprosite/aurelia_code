#!/bin/bash

set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
UNIT_TEMPLATE="$ROOT_DIR/scripts/aurelia.service"
UNIT_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/systemd/user"
UNIT_PATH="$UNIT_DIR/aurelia.service"
TMP_UNIT=$(mktemp)

cleanup() {
    rm -f "$TMP_UNIT"
}
trap cleanup EXIT

if [ ! -f "$UNIT_TEMPLATE" ]; then
    echo "❌ Unit template não encontrado em $UNIT_TEMPLATE"
    exit 1
fi

mkdir -p "$UNIT_DIR" "$HOME/.aurelia/logs"

escaped_root=${ROOT_DIR//&/\\&}
sed "s|__AURELIA_ROOT__|$escaped_root|g" "$UNIT_TEMPLATE" > "$TMP_UNIT"

echo "🔨 Construindo binário..."
"$ROOT_DIR/scripts/build.sh"

echo "⚙️ Instalando Unit de Usuário..."
install -m 0644 "$TMP_UNIT" "$UNIT_PATH"

echo "🔄 Recarregando systemd (--user)..."
systemctl --user daemon-reload
systemctl --user enable aurelia.service

echo "🚀 Reiniciando serviço..."
systemctl --user restart aurelia.service

echo "📊 Status operacional:"
systemctl --user status aurelia.service --no-pager -l
