---
title: Operational Governance — Incident Response Playbook
description: Runbooks para responder a service down, GPU OOM, ZFS degraded, Tunnel down, credential compromise
owner: codex
updated: 2026-03-20
---

# Operational Governance — Incident Response Playbook

**Purpose:** Procedimentos estruturados para triagem, comunicação e resolução de incidentes
**Authority:** ADR-20260319-Polish-Governance-All / OPERATIONAL_GOVERNANCE
**Updated:** 2026-03-20

---

## Incident Severity Levels

| Severidade | SLA Resposta | SLA Resolução | Exemplos |
|----------|---|---|---|
| **P1 (Critical)** | 5 min | 15 min | All services down, credential breach, data loss |
| **P2 (High)** | 15 min | 1 hour | Single service degraded, backup missing, VRAM >95% |
| **P3 (Medium)** | 1 hour | 4 hours | Slow API response, container restart loop, disk warning |
| **P4 (Low)** | 4 hours | 24 hours | Informational alerts, optimization opportunities |

---

## Incident Declaration

### Automatic Declaration
Severity auto-escalated when health-check detects:
- Container down → P1
- Qdrant/Postgres unreachable → P1 (if CRITICAL tables)
- Disk >95% → P1
- Backup missing >24h → P1
- Memory sync stale >1h → P2

### Manual Declaration
**Codex** declares incident when notified of service issue:
```
Title: [P{1-4}] {Service}: {Brief description}
Time: {HH:MM} UTC
Impact: {Users/Services affected}
Status: Investigating
```

---

## Incident Response Template

### Phase 1: Triage (First 5 minutes)

**Goal:** Confirm incident, assess severity, gather initial context

```
[ ] 1. Verify incident is real (not false alert)
       - Confirm service down: curl -v http://service:port/health
       - Check if isolated or widespread
       - Note exact time symptoms started

[ ] 2. Assess severity (P1, P2, P3, P4)
       - Is customer-facing? → Higher severity
       - Is production data at risk? → Higher severity
       - Can we degrade gracefully? → Lower severity

[ ] 3. Gather context
       - Check recent deployments, config changes
       - Review logs (last 10 minutes)
       - Check metrics (CPU, disk, VRAM at incident time)

[ ] 4. Document in incident tracker
       - Title, severity, impact, context
       - Start timestamp

[ ] 5. Send initial Telegram notification
       - "🚨 Incident declared: [P1] {Service} down (investigating)"
```

---

### Phase 2: Mitigation (First 15 minutes)

**Goal:** Reduce blast radius, restore partial functionality if possible

```
[ ] 1. Quick fixes (try in order)
       - Restart container: docker restart {service}
       - Reload config: systemctl reload {service}
       - Clear cache: redis-cli FLUSHALL (if applicable)
       - Roll back last change (if recently deployed)

[ ] 2. If service still down
       - Fail over to backup (if available)
       - Disable non-critical features
       - Switch to degraded mode (read-only, basic functionality)

[ ] 3. Escalate to humano if needed
       - "Still down after restart, escalating to humano"
       - Provide logs, diagnostics

[ ] 4. Update Telegram every 5 minutes
       - Status, ETA for resolution, what we're trying
```

---

### Phase 3: Root Cause Analysis (First 30 minutes)

**Goal:** Understand what broke and fix it

```
[ ] 1. Check logs for error patterns
       docker logs {service} --tail 100 | grep -i error
       tail -f ~/.aurelia/logs/{service}.log

[ ] 2. Check resource exhaustion
       - VRAM: nvidia-smi (if GPU service)
       - Disk: df -h /
       - CPU: top -n1 | head -20
       - Connections: netstat -an | grep ESTABLISHED | wc -l

[ ] 3. Check dependencies
       - Postgres: docker exec postgres-{db} pg_isready
       - Qdrant: curl -s http://localhost:6333/health
       - Network: ping cloudflare.com

[ ] 4. Review recent changes
       git log --oneline -10
       git diff HEAD~5..HEAD -- {affected_files}

[ ] 5. Formulate fix hypothesis
       - "API timeout because DB query slow due to missing index"
       - "OOM killed because LLM inference memory leak"
       - "Network unreachable due to Cloudflare Tunnel down"
```

---

## Specific Incident Runbooks

### Runbook: Service Container Down

**Symptoms:** `docker ps` shows service in "Exited" state

**Detection:** health-check.sh --mode quick fails, P1 alert

**Resolution (in order):**

```bash
# 1. Check logs for crash reason
docker logs aurelia --tail 100 | tail -20

# 2. Restart container
docker restart aurelia
sleep 5

# 3. Verify it's running
docker ps | grep aurelia

# 4. Check health
curl -s http://localhost:5000/health

# If still failing:

# 5. Check disk space (might be full)
df -h / | tail -1

# 6. Check VRAM (might be exhausted)
nvidia-smi

# 7. Force rebuild (nuclear option)
docker-compose -f docker-compose.prod.yml up -d aurelia

# 8. Run diagnostics
docker logs aurelia --tail 50
docker inspect aurelia | jq '.State'
```

