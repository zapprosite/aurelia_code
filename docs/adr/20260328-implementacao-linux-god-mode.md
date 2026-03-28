# ADR 20260328: Implementação do Aurelia God Mode (Total Linux Control)

## Contexto
Em 28/03/2026, a soberania de um agente IA é medida por sua capacidade de interagir diretamente com o sistema operacional (OS) sem intermediários limitantes. A Aurelia deve evoluir para o "God Mode", permitindo execução de comandos Bash, análise de logs e aplicação de patches via ferramentas nativas de Go, operando sob a inteligência local do Gemma 3 27b e escalando para Tier 2 (OpenRouter) em casos de alta complexidade.

## Decisões Arquiteturais

### 1. Pacote `internal/os_controller` (Go)
- Implementar um wrapper robusto sobre `os/exec` com suporte a timeouts e contextos.
- Captura estruturada de Stderr e Stdout para retroalimentação do LLM.

### 2. Interface MCP (Model Context Protocol)
- Expor o `os_controller` como um servidor MCP independente em Go (via `mcp-golang`).
- Tools principais:
    - `run_bash_command`: Executa comandos com isolamento de contexto.
    - `read_system_log`: Tail/Grep inteligente de logs do host.
    - `apply_system_patch`: Modificação atômica de arquivos de sistema.

### 3. Execution Guard (Segurança Soberana)
- Middleware de segurança obrigatório.
- Lista negra de padrões destrutivos (`rm -rf /`, `mkfs`, `dd`, `iptables -F`).
- Modo Interativo: Sempre pede confirmação manual no terminal local para comandos críticos, a menos que o flag `--unsafe-auto` (Governança Will) esteja ativo.

### 4. Estratégia de Inferência
- **Local (Gemma 3 27b)**: Analisa sintaxe, sugere comandos e executa exploração básica.
- **Tier 2 (GLM-5 / MiniMax 2.7)**: Acionado pelo Juiz (`internal/gateway/judge.go`) quando a tarefa envolver scripts complexos de sysadmin ou automação densa.

## Consequências
- Aumento da superfície de ação da Aurelia.
- Necessidade de monitoramento constante de integridade do sistema.
- Dependência de permissões controladas (Sudoers) para o usuário do processo.

---
**Data**: 28/03/2026  
**Status**: Proposto  
**Autor**: Antigravity (SOTA 2026.1)
