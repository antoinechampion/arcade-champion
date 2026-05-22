---
name: fightcade
description: >-
  Communicates with the Fightcade 2 platform — search games, join lobbies, challenge
  players, accept challenges, launch matches, training mode, and offline play. Use when
  the user mentions Fightcade, matchmaking, lobbies, challenges, netplay, GGPO, or
  wants to interact with the Fightcade back-end.
---

# Fightcade 2 integration

## Architecture overview

Fightcade 2 is a Nativefier-wrapped Electron app loading `https://web.fightcade.com/` with User-Agent `Fightcade2-OSX-v2.1.45`. Three backend services:

| Service | URL / Address | Protocol | Purpose |
|---------|---------------|----------|---------|
| **GGS lobby server** | `wss://ggs.fightcade.com/ws/` | WebSocket JSON | Auth, lobbies, matchmaking, chat |
| **REST API** | `https://www.fightcade.com/api/` | HTTP POST JSON | Replays, rankings, profiles, game info |
| **GGPO relay** | `ggpo.fightcade.com:7000-7001` | Raw UDP/TCP | NAT hole-punching, P2P netplay |

| **Status API** | `https://web.fightcade.com/fc2status/api/` | HTTP POST JSON | User online status reporting |

A local Python binary (`fcade`) handles emulator launching and NAT traversal via the `fcade://` URL scheme.

## WebSocket protocol (GGS)

All messages are JSON with a `req` field. Requests include `requestIdx` (integer) for correlating responses. Server-pushed events have no `result` field.

**Distinguish responses from events:** if `msg.result` exists → response to your request (matched by `requestIdx`). Otherwise → server-pushed event from another user's action.

### Authentication

```json
// Login with credentials
{"req": "login", "username": "...", "userpass": "...", "location": {}, "requestIdx": N}

// Login with session cookie
{"req": "autologin", "cookie": "<fcic_value>", "location": {}, "requestIdx": N}

// Response (both)
{"result": 200, "cookie": "<fcic>", "user": {"name": "...", "rank": 4, "token": "<session_token>", ...}}
```

The `token` is used by the `fcade` binary for Status API calls. The `cookie` (`fcic`) persists sessions.

### Lobby operations

```json
// Join channel
{"req": "join", "channelname": "sfiii3nr1", "status": "available", "away": false, "idx": -1, "requestIdx": N}
// Response: {result: 200, users: [...], emulator, gameid, ranked}

// Set status to "Looking to play" (REQUIRED after join — without this you appear as "Not available")
{"req": "stnoaway", "channelname": "sfiii3nr1", "requestIdx": N}

// Set status to "Not available"
{"req": "staway", "channelname": "sfiii3nr1", "requestIdx": N}

// Leave channel (fire-and-forget, requestIdx: -1)
{"req": "leave", "channelname": "sfiii3nr1", "requestIdx": -1}

// Search channels
{"req": "channels", "filter": "street fighter", "paginated": true, "page": 0, "requestIdx": N}

// Chat
{"req": "chat", "channelname": "...", "text": "...", "requestIdx": -1}
```

**Important:** Joining a channel alone does NOT make you "Looking to play". You must send `stnoaway` after joining to be challengeable. The `away` field in the join request is ignored by the server.

### Challenge → match flow

```
You                         GGS Server                    Opponent
 │── challenge ──────────────▶│──── challenge (event) ─────▶│
 │                             │◀─── accept ────────────────│
 │◀── accept (event) ─────────│                             │
 │◀── start (quark,port,delay)│──── start ─────────────────▶│
 │   Build fcade:// URL        │        Build fcade:// URL   │
 │   Open → fcade binary       │        Open → fcade binary  │
```

```json
// Challenge someone
{"req": "challenge", "username": "opponent", "channelname": "sfiii3nr1", "challengeid": 1, "ranked": true, "requestIdx": N}

// Accept incoming challenge
{"req": "accept", "username": "challenger", "channelname": "sfiii3nr1", "challengeid": 1, "ranked": true, "requestIdx": N}

// Reject / cancel
{"req": "reject", "username": "...", "channelname": "...", "challengeid": 1, "requestIdx": N}
{"req": "cancel", "username": "...", "channelname": "...", "challengeid": 1, "requestIdx": N}
```

