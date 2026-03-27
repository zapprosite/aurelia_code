---
name: add-subdomain
description: Adiciona novo subdomínio público {nome}.zappro.site via Terraform + Cloudflare Tunnel, atualizando todos os arquivos de governança obrigatórios.
phases: [E]
---

# 🌐 Add Subdomain: {nome}.zappro.site

Processo completo e auditável para expor um serviço local publicamente via Cloudflare Tunnel, com IaC em Terraform e atualização obrigatória dos arquivos de governança.

---

## 🎯 Pré-requisitos

Antes de qualquer edição, confirme:

| Item | Verificar |
|------|-----------|
| Serviço rodando localmente | `curl -s http://localhost:<PORTA>` retorna 200 |
| Porta não conflita | Checar `/srv/ops/ai-governance/PORTS.md` |
| Subdomínio não existe | Checar `/srv/ops/ai-governance/SUBDOMAINS.md` |
| Terraform disponível | `terraform -version` |
| Cloudflared ativo | `systemctl is-active cloudflared` |

---

## 📁 Arquivos que serão modificados

```
/srv/ops/terraform/cloudflare/variables.tf     ← mapa de services (Terraform IaC)
/srv/ops/terraform/cloudflare/main.tf          ← ingress rule do tunnel
/home/will/.cloudflared/config.yml             ← runtime do cloudflared
/srv/ops/ai-governance/NETWORK_MAP.md          ← OBRIGATÓRIO registrar porta + serviço
/srv/ops/ai-governance/SUBDOMAINS.md           ← OBRIGATÓRIO registrar subdomain público
```

---

## 🔧 Passo a Passo Completo

### Passo 1 — Adicionar ao mapa de services do Terraform

**Arquivo:** `/srv/ops/terraform/cloudflare/variables.tf`

Localizar o bloco `services = {` e adicionar ao final (antes do fechamento `}`):

```hcl
    <NOME_DO_SERVICO> = {
      url       = "http://localhost:<PORTA>"
      subdomain = "<NOME_SUBDOMAIN>"
    }
```

**Exemplo real (aurelia):**
```hcl
    aurelia = {
      url       = "http://localhost:3334"
      subdomain = "aurelia"
    }
```

---

### Passo 2 — Adicionar ingress rule no Terraform

**Arquivo:** `/srv/ops/terraform/cloudflare/main.tf`

Localizar o último `ingress_rule` antes do catch-all `http_status:404` e inserir:

```hcl
    ingress_rule {
      hostname = "${var.services.<NOME_DO_SERVICO>.subdomain}.${var.domain}"
      path     = ""
      service  = var.services.<NOME_DO_SERVICO>.url
    }
```

**Exemplo real (aurelia):**
```hcl
    ingress_rule {
      hostname = "${var.services.aurelia.subdomain}.${var.domain}"
      path     = ""
      service  = var.services.aurelia.url
    }
```

> ⚠️ O catch-all `ingress_rule { service = "http_status:404" }` deve SEMPRE ser o último.

---

### Passo 3 — Adicionar ao config.yml do cloudflared

**Arquivo:** `/home/will/.cloudflared/config.yml`

Inserir antes da linha `- service: http_status:404`:

```yaml
  - hostname: <SUBDOMAIN>.zappro.site
    service: http://localhost:<PORTA>
```

**Exemplo real:**
```yaml
  - hostname: aurelia.zappro.site
    service: http://localhost:3334
  - service: http_status:404
```

---

### Passo 4 — Atualizar NETWORK_MAP.md (OBRIGATÓRIO)

**Arquivo:** `/srv/ops/ai-governance/NETWORK_MAP.md`

**4a.** Na seção `Serviços Non-Docker` (ou na stack correta), adicionar linha:

```markdown
| <nome> (<tipo>) | :<PORTA> | ✅ UP | <descrição> → <subdomain>.zappro.site |
```

**4b.** Na seção `Subdomínios Públicos (Cloudflare Tunnel)`, adicionar linha:

```markdown
| `<subdomain>.zappro.site` | :<PORTA> | ✅ ativo |
```

---

### Passo 5 — Atualizar SUBDOMAINS.md (OBRIGATÓRIO)

**Arquivo:** `/srv/ops/ai-governance/SUBDOMAINS.md`

**5a.** Na tabela `Subdomínios Ativos`, adicionar linha:

