# CLAUDE.md

This file provides on how to work with this repository.

## Project Overview

Arcade Champion powers homemade arcade machines running Bazzite OS. Three components:
1. **Pre-install scripts** — configure Bazzite and install dependencies
2. **Front-end** (`front-end/`) — kiosk-mode game selector, navigable with an arcade stick
3. **Back-end** — game library manager supporting Steam, Fightcade, and classic MAME

## Tech Stack

- **Front-end:** Vue 3 (Composition API) + TypeScript + Vite + Pinia. Tailwind CSS with semantic tokens. CSS variables for theming.
- **Back-end:** Go
- **Target platform:** Bazzite OS (Fedora Atomic-based, runs on x86_64 arcade cabinets)

## Commands

### Front-end (`front-end/`)

```sh
npm install          # install dependencies
npm run dev          # dev server with hot reload
npm run build        # production build
npm run preview      # preview production build
```

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

## Testing Philosophy

- No coverage targets
- Few but meaningful tests covering realistic scenarios
- Don't try to cover every branch or code path

## Workflow

1. Plan steps and dependencies
2. Refine for consistency and good patterns
3. Complete implementation
4. Check if the codebase can be simplified post-change
5. Update tests
