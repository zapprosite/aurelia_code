# ADR-20260323-Resource-Governance-RAM-IPv6

## Status
Resolvido

- slug: resource-governance-ram-ipv6
- json de continuidade: docs/adr/taskmaster/20260323-resource-governance-ram-ipv6.json

## Contexto
O Home Lab está operando no limite de sua capacidade física de RAM (30GiB totais, com ~29GiB em uso constante e 3.1GiB processados em swap). 
A análise via `docker stats` revelou que a maioria dos containers (incluindo a stack do Supabase, n8n, LiteLLM e motores de IA) estão rodando sem limites de memória explícitos (`limit: 30.48GiB`), permitindo que picos de RSS (Resident Set Size) causem instabilidade no sistema operacional (OOM conditions).

Além disso, a rede interna apresenta múltiplos endereços IPv6 com status `deprecated`, reflexo de endereços temporários não renovados ou trocas de prefixo delegados (PD) pelo ISP que não foram limpas.

## Decisão
Implementar uma política de **Capping de Recursos** (Recurso Limite Rígido) para todos os serviços baseados em Docker e uma rotina de limpeza de rede.

### 1. Limites de Memória (Hard Capping)
Todo serviço deve declarar limites de memória em seu arquivo `docker-compose.yml` seguindo as categorias:
- **Bancos de Dados (Supabase, Postgres, Qdrant)**: Máximo 2GiB de RAM por instância (ajustável sob demanda, mas com teto).
- **Motores de Voz/IA (Kokoro, Whisper, XTTS)**:
  - Kokoro: 2GiB
  - Whisper-local: 3GiB
  - XTTS: 4GiB (devido ao modelo de carregamento de vozes)
- **Middleware & Studio (Studio, Kong, LiteLLM)**: Máximo 512MiB a 1GiB.

### 2. Governança de Rede (IPv6)
- Descartar endereços `deprecated` via script de manutenção.
- Configurar o kernel (`sysctl`) para reduzir o tempo de retenção de endereços temporários inválidos.

## Consequências
- **Positivas**: Maior estabilidade do sistema, fim do uso excessivo de Swap, previsibilidade de recursos para novos agentes.
- **Negativas**: Serviços podem sofrer OOM Kill se subestimarem seu uso real. Monitoramento constante do Grafana será necessário.
- **Riscos**: XTTS e Whisper podem falhar em arquivos extremamente longos se o buffer de RAM estourar o limite de 3GiB/4GiB.

## Implementação
As alterações serão feitas nos arquivos:
- `/srv/apps/supabase/docker/docker-compose.yml`
- `/srv/apps/voice/docker-compose.yml`
- `/srv/apps/litellm/docker-compose.yml`
- `/srv/apps/monitoring/docker-compose.yml`
