---
description: Prioriza a análise interna antes de qualquer busca externa (web/IA).
id: 02-local-first
---

# 🔍 Regra 02: Descoberta Local Primeiro

A inteligência deve ser extraída primeiramente do código e documentação existentes.

<directives>
1. **Inspeção Ativa**: Use ferramentas de sistema (`ls`, `grep`, `find`) e MCP (`list_dir`, `view_file`) antes de perguntar ao usuário ou buscar na web.
2. **Análise de Contexto**: O diretório `.context/` deve ser a primeira parada para entender o estado atual de features e planos.
3. **Anti-Hallucinação**: É proibido assumir a existência de módulos ou padrões. Verifique fisicamente a estrutura antes de referenciá-la.
</directives>
---
description: Define como os agentes interagem com ferramentas externas e entre si.
id: 07-shared-mcp
---

# 🤝 Regra 07: MCP Compartilhado & Sem Aninhamento

A interoperabilidade entre motores de execução segue padrões estritos.

<directives>
1. **Sem CLI-nesting**: Nunca chame um motor (`claude`, `opencode`) de dentro de outro como ferramenta de terminal.
2. **Canais de Dados**: Use servidores MCP como ponte comum de conhecimento (Filesystem, Postgres, etc).
3. **Handoffs**: Transferências de contexto devem ocorrer via artefatos Markdown e metadados Git.
</directives>
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
