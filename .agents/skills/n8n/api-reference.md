# n8n API Reference

## Autenticação
- Header: X-N8N-API-KEY: <sua_chave>
- Base URL: https://<seu-dominio>/api/v1

## Endpoints principais

### Workflows
- GET    /workflows                Lista todos os workflows
- GET    /workflows/:id            Busca workflow por ID
- POST   /workflows                Cria novo workflow
- PUT    /workflows/:id            Atualiza workflow
- DELETE /workflows/:id            Deleta workflow
- POST   /workflows/:id/activate   Ativa workflow
- POST   /workflows/:id/deactivate Desativa workflow

### Execuções
- GET    /executions               Lista execuções
- GET    /executions/:id           Busca execução por ID
- POST   /executions               Dispara execução manual
- DELETE /executions/:id           Deleta execução

### Credenciais
- GET    /credentials              Lista credenciais
- POST   /credentials              Cria credencial
- DELETE /credentials/:id          Deleta credencial

## Webhooks
Todo workflow com trigger Webhook gera uma URL no formato:
https://<dominio>/webhook/<uuid-do-workflow>

Em produção usar sempre HTTPS. Proteger com header de autenticação customizado quando possível.

## Dicas
- Disparar workflow via API: POST /workflows/:id/execute com body JSON
- Autenticar sempre via header, nunca via query string
- Execuções ficam salvas por padrão. Configure data pruning para não lotar o banco
