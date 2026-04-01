# Auditoria de Infraestrutura — Aurélia Sovereign 2026

## Data: 01/04/2026

---

## 1. STATUS GERAL

| Serviço | Status | Observação |
|---------|--------|------------|
| **Ollama** | ✅ rodando | qwen3.5:9b, qwen2.5vl:7b, nomic-embed |
| **Qdrant** | ✅ rodando | 2 collections |
| **LiteLLM** | ✅ rodando | porta 4000, 8484, 3334 |
| **TTS Kokoro** | ✅ rodando | porta 8012 |
| **STT Whisper** | ✅ rodando | porta 8020 |
| **Redis** | ⚠️ **RESTARTING** | erro RDB format v13 |
| **n8n** | ✅ rodando | |
| **Grafana** | ✅ rodando | monitor.zappro.site |
| **CapRover** | ✅ rodando | cap.zappro.site |
| **Cloudflare Tunnel** | ✅ rodando | will-zappro-homelab |
| **Tailscale** | ✅ rodando | 2 devices |
| **Supabase Local** | ❌ **DESINSTALADO** | ✅ correto |

---

## 2. CONTAINERS (15)

```
aurelia-redis-1         ⚠️ Restarting
aurelia-qdrant-1        ✅
aurelia-smart-router    ✅ (LiteLLM)
whisper-local           ✅ (STT)
aurelia-kokoro         ✅ (TTS)
captain-*               ✅ (CapRover)
grafana                 ✅
cadvisor                ✅
node-exporter           ✅
n8n + n8n-postgres      ✅
opencode-searxng        ✅
litellm-db              ✅
```

---

## 3. RECURSOS

| Recurso | Total | Usado | Livre |
|---------|-------|-------|-------|
| **VRAM** | 24GB | 1.2GB | 22.8GB |
| **RAM** | 30GB | 17GB | 13GB |
| **Disk (ZFS)** | 3.47TB | 47.9GB | 3.47TB |

---

## 4. SUBDOMÍNIOS (Cloudflare)

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

## 5. DOCUMENTAÇÃO

### Em /home/will/

| Arquivo | Descrição |
|---------|-----------|
| ARQUITETURA-DETALHADA.md | Arquitetura geral |
| DESKTOP-CONTROLE-APP.md | Controle desktop |
| TOP-10-MCP-2026.md | MCP servers |
| CLAUDE.md | Configuração Claude |

### Em /home/will/Documents/Obsidian/Aurelia/

| Vault | Descrição |
|-------|-----------|
| Governance | ADRs, políticas |
| Skills | Skills do agente |
| Training | Treinamento |

### Em /home/will/aurelia/docs/

| Diretório | Descrição |
|-----------|-----------|
| docs/adr/ | ADRs ( Architecture Decisions) |
| docs/governance/ | Políticas e contratos |
| docs/reports/ | Relatórios de saúde |

### Terraform

```
/home/will/dev/skills/homelab-cloud-governor/terraform/
└── homelab_dns.tf (DNS Cloudflare)
```

---

## 6. PROBLEMAS IDENTIFICADOS

### CRÍTICO

1. **Redis (aurelia-redis-1)** - Restarting
   - Erro: `Can't handle RDB format version 13`
   - Causa: AOF corrupt ou versão incompatível
   - Ação: Limpar dados e recriar

### WARNINGS

1. **Redis sem persistência confiável**
2. **某些文档 em inglês** (precisam traduzir)

---

## 7. AÇÕES RECOMENDADAS

### Imediato

- [ ] Fix Redis: remover volumes e recriar
- [ ] Atualizar docs deinfraestrutura

### Pós-Fix

- [ ] Testar todas as skills
- [ ] Verificar backup n8n
- [ ] Testar pipeline completo (STT → LLM → TTS)

---

## 8. MAPA FINAL

```
┌─────────────────────────────────────────────────────────────┐
│                    AURÉLIA SOBERANO 2026                   │
├─────────────────────────────────────────────────────────────┤
│  Hardware: RTX 4090 24GB | Ryzen 9 7900X | 30GB RAM        │
├─────────────────────────────────────────────────────────────┤
│  LLM Layer (Ollama)                                        │
│  ├── qwen3.5:9b (6.6GB) - Linguagem                         │
│  ├── qwen2.5vl:7b (6.0GB) - VL                             │
│  └── nomic-embed (274MB) - Embeddings                       │
├─────────────────────────────────────────────────────────────┤
│  Voice Layer                                                 │
│  ├── STT: Whisper GPU (:8020)                             │
│  ├── TTS: Kokoro GPU (:8012) + Edge (fallback)             │
│  └── Voz: Thalita PT-BR                                     │
├─────────────────────────────────────────────────────────────┤
│  Memória                                                     │
│  ├── Qdrant (:6333) - 2 collections                        │
│  ├── Redis (:6379) - ⚠️ PROBLEMA                           │
│  └── Obsidian - Vault Aurelia                              │
├─────────────────────────────────────────────────────────────┤
│  Orquestração                                                │
│  ├── LiteLLM (:4000) - Router                              │
│  ├── n8n (:5678) - Workflow                                │
│  └── CapRover (:3000) - Deploy                              │
├─────────────────────────────────────────────────────────────┤
│  Observabilidade                                             │
│  ├── Grafana (:3100) - Dashboards                          │
│  ├── Prometheus (:9090) - Métricas                          │
│  └── node-exporter, cadvisor, nvidia-gpu-exporter          │
├─────────────────────────────────────────────────────────────┤
│  Rede                                                        │
│  ├── Cloudflare Tunnel - DNS + HTTPS                        │
│  └── Tailscale - VPN mesh                                   │
└─────────────────────────────────────────────────────────────┘
```

---

*Auditoria realizada em 01/04/2026*
