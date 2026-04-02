# ADR 20260401 — Infraestrutura Pós-Sovereign (S-67, S-68)

## Status
Aceito ✅

## Contexto
Após a consolidação da arquitetura Sovereign 2026, iniciamos a otimização da infraestrutura local para suportar maior carga de modelos e ferramentas de busca autônoma. O armazenamento padrão dos modelos do Ollama estava competindo por espaço com o sistema operacional, e a Aurélia carecia de uma ferramenta de busca web local e privada.

## Decisões

### 1. Migração de Modelos Ollama (S-67)
Decidimos mover todos os modelos do Ollama para o diretório `/srv/models`, montado em um dataset ZFS dedicado (ZFS tank). 
- **Configuração**: Utilizamos um `override.conf` do systemd para injetar a variável `OLLAMA_MODELS=/srv/models`.
- **Benefício**: Isolamento de storage, snapshots ZFS para modelos e performance superior em IOPS.

### 2. Integração SearXNG (S-68)
Implementamos a ferramenta `search_web_local` utilizando o motor de busca **SearXNG** rodando localmente.
- **Implementação**: `internal/tools/searxng.go` consome a API JSON do SearXNG.
- **Configuração**: URL definida via env `SEARXNG_URL` com fallback seguro.
- **Privacidade**: Todo o tráfego de busca web agora é agregado e anonimizado localmente antes de sair para os motores externos.

## Consequências
- **Vantagem**: Maior resiliência do sistema operacional (espaço livre) e buscas web mais rápidas e privadas.
- **Monitoramento**: O status do Ollama deve ser verificado via `systemctl show ollama` para garantir que o override está ativo.
