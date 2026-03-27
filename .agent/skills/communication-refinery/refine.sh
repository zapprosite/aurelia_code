#!/bin/bash
# .agent/skills/communication-refinery/refine.sh
# Mestre Refinador v2 (SOTA 2026) - Multi-Bot & Sentiment Aware

TEXT=""
BOT="aurelia"
SENTIMENT="authoritative"
VOICE="af_sarah"
SPEED=1.0

# Parsing de argumentos
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -t|--text) TEXT="$2"; shift ;;
        -b|--bot) BOT="$2"; shift ;;
        -s|--sentiment) SENTIMENT="$2"; shift ;;
        *) TEXT="$1" ;; # Fallback para texto como primeiro argumento sem flag
    esac
    shift
done

if [ -z "$TEXT" ]; then
    echo "Uso: ./refine.sh -t \"Texto\" [-b bot] [-s sentimento]"
    echo "Bots: aurelia, vendas, agenda, infra"
    echo "Sentimentos: authoritative, empathy, focus, alarm"
    exit 1
fi

# 1. Mapeamento de Token do Bot (conforme .env)
case $BOT in
    vendas) TOKEN_VAR="TELEGRAM_TOKEN_VENDAS" ;;
    agenda) TOKEN_VAR="TELEGRAM_TOKEN_AGENDA" ;;
    infra)  TOKEN_VAR="TELEGRAM_TOKEN_INFRA" ;;
    caixa)  TOKEN_VAR="TELEGRAM_TOKEN_CAIXA" ;;
    obras)  TOKEN_VAR="TELEGRAM_TOKEN_OBRAS" ;;
    db)     TOKEN_VAR="TELEGRAM_TOKEN_DB" ;;
    *)      TOKEN_VAR="TELEGRAM_BOT_TOKEN" ;;
esac

# Carrega token do .env
if [ -f .env ]; then
    BOT_TOKEN=$(grep "^${TOKEN_VAR}=" .env | cut -d'=' -f2)
fi

# 2. Mapeamento de Sentimento (SOTA 2026)
case $SENTIMENT in
    empathy)
        VOICE="af_bella"
        SPEED=0.9
        ICON="🤍"
        ;;
    focus)
        VOICE="af_sarah"
        SPEED=1.15
        ICON="⚡"
        ;;
    alarm)
        VOICE="am_guilherme"
        SPEED=1.25
        ICON="🚨"
        ;;
    *)
        VOICE="af_sarah"
        SPEED=1.0
        ICON="🏛️"
        ;;
esac

# 3. Refinamento de Persona Dinâmico (SOTA 2026 - OLLAMA)
SOUL_FILE="$(dirname "$0")/SOUL_${BOT^^}.md"
if [ ! -f "$SOUL_FILE" ]; then
    SOUL_FILE="$(dirname "$0")/SOUL_AURELIA.md"
fi

echo "Efetuando reescrita de persona via Ollama ($BOT)..."
# Lê a alma do bot como prompt de sistema
SOUL_PROMPT=$(cat "$SOUL_FILE" | tr '\n' ' ' | sed 's/"/\\"/g')
CLEAN_TEXT=$(echo "$TEXT" | tr '\n' ' ' | sed 's/"/\\"/g')

REFINED_TEXT=$(curl -s -X POST "http://127.0.0.1:11434/api/generate" \
    -d "{
        \"model\": \"gemma3:12b\",
        \"prompt\": \"PERSONA: $SOUL_PROMPT\n\nREWRITE EXACTLY IN THE TONE OF THE PERSONA FOR TELEGRAM MARKDOWN. BE DIRECT, PREMIUM AND SENIOR. USE THE EMOJIS FROM PERSONA.\n\nORIGINAL TEXT: $CLEAN_TEXT\",
        \"stream\": false,
        \"options\": {\"temperature\": 0.3}
    }" | jq -r '.response')

if [ -n "$REFINED_TEXT" ] && [ "$REFINED_TEXT" != "null" ]; then
    echo "Texto Refinado: $REFINED_TEXT"
    FINAL_TEXT="$REFINED_TEXT"
else
    echo "Falha na reescrita dinâmica. Usando texto original."
    FINAL_TEXT="$ICON $TEXT"
fi

TEMP_AUDIO="/tmp/refine_${BOT}_output.opus"

# 4. Gerar Áudio Expressivo
python3 scripts/test_tts.py --text "$FINAL_TEXT" --voice "$VOICE" --speed "$SPEED" --output "$TEMP_AUDIO"

if [ $? -ne 0 ]; then
    echo "Erro na geração de áudio."
    exit 1
fi

# 5. Envio Multimodal com Feedback Loop (SOTA 2026)
REPLY_MARKUP="{\"inline_keyboard\":[[{\"text\":\"✅ Aprovar\",\"callback_data\":\"approve\"},{\"text\":\"✏️ Refinar\",\"callback_data\":\"refine\"}]]}"

BOT_TOKEN="$BOT_TOKEN" ./scripts/telegram-send-v2.sh "7220607041" "$FINAL_TEXT" "$TEMP_AUDIO" "$REPLY_MARKUP"

if [ $? -eq 0 ]; then
    echo "Sucesso! Refinamento enviado para o bot $BOT com botões de feedback."
else
    echo "Erro no envio para o bot $BOT."
fi
