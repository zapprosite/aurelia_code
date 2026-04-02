# ADR 20260401 — Industrialização Sovereign Pipeline (S-56..S-66)

## Status
Aceito ✅

## Contexto
O ecossistema Aurélia precisava de uma transição definitiva para o modelo "Sovereign 2026", eliminando dependências externas (Supabase) e reforçando a coordenação local entre múltiplos agentes através de rede local (Redis). Além disso, o pipeline de entrada do Telegram estava excessivamente complexo (> 500 linhas), o que dificultava a manutenção e auditoria de segurança.

## Decisões

### 1. Remoção do Supabase
Toda a infraestrutura de dados foi consolidada em **SQLite** (persistência relacional) e **Qdrant** (memória semântica). Referências ao Supabase foram eliminadas do código e testes para reduzir a superfície de ataque e dependência de nuvem externa.

### 2. Memória Compartilhada via Redis (`SharedMemory`)
Implementamos uma interface de memória compartilhada usando **Redis Pub/Sub**. Isso permite sincronização síncrona/assíncrona entre o enxame de agentes (Swarm), suportando os módulos KAIROS e Dream.

### 3. Refatoração do `input_pipeline.go`
O pipeline de entrada foi reestruturado em funções menores (< 100 linhas), com responsabilidades claras:
- **handleCommonInputPreProcess**: Sanitização, Porteiro e Memory Commands.
- **runExecutionAndDelivery**: Execução técnica e entrega final.

### 4. Sistema de Identidade Persistente
Extraímos os prompts de sistema para arquivos externos (`AURELIA.md`, `IDENTITY.md`, `HEARTBEAT.md`) montados como volumes, permitindo alterações de personalidade sem necessidade de recompilação.

## Consequências
- **Melhoria**: Maior desacoplamento e estabilidade nas retentativas de rede (anti-retry deduplication).
- **Eficiência**: Menos overhead de memória e latência com a remoção de requisições HTTPS externas para DB.
- **Governança**: Auditoria de segredos e ADRs garantem que o repositório permanece industrializado.