**Escalation:** If container crashes immediately after restart → humano review needed

---

### Runbook: GPU Out of Memory (OOM)

**Symptoms:** `nvidia-smi` shows 100% VRAM used, inference slow/failing

**Detection:** health-check.sh detects >95% VRAM, P2 alert

**Resolution (in order):**

```bash
# 1. Check current VRAM usage
nvidia-smi

# 2. Identify largest process
ps aux --sort=-%mem | head -10

# 3. Kill non-essential LLM processes (lowest priority first)
pkill -f "qwen3.5:9b" --oldest-first

# 4. Verify VRAM freed
sleep 5 && nvidia-smi

# 5. Monitor for memory leak
watch -n 5 nvidia-smi

# If memory keeps growing:

# 6. Restart LLM service
docker restart litellm

# 7. Review LLM logs for memory leak
docker logs litellm --tail 100 | grep -i memory

# If still failing:

# 8. Reduce batch size or max_tokens
# (Edit .env or Docker config)
docker-compose -f docker-compose.prod.yml up -d litellm
```

**Escalation:** If memory usage doesn't decrease → possible memory leak, requires code review

---

### Runbook: Disk Full

**Symptoms:** `df -h /` shows >95% used, writes failing

**Detection:** health-check.sh detects >95%, P1 alert

**Resolution (in order):**

```bash
# 1. Identify large files
du -sh /* | sort -rh | head -10

# 2. Check if backup dir is oversized
du -sh ~/.aurelia/backups/

# 3. Auto-cleanup old backups (>30d)
find ~/.aurelia/backups/ -name "*.sql.gz" -mtime +30 -exec rm -v {} \;

# 4. Check logs aren't too large
du -sh ~/.aurelia/logs/
find ~/.aurelia/logs/ -name "*.log" -mtime +30 -delete

# 5. Check Docker disk usage
docker system df

# 6. Cleanup unused Docker images/volumes
docker system prune -a --volumes

# 7. Check ZFS snapshots aren't consuming space
zfs list -t snapshot -o name,used

# If still full:

# 8. Archive old data to /srv/backups/
find ~/.aurelia/backups/ -name "*.sql.gz" -mtime +90 \
    -exec mv {} /srv/backups/archive/ \;

# 9. Contact humano for manual cleanup
```

**Escalation:** If still >80% after cleanup → manual intervention needed

---

### Runbook: Database Unreachable (PostgreSQL)

**Symptoms:** `docker exec postgres-n8n pg_isready` returns non-0

**Detection:** health-check.sh --mode services fails, P1 alert

**Resolution (in order):**

```bash
# 1. Check if container is running
docker ps | grep postgres-n8n

# 2. Restart container
docker restart postgres-n8n
sleep 10

# 3. Check logs
docker logs postgres-n8n --tail 50

# 4. Test connectivity
docker exec postgres-n8n psql -U root -d n8n -c "SELECT 1"

# 5. Check disk space (might be full)
docker exec postgres-n8n df -h /var/lib/postgresql

# If container is running but won't accept connections:

# 6. Check WAL files (might be corrupted)
docker exec postgres-n8n ls -lah /var/lib/postgresql/data/pg_wal/ | wc -l

# 7. Perform recovery
docker exec postgres-n8n pg_ctl status -D /var/lib/postgresql/data

# If recovery fails:

# 8. Restore from backup
LATEST=$(ls -t ~/.aurelia/backups/n8n-*.sql.gz | head -1)
gunzip -c "$LATEST" | docker exec -i postgres-n8n psql -U root
```

**Escalation:** If recovery fails → humano review of backup integrity

---

### Runbook: Cloudflare Tunnel Down

**Symptoms:** Domain unreachable (DNS resolves but connection hangs)

**Detection:** External health check fails, Telegram alert

**Resolution:**

```bash
# 1. Check tunnel status
cloudflared tunnel info

# 2. Check credentials are valid
cloudflared tunnel list

# 3. Restart tunnel
docker restart cloudflare-tunnel
sleep 10

# 4. Verify it's connected
cloudflared tunnel info

# 5. Check firewall rules
sudo ufw status | grep cloudflare  # or Tailscale port 7681

# 6. Manual tunnel test (if available)
cloudflared tunnel run --name {tunnel-name}

# 7. Check Cloudflare dashboard for issues
# (Log in to Cloudflare, check tunnel status)

# If still down:

# 8. Failover to Tailscale (if applicable)
# Users can access via Tailscale IP instead of domain
```

