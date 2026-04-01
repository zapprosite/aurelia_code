---
descricao: Visao com Qwen3.5-9B VL — Analise de telas e screenshots via Ollama local
---

# Skill: visao-qwen-vl

> **Autoridade**: Visao Local | **Data**: 01/04/2026
> **Modelo**: qwen3.5:9b (Ollama local, VL enabled, zero API)

## Proposito

Usar o modelo de visao Qwen3.5-9B VL para analisar:
- Screenshots da tela
- Imagens de paginas web
- Prints do browser
- Analise visual de interfaces

## Stack Atual (Sovereign 2026)

```
╔══════════════════════════════════════════════════════════════╗
║  AURÉLIA SOVEREIGN OS [v2026.04.01]                       ║
╠══════════════════════════════════════════════════════════════╣
║  HOST: will-zappro        GPU: RTX 4090 [24GB]            ║
║  STT: Whisper large-v3   [GPU] @ localhost:8020             ║
║  TTS: Edge TTS (GRÁTIS)  pt-BR-ThalitaMultilingualNeural   ║
║  TTS: Kokoro-82M ONNX    [GPU] @ localhost:8012 (fallback)║
║  LLM: Qwen3.5-9B VL     [Ollama] @ localhost:11434        ║
╚══════════════════════════════════════════════════════════════╝
```

## Comandos

| Comando | Funcao |
|---------|--------|
| `/vl` | Analisar screenshot atual |
| `/vl-url [url]` | Tirar print de site e analisar |
| `/vl-screenshot` | Capturar tela atual |
| `/vl-window` | Capturar janela ativa |
| `/vl-region` | Capturar regiao selecionada |

## Configuracao

### Ollama (porta 11434)
```bash
# Modelo ja disponivel
ollama list
# NAME                       ID              SIZE      MODIFIED   
# qwen3.5:9b                 6488c96fa5fa    6.6 GB    3 days ago    
```

## Fluxo de Execucao

### 1. Captura de Tela

```bash
# Screenshot total (scrot ja instalado)
scrot /tmp/vl-capture.png

# Janela ativa
scrot -u /tmp/vl-capture.png

# Regiao interativa
scrot -s /tmp/vl-capture.png
```

### 2. Envio para Qwen3.5-9B VL (API native)

```bash
# Via Ollama API
curl -s http://127.0.0.1:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen3.5:9b",
    "prompt": "Descreva o que voce ve nesta imagem em detalhes.",
    "images": ["'"$(base64 -w0 /tmp/vl-capture.png)"'"]
  }'
```

### 3. Resposta

```json
{
  "model": "qwen3.5:9b",
  "response": "A imagem mostra uma interface de terminal...",
  "done": true
}
```

## Scripts Helper

### vl-capturar.sh
```bash
#!/bin/bash
# Captura tela e salva em /tmp/vl-capture.png
scrot /tmp/vl-capture.png
echo "Captura: /tmp/vl-capture.png ($(file /tmp/vl-capture.png | cut -d: -f2))"
```

### vl-analisar.sh
```bash
#!/bin/bash
IMG="${1:-/tmp/vl-capture.png}"
PROMPT="${2:-Analise esta imagem e descreva o que voce ve.}"
curl -s http://127.0.0.1:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen3.5:9b",
    "prompt": "'"$PROMPT"'",
    "images": ["'"$(base64 -w0 "$IMG")"'"],
    "stream": false
  }' | jq -r '.response'
```

### vl-url.sh (com Playwright)
```bash
#!/bin/bash
URL="${1:-https://example.com}"
playwright screenshot --wait-for-timeout=2000 "$URL" /tmp/vl-url.png 2>/dev/null || \
curl -s "https://screenshotapi.net/api?url=$URL&wait=2000" -o /tmp/vl-url.png
echo "/tmp/vl-url.png"
```

## Exemplos de Uso

### Analisar tela atual
```
/vl
```
→ Captura screenshot → envia para Qwen3.5-9B VL → retorna descricao

### Analisar site
```
/vl-url https://github.com
```
→ Abre site com Playwright → captura → analiza

### Pergunta especifica
```
/vl o que esta escrito no botao vermelho?
```
→ Captura → pergunta especifica → resposta

## Requisitos

| Recurso | Status |
|---------|--------|
| Ollama | ✅ rodando na 11434 |
| qwen3.5:9b | ✅ ja instalado (6.6GB) |
| GPU RTX 4090 | ✅ 24GB VRAM |
| scrot | ✅ instalado |
| Playwright | Opcional para /vl-url |

## Environment

```bash
# .env
OLLAMA_URL=http://127.0.0.1:11434
OLLAMA_MODEL=qwen3.5:9b
VISION_ENABLED=true
```

## Health Check

```bash
# Verificar modelo
curl -s http://127.0.0.1:11434/api/tags | jq '.models[].name'

# Testar VL
curl -s http://127.0.0.1:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen3.5:9b",
    "prompt": "Hello",
    "stream": false
  }'
```

## Erros Comuns

| Erro | Solucao |
|------|---------|
| Model not found | `ollama pull qwen3.5:9b` |
| No GPU | Usar CPU (mais lento) |
| Image too large | Redimensionar: `convert img.png -resize 2048x2048 img.png` |
| Timeout | Reduzir prompt ou tamanho imagem |

---

*Visao local com Qwen3.5-9B VL — Soberano 2026*
