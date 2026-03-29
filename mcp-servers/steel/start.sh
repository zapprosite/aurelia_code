#!/bin/bash
# Steel Container Startup Script
# ADR: 20260328-container-steel-browser-isolation

set -e

echo "[Steel] Iniciando container..."

# Variables com defaults
export STAGEHAND_MODEL=${STAGEHAND_MODEL:-"aurelia-smart"}
export LITELLM_BASE_URL=${LITELLM_BASE_URL:-"http://litellm:4000/v1"}
export NODE_ENV=${NODE_ENV:-"production"}
export LOG_LEVEL=${LOG_LEVEL:-"info"}

echo "[Steel] Configurações:"
echo "  - Model: $STAGEHAND_MODEL"
echo "  - LiteLLM: $LITELLM_BASE_URL"
echo "  - Node Env: $NODE_ENV"

# Health check server simples
cat > /tmp/health-server.js << 'HEALTH'
const http = require('http');
const server = http.createServer((req, res) => {
    if (req.url === '/health') {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ status: 'ok', timestamp: new Date().toISOString() }));
    } else {
        res.writeHead(404);
        res.end('Not Found');
    }
});
server.listen(3000, () => console.log('[Health] Listening on :3000'));
HEALTH

# Inicia health server em background
node /tmp/health-server.js &
HEALTH_PID=$!

# Cleanup ao sair
trap "kill $HEALTH_PID 2>/dev/null" EXIT

# Aguarda LiteLLM estar disponível
echo "[Steel] Aguardando LiteLLM em $LITELLM_BASE_URL..."
for i in {1..30}; do
    if curl -sf "$LITELLM_BASE_URL/health" > /dev/null 2>&1; then
        echo "[Steel] LiteLLM disponível!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "[Steel] AVISO: LiteLLM não disponível após 30 tentativas"
        echo "[Steel] Continuando mesmo assim (Stagehand tentará conectar)..."
    fi
    sleep 1
done

# Inicia Stagehand MCP server
echo "[Steel] Iniciando Stagehand MCP server..."
cd /app

# Se build existe, usa; senão usa ts-node direto
if [ -f "dist/index.js" ]; then
    echo "[Steel] Usando build compilado..."
    node dist/index.js
else
    echo "[Steel] Usando ts-node (dev mode)..."
    npx ts-node --esm src/index.ts
fi
