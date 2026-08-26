# Rummy Client — Phaser 3

Phaser 3 web client for Romanian Tile Rummy (Nakama authoritative, Go backend). This is the **Day 24+ optional client** per `AGENTS.md:812` and `docs/daily-log.md` Phase 12 — a minimal, maintainable test client that proves the protocol without copying Remi Online branding.

**Current phase:** `Phase 1 — Foundation` through Day 5 `Local setup documentation` (scaffolding, Vite+TS+Phaser, Nakama JS `2.8.0`, lint/format `Day 4`, and this `Quick Start` with prerequisites, `VITE_NAKAMA_*` env, Docker startup, Console access, test execution). See `client/docs/roadmap.md` for the full Handmade Hero incremental plan (one small vertical slice per day).

## Prerequisites

- **Node** `>=20` (`node --version` `20.x` or `22.x` LTS; `npm` `>=10`).
- **Docker** `>=29` + `compose v5` — same as backend (`docker --version` `29.x`, `docker compose version` `v5.x`), tested on `darwin 25.6.0 arm64`.
- **Go backend** running via `docker compose up --build -d` from repo root (see `../README.md` `Quick Start`).
- No local `psql` needed — Postgres at `5433` via Docker.

## Environment Variables

Copy `client/.env.example` to `client/.env` (optional — `vite.config.ts` and `src/net/nakama.ts` have defaults via `import.meta.env?.VITE_* ?? "127.0.0.1:7350"`):

```bash
cd client
cp .env.example .env   # optional, defaults are sane for local dev
cat .env.example
# VITE_NAKAMA_HOST=127.0.0.1
# VITE_NAKAMA_PORT=7350
# VITE_NAKAMA_KEY=defaultkey
# VITE_NAKAMA_USE_SSL=false
# VITE_DEV_PORT=5173
```

`VITE_NAKAMA_KEY` must be `defaultkey` for local Nakama (not `defaulthttpkey` — see `../README.md` `Troubleshooting` `Server key invalid`). Never commit `client/.env` (already in `client/.gitignore`).

## Docker Startup (backend)

From **repo root** (not `client/`):

```bash
# Build Go plugin and start DB + Nakama (once, ~10s for DB init)
docker compose up --build -d

# Verify
docker compose ps                    # both healthy
docker compose logs nakama --tail=30 | grep -E "Rummy|modules|Registered"
# expect: Found runtime modules count 1 [rummy_backend.so]
# expect: Registered Go runtime RPC function invocation id health/version
# expect: Registered Go runtime Match creation function invocation name rummy
make smoke   # from repo root — SMOKE PASSED (pg_isready, healthcheck, InitModule, console 200, RPC health)
```

Compose project is `rummy_backend` (`compose.yml:1` `name:`) to isolate from `tinybot`. Host ports `5433→5432` for Postgres, `7350→7350` API, `7351→7351` console.

## Nakama Console Access

- **URL:** `http://127.0.0.1:7351` — user `admin` password `password` (local `nakama/data/local.yml` `DEBUG` only).
- **Health:** `docker inspect rummy_nakama --format '{{json .State.Health}}' | python3 -m json.tool` and `docker exec rummy_nakama /nakama/nakama healthcheck`.
- **DB:** `docker compose exec postgres psql -U postgres -d nakama -c "SELECT * FROM migrate;"` or `make db-shell` from repo root.

## Test Execution (client)

```bash
cd client
npm install          # once (120 packages, 4s) — creates node_modules/ (ignored)
npm run lint         # eslint src --ext .ts (Day 4, .eslintrc.json)
npm run fmt          # prettier --write src (Day 4, .prettierrc)
npm run typecheck    # tsc --noEmit (Day 4, tsconfig.json strict)
npm run build        # tsc --noEmit && vite build (Day 4, outputs dist/)
npm run dev          # Vite dev server http://127.0.0.1:5173 --host 127.0.0.1 --port 5173
# open http://127.0.0.1:5173
# check DevTools console: "Phaser 3 Rummy — Day 2 Vite + TypeScript + Phaser scaffold"
```

Or via root `Makefile` aliases (from repo root):

```bash
make client-lint       # cd client && npm run lint
make client-typecheck  # cd client && npm run typecheck
make client-build      # cd client && npm run build
make client-fmt        # cd client && npm run fmt
```

All `make client-*` are `Day 4` deliverables and must be `0` before next day.

## Quick Start (play)

```bash
# 1. Backend up (from repo root)
docker compose up --build -d && make smoke

# 2. Client dev (from client/)
cd client && npm install && npm run dev
# open http://127.0.0.1:5173

# 3. Play (2 browsers or 2 tabs at http://127.0.0.1:5173)
# - Alice: click "Create Match" (host, Seat 0, HasOpened false)
# - Bob: click "Join Match" (Seat 1)
# - Both see PrivateView.OwnRack vs PublicView counts (redaction per docs/protocol.md)
# - Use Draw / Discard / Meld / Extend / Pickup / Replace per docs/state-machine.md
```

**Backend API for client:** `http://127.0.0.1:7350` (device auth `defaultkey` via `src/net/nakama.ts` `Client("127.0.0.1","7350","defaultkey")`), `ws://127.0.0.1:7350/ws` (`Client.createSocket`), `http://127.0.0.1:7351` console `admin/password`.

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

Per `client/docs/roadmap.md` Phase 1 Day 5 done — `Quick Start` now covers prerequisites, `VITE_NAKAMA_*` env, `docker compose up --build -d`, Nakama Console `http://127.0.0.1:7351` `admin/password`, and `npm run lint`/`fmt`/`typecheck`/`build`/`dev`. Next is **Day 6 — Smoke test** (`client/docs/roadmap.md:15`): add `scripts/client-smoke.sh` or `make client-smoke` that verifies `npm run dev` serves `http://127.0.0.1:5173` `200`, `Preload` `complete`, and `nakama.ts` `Client` can `authenticateDevice` against `docker compose up` backend (like `make smoke` for Go).

---

*Stack: Phaser 3.80+, Vite 5, @heroiclabs/nakama-js 2.9+, TypeScript 5, Go backend `rummy_backend:local` at `127.0.0.1:7350`.*