**Escalation:** If tunnel is down and can't reconnect → requires manual Cloudflare credential check

---

### Runbook: Credential Compromise

**Symptoms:** Unusual API activity, unauthorized requests, secret-audit.sh finds plaintext

**Detection:** Automated secret-audit.sh, manual discovery

**Response (URGENT - P1, <5 min):**

```bash
# IMMEDIATE (Do within 1 minute):

# 1. Revoke exposed credentials
#    - Postgres password change: ALTER USER root WITH PASSWORD 'new-password'
#    - API keys: Delete old, generate new in provider dashboard
#    - OAuth tokens: Revoke in application settings

# 2. Remove plaintext credentials
rm -f ~/Desktop/rascunho-s.txt
rm -f ~/.env-old
grep -r "password=" ~/ --include="*.txt" --include="*.sh" | xargs -I{} rm {}

# 3. Notify affected services
# Update password in KeePassXC, distribute to services

# WITHIN 30 MINUTES:

# 4. Audit recent access logs
journalctl --since "2 hours ago" | grep -i password
cat ~/.aurelia/logs/*.log | grep -i unauthorized

# 5. Check for unauthorized changes
git log --oneline -20
git diff HEAD~20..HEAD -- config files

# 6. Document in incident report
cat > /tmp/breach-report-$(date +%s).txt <<EOF
- Compromised credential: {type}
- Exposed at: {location}
- Revoked at: {time}
- Services affected: {list}
- Mitigation: {actions taken}
EOF

# WITHIN 24 HOURS:

# 7. Full audit trail review
# Check for unauthorized API calls, database changes

# 8. File ADR amendment
# Document how this credential leaked (e.g., plaintext in file)
# Add preventive measures to SECRETS_GOVERNANCE
```

**Escalation:** Immediately notify humano, may require CISO review

---

## Post-Incident (Blameless Postmortem)

**Within 24 hours of resolution:**

```markdown
## Incident Report

**Title:** [P{1-4}] {Service}: {Description}
**Date:** {YYYY-MM-DD HH:MM UTC}
**Duration:** {minutes}
**Impact:** {Downtime, data loss, users affected}

### Timeline

| Time | Event |
|------|-------|
| HH:MM | Alert triggered |
| HH:MM | Incident declared |
| HH:MM | Root cause identified |
| HH:MM | Fix applied |
| HH:MM | Service restored |
| HH:MM | All systems green |

### Root Cause

{What actually happened, not who caused it}
- Contributing factors
- Cascading failures
- System limitations exposed

### Resolution

{What we did to fix it}
1. Step 1
2. Step 2
...

### Prevention

{How to prevent next time}
- Code changes needed
- Process improvements
- Monitoring enhancements
- Documentation updates

### Action Items

- [ ] Item 1 (assigned to X, due date Y)
- [ ] Item 2 (assigned to X, due date Y)

### Lessons Learned

{What we learned, technology decisions to revisit}
```

---

## Escalation Matrix

| Severity | Responder | Escalate To | Time |
|----------|-----------|-------------|------|
| **P1** | codex | humano | Immediately (5 min) |
| **P2** | codex | humano | 30 min if not resolved |
| **P3** | codex | - | 4 hours if not resolved |
| **P4** | codex | - | Next business day |

**Contact:** Telegram bot (instant), email (backup)

---

## Communication Template

### Initial Alert (T+0)
```
🚨 [P1] {Service} DOWN
Status: Investigating
ETA: 15 minutes
Details: {Brief description}
```

### Status Update (Every 5 minutes if P1, every 15 min if P2)
```
⚠️ [P1] {Service} — {Status}
What we tried: {Actions taken}
Next step: {Plan}
ETA: {Time}
```

### Resolution (When fixed)
```
✅ [P1] {Service} RESOLVED
Duration: {X minutes}
Root cause: {Brief explanation}
Postmortem: {Link to report}
```

---

## Validation

Test incident response quarterly:

```bash
# Q1: Simulate container crash
docker kill aurelia  # Trigger auto-restart
# Verify health-check alerts within 5 min
# Verify container restarts automatically

# Q2: Simulate disk full
dd if=/dev/zero of=/tmp/filltest bs=1M count=10000  # Simulate full disk
# Verify cleanup script runs
# Verify alert triggered

# Q3: Simulate database down
docker stop postgres-n8n
# Verify health-check alerts
# Verify recovery procedure works

# Q4: Simulate credential compromise
echo "password=secret123" > ~/test-plaintext.txt
bash scripts/secret-audit.sh  # Should detect
# Verify alert sent
```

---

## Links

- [Health Checks](./operational-governance-health-checks.md)
- [Backup Verification](./operational-governance-backup-verification.md)
- [NETWORK_MAP.md](/srv/ops/ai-governance/NETWORK_MAP.md) (Ports, services)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
