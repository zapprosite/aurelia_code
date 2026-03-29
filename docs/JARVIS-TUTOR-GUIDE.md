# Jarvis Tutor - Guia Rápido

## Pré-requisitos

1. `.env` configurado com:
```bash
# Obrigatório
TELEGRAM_BOT_TOKEN=seu_token
TELEGRAM_USER_ID=seu_id
TELEGRAM_CHAT_ID=chat_id

# STT (pelo menos um)
GROQ_API_KEY=sua_chave_groq
# OU
WHISPER_LOCAL_MODEL=medium
```

2. Build:
```bash
cd /home/will/aurelia
go build -o aurelia .
```

---

## Teste Rápido

```bash
./aurelia tutor
```

Se não tiver IDs, vai dar erro pedindo.

---

## Instalação BK 24/7

### 1. Configure .env
```bash
# Edite /home/will/aurelia/.env
nano /home/will/aurelia/.env

# Adicione:
TELEGRAM_USER_ID=123456789
TELEGRAM_CHAT_ID=123456789
JARVIS_THRESHOLD=500
```

### 2. Instale Service
```bash
sudo cp docs/systemd/jarvis-tutor.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now jarvis-tutor
```

### 3. Verifique
```bash
sudo systemctl status jarvis-tutor
sudo journalctl -u jarvis-tutor -f
```

---

## Troubleshooting

### Jarvis não inicia
```bash
# Verificar logs
sudo journalctl -u jarvis-tutor -n 50

# Verificar .env
cat /home/will/aurelia/.env | grep -E "TELEGRAM|GROQ|WHISPER"
```

### Erro de STT
```bash
# Groq OK?
curl -s https://api.groq.com/openai/v1/models | head -c 100

# Ollama OK?
curl -s http://localhost:11434/api/tags | jq '.models | length'
```

### Áudio não capturado
```bash
# Microfone OK?
arecord -l

# Teste gravação
arecord -d 3 /tmp/test.wav && aplay /tmp/test.wav
```

---

## Comandos

```bash
# Iniciar
sudo systemctl start jarvis-tutor

# Parar
sudo systemctl stop jarvis-tutor

# Reiniciar
sudo systemctl restart jarvis-tutor

# Logs
sudo journalctl -u jarvis-tutor -f

# Status
sudo systemctl status jarvis-tutor
```

---

## Atualização

```bash
cd /home/will/aurelia
git pull
go build -o aurelia .
sudo systemctl restart jarvis-tutor
```
