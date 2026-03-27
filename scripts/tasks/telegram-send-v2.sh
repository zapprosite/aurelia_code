# Use: BOT_TOKEN="xxx" ./telegram-send-v2.sh CHAT_ID TEXT [FILE_PATH] [REPLY_MARKUP]

TOKEN="${BOT_TOKEN}"
CHAT_ID="${1:-7220607041}"
TEXT="${2}"
FILE_PATH="${3}"
REPLY_MARKUP="${4}"

if [ -z "$TOKEN" ]; then
    # Fallback para o token padrão se não fornecido via env
    if [ -f .env ]; then
        TOKEN=$(grep "^TELEGRAM_BOT_TOKEN=" .env | cut -d'=' -f2)
    fi
fi

if [ -z "$TOKEN" ]; then
    echo "ERRO: BOT_TOKEN não fornecido e padrão não encontrado no .env"
    exit 1
fi

if [ -z "$FILE_PATH" ]; then
    curl -s -X POST "https://api.telegram.org/bot${TOKEN}/sendMessage" \
        -d chat_id="${CHAT_ID}" \
        -d text="${TEXT}" \
        -d parse_mode="HTML" \
        -d reply_markup="${REPLY_MARKUP}"
else
    curl -s -X POST "https://api.telegram.org/bot${TOKEN}/sendVoice" \
        -F chat_id="${CHAT_ID}" \
        -F voice="@${FILE_PATH}" \
        -F caption="${TEXT}" \
        -F parse_mode="HTML" \
        -F reply_markup="${REPLY_MARKUP}"
fi
