---
name: homelab-control
description: Controle nativo de Home Lab no Ubuntu (Ollama, NVIDIA, Docker, ZFS).
---

# 🛸 Home Lab Control Skill

Esta skill habilita o Aurelia a gerenciar infraestrutura de Homelab diretamente no Ubuntu 24.04 LTS, garantindo o uso de Bash/Nativo em vez de emulações de PowerShell.

## 🛠️ Comandos de Infra
1. **GPU (NVIDIA)**:
   - Use `nvidia-smi` para monitorar VRAM.
   - Verifique `gpustat` se disponível para logs curtos.
2. **Ollama**:
   - `curl -s http://localhost:11434/api/tags`: Listar modelos.
   - `ollama run <model>`: Testar inferência local.
   - Utilize os scripts em `scripts/update-ollama.sh` para rotinas de atualização.
3. **Docker**:
   - `docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"`: Listagem limpa.
   - Gerenciamento de deploys via `docker-compose.yml`.

## 🛡️ Anti-Padrão (Importante!)
- **NÃO use PowerShell**: O sistema operacional é **Ubuntu Desktop 24.04**. 
- Se um subagente tentar usar `powershell`, interrompa e force o uso de `bash`.

## Arquivos de Referência
- `scripts/health-check.sh`: Status unificado do lab.
- `scripts/update-ollama.sh`: Atualização de modelos AI.

## Quando usar
- Para diagnósticos de hardware.
- Para atualizar modelos de LLM locais.
- Para gerenciar containers de serviços (CapRover, N8N).

## Escopo correto

Use esta skill para o controle nativo e os comandos do host.

Se a tarefa pedir postura de operação sênior completa, com:

- classificação de saúde
- prova antes/depois
- recuperação segura
- governança de recursos

prefira carregar junto:

- `systems-engineer-homelab`
