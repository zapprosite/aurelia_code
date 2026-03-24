# Post-Reboot Validation - 2026-03-19

## Deployment Snapshot

- Deployed worktree: `/home/will/aurelia-24x7`
- Branch: `feat/24x7-system-service`
- Commit: `533a380`
- Systemd unit: `/etc/systemd/system/aurelia.service`
- Service working directory: `/home/will/aurelia-24x7`

## Boot Validation Evidence

- Current boot reached a graphical session for user `will` on `2026-03-19 01:50:29 -03`.
- `GDM3` autologin is enabled in `/etc/gdm3/custom.conf` with `AutomaticLoginEnable=true` and `AutomaticLogin=will`.
- Journal evidence for this boot shows `pam_unix(gdm-autologin:session): session opened for user will(uid=1000) by will(uid=0)`.
- `ssh.service` is `enabled` and `active`, listening on `0.0.0.0:22` and `[::]:22`.
- `aurelia.service` entered `active (running)` at `2026-03-19 01:50:34 -03`.
- `http://127.0.0.1:8484/health` returned `200 OK` with `{\"status\":\"ok\"}` during validation.

## Keyring Issue And Fix

- Initial boot logs also showed `gkr-pam: no password is available for user` and `gkr-pam: couldn't unlock the login keyring.` under `gdm-autologin`.
- The login keyring contained 6 items used by GitHub and Chromium/Chrome safe storage.
- A backup of the original keyring files was stored at:
  - `/home/will/.local/share/keyrings/backup-antigravity-20260319T015938`
  - `/home/will/.local/share/keyrings/backup-antigravity-files-20260319T020104`
- The keyring was migrated to a no-password `login.keyring`, preserving all 6 items.
- `default` and `login` now both resolve to `/org/freedesktop/secrets/collection/login` after restarting `gnome-keyring-daemon.service`.
- Post-fix validation showed the standard Secret Service `Unlock` call returned `prompt=/`, confirming unlock without a password dialog.

## Follow-Up

- Service health remained intact after the keyring migration.
- The next reboot should be used to confirm the GNOME password dialog no longer appears during autologin.

## Reversal - 2026-03-19 02:12 -03

- On user request, the no-password keyring change was reverted before the next reboot.
- Restored files from `/home/will/.local/share/keyrings/backup-antigravity-20260319T015938`.
- Created a rollback safety copy at `/home/will/.local/share/keyrings/backup-antigravity-restore-20260319T021237`.
- Removed the temporary `default` alias file created during the no-password migration.
- After restart, `login.keyring` is back to the original password-protected GNOME keyring format.
- Validation after the restore showed:
  - `Locked=true` before unlock
  - `Unlock` returned `prompt=/org/freedesktop/secrets/prompt/u1`
  - `Locked=true` remained after the non-interactive unlock attempt
- This confirms the password prompt behavior is back for the next boot.

## Reboot Revalidation And Safe Reapply - 2026-03-19 02:20 -03

- Current reboot validation confirmed `gdm-autologin` succeeded for user `will` at `2026-03-19 02:16:15 -03`.
- `ssh.service` was still `active` and listening on port `22`.
- `aurelia.service` was `active (running)` since `2026-03-19 02:16:20 -03`.
- `http://127.0.0.1:8484/health` returned `200 OK` during this run.
- The first failure in this boot remained the keyring path:
  - `gkr-pam: no password is available for user`
  - `gkr-pam: couldn't unlock the login keyring`
- Root cause remained unchanged: GDM autologin does not provide a login password to PAM, so a password-protected GNOME login keyring cannot unlock automatically during session start.
- Applied the smallest known-good correction by restoring the previously validated no-password keyring state from:
  - `/home/will/.local/share/keyrings/backup-antigravity-restore-20260319T021237/login.keyring.before-restore`
  - `/home/will/.local/share/keyrings/backup-antigravity-restore-20260319T021237/default.before-restore`
- Created a fresh rollback backup before reapplying:
  - `/home/will/.local/share/keyrings/backup-antigravity-reapply-20260319T022047`
- Post-fix validation after restarting `gnome-keyring-daemon.service` showed:
  - `login.keyring` is now the ASCII no-password format again
  - `default` and `login` aliases both resolve to `/org/freedesktop/secrets/collection/login`
  - `Unlock` returned `prompt=/`
  - A forced `Lock` followed by `Unlock` also returned `prompt=/`, with `Locked=true` then `Locked=false`
- This re-establishes the non-interactive keyring behavior required for autologin boots.

## Live Validation After User Reboot Report - 2026-03-19 02:50 -03

- Deploy target still matches the expected runtime:
  - Worktree: `/home/will/aurelia-24x7`
  - Branch: `feat/24x7-system-service`
  - Commit: `533a380`
- Current boot started at `2026-03-19 02:41 -03`.
- Active graphical session for `will` is session `2` with:
  - `Service=gdm-autologin`
  - `Type=x11`
  - `State=active`
- `aurelia.service` is loaded from `/etc/systemd/system/aurelia.service`, `enabled`, and `active (running)` since `2026-03-19 02:41:46 -03`.
- The unit is configured with:
  - `User=will`
  - `WorkingDirectory=/home/will/aurelia-24x7`
  - `WantedBy=multi-user.target`
- `systemd-analyze critical-chain aurelia.service` shows it starting after `network-online.target`, which is consistent with boot-time autonomous startup.
- `http://127.0.0.1:8484/health` returned `HTTP/1.1 200 OK` with `{\"status\":\"ok\"}` during this validation.

## Keyring State In This Boot - 2026-03-19 02:50 -03

- Boot journal still contains `gkr-pam: couldn't unlock the login keyring.` under `gdm-autologin`.
- Despite that PAM line, the live Secret Service state is currently healthy and non-interactive:
  - `org.freedesktop.secrets` is available on the user session bus
  - alias `login` resolves to `/org/freedesktop/secrets/collection/login`
  - collection `Locked=false`
- No active prompt helper process was found (`gcr-prompter`, `pinentry`, `seahorse` absent).
- No journal evidence of a keyring unlock dialog or prompt helper launch was found in this boot.
- `gnome-keyring-daemon` is running for `will`, and later journal entries show item registration in the `login` collection without any recorded prompt.
- Conclusion for this boot: there is strong evidence that no GUI password prompt was needed after login and that Secret Service access remained usable when tools registered secrets.

## SSH Readiness Caveat - 2026-03-19 02:51 -03

- `ssh.service` is `active`, and `ss -ltnp` confirms listeners on `0.0.0.0:22` and `[::]:22`.
- Local non-interactive SSH self-test to `localhost` failed with `Permission denied (publickey,password)`.
- Root cause on this host: `/home/will/.ssh/authorized_keys` exists but is empty.
- Therefore:
  - port `22` is up and reachable locally
  - SSH daemon recovery after reboot is validated
  - passwordless SSH reconnection is **not** proven from this host alone and is not currently guaranteed by the local `will` account configuration unless an external client key is already authorized elsewhere
