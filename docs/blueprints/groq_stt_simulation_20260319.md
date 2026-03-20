# Groq STT Simulation

## Estado Atual do Código

A Aurelia já suporta Groq como `STT provider`:

- [pkg/stt/groq.go](/home/will/aurelia/pkg/stt/groq.go)
- [cmd/aurelia/app.go](/home/will/aurelia/cmd/aurelia/app.go)
- [internal/config/config.go](/home/will/aurelia/internal/config/config.go)

Ponto de entrada de config:

- `groq_api_key`
- `stt_provider = groq`

## Onde o Token Entra

Arquivo de config da instância:

```json
{
  "stt_provider": "groq",
  "groq_api_key": "COLE_AQUI"
}
```

Local esperado pela Aurelia:

- `~/.aurelia/config/app.json`

## Curl Exato

```bash
curl -sS -X POST "https://api.groq.com/openai/v1/audio/transcriptions" \
  -H "Authorization: Bearer $GROQ_API_KEY" \
  -F "file=@/caminho/para/audio.wav" \
  -F "model=whisper-large-v3-turbo" \
  -F "language=pt" \
  -F "temperature=0" \
  -F "response_format=json"
```

## Smoke Script

Script criado:

- [groq-stt-curl-smoke.sh](/home/will/aurelia/scripts/groq-stt-curl-smoke.sh)

Uso sem token:

```bash
bash /home/will/aurelia/scripts/groq-stt-curl-smoke.sh /caminho/para/audio.wav
```

Uso com token:

```bash
GROQ_API_KEY=SEU_TOKEN \
bash /home/will/aurelia/scripts/groq-stt-curl-smoke.sh /caminho/para/audio.wav
```

## Comportamento Esperado

Sem token:

- imprime o comando final em modo `dry-run`

Com token:

- faz a chamada real na Groq
- retorna JSON com o campo `text`

## Próximo Passo

Quando você trouxer o token, o caminho é:

1. colocar `groq_api_key` em `~/.aurelia/config/app.json`
2. rodar o smoke com um `.wav` curto em PT-BR
3. validar a resposta do JSON
4. conectar persistência em `Supabase`
5. conectar indexação em `Qdrant`
