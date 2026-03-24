---
name: Api Design
description: Design de APIs robustas seguindo o contrato Zod-First e tRPC (Soberano 2026).
phases: [P, R]
---

# 🔌 API Design: Zod-First & tRPC (Sovereign 2026)

Esta skill habilita o Antigravity a projetar interfaces de comunicação seguras e tipadas entre o Backend (Go/tRPC) e o Frontend (Next.js/React).

## 🏛️ Princípios de Design (Industrial)

### 1. Contrato Zod-First
- **Regra**: Todos os esquemas de dados devem residir exclusivamente em `packages/zod-schemas/`.
- **Compartilhamento**: É proibido duplicar lógica de validação. O esquema definido no package deve ser consumido tanto pelo backend quanto pelo frontend.
- **Tipagem**: Utilize inferência automática de tipos a partir do Zod (`z.infer<typeof schema>`).

### 2. tRPC & Protocolo
- **Soberania**: O tRPC é o protocolo padrão para comunicação interna. Evite REST puro para lógica de aplicação complexa.
- **Estrutura**: Agrupe rotas em `internal/gateway/router/` seguindo a lógica de domínios (infra, agent, voice).
- **Segurança**: Todo end-point que altera estado deve exigir auditoria de permissões (Middleware de Auth).

### 3. Observabilidade (Logging Estruturado)
- Todo end-point tRPC deve implementar logging estruturado via `slog`.
- Inclua metadados como `conversation_id`, `user_id` e `trace_id` em todas as requisições para facilitar o debug remoto.

## 🚀 Workflow de Design
1. **Definição do Schema**: Crie o arquivo `.ts` em `packages/zod-schemas/`.
2. **Registro do Router**: Implemente o handler em Go seguindo o contrato definido.
3. **Drafting Frontend**: Utilize os hooks do tRPC no React para consumo seguro.
4. **Validação**: Realize testes de fumaça (Smoke Tests) antes de liberar o end-point.

## 📍 Quando usar
- Ao criar novos módulos no Dashboard.
- Ao adicionar novas ferramentas (Tools) que precisam expor estado via API.
- Ao refatorar comunicações legadas entre serviços.

## 🚫 Anti-Padrões
- Definir tipos manuais que podem divergir do esquema de validação.
- Expor dados sensíveis (tokens, PII) em retornos de API sem filtragem.
- Criar rotas tRPC "gigantes". Prefira routers atômicos e especializados.