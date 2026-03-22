> [!NOTE]
> Status: ✅ Arquivado / Concluído em 22/03/2026

# ADR-20260322-dashboard-cockpit-cmdk

## Status
Proposta

## Contexto
O Dashboard Ultratrink atual é majoritariamente passivo (apenas leitura via SSE). Para transformar o dashboard em um verdadeiro "Cockpit", é necessária uma interface de comando ativa que permita ao operador humano disparar ações rápidas no enxame de agentes sem sair da UI. O padrão da indústria para isso é o atalho `CMD+K` (Command Palette).

## Decisão
Implementaremos um sistema de comandos bidirecional:
1.  **Backend (Go)**: Novo endpoint `/api/commands` (POST) que recebe um `CommandRequest` e o despacha para os serviços internos (Agent, Memory, Config).
2.  **Frontend (React)**: Componente `CommandMenu` utilizando `framer-motion` para animações suaves e `lucide-react` para ícones. O atalho global `Ctrl+K` ou `Cmd+K` abrirá a interface.
3.  **Feedback Visual**: As ações disparadas pelo Cockpit emitirão eventos SSE que retornarão ao Feed, fechando o loop de feedback.

## Comandos Iniciais
- `Sync Knowledge`: Força a sincronização do ai-context.
- `Switch Model`: Altera o modelo LLM ativo (Ollama/Gemini).
- `Clear Memory`: Reseta o contexto do Agent Loop atual.
- `Restart Services`: Reinicia serviços do Homelab (Ollama, Qdrant).

## Consequências
- Aumenta a autonomia do operador na interface web.
- Introduz necessidade de autenticação/segurança futura (por enquanto, o homelab opera em rede confiável).
- Melhora a observabilidade ativa.
