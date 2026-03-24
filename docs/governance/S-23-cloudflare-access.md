# S-23: Cloudflare Access — Cloud-Only Configuration

**Status:** Cloud-only config. No code implementation required in this repository.

## Overview

Cloudflare Access protects the Aurelia dashboard endpoint (`https://aurelia.kocar.dev` or equivalent)
behind Cloudflare Zero Trust, requiring authentication before traffic reaches the origin.

## Setup Steps (Manual, Cloudflare Dashboard)

1. **Create an Application** in Cloudflare Zero Trust > Access > Applications.
   - Type: Self-hosted
   - Application domain: `aurelia.yourdomain.com`
   - Session duration: 24h (recommended)

2. **Create a Policy** tied to the application:
   - Policy name: `aurelia-owners`
   - Action: Allow
   - Include rule: Emails → `your@email.com`

3. **Add a Service Token** (optional, for API access):
   - Zero Trust > Access > Service Auth > Service Tokens
   - Use `CF-Access-Client-Id` + `CF-Access-Client-Secret` headers in automated requests.

4. **DNS & Tunnel** (if using Cloudflare Tunnel):
   ```bash
   cloudflared tunnel create aurelia
   cloudflared tunnel route dns aurelia aurelia.yourdomain.com
   ```
   Configure `~/.cloudflared/config.yml`:
   ```yaml
   tunnel: <tunnel-id>
   credentials-file: /home/user/.cloudflared/<tunnel-id>.json
   ingress:
     - hostname: aurelia.yourdomain.com
       service: http://localhost:3334
     - service: http_status:404
   ```

5. **Origin validation** (optional hardening):
   - In the dashboard Go backend, verify `Cf-Access-Jwt-Assertion` header against Cloudflare public keys.
   - Library: `github.com/cloudflare/cloudflare-go` or manual JWT validation.

## No Backend Code Needed

The protection is enforced at the Cloudflare edge. The Go backend at `localhost:3334` does not need
to implement any authentication logic for Cloudflare Access to work.

If origin validation is desired in the future, add middleware in `internal/dashboard/dashboard.go`
to verify the CF-Access JWT assertion before serving routes.

## References

- https://developers.cloudflare.com/cloudflare-one/applications/configure-apps/self-hosted-apps/
- https://developers.cloudflare.com/cloudflare-one/identity/service-tokens/
