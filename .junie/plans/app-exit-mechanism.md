---
sessionId: session-260923-171140-1irt
---

# Requirements

### Overview & Goals
Provide a simple, reliable, and foolproof mechanism to exit the application when running in kiosk mode on arcade cabinets or development environments. The solution must work without relying on problematic Wayland/OS-level global shortcut portals.

### Scope
- **In Scope:**
  - A lightweight Tauri command (`exit_app`) that terminates the application cleanly (`app.exit(0)`).
  - A global keyboard listener inside the webview for `Ctrl+Shift+Delete`.
  - A dedicated "Quit App" button in the Back Office header (`/backoffice`) for direct UI access.
  - Removal of the unnecessary `tauri-plugin-global-shortcut` dependency.
- **Out of Scope:**
  - System-level background hotkeys intercepted when the application is unfocused or minimized.
  - Operating system power/shutdown integration.

### User Stories
- As an arcade cabinet operator, I want to press `Ctrl+Shift+Delete` while inside the app so that I can exit the kiosk interface cleanly.
- As an administrator in the Back Office, I want a "Quit App" button so that I can close the application using a mouse, touch, or navigation stick without needing a physical keyboard attached.
- As a developer, I want the exit mechanism to be lightweight and zero-configuration across Linux, macOS, and dev mode.

### Functional Requirements
- Pressing `Ctrl+Shift+Delete` anywhere inside the application window triggers `exit_app`.
- Clicking the "Quit App" button in the Back Office interface triggers `exit_app`.
- In standard browser development mode (outside Tauri), invoking exit gracefully falls back without throwing runtime errors.

# Technical Design

### Current Implementation
- `tauri/Cargo.toml` and `tauri/src/main.rs` previously attempted to use `tauri-plugin-global-shortcut`. On Linux (Wayland / WebKitGTK / Bazzite OS), OS-level global shortcut interception often fails silently due to missing desktop portal permissions or compositor restrictions.
- The frontend (`front-end/src/main.ts`) already captures window-level keyboard navigation but has no exit command binding.
- The Back Office header (`front-end/src/pages/BackOfficePage.vue`) provides navigation to Settings and Game creation, but lacks an exit action.

### Key Decisions
- **In-Webview Keydown Handler vs. OS-Level Global Shortcut Plugin:**
  - *Decision:* Use standard webview `window.addEventListener('keydown')` and invoke a Tauri command.
  - *Rationale:* The kiosk runs full-screen and has focus; capturing the keypress inside the webview is 100% reliable, cross-platform, and avoids Wayland permission roadblocks and extra plugin overhead.
- **Dual Access (Keyboard Shortcut + UI Button):**
  - *Decision:* Support both `Ctrl+Shift+Delete` and an "Quit App" button in the Back Office.
  - *Rationale:* Maximizes utility for both keyboard users and cabinet/mouse operators.
- **Tauri IPC Command (`exit_app`):**
  - *Decision:* Use a native `#[tauri::command]` in `main.rs` registered via `tauri::generate_handler![exit_app]`.
  - *Rationale:* Cleanest and standard Tauri v2 approach, requiring no extra crates.

### Proposed Changes

#### Tauri (`tauri/`)
1. `tauri/Cargo.toml`:
   - Remove `tauri-plugin-global-shortcut`.
2. `tauri/src/main.rs`:
   - Add `#[tauri::command] fn exit_app(app: tauri::AppHandle) { app.exit(0); }`.
   - Register the command with `.invoke_handler(tauri::generate_handler![exit_app])`.
   - Remove shortcut plugin initialization and imports.

#### Frontend (`front-end/`)
1. `front-end/src/api/client.ts` (or `front-end/src/utils/app.ts`):
   - Add `exitApp()` helper:
     ```ts
     export async function exitApp(): Promise<void> {
       if (typeof window !== 'undefined' && '__TAURI_INTERNALS__' in window) {
         // @ts-expect-error Tauri internal invoke
         await window.__TAURI_INTERNALS__.invoke('exit_app')
       } else if (typeof window !== 'undefined') {
         window.close()
       }
     }
     ```
2. `front-end/src/main.ts`:
   - Listen for `Ctrl+Shift+Delete` in the global `keydown` event listener and invoke `exitApp()`.
3. `front-end/src/pages/BackOfficePage.vue`:
   - Add a `<button class="quit-btn" @click="exitApp">Quit App</button>` to the header next to Settings and Add Game.

### Architecture Diagram
```mermaid
graph LR
    KB[Keyboard Ctrl+Shift+Delete] -->|keydown event| FE[Frontend App]
    UI[Back Office 'Quit App' Button] -->|click event| FE
    FE -->|invoke 'exit_app'| IPC[Tauri IPC]
    IPC -->|app.exit(0)| App[Tauri Backend Process]
```

### File Structure
- `tauri/Cargo.toml` (modify)
- `tauri/src/main.rs` (modify)
- `front-end/src/api/client.ts` (modify)
- `front-end/src/main.ts` (modify)
- `front-end/src/pages/BackOfficePage.vue` (modify)

# Testing

### Validation Approach
Verify exit command functionality at both unit level (Rust and TypeScript/Vue) and integration level.

### Key Scenarios
1. **Tauri Command Registration:**
   - Verify `exit_app` handler is compiled and registered in `tauri/src/main.rs`.
   - Verify `cargo test --manifest-path tauri/Cargo.toml` passes.
2. **Back Office UI Button:**
   - Mount `BackOfficePage.vue` with Vitest and verify the Quit button is rendered.
   - Trigger click on the Quit button and verify that `exitApp()` is invoked.
3. **Keyboard Shortcut:**
   - Dispatch `Ctrl+Shift+Delete` keydown event in frontend tests and verify `exitApp()` is called.
   - Verify other keys (like arrow keys and standard typing) are not blocked or misrouted.
4. **Browser Fallback:**
   - Verify calling `exitApp()` outside Tauri environment does not crash or throw unhandled exceptions.

# Delivery Steps

### ✓ Step 1: Implement Tauri exit_app command and clean up Rust dependencies
Update the Tauri backend to handle application exit via direct IPC command and remove the global shortcut plugin dependency.

- Remove `tauri-plugin-global-shortcut` from `tauri/Cargo.toml` to eliminate unneeded dependencies and portal permission issues on Linux/Wayland.
- Define a `#[tauri::command]` named `exit_app` in `tauri/src/main.rs` that calls `app.exit(0)`.
- Register `exit_app` with `tauri::generate_handler![exit_app]` in `tauri::Builder`.
- Add unit tests in `tauri/src/main.rs` to verify the Tauri application builder and command registration.

### ✓ Step 2: Add frontend exit utility and keyboard shortcut handler
Create a shared frontend utility function and global keyboard listener to trigger the application exit.

- Create an `exitApp` helper in `front-end/src/api/client.ts` (or `front-end/src/utils/app.ts`) that invokes the `exit_app` Tauri command when running under Tauri, with a fallback for browser environments.
- Add an in-window `keydown` listener in `front-end/src/main.ts` for `Ctrl+Shift+Delete` to trigger `exitApp()`.
- Add unit tests verifying the exit invocation behavior and keyboard event matching.

### ✓ Step 3: Add Quit button to Back Office header
Add a visible Exit/Quit button to the Back Office interface for mouse, touch, and UI-driven exits.

- Update `front-end/src/pages/BackOfficePage.vue` to include a prominent "Quit App" button in the header action bar.
- Wire the button click handler to `exitApp()`.
- Add style definitions for the exit button adhering to the arcade theme.
- Add component tests to verify that clicking the button triggers the exit function.