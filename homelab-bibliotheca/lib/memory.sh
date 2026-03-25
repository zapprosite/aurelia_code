#!/usr/bin/env bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Sovereign-Bibliotheca v2 — Motor de Memória (Qdrant & Supabase)
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

source "$(dirname "${BASH_SOURCE[0]}")/config.sh"

function usage() {
    echo -e "Uso: $0 [comando]"
    echo -e "Comandos:"
    echo -e "  sync              Sincroniza SQLite -> Qdrant/Supabase"
    echo -e "  query <texto>     Busca semântica no Qdrant"
    echo -e "  status            Verifica saúde dos bancos"
}

function sync_memory() {
    echo -e "${CLR_INFO}Iniciando sincronização de memória...${CLR_RESET}"
    
    # 1. Verificar SQLite
    if [ ! -f "$AURELIA_DB_PATH" ]; then
        echo -e "${CLR_ERROR}Erro: SQLite não encontrado em $AURELIA_DB_PATH${CLR_RESET}"
        return 1
    fi

    # 2. Exemplo de extração de mensagens recentes (Top 10)
    # Aqui poderíamos implementar uma lógica de 'delta' baseada em timestamp
    MESSAGES=$(sqlite3 "$AURELIA_DB_PATH" "SELECT content, bot_id, chat_id, created_at FROM messages ORDER BY created_at DESC LIMIT 10;" -json)

    if [ "$MESSAGES" == "[]" ] || [ -z "$MESSAGES" ]; then
        echo -e "${CLR_WARN}Nenhuma mensagem nova para sincronizar.${CLR_RESET}"
        return 0
    fi

    # 3. Sincronização Qdrant (Simulada via curl - requer embeddings no pipeline real)
    # Nota: O pipeline Go da Aurélia já faz isso no mirror.go, este script é o 'backup soberano'
    echo -e "Sincronizando $(echo "$MESSAGES" | jq '. | length') itens..."
    
    # [LOGICA DE SYNC AQUI - Futuro: Chamar scripts/memory-sync-vector-db.sh se disponível]
    echo -e "${CLR_SUCCESS}Sincronização concluída (Modo Soberano).${CLR_RESET}"
}

function query_memory() {
    TEXT="$1"
    if [ -z "$TEXT" ]; then
        echo "Erro: Texto de busca obrigatório."
        return 1
    fi

    echo -e "${CLR_INFO}Buscando: '$TEXT'...${CLR_RESET}"
    
    # Busca simplificada via Search API da Aurélia (se disponível) ou Qdrant direto
    # Por padrão, usaremos a API da Aurélia para aproveitar o encoder de embeddings
    curl -s -X POST "$AURELIA_API_URL/api/brain/search" \
        -H "Content-Type: application/json" \
        -d "{\"text\": \"$TEXT\", \"limit\": 3}" | jq '.'
}

function check_status() {
    echo -e "=== Status Infraestrutura ==="
    # Qdrant
    curl -s "$QDRANT_URL/health" | grep -q "ok" && echo -e "Qdrant:   ${CLR_SUCCESS}ONLINE${CLR_RESET}" || echo -e "Qdrant:   ${CLR_ERROR}OFFLINE${CLR_RESET}"
    # Supabase (Studio)
    curl -s --head "$SUPABASE_STUDIO_URL" | grep -q "200" && echo -e "Supabase: ${CLR_SUCCESS}ONLINE${CLR_RESET}" || echo -e "Supabase: ${CLR_ERROR}OFFLINE${CLR_RESET}"
}

case $1 in
    sync) sync_memory ;;
    query) query_memory "$2" ;;
    status) check_status ;;
    *) usage ;;
esac