`challengeid` is a client-generated incrementing integer per channel.

### Start event (match begins)

Pushed to both players after accept:

```json
{
  "req": "start",
  "quarkid": "1234567890-1234",
  "playerid": 0,
  "port": 6000,
  "delay": 2,
  "ranked": true,
  "token": "<session_token>",
  "user": {"name": "opponent", "rank": 3}
}
```

### Other server-pushed events

| Event `req` | Meaning |
|-------------|---------|
| `challenge` | Incoming challenge from another player |
| `accept` / `reject` | Response to your challenge |
| `start` | Match starting — contains quark parameters |
| `join` / `leave` | User entered/left the lobby |
| `staway` / `stnoaway` | User changed availability status |
| `chat` | Lobby chat message |
| `msg` | Private message |

## Launching games via `fcade://` URLs

After receiving a `start` event, build a URL and invoke it via `open` (macOS) / `xdg-open` (Linux):

```
# Proxy-based emulators (fbneo, snes9x, fc1)
fcade://served/<emulator>/<gameid>/<quarkid>.<playerid>,<port>,<delay>,<ranked>

# Native GGPO emulators (flycast, duckstation, custom)
fcade://launch/<emulator>/<gameid>/<quarkid>.<playerid>,<port>,<delay>,<ranked>/<token>

# Offline test mode
fcade://play/<emulator>/<gameid>

# Training mode (FBNeo Lua training script)
fcade://training/<emulator>/<gameid>
```

Native GGPO emulators: `flycast`, `duckstation`, `custom`. Everything else uses the proxy path (`served`).

## REST API

Single POST endpoint `https://www.fightcade.com/api/` — JSON body with `req` field:

| `req` | Key params | Purpose |
|-------|-----------|---------|
| `searchquarks` | `gameid?`, `username?`, `best`, `since`, `offset`, `limit` | Search replays |
| `getuser` | `username` | User profile |
| `searchrankings` | `gameid`, `offset`, `limit`, `byElo`, `recent` | Rankings |
| `gameinfo` | `gameid` | Game details |
| `searchevents` | `gameid?`, `offset`, `limit` | Tournaments |

## Becoming "Looking to play" (required for challenges)

Two things are needed to appear as "Looking to play" to other players:

### 1. Status API call (after authentication, before joining)

POST to the Status API to report your client as online:

```
POST https://web.fightcade.com/fc2status/api/
User-Agent: fcade
Content-Type: application/json
```

```json
{
  "req": "userstatus",
  "token": "<session_token from login>",
  "userstatus": "stcable",
  "uuid": "<persistent device UUID>",
  "guid": "<machine GUID (same as uuid on non-Windows)>",
  "huid": "<MD5 of uuid>",
  "version": "6",
  "hash": "<MD5 of salted concatenation>"
}
```

- `userstatus`: `"stwlan"` (WiFi) or `"stcable"` (wired)
- `uuid`: zlib-compressed UUID stored in `~/.fcuid`
- `hash`: `MD5("3jedoQ" + token + "qmkq0" + uuid + "dsnds" + guid + "sec or" + version + "2jden3" + userstatus + "llNjjha" + huid)`

### 2. Send `stnoaway` (after joining a channel)

After joining a channel, you must explicitly send `{"req": "stnoaway", "channelname": "..."}` to set your status to "Looking to play". Without this, the server defaults you to "Not available" regardless of the `status` and `away` fields in the join request.

## Rank system

7 levels: `["?", "E", "D", "C", "B", "A", "S"]` (indices 0–6).

## Emulators

| Key | Native GGPO | Protocol |
|-----|------------|----------|
| `fbneo` | No | UDP proxy |
| `flycast` | **Yes** | Direct GGPO + Dojo config |
| `duckstation` | **Yes** | Direct GGPO |
| `snes9x` | No | UDP proxy |
| `custom` | **Yes** | TCP proxy |
| `nulldc` | No | TCP proxy + reliable UDP |
| `fc1` (legacy) | No | UDP proxy |

For detailed protocol docs (NAT hole-punching, UDP/TCP protocols, device fingerprinting, decompiled Python sources), see [reference.md](reference.md).
