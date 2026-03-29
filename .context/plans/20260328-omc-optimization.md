# Plan: OMC Infra Optimization 2026

## Goal
Optimize Home Lab resource usage and fix healthcheck inconsistencies using the `oh-my-claude-sisyphus` (OMC) tools.

## Worktree Context
- **Path**: `../20260328-omc-optimization`
- **Branch**: `feature/omc-optimization`

## Tasks
1. **LiteLLM Healthcheck Alignment**:
   - Verify why Docker reports `unhealthy` despite API response.
   - Adjust `HEALTHCHECK` command in `docker-compose.yml` or CapRover config.
2. **Whisper VRAM Optimization**:
   - Migrate `speaches` from `faster-whisper-large-v3` to `distil-large-v3`.
   - Goal: Reduce VRAM usage from 4.2GB to approx 1.5GB.
3. **Daemon Logging Refinement**:
   - Use OMC `ralph` mode to audit `aurelia.service` logs.
   - Implement auto-healing logic for the daemon.

## Status
- [ ] LiteLLM Health Fix
- [ ] Whisper VRAM Optimization
- [ ] Logging & Auto-healing
