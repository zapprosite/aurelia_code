# ADR-20260326: Acessibilidade e Portabilidade Universal

## Status
Proposto

## Contexto
O ecossistema Aurelia foi desenvolvido com foco em "Soberania Industrial", exigindo hardware de ponta (RTX 4090, etc.) e conhecimento técnico sênior para operação. Para expandir o uso para usuários menos técnicos ("leigos") e permitir a portabilidade para hardwares mais modestos (Macs, Laptops sem GPU dedicada), precisamos simplificar a entrada no sistema.

## Decisão
1. **Script de Setup Único (`iniciar.sh`)**: Criar um script amigável no root que abstraia a complexidade do Go, Docker e dependências de sistema.
2. **Modo "Lite" Explícito**: Implementar uma configuração `LITE_MODE=true` que desativa o processamento de voz local pesado e busca semântica GPU-intensive em favor de APIs Cloud (serviços externos) de baixo custo ou processamento CPU.
3. **Documentação para Iniciantes**: Criar `docs/BEM-VINDO.md` com linguagem não técnica, focando no valor entregue (o que o bot faz) em vez de como ele funciona internamente.
4. **Modularidade de Soberania**: Permitir que o usuário escolha entre "Soberania Total" (Local) e "Conveniência" (Cloud) via flag de setup.

## Consequências
- **Prós**: Maior base de usuários; facilidade de demonstração; redução da barreira de entrada.
- **Contras**: Maior complexidade na manutenção de caminhos de execução alternativos (Cloud vs Local).
