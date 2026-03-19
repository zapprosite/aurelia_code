# Architecture: Aurelia

Aurelia segue um padrão de **Monolito Modular** em Go, organizado para separação clara de responsabilidades entre interface, domínio e infraestrutura.

## Estrutura de Pastas
- `cmd/aurelia/`: Ponto de entrada, composição da aplicação e fiação (wiring) de dependências.
- `internal/agent/`: Lógica central do loop ReAct, registro de ferramentas e gestão de times.
- `internal/runtime/`: Resolução de caminhos, bootstrap de diretórios e lock de instância única.
- `internal/observability/`: Logging estruturado (`slog`), expurgação de dados e métricas.
- `internal/telegram/`: Adaptadores de entrada e saída para a API do Telegram.
- `internal/memory/`: Persistência de estado operacional e camadas de memória em SQLite.
- `internal/persona/`: Gestão de identidade e construção de prompts.
- `pkg/llm/`: Provedores de modelos de linguagem (Anthropic, Google, OpenAI, etc.).
- `pkg/stt/`: Serviços de transcrição de áudio (Groq Whisper).

## Ciclo de Vida do Processo
1. **Bootstrap**: Garante que `~/.aurelia/` exista.
2. **Lock**: Adquire `~/.aurelia/instance.lock` gravando PID, Comando e Timestamp.
3. **Configure**: Inicializa logging estruturado e bridge de stdlog.
4. **Wiring**: Instancia componentes e registra ferramentas.
5. **Run**: Inicia loops de I/O (Telegram) e agendamentos.
6. **Shutdown**: Captura sinais de interrupção, finaliza processos pendentes e libera o lock.

## Gestão de Estado
- O estado operacional é persistido em bancos SQLite localizados em `~/.aurelia/data/`.
- Configurações canônicas residem em `~/.aurelia/config/app.json`.
