# Rummy Client — Phaser 3

Phaser 3 web client for Romanian Tile Rummy (Nakama authoritative, Go backend). This is the **Day 24+ optional client** per `AGENTS.md:812` and `docs/daily-log.md` Phase 12 — a minimal, maintainable test client that proves the protocol without copying Remi Online branding.

**Current phase:** `Phase 1 — Foundation` (scaffolding only, no gameplay yet). See `client/docs/roadmap.md` for the full Handmade Hero incremental plan (one small vertical slice per day).

## Quick Start (local)

```bash
# 1. Prerequisites: Node 20+ and Go backend running
node --version   # >=20
npm --version    # >=10
docker compose up --build -d   # from repo root, starts Nakama 7350/7351 + Postgres 5433
make smoke       # check Nakama healthy

# 2. Install and run client
cd client
npm install
npm run dev      # Vite dev server http://127.0.0.1:5173
# open http://127.0.0.1:5173

# 3. Play (2 browsers or 2 tabs)
# - Alice: click "Create Match" (host, Seat 0, HasOpened false)
# - Bob: click "Join Match" (Seat 1)
# - Both see PrivateView.OwnRack vs PublicView counts (redaction)
# - Use Draw / Discard / Meld / Extend / Pickup / Replace per docs/protocol.md
```

**Backend API:** `http://127.0.0.1:7350` (device auth `defaultkey`), `ws://127.0.0.1:7350/ws` (Nakama JS `Client` + `Socket`), `http://127.0.0.1:7351` console `admin/password`.

## Project Structure

```
client/
  docs/
    roadmap.md       # Handmade Hero incremental plan (this file's source)
  src/
    main.ts          # Phaser 3 `new Phaser.Game` with `Preload`/`TableScene`/`RackScene`
    scenes/
      Preload.ts     # loads tile sprites, joker, table, discard, stock
      TableScene.ts  # renders PublicView: TableMelds, DiscardRow, StockCount, CurrentSeat
      RackScene.ts   # renders PrivateView.OwnRack only (redaction)
    net/
      nakama.ts      # Device auth via `defaultkey`, `Client` + `Socket`, `CreateMatch`/`JoinMatch`, `SendMatchState` with Envelope {v,op,requestId,payload}
      protocol.ts    # typed wrappers for opcodes 1..9/100..199 per docs/protocol.md
    state/
      snapshot.ts    # PublicSnapshot/PrivateSnapshot types and `PublicView`/`PrivateView` projection (mirrors Go visibility.go)
      sync.ts        # reconnection: store last PrivateSnapshot per Seat, rehydrate on `OpServerState 100`
    ui/
      TurnIndicator.ts
      ErrorToast.ts  # shows OpServerError {code, message, requestId, op}
      Rack.ts        # drag/drop, click selection, sort by colour/rank
  index.html
  vite.config.ts
  package.json       # phaser 3.80+, @heroiclabs/nakama-js 2.9+, vite 5
```

Adapted from `AGENTS.md:133` but with `client/src` instead of `src/` to keep backend `internal/` separate.

## How to Add a New Command (client)

1. Add `OpClient*` to `src/net/protocol.ts` (never reuse `1..9`, keep `Version 1` per `docs/protocol.md`).
2. Add payload type and `ValidatePayload` helper (mirror Go `validator.go:12`).
3. Add `AllowedOps` check in `src/state/sync.ts` (mirror `phases.go:15`).
4. Implement handler in `src/scenes/RackScene.ts` or `TableScene.ts` (keep scenes thin, pure validation in `src/state/snapshot.ts`).
5. Wire to `Socket.sendMatchState` with `Envelope{v,op,requestId,payload}` and handle `OpServerError` `102`/`OpServerEvent` `103`/`OpServerState` `100`.
6. Add `src/ui` feedback for `bad_payload`/`not_your_turn`/`wrong_phase` and test `PrivateView` redaction.

See `docs/roadmap.md` Phase 5–8 for networking/state sync.

## Docs

- `client/docs/roadmap.md` — full incremental plan (this file's source)
- `../docs/protocol.md` — opcodes, envelope, payload schemas, snapshots, errors (shared with backend)
- `../docs/state-machine.md` — `GamePhase`/`TurnPhase` and `AllowedOps` matrix (shared)
- `../docs/testing.md` — deterministic harness and `make cli` (backend CLI) vs this Phaser client
- `../AGENTS.md` — product spec, `Day 24` optional client acceptance criteria

## Next Steps

Per `client/docs/roadmap.md` Phase 1 Day 1: `npm init` + `phaser` + `vite` + minimal `Preload` scene that loads a single tile sprite and logs `InitModule` for `make dev` health check.

---

*Stack: Phaser 3.80+, Vite 5, @heroiclabs/nakama-js 2.9+, TypeScript 5, Go backend `rummy_backend:local` at `127.0.0.1:7350`.*
