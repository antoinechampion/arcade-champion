# CLAUDE.md

This file provides guidance on how to work with this repository.

## Project Overview

Arcade Champion powers homemade arcade machines running Bazzite OS. Three components:
1. **Front-end** (`front-end/`) — kiosk-mode game selector, navigable with arcade stick/gamepad
2. **Back-end** (`back-end/`) — Go game library manager supporting Steam, Fightcade, and MAME
3. **Shell** (planned: Tauri, `tauri/`) — bundles front-end + back-end sidecar into a kiosk app with window management for game launch/refocus

## Tech Stack

- **Front-end:** Vue 3 (Composition API) + TypeScript + Vite + vue-router. Tailwind CSS v4. CSS variables for theming. WebKitGTK on Linux (via Tauri), WKWebView on macOS.
- **Back-end:** Go + gorilla/mux — runs as a Tauri sidecar
- **Shell:** Tauri v2 (Rust) — kiosk window, sidecar lifecycle, window hide/show for game focus
- **Target platform:** Bazzite OS (Fedora Atomic-based, x86_64 arcade cabinets). Dev on macOS.
- **Navigation:** Gamepad API + custom spatial navigation system (`composables/navigation.ts`). Built-in virtual keyboard for text input.

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