#!/bin/bash
# bot-cli.sh — CLI para gerenciar o time de bots da Aurélia (S-32)
# Uso: ./scripts/bot-cli.sh <comando> [args]

API="http://localhost:3334"
IMPERSONATE="http://localhost:8484/v1/telegram/impersonate"

cmd=$1
shift

case $cmd in
  list)
    echo "=== Bots no time ==="
    curl -s "$API/api/bots" | python3 -c "
import sys, json
bots = json.load(sys.stdin)
for b in bots:
  status = '🟢' if b.get('running') else '🔴'
  print(f\"{status}  {b['id']:20} {b['name']:25} persona={b.get('persona_id','—')}\")
print(f'\nTotal: {len(bots)} bot(s)')
"
    ;;

  personas)
    echo "=== Personas disponíveis ==="
    curl -s "$API/api/personas" | python3 -c "
import sys, json
for p in json.load(sys.stdin):
  print(f\"  {p['id']:25} {p['name']} — {p['description'][:60]}...\")
"
    ;;

  add)
    # ./bot-cli.sh add <id> <nome> <token> <persona_id> [focus_area] [user_ids]
    ID=$1; NAME=$2; TOKEN=$3; PERSONA=$4; FOCUS="${5:-}"; USERS="${6:-}"
    if [ -z "$ID" ] || [ -z "$TOKEN" ]; then
      echo "Uso: bot-cli.sh add <id> <nome> <token> <persona_id> [focus_area] [user_ids_csv]"
      exit 1
    fi
    PAYLOAD=$(python3 -c "
import json, sys
d = {'id': '$ID', 'name': '$NAME', 'token': '$TOKEN', 'persona_id': '$PERSONA',
     'focus_area': '$FOCUS', 'enabled': True}
if '$USERS':
    d['allowed_user_ids'] = [int(x.strip()) for x in '$USERS'.split(',') if x.strip()]
print(json.dumps(d))
")
    echo "Criando bot '$ID'..."
    curl -s -X POST "$API/api/bots/create" \
      -H "Content-Type: application/json" \
      -d "$PAYLOAD" | python3 -c "
import sys, json
try:
  d = json.load(sys.stdin)
  print(f'✅ Bot criado: {d[\"id\"]} — {d[\"name\"]}')
except:
  print('Resposta:', sys.stdin.read())
" 2>/dev/null || echo "✅ Bot criado"
    ;;

  remove)
    ID=$1
    echo "Removendo bot '$ID'..."
    curl -s -X DELETE "$API/api/bots/remove?id=$ID"
    echo "✅ Removido"
    ;;

  ping)
    MSG="${1:-Aurélia, qual é o status do time de bots?}"
    echo "Enviando: $MSG"
    curl -s -X POST "$IMPERSONATE" \
      -H "Content-Type: application/json" \
      -d "{\"text\": \"$MSG\"}" | python3 -c "
import sys, json
d = json.load(sys.stdin)
print(f'→ {d[\"status\"]}: {d[\"message\"]}')
"
    ;;

  qdrant)
    echo "=== Collections Qdrant ==="
    APIKEY=$(grep QDRANT_API_KEY /home/will/.aurelia/config/secrets.env 2>/dev/null | cut -d= -f2)
    curl -s -H "api-key: $APIKEY" http://localhost:6333/collections | python3 -c "
import sys, json
d = json.load(sys.stdin)
for c in d['result']['collections']:
  print(f\"  {c['name']}\")
"
    ;;

  *)
    echo "Comandos: list | personas | add | remove | ping | qdrant"
    ;;
esac
