---
name: master-response-2026
description: Master Response 2026 — Força Markdown profissional + Kokoro TTS local GPU
---

# Skill: Master Response 2026

> **Autoridade**: Soberania de Output | **Ativação**: `/mr2026` | **Data**: 01/04/2026

## Propósito

Esta skill garante que toda resposta do agente seja entregue em:
1. **Markdown 2026 profissional** — formato limpo, curto e direto
2. **Áudio TTS via Kokoro (GPU local)** — sem APIs de nuvem, voz pt-BR

## Contrato de Output

### Markdown 2026 (Obligatório)

Toda resposta DEVE seguir este contrato (`markdown2026Contract`):

```
## Saida Obrigatoria
- Entregue texto em Markdown limpo, curto e direto.
- Use listas planas; nao use JSON, YAML ou blocos soltos como resposta final.
- Use tabela apenas quando houver comparacao real entre multiplos itens.
- Se algo nao puder ser confirmado, diga isso explicitamente em uma linha objetiva.
```

**Regras**:
- Headers com `#`, `##`, `###` limpos
- Listas com `-` ou `1.` (sem asteriscos)
- Blocos de código apenas quando necessário (nunca como output final)
- Tabelas apenas para comparações reais
- Sem Emojis decorativos excessivos (apenas quando relevante)

### TTS Kokoro (Obligatório)

Após gerar a resposta Markdown:

1. **Endpoint**: `http://127.0.0.1:8011/v1/audio/speech`
2. **Modelo**: `kokoro` (ONNX 82M local)
3. **Voz**: Feminina pt-BR (Aurélia official)
4. **Formato**: `opus` (ou `mp3` como fallback)
5. **Velocidade**: `1.0`

**Prompt de Sanitização** (aplicar ANTES de enviar ao TTS):

```
ATENÇÃO — MODO VOZ: Esta mensagem será sintetizada em voz. 
Regras obrigatórias:
1) Sem markdown — nada de asteriscos, underlines, cerquilhas, colchetes ou backticks.
2) Sem listas numeradas ou com marcadores — use frases conectadas ("Primeiro... depois... por fim...").
3) Sem tabelas.
4) Números por extenso quando for natural ("dois mil e vinte e seis", "trinta por cento").
5) Frases curtas e diretas.
6) Se precisar enumerar itens, separe por vírgula ou use "e" entre eles.
```

## Comandos

| Comando | Ação |
|---------|------|
| `/mr2026` | Ativar modo Master Response 2026 (força markdown + audio) |
| `/mr2026-audio-only` | Gerar áudio da última resposta sem novo markdown |
| `/mr2026-off` | Desativar modo Master Response 2026 |

## Fluxo de Execução

```
1. Recebe input do usuário
2. Processa normalmente (raciocínio, ferramentas, etc)
3. Aplica markdown2026Contract na resposta final
4. Sanitiza texto para modo voz (remove markdown decoration)
5. Envia para Kokoro TTS (POST /v1/audio/speech)
6. Entrega: Markdown + Audio (ambos simultâneos ou audio como follow-up)
```

## Configuração Técnica

- **TTS Provider**: OpenAI Compatible (`openai_compatible`)
- **TTS Base URL**: `http://127.0.0.1:8011`
- **TTS Model**: `kokoro`
- **Voz**: `af_sarah` (ou default feminine pt-BR)
- **Timeout**: 30s (GPU local é rápido)
- **Max chars**: 50000 (limite Kokoro)

## Casos de Erro

| Cenário | Tratamento |
|---------|-------------|
| `voice-proxy` (8011) offline | Retornar apenas markdown, logar erro |
| Texto > 50k chars | Truncar para 50k, avisar no output |
| TTS falhar | Retry 1x, depois fallback markdown only |
| JSON/YAML detectado na resposta | Converter para markdown limpo automaticamente |

## Referências

- `internal/telegram/bot_governance.go` — `markdown2026Contract`
- `internal/telegram/input_pipeline.go:483-484` — voice system prompt
- `pkg/voice/tts/openai_compatible.go` — TTS client

---

*Assinado: Aurélia Master Response 2026 | Soberano de Output*