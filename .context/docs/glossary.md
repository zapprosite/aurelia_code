# Glossary

## Termos de Runtime
- **AURELIA_HOME**: Diretório oculto (geralmente `~/.aurelia`) que contém configurações, logs e estado persistente.
- **Instance Lock**: Arquivo de trava (`instance.lock`) que impede a execução de múltiplas instâncias do bot simultaneamente.
- **Daemon de Usuário**: Serviço gerenciado pelo `systemd --user`, permitindo que o bot rode em background sem privilégios de root.
- **Bootstrap**: Processo inicial de criação da estrutura de pastas e verificação de dependências.

## Termos de Agente
- **ReAct Loop**: Ciclo "Reasoning and Acting" onde o modelo decide quais ferramentas usar para atingir um objetivo.
- **Teams**: Estrutura de coordenação entre o agente mestre e subagentes especialistas.
- **Skills**: Funcionalidades modulares em Go que podem ser carregadas dinamicamente pelo runtime.
- **MCP (Model Context Protocol)**: Protocolo aberto para conectar modelos de IA a fontes de dados e ferramentas externas.

## Termos de Observabilidade
- **Structured Logging (slog)**: Formato de log (JSON ou Texto estruturado) que facilita a busca e análise automática.
- **Redaction (Redação/Expurgação)**: Prática de remover ou mascarar dados sensíveis nos logs para proteger a privacidade.
- **Bridge**: Mecanismo que encaminha logs da biblioteca padrão (`log`) para o sistema de log estruturado.
