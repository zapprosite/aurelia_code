# ADR-20260329-os-native-god-mode

## Status
Proposto / SOTA 2026.1

## Contexto
Até agora, a Aurélia operava majoritariamente via Browser ou CLI básica. Para um assistente soberano, é necessário controle total sobre o host Ubuntu, permitindo gerenciar processos, rede, storage e interface gráfica de forma segura e autônoma.

## Decisão
Implementamos a soberania total via **MCP OS Controller**:
1. **Host-Native MCP**: Servidores MCP rodando localmente com permissões elevadas (auditadas pelo Porteiro).
2. **GNOME Integration**: Controle de janelas, brilho, volume e notificações via DBus/GSettings.
3. **Infrastructure Control**: Gestão nativa de Docker, ZFS e NVIDIA-SMI para monitoramento de saúde do Homelab.
4. **Shell Sovereign**: Um subset de comandos Bash "Safe-by-Design" que o Jarvis pode executar para auto-correção de sistema.

## Consequências
- **Poder Total**: O Jarvis agora pode "se consertar" (ex: reiniciar o Ollama se houver hang).
- **Segurança**: Risco de comandos destrutivos. Implementamos o `Sovereign Guard` que exige aprovação manual para comandos `rm -rf /` ou similares, independente do modo.
- **Observabilidade**: Todos os comandos do "God Mode" são logados no dashboard local para auditoria posterior.
