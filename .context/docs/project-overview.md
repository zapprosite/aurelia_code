# Project Overview: Aurelia

Aurelia é um agente de codificação autônomo, "local-first", construído em Go. Ele foi projetado para operar com alta autonomia em ambientes de monorepo e homelab.

## Missão
Prover uma interface inteligente e resiliente para automação de tarefas de desenvolvimento, utilizando múltiplos modelos de linguagem (LLMs) e ferramentas locais.

## Pilares Operacionais
1. **Disciplina de Runtime**: Uso de lockfile para instância única e daemon via systemd (usuário).
2. **Observabilidade Total**: Logging estruturado com `slog` e proteção de dados sensíveis.
3. **Arquitetura Modular**: Monolito modular em Go facilitando a extensão via "Skills" e "MCP".
4. **Segurança Local**: Execução controlada, sem necessidade de privilégios de root para o dia a dia.

## Componentes Principais
- **ReAct Loop**: Motor de raciocínio baseado em pensamento e ação.
- **Teams**: Orquestração de subagentes para tarefas complexas.
- **Memory**: Camadas de memória (identidade, fatos, notas, episódica) em SQLite.
- **Interfaces**: Telegram Bot API como canal principal de I/O.
