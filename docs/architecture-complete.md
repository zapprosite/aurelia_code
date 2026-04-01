# Arquitetura Completa — Aurélia Sovereign 2026

> **Autoridade**: Supremacia Arquitetural | **Atualizado**: 01/04/2026

## Hardware

| Componente | Especificação |
|------------|---------------|
| **CPU** | AMD Ryzen 9 7900X — 12c/24t |
| **RAM** | 32 GB DDR5 |
| **GPU** | RTX 4090 — 24 GB VRAM · Driver 580 · CUDA 13 |
| **Disco SO** | nvme1n1 — 931 GB (Kingston Gen3) |
| **Disco Dados** | nvme0n1 — 3,64 TB (Crucial Gen5) → ZFS pool "tank" |

---

## Stack Completa (01/04/2026)

```
╔══════════════════════════════════════════════════════════════════════╗
║  AURÉLIA SOVEREIGN OS [v2026.04.01]                               ║
╠══════════════════════════════════════════════════════════════════════╣
║  HOST: will-zappro        GPU: RTX 4090 [24GB]                     ║
║  STT: Whisper large-v3   [GPU] @ localhost:8020                     ║
║  TTS: Edge TTS (GRÁTIS)  pt-BR-ThalitaMultilingualNeural           ║
║  TTS: Kokoro-82M ONNX    [GPU] @ localhost:8012 (fallback)         ║
║  LLM: Qwen3.5-9B VL     [Ollama] @ localhost:11434                 ║
║  VL: Qwen3.5-9B Vision  [Ollama] @ localhost:11434                  ║
╚══════════════════════════════════════════════════════════════════════╝
```

---

## Serviços e Containers

### Principais

| Serviço | Container | Porta | Status |
|---------|-----------|-------|--------|
| **Aurelia Bot** | aurelia | 8484 | ✅ |
| **Redis (Main)** | aurelia-redis-main | 6379 | ✅ |
| **Redis (LiteLLM)** | litellm-redis | 6380 | ✅ |
| **Redis (n8n)** | n8n-redis | 6381 | ✅ |
| **Qdrant** | aurelia-qdrant-1 | 6333, 6334 | ✅ |
| **LiteLLM** | aurelia-smart-router | 4000, 8484, 3334 | ✅ |
| **Whisper (STT)** | whisper-local | 8020 | ✅ |
| **Kokoro (TTS)** | aurelia-kokoro | 8012 | ✅ |
| **n8n** | n8n | 5678 | ✅ |
| **n8n Postgres** | n8n-postgres | 5432 | ✅ |
| **Grafana** | grafana | 3100 | ✅ |
| **Prometheus** | prometheus | 9090 | ✅ |
| **node-exporter** | node-exporter | 9100 | ✅ |
| **cadvisor** | cadvisor | 9250 | ✅ |

### CapRover

| Serviço | Container | Porta |
|---------|-----------|-------|
| **nginx** | captain-nginx | 80, 443 |
| **captain** | captain-captain | 3000 |

---

## Pipeline de Voz

```
[Áudio .wav]
    │ POST :8020/v1/audio/transcriptions
    ▼
[Whisper large-v3 — STT GPU]
    │ texto
    ▼
[BGE-M3 :11434] → embed → [Qdrant :6333] → contexto
    │
    ▼ POST :11434/api/chat
[Qwen 3.5 — 262K ctx, thinking]
    │ resposta
    ▼ POST :8012/v1/audio/speech
[Kokoro TTS GPU / Edge TTS]
    │
    ▼
[Áudio .wav saída]
```

---

## Ollama (systemd)

**API:** http://localhost:11434

| Modelo | Params | Quant | VRAM | Contexto |
|--------|--------|-------|------|----------|
| `qwen3.5:9b` | 9.65B | Q4_K_M | ~6.5 GB | **262K tokens** |
| `qwen2.5vl:7b` | 7B | Q4_K_M | ~6.0 GB | Vision |
| `nomic-embed-text` | 566M | F16 | ~1.2 GB | Embeddings |

---

## Budget VRAM — RTX 4090 (24 GB)

