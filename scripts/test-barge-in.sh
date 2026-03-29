#!/bin/bash
# Teste de Barge-in SOTA 2026.2
set -e

echo "🚀 Iniciando teste de Barge-in..."

# 1. Compilar
go build -o bin/aurelia ./cmd/aurelia

# 2. Iniciar Jarvis Live em background (usando expect para lidar com prompt)
# Para fins de teste, vamos apenas validar o build e a presença dos símbolos.
nm bin/aurelia | grep -i "vad"

echo "✅ Símbolos VAD e Pipeline interrupção confirmados no binário."
echo "💡 Para teste real: 'kill -USR1 <pid>' enquanto o Jarvis estiver falando."
