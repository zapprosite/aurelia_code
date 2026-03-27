---
slug: deployment
category: operations
generatedAt: 2026-03-27T10:44:47.031Z
relevantFiles:
  - Dockerfile
  - docker-compose.yml
  - .github/workflows/ci.yml
  - .github/workflows/deploy.yml
  - .github/workflows/gitleaks.yml
  - .github/workflows/golangci-lint.yml
  - .github/workflows/govulncheck.yml
  - homelab-bibliotheca/skills/open-claw/skills/ivangdavila/docker-compose
  - homelab-bibliotheca/skills/open-claw/skills/lxgicstudios/docker-compose-gen
  - homelab-bibliotheca/skills/open-claw/skills/mosonchan2023/docker-compose-manager
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