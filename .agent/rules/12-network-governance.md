---
description: Governança obrigatória de portas e subdomínios — leitura antes de qualquer ação de rede.
id: 12-network-governance
---

# Regra 12: Governança de Rede — Portas e Subdomínios

> **Lei imutável**: Nenhum agente pode sugerir, abrir, mapear ou documentar uma porta ou subdomínio
> sem primeiro consultar os registros canônicos abaixo.

---

## Registros Autoritativos (leitura obrigatória)

| O que precisar saber | Onde ler |
|---------------------|----------|
| Portas ativas, reservadas e livres | `/srv/ops/ai-governance/PORTS.md` |
| Subdomínios ativos, pendentes e removidos | `/srv/ops/ai-governance/SUBDOMAINS.md` |
| Topologia de rede, tunnel, Terraform | `/srv/ops/ai-governance/NETWORK_MAP.md` |

**Leia SEMPRE antes de:**
- Sugerir porta para novo serviço
- Editar `docker-compose.yml`, `.env`, qualquer config de rede
- Adicionar entrada no Cloudflare Tunnel
- Criar deploy no Coolify
- Responder "qual porta usar para X"

---

## Estado Atual — Portas Reservadas (2026-04-02)

### ❌ Portas PROIBIDAS — já em uso

| Porta | Serviço |
|-------|---------|
| **3000** | Open WebUI (Coolify) — reservada |
| **3100** | Grafana |
| **3334** | aurelia-smart-router (LiteLLM UI) |
| **4000** | LiteLLM proxy (api.zappro.site / llm.zappro.site) |
| **4001** | OpenClaw Bot (Coolify) — reservada |
| **5678** | n8n |
| **6001/6002** | coolify-realtime (soketi) |
| **6333/6334** | Qdrant |
| **6379** | Redis (aurelia) |
| **8000** | Coolify PaaS |
| **8012** | Kokoro TTS |
| **8080** | aurelia-api |
| **8484** | aurelia health |
| **8888** | SearXNG |
| **9090** | Prometheus |
| **9100** | node-exporter |
| **9250** | cAdvisor |
| **11434** | Ollama |

### ✅ Portas LIVRES (confirmado)

`4002–4099` · `8443` · `9000` · `3001` · `3002`

---

## Estado Atual — Subdomínios (2026-04-02)

### Ativos (`*.zappro.site`)

| Subdomínio | Porta | Status |
|------------|-------|--------|
| `api` | :4000 | ✅ LiteLLM com auth |
| `llm` | :4000 | ✅ LiteLLM alias |
| `aurelia` | :8080 | ✅ Aurelia API |
| `coolify` | :8000 | ✅ Coolify PaaS |
| `monitor` | :3100 | ✅ Grafana |
| `n8n` | :5678 | ✅ n8n |
| `qdrant` | :6333 | ✅ Qdrant |

### Reservados — deploy pendente

| Subdomínio | Porta | Serviço |
|------------|-------|---------|
| `bot` | :4001 | OpenClaw Bot (Coolify) |
| `chat` | :3000 | Open WebUI (Coolify) |

---

## Regras Anti-Conflito

```
❌ NUNCA hardcodar :8000  → Coolify
❌ NUNCA hardcodar :4000  → LiteLLM produção
❌ NUNCA usar :3000       → reservada Open WebUI
❌ NUNCA usar :4001       → reservada OpenClaw Bot
✅ Dev local monorepo     → PORT=4002+ ou PORT=5173 (Vite)
✅ Novo microserviço      → faixa 4002–4099
```

---

## Protocolo ao Adicionar Porta ou Subdomínio

```
1. ler PORTS.md → confirmar porta livre
2. ler SUBDOMAINS.md → confirmar subdomínio disponível
3. ss -tlnp | grep :PORTA → verificar no host
4. Adicionar ao docker-compose ou config do serviço
5. Atualizar PORTS.md com o novo registro
6. Se subdomínio público:
   a. Atualizar SUBDOMAINS.md
   b. Atualizar ~/.cloudflared/config.yml (espelho local)
   c. Atualizar tunnel remoto via Cloudflare API
   d. Adicionar resource em homelab_dns.tf
   e. terraform apply (com tfvars das credenciais do .env)
7. Registrar ADR em docs/adr/ se for mudança estrutural
```

---

## Terraform — Como Aplicar

```bash
cd ~/dev/skills/homelab-cloud-governor/terraform

# Criar tfvars temporário (NUNCA commitar)
cat > terraform.tfvars << EOF
cloudflare_api_token = "$(grep CLOUDFLARE_API_TOKEN ~/aurelia/.env | cut -d= -f2)"
zone_id              = "$(grep CLOUDFLARE_ZONE_ID ~/aurelia/.env | cut -d= -f2)"
tunnel_id            = "$(grep CF_TUNNEL_ID ~/aurelia/.env | cut -d= -f2)"
EOF

terraform plan
terraform apply -auto-approve
rm terraform.tfvars   # limpar secrets imediatamente
```

## Tunnel Remoto — Como Atualizar

> O tunnel usa config **remota** (Cloudflare API), NÃO o config.yml local.

```bash
CF_TOKEN=$(grep CLOUDFLARE_API_TOKEN ~/aurelia/.env | cut -d= -f2)
CF_ACCOUNT=$(grep CLOUDFLARE_ACCOUNT_ID ~/aurelia/.env | cut -d= -f2)
CF_TUNNEL="8c55fcb7-fc8a-4ddc-a332-bff80012f11f"

curl -X PUT \
  "https://api.cloudflare.com/client/v4/accounts/${CF_ACCOUNT}/cfd_tunnel/${CF_TUNNEL}/configurations" \
  -H "Authorization: Bearer $CF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"config": {"ingress": [ ...entradas... ]}}'
```

---

**Referências:** [PORTS.md](/srv/ops/ai-governance/PORTS.md) | [SUBDOMAINS.md](/srv/ops/ai-governance/SUBDOMAINS.md) | [NETWORK_MAP.md](/srv/ops/ai-governance/NETWORK_MAP.md)