```
Desktop (Xorg + GNOME)       ~1,0 GB   fixo
Whisper STT                   ~4,0 GB   ativo
Kokoro TTS                     ~5,0 GB   ativo
────────────────────────────────────────────
Em uso:                       ~10 GB    → ~14 GB livres

+ Qwen 3.5 (sob demanda)     ~6,5 GB
+ Qwen VL (sob demanda)       ~6,0 GB
────────────────────────────────────────────
Pior caso:                   ~22.5 GB  → ~1.5 GB livres
```

---

## Armazenamento ZFS "tank"

```
nvme0n1 (3,64 TB) → 3,47 TB livre
├── docker-data  → /srv/docker-data   [45 GB]
├── monorepo     → /srv/monorepo     [256 MB]
├── backups      → /srv/backups      [194 MB]
├── models       → /srv/models       [1.75 GB]
├── n8n          → /srv/data/n8n     [10 MB]
└── monitoring   → /srv/data/monitoring
```

---

## Subdomínios Cloudflare

| Subdomínio | Target |
|------------|--------|
| n8n.zappro.site | cf-tunnel |
| qdrant.zappro.site | cf-tunnel |
| cap.zappro.site | cf-tunnel |
| llm.zappro.site | cf-tunnel |
| monitor.zappro.site | cf-tunnel |
| aurelia.zappro.site | cf-tunnel |
| supabase.zappro.site | cf-tunnel |
| studio.zappro.site | cf-tunnel |

---

## Comandos Rápidos

```bash
# Smoke test completo
./scripts/smoke-test.sh

# Rate limits
./scripts/rate-limit-check.sh

# Status containers
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# VRAM
nvidia-smi --query-gpu=memory.used,memory.free --format=csv,noheader

# ZFS
zpool status tank
```

---

## Skills Disponíveis

| Skill | Caminho | Descrição |
|-------|---------|-----------|
| `aurelia_code` | `.agent/skills/aurelia_code/` | Líder DevOps - Agent Swarm |
| `visao-qwen-vl` | `.agent/skills/visao-qwen-vl/` | Análise visual |
| `master-response-2026` | `.agent/skills/master-response-2026/` | Markdown + TTS forçado |
| `gerente-vendas-agencia` | `.agent/skills/gerente-vendas-agencia/` | Vendas PT-BR |

---

## Mapa de Portas Atualizado (01/04/2026)

| Porta | Serviço | Status |
|-------|---------|--------|
| 80, 443 | captain-nginx (CapRover) | ✅ |
| 3000 | CapRover Dashboard | ✅ |
| 3100 | Grafana | ✅ |
| 5432 | n8n-postgres | ✅ |
| 5433, 6543 | Supabase PgBouncer | ✅ |
| 5435 | Supabase Postgres (MCP) | ✅ |
| 5678 | n8n | ✅ |
| 6333, 6334 | Qdrant | ✅ |
| 6379 | Redis Main (Aurelia) | ✅ |
| 6380 | Redis LiteLLM | ✅ |
| 6381 | Redis n8n | ✅ |
| 8000 | Supabase Kong | ✅ |
| 8010 | Voice proxy → STT | ✅ |
| 8011 | Voice proxy → TTS | ✅ |
| 8012 | Kokoro TTS | ✅ |
| 8020 | Whisper STT | ✅ |
| 9090 | Prometheus | ✅ |
| 9100 | node-exporter | ✅ |
| 9250 | cAdvisor | ✅ |
| 11434 | Ollama | ✅ |
| 54323 | Supabase Studio | ✅ |

---

## Referências Desktop

| Documento | Descrição |
|-----------|-----------|
| `/home/will/Desktop/SYSTEM_ARCHITECTURE.md` | Arquitetura detalhada (mar/2026) |
| `/home/will/Desktop/guia-homelab-estavel.md` | Guia de manutenção homelab |
| `/home/will/Desktop/guide-audio-tts-stt.md` | Guia de voz |
| `/home/will/Desktop/guide-antigravity.md` | Guia Antigravity |
| `/home/will/Desktop/GUIA-AURELIA-CAPROVER-TERRAFORM.md` | CapRover + Terraform |

---

*Atualizado: 01/04/2026 | Aurélia Sovereign*
