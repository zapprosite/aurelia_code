# S-23: Cloudflare Access — Zero Trust via Terraform

**Status:** ✅ Implementado via IaC
**Terraform:** `/srv/ops/terraform/cloudflare/access.tf`

## Recursos criados

- `cloudflare_access_application` — uma aplicação Zero Trust por serviço (n8n, qdrant, caprover, supabase, studio, monitor, aurelia)
- `cloudflare_access_policy` — política `owners` com `decision = allow` por email

## Como usar

```bash
cd /srv/ops/terraform/cloudflare

# 1. Adicionar email autorizado no terraform.tfvars:
# allowed_emails = ["will@example.com"]

# 2. Aplicar
terraform plan
terraform apply
```

## Variável

```hcl
variable "allowed_emails" {
  type    = list(string)
  default = []
}
```

## Serviços protegidos

| Serviço   | Domínio                |
|-----------|------------------------|
| aurelia   | aurelia.zappro.site    |
| n8n       | n8n.zappro.site        |
| qdrant    | qdrant.zappro.site     |
| caprover  | cap.zappro.site        |
| supabase  | supabase.zappro.site   |
| studio    | studio.zappro.site     |
| monitor   | monitor.zappro.site    |

## Sem código no backend

A proteção é enforced na edge Cloudflare. O backend Go não precisa de middleware de auth.
Se no futuro quiser validação de JWT CF-Access, adicionar middleware em `internal/dashboard/dashboard.go`.
