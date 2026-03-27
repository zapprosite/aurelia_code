#!/bin/bash
echo "🖥️ Iniciando Deploy do Desktop Controller..."

# Placeholder para o seu setup de Docker/VNC/RustDesk
if [ -f "docker-compose.yml" ]; then
    docker compose up -d
    echo "✅ Containers iniciados via Docker Compose."
else
    echo "⚠️ Nenhum docker-compose.yml encontrado no diretório atual."
    echo "Sugestão: Crie um compose para o RustDesk ou Desktop-Native."
fi

echo "--------------------------------"
