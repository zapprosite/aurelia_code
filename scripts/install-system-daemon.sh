#!/bin/bash

set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
UNIT_TEMPLATE="$ROOT_DIR/scripts/aurelia.system.service"
UNIT_PATH="/etc/systemd/system/aurelia.service"
TARGET_BIN="/usr/local/bin/aurelia"
TMP_UNIT=$(mktemp)
USER_UNIT_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/systemd/user"
USER_UNIT_PATH="$USER_UNIT_DIR/aurelia.service"

AURELIA_USER=${AURELIA_USER:-$(id -un)}
AURELIA_GROUP=${AURELIA_GROUP:-$(id -gn)}
AURELIA_USER_HOME=${AURELIA_USER_HOME:-$(getent passwd "$AURELIA_USER" | cut -d: -f6)}
AURELIA_HOME=${AURELIA_HOME:-"${AURELIA_USER_HOME}/.aurelia"}

cleanup() {
    rm -f "$TMP_UNIT"
}
trap cleanup EXIT

if [ ! -f "$UNIT_TEMPLATE" ]; then
    echo "❌ Unit template não encontrado em $UNIT_TEMPLATE"
    exit 1
fi

escaped_root=${ROOT_DIR//&/\\&}
escaped_user=${AURELIA_USER//&/\\&}
escaped_group=${AURELIA_GROUP//&/\\&}
escaped_home=${AURELIA_HOME//&/\\&}
escaped_user_home=${AURELIA_USER_HOME//&/\\&}

sed \
    -e "s|__AURELIA_ROOT__|$escaped_root|g" \
    -e "s|__AURELIA_USER__|$escaped_user|g" \
    -e "s|__AURELIA_GROUP__|$escaped_group|g" \
    -e "s|__AURELIA_HOME__|$escaped_home|g" \
    -e "s|__AURELIA_USER_HOME__|$escaped_user_home|g" \
    "$UNIT_TEMPLATE" > "$TMP_UNIT"

echo "🔨 Construindo binário Linux..."
"$ROOT_DIR/scripts/build.sh" aurelia

echo "📁 Preparando runtime em $AURELIA_HOME..."
sudo mkdir -p "$AURELIA_HOME"/{config,data,logs}
sudo chown -R "$AURELIA_USER:$AURELIA_GROUP" "$AURELIA_HOME"

if systemctl --user list-unit-files aurelia.service >/dev/null 2>&1; then
    echo "🧹 Desativando serviço systemd de usuário para evitar conflito..."
    systemctl --user stop aurelia.service || true
    systemctl --user disable aurelia.service || true
    if [ -f "$USER_UNIT_PATH" ]; then
        backup_user_unit="${USER_UNIT_PATH}.bak.$(date +%Y%m%d-%H%M%S)"
        mv "$USER_UNIT_PATH" "$backup_user_unit"
        echo "   ↳ unit de usuário movida para $backup_user_unit"
    fi
    systemctl --user daemon-reload || true
fi

echo "🧹 Parando serviço systemd de sistema anterior (se existir)..."
sudo systemctl stop aurelia.service || true

echo "⚙️ Instalando binário em $TARGET_BIN..."
sudo install -m 0755 "$ROOT_DIR/aurelia" "$TARGET_BIN"

echo "⚙️ Instalando unit systemd em $UNIT_PATH..."
sudo install -m 0644 "$TMP_UNIT" "$UNIT_PATH"

echo "🔄 Recarregando systemd..."
sudo systemctl daemon-reload
sudo systemctl enable aurelia.service

echo "🚀 Reiniciando serviço..."
sudo systemctl restart aurelia.service

echo "📊 Status operacional:"
sudo systemctl status aurelia.service --no-pager -l

echo ""
echo "🧪 Health local:"
curl -fsS http://127.0.0.1:8484/health || true
