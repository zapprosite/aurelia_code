#!/bin/bash
# Jarvis Tutor - Script de Instalação 24/7
# Uso: ./scripts/install-jarvis-tutor.sh USER_ID CHAT_ID

set -e

USER_ID=${1:-}
CHAT_ID=${2:-}

if [ -z "$USER_ID" ] || [ -z "$CHAT_ID" ]; then
    echo "Uso: $0 USER_ID CHAT_ID"
    echo "Exemplo: $0 123456789 987654321"
    exit 1
fi

echo "🤖 Instalando Jarvis Tutor 24/7"
echo "   User ID: $USER_ID"
echo "   Chat ID: $CHAT_ID"
echo ""

# 1. Build
echo "📦 Fazendo build..."
cd /home/will/aurelia
go build -o aurelia .

# 2. Copiar service
echo "📋 Instalando systemd service..."
sudo cp docs/systemd/jarvis-tutor.service /etc/systemd/system/

# 3. Substituir USER_ID e CHAT_ID
sudo sed -i "s/USER_ID/$USER_ID/g" /etc/systemd/system/jarvis-tutor.service
sudo sed -i "s/CHAT_ID/$CHAT_ID/g" /etc/systemd/system/jarvis-tutor.service

# 4. Reload systemd
echo "🔄 Reload systemd..."
sudo systemctl daemon-reload

# 5. Habilitar e iniciar
echo "🚀 Iniciando Jarvis Tutor..."
sudo systemctl enable --now jarvis-tutor

# 6. Status
echo ""
echo "✅ Jarvis Tutor instalado!"
echo ""
sudo systemctl status jarvis-tutor

echo ""
echo "📝 Comandos úteis:"
echo "   Status:  sudo systemctl status jarvis-tutor"
echo "   Logs:    sudo journalctl -u jarvis-tutor -f"
echo "   Parar:   sudo systemctl stop jarvis-tutor"
echo "   Reiniciar: sudo systemctl restart jarvis-tutor"
