# ADR 20260330: Integração Formal do Smart Router ao Docker Compose

## Status
Proposto

## Contexto
O serviço `aurelia-smart-router` (LiteLLM) estava sendo iniciado de forma manual ou por mecanismos fora do `docker-compose.yml` principal. Isso resultou em:
1. **Falha na Injeção de Segredos**: O container não recebia chaves de API críticas (como `GOOGLE_API_KEY`) definidas no sistema.
2. **Falha de Conectividade**: O roteador tentava acessar o Ollama via `localhost`, o que dentro do container isolado fallhava.
3. **Conflito de Portas**: Falta de orquestração centralizada causava tentativas de ocupação de portas já em uso (`8484`, `3334`).

## Decisão
Mover a definição do serviço `aurelia-smart-router` para o arquivo `docker-compose.yml` na raiz do projeto.

Configurações específicas:
- **Imagem**: `ghcr.io/berriai/litellm:main-stable`
- **Ambiente**: Uso de `env_file: .env` para garantir paridade com o ecossistema.
- **Portas**: Mapeamento fixo para `4000:4000` (API), `8484` (Health) e `3334` (Dashboard/UI).
- **Rede**: Adição de `extra_hosts` para `host.docker.internal:host-gateway` para permitir acesso ao Ollama no host.
- **Configuração**: Mapeamento do arquivo local `configs/litellm/config.yaml` para `/app/config.yaml`.

## Consequências
- **Positivas**: Centralização da governança, injeção automática de segredos do `.env`, e maior estabilidade de rede.
- **Neutras**: Necessidade de atualizar o `config.yaml` para usar `host.docker.internal`.
