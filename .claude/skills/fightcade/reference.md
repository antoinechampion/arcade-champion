# Fightcade 2 — detailed protocol reference

This document covers the low-level protocols, NAT traversal, device fingerprinting, and emulator configuration details. For the high-level overview and WebSocket API, see [SKILL.md](SKILL.md).

## Application stack

```
┌─────────────────────────────────────────────────┐
│        Electron Shell (Nativefier 8.0.7)        │
│     Vue.js SPA @ https://web.fightcade.com/     │
│     User-Agent: Fightcade2-OSX-v2.1.45          │
└────────────────────┬────────────────────────────┘
                     │
         ┌───────────┼───────────┐
         ▼           ▼           ▼
    REST API    WebSocket    fcade:// URL
    (search,    (lobby,      scheme handler
    rankings)   matchmaking)      │
                     │            ▼
                     │    ┌────────────────┐
                     │    │ fcade (Python)  │
                     │    │ NAT punching   │
                     │    │ Emulator proxy │
                     │    └────────┬───────┘
                     │             ▼
                     │    UDP NAT hole-punch
                     │    Emulator launch
                     │
              ┌──────┴──────┐
              │  GGS Server │  wss://ggs.fightcade.com/ws/
              │  GGPO Relay │  ggpo.fightcade.com:7000-7001 (UDP)
              └─────────────┘
```

## UDP NAT hole-punching protocol

### Port map

| Port | Protocol | Purpose |
|------|----------|---------|
| 6000–6009 | UDP | GGPO netplay (emulator ↔ emulator) |
| 7000–7001 | UDP | Master server communication |
| 6004 | UDP | Default client bind for hole punching |
| 6006 | UDP | Fallback client bind |
| 26004 | UDP | Secondary fallback for restricted NAT |
| 7001 | UDP | Local proxy (127.0.0.1) for emulator |

### Registration & hole-punching flow

```
Client binds UDP socket (port 6004 or auto)

CLIENT → MASTER (UDP port 7000 or 7001):
  "<quarkid>.<side>/<emulator_port>"
  e.g. "1234567890-1234.0/7001"

MASTER → CLIENT:
  "ok <quarkid>.<side>"

CLIENT → MASTER:
  "ok"

MASTER → CLIENT (6 bytes):
  [4-byte IP][2-byte port] — peer's public address

CLIENT → PEER (repeated 8 iterations, ~4 seconds):
  "<random_float> _"

PEER → CLIENT:
  "<peer_token> <client_token>"

CLIENT → PEER:
  "<client_token> <peer_token> ok"

Direct UDP connection established.
```

### Token exchange example

```
Client A → B: "0.5432 _"
B → A: "0.2891 0.5432"
A → B: "0.5432 0.2891 ok"
```

### Restricted NAT handling

- If punch fails on port 6004–6009, tries up to 512 alternate ports
- If peer reports restricted NAT, attempts aggressive port bombing
- Falls back to `useports/<quark>` to request server relay

### UDP proxy (post hole-punch)

After hole-punching succeeds, `fcade` runs a local proxy:

```
Emulator (127.0.0.1:7001)
         ↓
fcade UDP proxy (localhost:7001 → peer's punched address)
         ↓
Peer Emulator
```

Proxy filters out hole-punch artifacts (packets containing ` ok` or ` _`).

## TCP hole-punching (NullDC / custom emulators)

Length-prefixed binary format: `[4-byte big-endian length][payload]`

```
Client MSG1 (to master):
  <quarkid>.<side>|<private_ip>:<private_port>

Master ACK:
  <public_ip>:<public_port>

Master → Peer Data:
  <peer_public>:<port>|<peer_private>:<port>
```

Both peers connect to master via TCP port 7000, exchange addresses, then establish direct TCP or fall back to relay.

## Status API (`https://web.fightcade.com/fc2status/api/`)

Used by the `fcade` binary (not the web frontend). Headers: `User-Agent: fcade`, `Content-Type: application/json`.

