---
name: homelab-control
description: Controle nativo de Home Lab no Ubuntu (Ollama, NVIDIA, Docker, ZFS).
---

# 🛸 Home Lab Control: Sovereign 2026

Habilita o gerenciamento operacional direto da infraestrutura do Home Lab sob a governança da Aurélia e Antigravity.

## 🛠️ Comandos de Infraestrutura (Ubuntu 24.04 Native)

### 1. GPU & CUDA (NVIDIA RTX 4090)
- **Monitoramento**: `nvidia-smi` é a fonte primária. Verifique se o `gemma3:12b` (Ollama) está consumindo os ~795MB necessários de VRAM.
- **Troubleshooting**: Se o `nvidia-smi` travar, verifique processos zumbis ou atualizações pendentes do `needrestart`.
- **Ref**: `docs/adr/ADR-20260320-politica-modelos-hardware-vram.md`.

### 2. Ollama & Sovereign Models
- **Gestão**: `curl -s http://localhost:11434/api/tags` para listar modelos ativos.
- **Modelos Pinados (2026)**:
  - `gemma3:12b`: Modelo residente principal (Sovereign fallback).
  - `qwen3.5:9b`: Alternativa para raciocínio local.
  - `bge-m3`: Modelo de embedding para o vetor DB (Qdrant).

### 3. Docker & Orquestração
- **Health Check**: `docker ps --format "table {{.Names}}\t{{.Status}}\t{{.CPUPerc}}\t{{.MemUsage}}"`
- **Deploy**: Use comandos `docker compose` apenas no diretório da aplicação alvo.
- **Segurança**: Auditoria periódica de portas expostas via `netstat -tulpn`.

### 4. Storage & ZFS
- **Status**: `zpool status` para verificar integridade dos pools de dados da Aurélia.
- **Snapshot**: Garanta snapshots regulares antes de upgrades críticos de sistema.

## 🛡️ Protocolo de Segurança (Antigravity Exclusive)
- **Bash Only**: Reprise total contra PowerShell. Se detectar scripts `.ps1`, converta para `.sh`.
- **Sudo=1**: Autoridade total habilitada. Sempre forneça prova de sucesso via logs.
- **Dry-Run**: Comandos de sistema (`reboot`, `systemctl stop`) exigem validação de impacto.

## 📍 Quando usar
- Para manter os serviços core (CapRover, N8N, Ollama) saudáveis.
- Para realizar o "Watchdog" de containers.
- Para gerenciar a alocação de hardware para squads de agentes.

## 🚫 Anti-Padrões
- Ignorar o custo de memória RAM total (7900x/RTX 4090 caps).
- Reiniciar o daemon `aurelia` sem verificar logs de erro primeiro.
