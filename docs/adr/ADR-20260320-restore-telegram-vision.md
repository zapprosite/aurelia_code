# ADR-20260320-restore-telegram-vision: Restauração da Funcionalidade de Visão e OCR

## Status
**Accepted / Concluído**

## Contexto
Durante a auditoria de integridade do monorepo, identificamos que as rotas de tratamento de imagem (`handlePhoto`, `handleDocument` de imagem, `Album` storage) do bot do Telegram (`internal/telegram`) estavam referenciadas mas fisicamente ausentes do sistema de arquivos. Para restaurar o build "Verde", as rotas foram comentadas na [ADR-20260320-repository-clean-audit-restore.md](file:///home/will/aurelia/docs/adr/ADR-20260320-repository-clean-audit-restore.md).

Este slice visa recuperar e reintegrar essas funcionalidades de forma robusta e otimizada (Visão por GPU).

## Decisão
1.  **Arqueologia de Código**: Investigar o histórico de commits para localizar os arquivos deletados (ex: `bot_photo.go`, `media.go`).
2.  **Restauração Funcional**: Reintegrar os métodos `handlePhoto`, `storeAlbumPhoto`, `flushAlbumPhotos` e helpers de detecção de MIME.
3.  **Reativação Gradual**:
    - Re-habilitar `OnPhoto` no `bot_middleware.go`.
    - Des-comentar e atualizar testes em `input_test.go`.
4.  **Fallback de Visão**: Garantir que o sistema reporte corretamente se o modelo LLM selecionado não suportar visão (ex: Ollama sem LLaVA).

## Resumo Executivo (Arquitetura Sênior)

A funcionalidade de **Visão Multimodal** foi totalmente restaurada e **melhorada**.
1. **Álbuns**: Implementada detecção de álbuns (`MediaGroupID`) com debounce de 2s para processamento em lote.
2. **Contexto**: A 'inputSession' agora transporta 'agent.Message' completo, permitindo que o LLM receba imagens nativamente.
3. **Reuso**: Implementado `recentMedia` (cache de 3min) para responder a comandos como "O que tem nesta imagem?" após o envio.

> [!IMPORTANT]
> O build global está verde e todos os testes do pacote `telegram` estão passando.

## Consequências
- **Positivo**: Recuperação de feature crítica de visão e percepção multimodal.
- **Riscos**: Conflitos de dependência entre pacotes de 'agent' e 'telegram' (débito de acoplamento).

## Verificação Técnica
- `go test -v ./internal/telegram` (Validar OCR e Álbuns).
- Log de simulação: Envios de fotos via bot de teste.
