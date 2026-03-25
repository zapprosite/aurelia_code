#!/usr/bin/env bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Sovereign-Bibliotheca v2 — Script de Unificação (Entrypoint)
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

source "$(dirname "${BASH_SOURCE[0]}")/lib/config.sh"

echo -e "${CLR_INFO}🚀 Iniciando Sovereign-Bibliotheca v2...${CLR_RESET}"

# 1. Sincronizar Memória (SQLite -> Qdrant/Supabase)
"$LIB_DIR/memory.sh" sync

# 2. Gerar Manifesto de Skills
"$LIB_DIR/skills.sh" manifest

# 3. Notificar Aurélia
"$LIB_DIR/comms.sh" send "Sovereign-Bibliotheca v2 sincronizada com sucesso. Qdrant, Supabase e Skills OpenClaw estão operacionais."

echo -e "${CLR_SUCCESS}✅ Pronto para operação multi-agente.${CLR_RESET}"
