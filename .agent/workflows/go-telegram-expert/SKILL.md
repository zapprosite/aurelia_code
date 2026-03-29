---
name: go-telegram-expert
description: Especialista em arquitetura sênior de bots Telegram em Go, focado em alta concorrência e mídias ricas (Soberano 2026).
---

# 🤖 Go Telegram Expert: Sovereign 2026

Habilita o Antigravity a projetar e implementar handlers de alta performance para o bot Telegram da Aurélia, utilizando os padrões de concorrência e tipos do Go 1.22+.

## 🏗️ Arquitetura de Comunicação (Triple-Tier)
1. **Premium (MiniMax 2.7)**: Use para gerar respostas empáticas, criativas ou realizar análise de sentimentos profunda em mensagens de voz/texto.
2. **Structured (DeepSeek 3.1)**: Use para parsing de comandos complexos, extração de entidades e roteamento de lógica de fluxos (FSM).
3. **Local (Qwen 3.5 (VL))**: Use para validação de sintaxe, respostas curtas de status e filtragem de spam.

## 🚀 Padrões de Implementação (Industrial)

### 1. Concorrência e Contexto
- **Goroutines**: Todo processamento de mídia pesada (foto/áudio) **deve** rodar em goroutines separadas com `context.WithTimeout`.
- **Worker Pool**: Use pools para evitar explosão de memória ao receber álbuns de imagens ou vídeos simultâneos.

### 2. Visão e OCR (Sovereign Vision)
- **Recuperação**: Todo handler de foto deve suportar múltiplos tipos (`PhotoSize`) e álbuns (`MediaGroup`).
- **OCR**: Realize o processamento inicial de texto em imagens localmente antes de enviar para análise de alto nível.

### 3. Voz e Áudio (MiniMax/Kokoro)
- **Transcrição**: Utilize o Tier 1 para transcrições precisas de áudios técnicos.
- **Síntese**: Integre respostas via Kokoro-TTS para voz premium brasileira.

### 4. Gestão de Estado (FSM)
- Evite IFs aninhados. Utilize máquinas de estado explícitas para fluxos de configuração ou onboarding em `internal/telegram/fsm/`.

## 🛡️ Guardrails
- **Rate Limiting**: Respeite os 30 msgs/seg do Telegram. Use buffers se necessário.
- **Sanitização**: Todo conteúdo `MarkdownV2` deve ser escapado corretamente para evitar crashes no parse da API.
- **Sudo=1**: Ao executar comandos via Telegram, exija autenticação de `admin_id` e valide o impacto (Dry-Run).

## 📍 Quando usar
- Para criar novos comandos (`/start`, `/status`, `/infra`).
- Para refatorar a lógica de recebimento de mídias ricas.
- Para otimizar a latência de resposta do bot.

## 🚫 Anti-Padrões
- Bloquear o loop principal do Telegram com operações de I/O pesadas.
- Usar globais para gerenciar o estado do usuário (Use DB ou Cache).
- Ignorar o tratamento de erros em `SendMessage`.