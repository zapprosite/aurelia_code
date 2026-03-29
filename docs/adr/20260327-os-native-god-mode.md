# ADR-20260327-os-native-god-mode

## Status
Proposto (SOTA 2026.1)

## Contexto
O controle atual é focado em Browser (Playwright). O "God Mode" exige manipulação direta de janelas, processos, arquivos e infraestrutura de hardware (NVIDIA/Docker) sem intermediários.

## Decisão
Consolidar o `os_controller` como o driver mestre de sistema:
1.  **Native Bridge**: Integração direta com APIs de desktop (X11/Wayland para janelas) e CLI de sistema (systemd, docker, zfs).
2.  **Autonomous Ops**: Capacidade de realizar deploys, troubleshoot de drivers e gestão de recursos de forma autônoma.
3.  **Audit Log**: Todo comando de sistema é registrado em `logs/sovereign/audit.log` com análise semântica pré-execução (Porteiro).

## Consequências
- **Pró**: Autonomia total sobre a máquina host (Ubuntu Sovereign).
- **Contra**: Risco de segurança elevado (mitigado pelo Porteiro e Hard-Lock).
