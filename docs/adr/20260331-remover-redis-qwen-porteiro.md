# ADR-20260331: Substituição do Redis e Qwen 0.5b por Deduplicação In-Memory e Regex no Porteiro Middleware

## Status
Aceito

## Contexto
O `PorteiroMiddleware` da Aurélia (chatbot Go + Telegram + modelos locais, home lab) operava com duas dependências críticas no hot path de cada mensagem:
1. **Deduplicação via Redis SetNX**: Evitava que mensagens duplicadas do Telegram fossem processadas. Contudo, em casos de indisponibilidade do Redis, o método `Deduplicate()` realizava *fail-open*, retornando sempre `isDupe = false`, o que permitia o processamento duplicado de todas as mensagens subsequentes.
2. **Análise de Segurança via Qwen 0.5b**: Utilizado através do LiteLLM para polimento de output e proteção contra prompt injections. Esse processamento via LLM adicionava aproximadamente 800ms de latência extra por mensagem.

Com o amadurecimento do ecossistema e a busca por resiliência e baixa latência (focando em single-instance num servidor com RTX 4090 e Ubuntu Server), fez-se necessário revisitar estas dependências para reduzir a complexidade operacional, eliminar latência desnecessária e operar independentemente do Redis.

## Decisão
Removemos o Redis e o LLM Qwen 0.5b do Porteiro Middleware. A nova implementação baseia-se em soluções locais, in-memory e determinísticas:
- **Deduplicação In-Memory**: Utilização de `sync.Map` armazenando uma struct `dupeEntry{expiresAt time.Time}`, complementado por uma goroutine dedicada para *cleanup* (rodando a cada 60s). O TTL do *lock* foi ajustado para 15 segundos (anteriormente 10s no Redis) para fornecer margem extra para retries do Telegram.
- **Validação de Injections (IsSafe)**: Substituição do LLM por validação *guardrail* baseada em *regex* puro. Inclui verificação contra *ignore instructions*, tentativas de *jailbreak* e *DAN mode*.
- **Polimento e Prevenção de Vazamentos (IsOutputSafe e PolishOutput)**: A função de *polish* (PolishOutput) agora simplesmente aplica `strings.TrimSpace(content)` em vez de reescrever via LLM. A função `IsOutputSafe` mantém a varredura via *regex* para proteção estrita de segredos, abordando chaves da OpenAI, tokens do GitHub e tokens do Telegram.

## Alternativas Consideradas

| Alternativa | Prós | Contras | Motivo da Rejeição |
|-------------|------|---------|--------------------|
| Manter Redis com reconexão automática | Deduplicação persiste a restarts do bot | Adiciona complexidade operacional e overhead de rede local | Não entrega benefício tangível para cenários single-instance de home lab |
| Substituir Redis por SQLite para dedup persistente | Persistência entre reinicializações sem dependência de container externo | I/O contínuo no disco para um lock efêmero (TTL de 15s) | Overhead e complexidade desnecessária dadas as características de TTL curto |
| Manter Qwen 0.5b como configuração opcional | Mantém a capacidade de análise semântica avançada de segurança | Aumenta a superfície de configuração sem ganho prático comprovado | Adição recorrente de 800ms à latência do hot path é injustificável |

## Consequências
- **Positivas**: Redução da stack infraestrutural via exclusão de 1 container (Redis). Redução expressiva na latência (~800ms economizados por mensagem em média). Deduplicação executada de forma ágil e determinística, sem gargalos de rede (network hops). Implementação de zero downtime com design enxuto.
- **Negativas**: A deduplicação in-memory não sobrevive a reinicializações completas do processo da Aurélia.
- **Riscos Adicionais**: Retrying de mensagens durante o restart do container não possivelmente dedupadas nos primeiros instantes. No entanto, para a realidade atual, o fato do restart ocorrer em menos de 5 segundos contra um TTL de 15s torna o problema gerenciável em termos operacionais.

## Tech Debt
- **Sincronização Distribuída**: O mecanismo `sync.Map` é stateful apenas para a instância atual. Se no futuro houver exigência de múltiplos réplicas (pods) da aplicação processando mensagens concorrentemente, a deduplicação inevitavelmente apresentará falhas sem um state store centralizado (Redis, Etcd, etc). Para esse cenário *scale-out* futuro, a reinserção do Redis ou implementação de mensageria focada será necessária. O estado atual deve ser tratado documentado como tech debt não-urgente nesse tocante.

## Referências
- `internal/middleware/porteiro.go`: Rewrite completo substituindo Redis/LLM.
- `docker-compose.yml`: Remoção do service `redis`.
- `.env.example`: Remoção das variáveis `REDIS_*`.
- Documentação de design e arquitetura Soberano 2026 (`docs/architecture-2026.md`).
