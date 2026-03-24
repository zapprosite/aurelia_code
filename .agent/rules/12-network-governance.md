---
description: Governança obrigatória de rede, portas e subdomínios antes de qualquer mudança de infraestrutura.
id: 12-network-governance
---

# 🌐 Regra 12: Governança de Rede, Portas e Subdomínios

Antes de qualquer mudança que envolva **rede, porta, serviço, container ou subdomínio**, o agente DEVE consultar os documentos de governança de rede.

<directives>
1. **Leitura obrigatória** antes de tocar rede/serviços:
   - `/srv/ops/ai-governance/NETWORK_MAP.md` — mapa canônico de portas, subdomínios, GPU e serviços do homelab
   - `/srv/ops/ai-governance/SUBDOMAINS.md` — registro oficial de subdomínios (se existir)

2. **Adicionar subdomínio**: use a skill `/add-subdomain` — ela executa o processo completo via Terraform + Cloudflare Tunnel e atualiza todos os arquivos de governança obrigatórios.

3. **Cloudflare Access (Zero Trust)**: governado por `/srv/ops/terraform/cloudflare/access.tf`.
   - Ver `docs/governance/S-23-cloudflare-access.md` para o processo IaC.
   - Toda nova exposição pública deve ter uma entrada no Access Application correspondente.

4. **Expor porta pública** sem atualizar NETWORK_MAP.md é **PROIBIDO** (§ GUARDRAILS.md).

5. **Tier C obrigatório** (log + aprovação humana):
   - Mudanças em firewall/UFW
   - Novos túneis Cloudflare
   - Roteamento de subdomínios novos
   - `terraform apply` em `/srv/ops/terraform/cloudflare/`
</directives>

## Referências

- [`/srv/ops/ai-governance/NETWORK_MAP.md`](/srv/ops/ai-governance/NETWORK_MAP.md)
- [`/srv/ops/terraform/cloudflare/access.tf`](/srv/ops/terraform/cloudflare/access.tf)
- [`docs/governance/S-23-cloudflare-access.md`](../../docs/governance/S-23-cloudflare-access.md)
- [Skill: add-subdomain](../skills/add-subdomain/SKILL.md)
