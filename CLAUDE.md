# CLAUDE.md

This file provides guidance on how to work with this repository.

## Project Overview

Arcade Champion powers homemade arcade machines running Bazzite OS. Four components:
1. **Front-end** (`front-end/`) — kiosk-mode game selector, navigable with arcade stick/gamepad
2. **Back-end** (`back-end/`) — Go game library manager supporting Steam, Fightcade, and MAME
3. **Shell** (`tauri/`) — Tauri v2 app: fullscreen, undecorated window that loads the built front-end
4. **Boot/session** (`boot/`, `KIOSK_SETUP.md`) — Plymouth theme, initramfs trimming, and the KWin kiosk session

## Architecture

```
KWin (Wayland + XWayland) ── app-session.sh ─┬─ steam -silent
                                             ├─ back-end  (HTTP :8080)
                                             └─ Tauri shell (WebKitGTK) ── front-end (Vue)
                                                        front-end ──/api, /images──▶ back-end
                                                        back-end ──launches──▶ Steam / Fightcade / MAME
```

- The back-end is a standalone process (not a Tauri sidecar), started by the session script before the shell. It inherits the compositor's display env so launched games do too.
- Data lives in `<UserConfigDir>/arcade-champion/`: SQLite `arcade-champion.db` (`games`, `settings` tables; secrets such as Fightcade credentials are encrypted), `arcade-champion.log` (truncated at each start), and cover/banner images served at `/images/`.
- **Front-end routes:** `/` home (featured, recently played, all games), `/playing`, `/backoffice` (game CRUD, image cropper, platform search), `/backoffice/add`, `/backoffice/edit/:id`, `/backoffice/settings` (platform paths/credentials).
- **Launch flow:** front-end `POST /api/launch` → `platform.Get(name)` → `Platform.Launch`. Steam uses `steam -applaunch <id>` on Linux, MAME runs the configured command, Fightcade has its own session/matchmaking/WebSocket client (`platform/fightcade/`).
- **Exit:** `Ctrl+Shift+Delete` or the Back Office "Quit App" button invokes the Tauri `exit_app` command (in-webview handler, since global shortcuts are unreliable on Wayland). See `.junie/plans/app-exit-mechanism.md`.
- **Tauri quirks:** window starts hidden and is shown on page load (so the front-end can hide the cursor first); `__NV_DISABLE_EXPLICIT_SYNC=1` is set on Linux to avoid a WebKitGTK/NVIDIA Wayland crash; the webview runs with `--disable-web-security` (cross-origin calls to the back-end, which sends permissive CORS headers).

## Kiosk Configuration (Bazzite)

Full guide: `KIOSK_SETUP.md`. Summary:
- Minimal KWin session, not a full Plasma desktop. Three scripts in `/var/home/arcade/`: `session-launcher.sh` (dbus-run-session) → `session-launcher-inner.sh` (`kwin_wayland --drm --xwayland --exit-with-session=app-session.sh`) → `app-session.sh` (starts Steam silently, back-end, then the Tauri binary; cleans up on exit).
- Never start the back-end before KWin: Steam needs XWayland's `DISPLAY`/`XAUTHORITY`.
- `/usr` is immutable, so the Wayland session entry (`arcade.desktop`) and the Plymouth theme are packaged as RPMs (built in a toolbox container) and layered with `rpm-ostree install`. Then enable autologin (not relogin) for the Arcade Champion session in Plasma Login settings.
- Cursor: KWin's Hide Cursor effect (`kwriteconfig6`, 10 s inactivity) complements the front-end's `cursor: none`.
- Boot speed: Plymouth theme in `boot/themes`, RPM spec in `boot/rpmbuild`, `boot/dracut/99-zz-arcade.conf` to shrink the initramfs (hardware-specific, regenerate if hardware changes), BIOS Fast Boot, disable `NetworkManager-wait-online`, auto-hide GRUB.

## Tech Stack

- **Front-end:** Vue 3 (Composition API) + TypeScript + Vite + vue-router. Tailwind CSS v4. CSS variables for theming. WebKitGTK on Linux (via Tauri), WKWebView on macOS.
- **Back-end:** Go + gorilla/mux + SQLite (`modernc.org/sqlite`, pure Go)
- **Shell:** Tauri v2 (Rust) — kiosk window and exit command
- **Target platform:** Bazzite OS (Fedora Atomic-based, x86_64 arcade cabinets). Dev on macOS.
- **Navigation:** Gamepad API + custom spatial navigation system (`composables/navigation.ts`). Built-in virtual keyboard for text input.

## Commands

- Back-end: `cd back-end && go run .` (serves `:8080`); tests: `go test ./...`
- Front-end: `cd front-end && npm install && npm run dev` (proxy target in `.env`: `VITE_BACKEND_URL`); tests: `npm test`
- Shell: `cd tauri && cargo tauri dev` / `cargo tauri build` (runs the front-end dev/build via `tauri.conf.json`)

## API Conventions

- All back-end routes live under `/api` (e.g. `/api/search`)
- Front-end fetches `/api/*`; Vite dev proxy forwards to `localhost:8080`
- JSON keys use camelCase (use struct tags: `json:"camelCase"`)

## Design Principles

- KISS — simplest solution that works
- Security must not degrade user/dev experience (runs locally, not exposed to the internet)
- Intuitive, maintainable code over clever code
- No defensive programming — trust internal code, only validate at system boundaries
- Minimal comments — only when the *why* is non-obvious

## Front-end Conventions

- Vue Composition API (`<script setup lang="ts">`) exclusively
- Strict TypeScript — never use `any`
- CSS variables for theming; prefer Tailwind semantic tokens
- `@` alias maps to `front-end/src/`

## Back-end Conventions

- Platform interface pattern: each platform (Steam, Fightcade, MAME) implements `platform.Platform`
- Factory function `platform.Get()` returns the implementation by name
- Handlers in `handlers/`, platform logic in `platform/`

## Testing Philosophy

- No coverage targets
- Few but meaningful tests covering realistic scenarios
- Don't try to cover every branch or code path

## Workflow

**You MUST follow these steps in order. Do NOT write code before completing steps 1–2. Do NOT consider the task done before completing steps 4–5.**

1. **Plan** — state the steps, files involved, and dependencies.
2. **Refine** — check the plan for consistency with existing patterns and conventions in this file.
3. **Implement** — write the code.
4. **Simplify** — review the codebase post-change. Remove dead code, reduce duplication, flatten unnecessary abstractions.
5. **Test** — add or update tests covering the new behavior.

Always apply KISS principle. If a change may result in bloated code, warn the user and offer alternatives.