```json
// User status (sent periodically)
{
  "req": "userstatus",
  "token": "<session_token>",
  "userstatus": "stwlan",
  "uuid": "<device_uuid>",
  "guid": "<machine_guid>",
  "huid": "<hardware_uid>",
  "version": "6",
  "hash": "<md5_hash>"
}

// Match status
{
  "req": "quarkstatus",
  "quarkstatus": "start",
  "quarkid": "<quark_id>",
  "token": "<session_token>"
}
```

- `stwlan` = WiFi, `stcable` = wired
- UUID stored in `~/.fcuid` (zlib-compressed)
- All hardware identifiers are MD5-hashed before transmission

## Device fingerprinting (anti-cheat)

| Platform | Identifiers |
|----------|-------------|
| macOS | `IOPlatformSerialNumber`, `IOPlatformUUID`, disk serial via `ioreg` |
| Windows | `wmic csproduct get uuid`, registry `MachineGuid`, disk serial |
| Linux | `/etc/machine-id`, `lsblk` hardware serial |

## Quark format

```
<quark_id>.<side>,<port>,<delay>,<ranked>

Examples:
1234567890-1234.0,6000,2,1  → Player 1, port 6000, 2-frame delay, ranked
1234567890-1234.1,6000,2,1  → Player 2, same match
```

## Emulator configuration details

### Flycast launch flags

```
-config dojo:Enable=yes
-config network:GGPO=yes
-config network:ActAsServer=yes|no
-config network:server=<ip>
-config network:GGPOPort=<port>
-config network:GGPORemotePort=<port>
-config network:GGPODelay=<n>
-config dojo:Transmitting=yes          (spectating)
-config dojo:SpectatorIP=<server>
-config dojo:TransmitScore=yes         (ranked)
-config dojo:FirstTo=<n>               (ranked match sets)
```

### Offline modes

| Mode | FBNeo | Flycast | DuckStation |
|------|-------|---------|-------------|
| Test game (`play`) | ROM launched directly | Dojo disabled explicitly | `-portable <game>` |
| Training (`training`) | ROM + `fbneo-training-mode/fbneo-training-mode.lua` | Same as test game | Same as test game |

FBNeo training mode auto-loads savestate from `savestates/<game>_fbneo.fs` if present.

## `fcade://` URL scheme — complete list

| URL | Action |
|-----|--------|
| `fcade://userstatus/stwlan/<token>` | Report user online (WiFi) |
| `fcade://userstatus/stcable/<token>` | Report user online (wired) |
| `fcade://launch/<emu>/<game>/<quark>/<token>` | Launch match (native GGPO) |
| `fcade://served/<emu>/<game>/<quark>` | Serve match (proxy emulators) |
| `fcade://stream/<emu>/<game>/<quark>` | Spectate match |
| `fcade://play/<emu>/<game>` | Test game offline |
| `fcade://training/<emu>/<game>` | Training mode |
| `fcade://direct/<emu>/<game>/<ip>/<side>` | Direct P2P connection |
| `fcade://checkrom/<emu>/<rom>` | Verify ROM via `frm` |
| `fcade://killemu` | Kill all emulator processes |
| `fcade://autoupdate` | Run updater |

## String deobfuscation

The `fcade` Python binary uses repeating-key XOR (key = ASCII `GetAddrFromArgs`):

```python
key = [71, 101, 116, 65, 100, 100, 114, 70, 114, 111, 109, 65, 114, 103, 115]

def decode(data: bytes) -> str:
    return ''.join(chr(b ^ key[i % 15]) for i, b in enumerate(data))
```

## Static assets

| URL | Purpose |
|-----|---------|
| `https://web.fightcade.com/static/previews/<gameid>.png` | Game thumbnails |
| `https://replay.fightcade.com/<emulator>/<gameid>/<quarkid>` | Replay download |
| `https://www.fightcade.com/id/<username>` | User profile page |
| `https://www.fightcade.com/game/<gameid>` | Game info page |

