---
name: homelab-stabilizer
description: >
  Skill para estabilização de infraestrutura Homelab no Ubuntu Desktop (SOTA 2026).
  Gerencia a stack Qdrant, CapRover, LiteLLM, Ollama, Grafana e n8n.
  Garante governança de portas, subdomínios (ADR 20240401) e saúde dos serviços.
  Use quando o usuário pedir para "estabilizar humberto", "checar saúde do homelab", "configurar stack IA" ou "gerenciar portas".
---

# Homelab Stabilizer (Estabilizador de Infraestrutura)

Esta skill é o orquestrador central para a saúde e governança do seu Homelab no Ubuntu Desktop.

## Quando Usar
- **Check-up Diário**: Para garantir que todos os serviços de IA e monitoramento estão ativos.
- **Configuração de Portas**: Para evitar conflitos entre Grafana (3001) e CapRover (3000).
- **Recuperação de Serviços**: Para reiniciar containers que falharam ou habilitar serviços no boot.

## Governança de Portas (ADR 20240401)
A skill segue o contrato de portas:
- **3000**: CapRover Dashboard
- **3001**: Grafana (Alterado do padrão p/ estabilidade)
- **4000**: LiteLLM Proxy
- **6333**: Qdrant REST
- **11434**: Ollama Engine
- **5678**: n8n Workflow

## Comandos Suportados

### 1. Health Check (Saúde do Sistema)
Verifica se as portas estão abertas e os containers Docker estão respondendo.
- **Gatilho**: "Checar saúde", "health status", "verificar homelab".
- **Ação**: Executar `scripts/homelab_health_check.sh`.

### 2. Estabilizar Serviços (Stabilize)
Configura auto-restart nos containers e integra serviços com o boot do sistema (systemd).
- **Gatilho**: "Estabilizar ubuntu", "fix stability", "habilitar boot".
- **Ação**: Executar `scripts/service_manager.sh`.

## Dependências
- Docker e Docker Compose
- CapRover CLI (opcional)
- `ss` ou `netstat` (para monitoramento de portas)

---
- ADR: `20240401-governanca-homelab.md`
- Walkthrough: `walkthrough.md`
