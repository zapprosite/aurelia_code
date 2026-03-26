# ADR 20260326-Zero-Hardcode-Policy

## Contexto
O ecossistema Aurélia estava sofrendo com a persistência de segredos em texto claro em arquivos de configuração (`app.json`) e potencial hardcode no código-fonte. Isso viola os princípios de infraestrutura soberana e segurança industrial.

## Decisão
Implementar a política **Zero Hardcode** em todo o monorepo:
1. **Mascaramento Automático**: O motor de configuração (`internal/config/config.go`) deve substituir segredos pelo placeholder `{chave-para-env}` antes de qualquer persistência em disco.
2. **Placeholders no Código**: É proibido o uso de strings reais de tokens/chaves no código. Devem ser usados placeholders para documentação e variáveis de ambiente para execução.
3. **Paridade Estrita**: O arquivo `.env.example` deve ser um espelho exato do `.env`, garantindo que todas as chaves necessárias estejam mapeadas.

## Consequências
- Maior segurança contra vazamentos acidentais via logs ou arquivos de configuração.
- Facilidade de portabilidade e setup em novos ambientes.
- Necessidade de gerenciar segredos exclusivamente via `.env` ou variáveis de sistema.
