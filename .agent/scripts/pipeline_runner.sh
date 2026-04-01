#!/bin/bash
# Pipeline Runner — executa a fila de .agent/tasks.json
# Uso: bash .agent/scripts/pipeline_runner.sh

set -e
TASKS_FILE=".agent/tasks.json"
DATE=$(date +%d/%m/%Y)

echo "🚀 Pipeline Soberano iniciado — $DATE"
echo "Lendo fila de $TASKS_FILE..."

# O agente lê o tasks.json e executa cada slice em ordem
# Este script é o contrato — o agente implementa a lógica
echo "PIPELINE_READY"
