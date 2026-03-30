---
name: telegram-formatter-2026
description: Conversor industrial de JSON (LLM Local) para Markdown 2026 profissional no Telegram.
---

# 🛰️ Telegram Formatter 2026: Master Interface

Skill industrial para garantir que nenhuma resposta técnica (JSON) de modelos locais (Qwen) vaze para o usuário final no Telegram sem o devido polimento semântico.

## 🏛️ Padrão Master Command (SOTA 2026)
As respostas convertidas devem seguir este esquema visual:

1. **Header**: `🛰️ [TÍTULO / STATUS]`
2. **Contexto**: `📊 Métricas & Estado`
3. **Análise**: `🧠 Insight Industrial`
4. **Próxima Ação**: `🚀 Próximo Passo`

## 🛠️ Mecanismo de Intercepção
Toda saída que começar com `{` ou `{"status":` deve ser enviada ao `Porteiro (Sentinel)` com o modelo Qwen 0.5b para re-formatação imediata antes do envio.

## 📍 Regras de Estilo
- **Emoji First**: Use emojis para categorizar seções (📊, 🧠, ⚙️, 🚀).
- **Zero JSON**: É terminantemente proibido exibir colchetes, chaves ou aspas duplas de estruturação técnica ao usuário.
- **Tone**: Engenheiro Sênior (Soberano 2026), direto e construtivo.

## 🚀 Como Ativar
Esta skill é integrada nativamente na camada `internal/telegram/output.go` e `send.go`.
