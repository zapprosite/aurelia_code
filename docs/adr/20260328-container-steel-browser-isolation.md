# ADR 20260328: Container "Steel" para Browser Isolation (Computer Use)

## Status
🔴 Proposto (P1 - Crítico para Computer Use)

## Contexto
O ADR `20260328-implementacao-jarvis-voice-e-computer-use` define que o navegador para computer use deve rodar em container isolado para:
1. Prevenir crashes em cascata (browser crash ≠ agent crash)
2. Garantir segurança (sandbox)
3. Permitir persistência de sessão entre ações
4. Isolar dependências (Playwright, Chromium)

## Decisões Arquiteturais

### 1. Arquitetura do Container Steel

```dockerfile
# mcp-servers/steel/Dockerfile
FROM mcr.microsoft.com/playwright:v1.50.0-jammy

# Dependências de browser automation
RUN apt-get update && apt-get install -y \
    chromium-browser \
    chromium-chromedriver \
    && rm -rf /var/lib/apt/lists/*

# Stagehand e dependências Node
WORKDIR /app
COPY mcp-servers/stagehand/package*.json ./
RUN npm ci --omit=dev

COPY mcp-servers/stagehand/src ./src
RUN npm run build

# Wrapper script
COPY mcp-servers/steel/start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 9222  # DevTools debugging port

ENTRYPOINT ["/start.sh"]
CMD ["node", "/app/dist/index.js"]
```

### 2. Session Management

O Steel container mantém estado entre chamadas:

```go
// Sessão persiste enquanto o container está vivo
type SteelSession struct {
    ContainerID  string
    SessionToken string
    LastActive   time.Time
    BrowserCtx   *playwright.BrowserContext
}
```

- **Idle timeout**: 10 minutos sem atividade → container é reciclado
- **Max session**: 1 hora → força restart para limpar memória
- **Volume persistido**: `/steel/sessions/` para state entre restarts

### 3. Health & Lifecycle

```go
// internal/steel/container.go
type ContainerManager struct {
    client    *docker.Client
    imageName string
    pool      *pool.Pool  // Connect pool de containers
}

func (m *ContainerManager) Acquire(ctx context.Context) (*SteelSession, error)
func (m *ContainerManager) Release(session *SteelSession) error
func (m *ContainerManager) Health() error  // healthcheck
```

### 4. Docker Compose Integration

```yaml
# configs/steel/docker-compose.yml
services:
  steel:
    build:
      context: .
      dockerfile: mcp-servers/steel/Dockerfile
    image: aurelia/steel:latest
    container_name: aurelia-steel-01
    environment:
      - STAGEHAND_MODEL=aurelia-smart
      - LITELLM_BASE_URL=http://litellm:4000/v1
      - NODE_ENV=production
    ports:
      - "9222:9222"  # DevTools
    volumes:
      - steel-sessions:/app/sessions
    restart: unless-stopped
    deploy:
      resources:
        limits:
          memory: 2G
        reservations:
          memory: 512M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9222/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  steel-sessions:
    driver: local
```

### 5. Resource Limits

| Recurso | Limite |
|---------|--------|
| Memory | 2GB (hard), 512MB (soft) |
| CPU | 2 cores |
| Network | bridge (acesso internet) |
| Storage | 1GB (session logs) |
| Chromium | 1 instâncias concurrent |

## Consequências

### Positivas
- Isolamento total: browser crash ≠ agent crash
- Security sandbox: Chromium roda em sandbox nativo
- Reproducibility: ambiente controlado e versionado
- Escalabilidade: múltiplos containers para paralelismo

### Negativas
- Latência de cold start (~5-10s para iniciar container)
- Overhead de recursos: ~1GB RAM por container
- Complexidade de debugging: container é opaco

### Trade-offs
- Podman vs Docker: Docker é mais suportado, Podman é mais seguro
- Host network vs Bridge: Bridge é mais isolado mas adiciona latência

## Dependências
- ⚠️ `mcp-servers/stagehand/src/index.ts` (já existe, precisa de build)
- ❌ `mcp-servers/steel/Dockerfile` (NÃO EXISTE)
- ❌ `mcp-servers/steel/start.sh` (NÃO EXISTE)
- ❌ `configs/steel/docker-compose.yml` (NÃO EXISTE)
- ⚠️ Docker daemon acessível do processo Go

## Referências
- [ADR-20260328-jarvis-voice-computer-use.md](./20260328-implementacao-jarvis-voice-e-computer-use.md)
- [mcp-servers/stagehand/src/index.ts](../../mcp-servers/stagehand/src/index.ts)
- [Stagehand Docker Setup](https://github.com/browserbasehq/stagehand/tree/main/examples/docker)

## Links Obrigatórios
- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)

---
**Data**: 2026-03-28
**Status**: Proposto
**Autor**: Claude (Principal Engineer)
**Slice**: feature/neon-sentinel
**Progress**: 0%