```markdown
| `<subdomain>.zappro.site` | <serviço> | :<PORTA> | ✅ PÚBLICO | <descrição breve> |
```

**5b.** No diagrama de topologia (bloco de código), adicionar linha:

```
    │       └── <subdomain>.zappro.site  → localhost:<PORTA>
```

---

### Passo 6 — Terraform plan + apply

```bash
cd /srv/ops/terraform/cloudflare

# Ver o que será criado (1 DNS CNAME + 1 ingress rule)
terraform plan

# Aplicar (cria/atualiza DNS no Cloudflare)
terraform apply -auto-approve
```

**Saída esperada:**
```
Plan: 1 to add, 1 to change, 0 to destroy.
Apply complete! Resources: 1 added, 1 changed, 0 destroyed.
```

---

### Passo 7 — Recarregar cloudflared

```bash
# Tenta reload gracioso primeiro
sudo systemctl reload cloudflared || sudo systemctl restart cloudflared

# Verificar que está ativo
systemctl is-active cloudflared
```

> ⚠️ **Requer aprovação** por política de governança (`GUARDRAILS.md` — Service restart).

---

### Passo 8 — Verificação end-to-end

```bash
# 1. DNS propagado?
nslookup <subdomain>.zappro.site

# 2. Tunnel respondendo? (aguarde 10-30s após reload)
curl -sI https://<subdomain>.zappro.site | head -5

# 3. Conteúdo correto?
curl -s https://<subdomain>.zappro.site/<endpoint>
```

**Critério de aceite:**
- `HTTP/2 200` no curl -sI
- `server: cloudflare` nos headers
- Conteúdo do serviço local retornado corretamente

---

## 🗂️ Referência: Estado atual dos arquivos de governança

| Arquivo | Localização | O que registrar |
|---------|-------------|-----------------|
| `NETWORK_MAP.md` | `/srv/ops/ai-governance/` | Porta + serviço + status + subdomain |
| `SUBDOMAINS.md` | `/srv/ops/ai-governance/` | Subdomain + porta + observação de segurança |
| `variables.tf` | `/srv/ops/terraform/cloudflare/` | Entrada no mapa `services` |
| `main.tf` | `/srv/ops/terraform/cloudflare/` | `ingress_rule` no `cloudflare_tunnel_config` |
| `config.yml` | `/home/will/.cloudflared/` | Regra `hostname + service` antes do catch-all |

---

## 🔒 Segurança — Checklist antes de expor

- [ ] O serviço exige autenticação? Se não → considerar Cloudflare Access Zero Trust
- [ ] Expõe dados sensíveis (DB, métricas, secrets)? → **NÃO expor** sem auth
- [ ] Porta já listada em `PORTS.md` como interna? → Revisar antes de tornar pública
- [ ] Subdomain adicionado em `SUBDOMAINS.md`? → **Obrigatório**
- [ ] Porta adicionada em `NETWORK_MAP.md`? → **Obrigatório**

---

## 🔑 Contexto da Infraestrutura

| Recurso | Valor |
|---------|-------|
| Domínio base | `zappro.site` |
| Tunnel name | `will-zappro-homelab` |
| Tunnel ID | `8c55fcb7-fc8a-4ddc-a332-bff80012f11f` |
| Tunnel CNAME | `8c55fcb7-fc8a-4ddc-a332-bff80012f11f.cfargotunnel.com` |
| Zone ID | `c0cf47bc153a6662f884d0f91e8da7c2` |
| Terraform dir | `/srv/ops/terraform/cloudflare/` |
| Cloudflared config | `/home/will/.cloudflared/config.yml` |
| Provider version | `cloudflare/cloudflare ~> 4.0` |

---

## ⚠️ Anti-Padrões

- **Nunca** editar só o `config.yml` sem atualizar o Terraform — estado ficará dessincronizado
- **Nunca** adicionar subdomain sem atualizar `NETWORK_MAP.md` + `SUBDOMAINS.md`
- **Nunca** colocar o catch-all `http_status:404` antes das regras de hostname
- **Nunca** usar `terraform apply` sem `terraform plan` antes
- **Nunca** expor serviços internos (postgres, prometheus, Ollama) sem Cloudflare Access

---

## 📍 Quando usar esta skill

- Qualquer novo serviço do homelab que precise de acesso externo
- Ao configurar dashboard, API, ou UI de qualquer stack nova
- Ao mover serviço de Tailscale-only para público
- Para replicar o processo em outro domínio/tunnel