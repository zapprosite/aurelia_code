# Steel Container - Browser Isolation for Computer Use
# ADR: 20260328-container-steel-browser-isolation
# Stagehand + Playwright em container isolado

FROM mcr.microsoft.com/playwright:v1.50.0-jammy

# Non-root user para segurança
RUN groupadd -r steel && useradd -r -g steel steel

# Dependências de sistema
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    wget \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Stagehand e dependências Node
WORKDIR /app

# Copia package files
COPY mcp-servers/stagehand/package*.json ./

# Instala dependências (sem devDependencies para减小 imagem)
RUN npm ci --omit=dev --ignore-scripts

# Copia source do Stagehand
COPY mcp-servers/stagehand/src ./src

# Build TypeScript
RUN npm run build || echo "Build pode falhar se tsconfig precisar de ajustes"

# Permissões
RUN chown -R steel:steel /app

# Wrapper script
COPY mcp-servers/steel/start.sh /start.sh
RUN chmod +x /start.sh

# Ports
EXPOSE 9222  # DevTools
EXPOSE 3000  # Health check

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

# User steel
USER steel

# Entrypoint via wrapper
ENTRYPOINT ["/start.sh"]
