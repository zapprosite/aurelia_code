# ADR-20260330-enterprise-hardening-ctu

## Contexto
O ecossistema Aurélia evoluiu focado em prototipagem rápida. Para tornar o Jarvis um assistente de uso diário (CTU - Computer to Use) no Ubuntu 24.04 LTS, precisamos sair do modo interativo para um modo de **Serviço de Sistema** (Daemon) com governança industrial, observabilidade total e resiliência.

## Decisão
Adotar o padrão **Enterprise Sovereign 2026**:
1.  **Orquestração via Systemd**:
    - `aurelia.service`: Gerencia o binário Go principal.
    - `aurelia-voice-gateway.service`: Gerencia o bridge de áudio Python/VAD.
2.  **Endurecimento de Operação**:
    - Uso de `EnvironmentFile` para segredos via `governance-polish`.
    - Logs estruturados em `internal/logs` e espelhados no `journalctl`.
3.  **Processo Spec-Driven (SDD)**:
    - Inicialização do `spec-kit` no projeto (`specify init`) para governar futuras fatias (Slices).
4.  **Integração Her-Mode**:
    - Ativação do VAD Monitor como interface padrão de entrada.

## Consequências
- **Pró**: Jarvis torna-se um componente nativo do sistema operacional, inicia no boot e recupera-se de falhas automaticamente.
- **Soberania**: Independência total de execução manual; o Homelab torna-se uma "Appliance" de IA.
