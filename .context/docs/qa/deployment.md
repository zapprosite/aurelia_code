---
slug: deployment
category: operations
generatedAt: 2026-03-31T07:48:20.732Z
relevantFiles:
  - Dockerfile
  - docker-compose.yml
  - .github/workflows/ci.yml
  - .github/workflows/deploy.yml
  - .github/workflows/gitleaks.yml
  - .github/workflows/golangci-lint.yml
  - .github/workflows/govulncheck.yml
  - mcp-servers/steel/Dockerfile
  - mcp-servers/steel/docker-compose.yml
  - services/aurelia-api/Dockerfile
---

# How do I deploy this project?

## Deployment

### Docker

This project includes Docker configuration.

```bash
docker build -t app .
docker run -p 3000:3000 app
```

### CI/CD

CI/CD pipelines are configured for this project.
Check `.github/workflows/` or equivalent for pipeline configuration.