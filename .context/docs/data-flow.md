# Data Flow

O fluxo de dados no Aurelia é centrado no ciclo de percepção-ação (ReAct).

1. **Input**: O usuário envia texto, imagem ou áudio via Telegram.
2. **Preprocessing**: Adaptadores em `internal/telegram` capturam e pré-processam (STT via Groq se áudio).
3. **Context Assembly**: `internal/persona` e `internal/memory` montam o prompt atual (histórico, identidade, fatos).
4. **LLM Reasoning**: O loop ReAct em `internal/agent` interage com o LLM provider.
5. **Tool Execution**: Se o LLM solicita, ferramentas são executadas em `internal/tools` ou `internal/skill`.
6. **Output**: A resposta final é sintetizada e enviada de volta via Telegram (MarkdownV2).
7. **Persistence**: Mensagens e fatos aprendidos são salvos em SQLite para continuidade.
