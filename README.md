# Rummy Backend — Romanian Tile Rummy (Nakama authoritative, Go)

Server-authoritative multiplayer Romanian Tile Rummy for 2–4 players, implemented as a Nakama Go runtime plugin. This is the **Go** version — migrated from TypeScript on 2026-08-25 (`55c7f3b`). See `docs/project-baseline.md` §13 for migration rationale.

The game follows Romanian Tile Rummy (106 tiles, 2 jokers, 50-point opening meld with at least one run, anticlockwise turns) per `AGENTS.md` and is being built incrementally (“Handmade Hero” vertical slices). Current phase is **Phase 11 — Deterministic simulation and minimal client**: validated runs/sets with jokers, 50-point opening, `MELD_INITIAL`/`MELD_NEW`/`EXTEND_MELD`, `DRAW_PREVIOUS`/`PICKUP`/`REPLACE_JOKER`, `ROUND_COMPLETE` win detection, centralized `PublicView`/`PrivateView` versioned `1` with reconnection, plus `TestDeterministicSimulation` harness and minimal CLI `cmd/rummy-cli` (`make cli`) showing `PrivateView` vs `PublicView` redaction (see `docs/rules-decisions.md` and `docs/daily-log.md`).

## Prerequisites

- Docker Desktop 4.x+ (`docker --version` 29.x, `docker compose version` v5.x) — tested on `darwin 25.6.0 arm64`.
- Go 1.23.x (`go version` — `go.mod` pins `1.23.5` to match `heroiclabs/nakama:3.26.0` runtime `go1.23.5`; local `1.27` works but `go vet` should pass on `1.23`).
- Git.
- No local `psql` required — Postgres runs in Docker at `5433`. Host `5432` is left for the unrelated `tinybot` stack.

## Quick Start

```bash
# 1. Clone (already done)
git clone git@github.com:gabriel-d0/rummy_backend.git
cd rummy_backend

# 2. Copy env defaults (non-secret, local only)
cp .env.example .env   # optional — compose.yml has defaults via ${VAR:-default}

# 3. Build Go plugin and start DB + Nakama
docker compose up --build -d

# 4. Watch startup (should end with "Startup done" and Go module log)
docker compose logs -f nakama --tail=50
# expected: Rummy backend InitModule — Romanian Tile Rummy v0.1.0 Go Day 3 skeleton starting
# expected: Found runtime modules count 1 [rummy_backend.so]
# expected: Registered Go runtime RPC function invocation id health/version

# 5. Verify console and API
open http://127.0.0.1:7351   # Nakama console, admin/password (local)
curl -s http://127.0.0.1:7350/ | head   # API gateway (empty 200/404 is ok, not connection refused)
```

**First-time DB init** takes ~10s for Postgres healthcheck + Nakama `migrate up`. Re-running `docker compose up -d` without `--build` reuses the existing `rummy_backend_pgdata` volume and skips migrations (`Successfully applied migration count 0`).

## Services & Ports

| Service | Container | Image / Build | Host port → Container | Purpose |
|---|---|---|---|---|
| `postgres` | `rummy_postgres` | `postgres:15-alpine` | `5433 → 5432` | Nakama DB — avoids `tinybot`’s `5432` (see `docs/project-baseline.md:6`). Uses volume `rummy_backend_pgdata`. |
| `nakama` | `rummy_nakama` | Built `rummy_backend:local` via `Dockerfile` (`nakama-pluginbuilder:3.26.0` → `nakama:3.26.0`) | `7350 → 7350` (HTTP API), `7351 → 7351` (Console), `7349 → 7349` (gRPC) | Authority, runs `main.go:12` `InitModule` → `health`/`version` RPCs and `rummy` match handler. Plugin baked as `/nakama/data/modules/rummy_backend.so`; only `local.yml` is host-mounted. |

Compose project name is `rummy_backend` (`compose.yml:1` `name:`) to isolate volumes/networks from `tinybot`.

## Development Workflow

### Rebuild after Go changes

Go plugin is **baked into the image**, not volume-mounted like the previous TS `build/index.js` flow (TS mounted `./nakama/data` and hid the image; Go failed with `plugin was built with ... pragma` until we baked). After editing `main.go` or any Go file:

```bash
docker compose up --build -d
docker compose logs nakama --tail=40 | grep -E "Rummy|modules|Registered"
# expect: Found runtime modules count 1 [rummy_backend.so]
```

**Do NOT** run `go build --buildmode=plugin` on macOS host — the `.so` would be `darwin` and fail inside `linux` container. Always build via `Dockerfile`.

### Start / Stop / Logs / Clean

```bash
docker compose up -d                # start (with existing image)
docker compose up --build -d        # start + rebuild Go plugin
docker compose ps                   # status, health (both should be healthy)
docker compose logs nakama --tail=100   # Nakama logs
docker compose logs postgres --tail=50
docker compose restart nakama        # reload after local.yml edit (no rebuild)
docker compose down                  # stop, preserve PG volume + host plugin baked in image
docker compose down -v               # stop AND wipe DB (pgdata volume) — next up will initdb
docker compose exec postgres psql -U postgres -d nakama -c "SELECT * FROM migrate;"  # DB inspect
```

### Go commands (on host, for lint/vet/test before building image)

```bash
go vet ./...          # static checks (required before commit)
go test ./...         # deterministic unit + match flow tests with seeded shuffle
go fmt ./...          # format (use go fmt, not prettier)
go mod tidy           # sync deps after editing go.mod (keeps go 1.23.5, nakama-common v1.36.0, protobuf v1.36.4)
make vet              # alias via Makefile
make test             # alias
make fmt              # alias
make check            # vet + fmt-check + test (CI baseline)
```

Protobuf is pinned to `v1.36.4` — `v1.36.6` caused `plugin was built with a different version of pragma` crash on `heroiclabs/nakama:3.26.0` (`go1.23.5`). Keep `go 1.23` line in `go.mod` aligned with builder.

### RPC smoke test (proves runtime loaded)

Nakama HTTP key is `defaultkey` (not `defaulthttpkey` — the latter gives `Server key invalid`).

```bash
UUID=$(uuidgen | tr '[:upper:]' '[:lower:]')
TOKEN=$(curl -s -X POST "http://127.0.0.1:7350/v2/account/authenticate/device?create=true" \
  -H "Content-Type: application/json" --user "defaultkey:" --data "{\"id\":\"$UUID\"}" | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")

curl -s -X POST "http://127.0.0.1:7350/v2/rpc/health" -H "Authorization: Bearer $TOKEN" --data '""' | python3 -m json.tool
# {"payload": "{\"status\":\"ok\",\"version\":\"0.1.0-go-day3-skeleton\",...}"}

curl -s -X POST "http://127.0.0.1:7350/v2/rpc/version" -H "Authorization: Bearer $TOKEN" --data '""' | python3 -m json.tool
# {"payload": "{\"runtime\":\"go\",\"phase\":\"1-day3-go\",...}"}
```

Or use Nakama Console → API Explorer → RPC → `health` with the same `$TOKEN`.

## Project Structure

```
.
├── main.go                # Nakama InitModule (Go) — registers health/version RPCs + rummy match (main.go:12)
├── go.mod / go.sum        # module github.com/gabriel-d0/rummy_backend, nakama-common v1.36.0
├── Dockerfile             # multi-stage pluginbuilder:3.26.0 → nakama:3.26.0, outputs backend.so
├── compose.yml            # name rummy_backend, postgres 5433 + nakama 7350/7351, builds Dockerfile
├── nakama/data/local.yml  # minimal DEBUG config, mounted ro
├── .env.example           # non-secret local defaults (5433, 7350/7351, admin/password)
├── Makefile               # dev helpers: vet/fmt/test/check + compose shortcuts (make help)
├── internal/
│   ├── match/             # authoritative state: RoundState (Kind, Winner), phases, turn advance, visibility, handlers (discard, draw, meld_initial, meld_new, extend_meld, draw_previous, pickup_for_meld, replace_joker, win), meld_common
│   ├── rules/
│   │   ├── tile/          # TileColour, Rank, TileInstance, Joker
│   │   ├── meld/          # Meld, ValidateRun/Set, joker ratio, Ace edge cases
│   │   └── scoring/       # ScoreTile/Run/Set, Batch, ValidateInitialBatch/HasRun/Score/Ownership
│   ├── setup/             # deck 106, seeded Rand, Shuffle, Deal, NewRoundState (15/14/stock)
│   └── protocol/          # opcodes, envelope, validator (Pickup/Replace jokerReps), errors (102 OpServerError)
├── docs/
│   ├── project-baseline.md
│   ├── rules-decisions.md # binding product decisions + TODO(product) ambiguities
│   ├── terminology.md
│   ├── protocol.md        # opcodes, envelope, payload schemas, snapshots, errors
│   ├── state-machine.md   # GamePhase/TurnPhase and allowed ops
│   ├── testing.md         # how to run deterministic tests
│   └── daily-log.md       # Handmade Hero daily slices (git log companion)
└── AGENTS.md              # product spec — source of truth, includes 24-day plan
```

## Architecture Overview

Small, explicit, server-authoritative. Match handler orchestrates; rules modules are pure.

- **`internal/match`** owns `RoundState` (`Players`, `Racks` private, `Stock`, `DiscardRow` ordered, `TableMelds` public with `Kind`, `CurrentSeat`, `GamePhase`, `TurnPhase`, `Winner`), `AdvanceTurn` anticlockwise `(current+1)%n`, `AllowedOps` matrix (`RoundComplete` none), `ValidateActivePlayer`/`ValidatePhaseOp`, visibility projection `PublicView`/`PrivateView`, and handlers `handleOpeningDiscard`, `handleDrawStock`, `handleNormalDiscard`, `handleMeldInitial`, `handleMeldNew`, `handleExtendMeld`, `handleDrawPreviousDiscard`, `handlePickupDiscardForMeld`, `handleReplaceJoker` plus shared `meld_common.go` (`buildMeldsFromPayload`, `applyMeldBatch`) and `win.go` (`checkWinAndComplete` on `rack==0` after meld/discard). Never trusts client meld validity.
- **`internal/rules/meld`** validates `Run` (same colour, consecutive, `1-2-3` and `12-13-1` ok, `13-1-2` rejected, joker reps explicit, `real>=2*joker`) and `Set` (same rank, distinct colours, 3–4 tiles, joker reps) with structured `ValidationError`.
- **`internal/rules/scoring`** scores tiles in context (`2–9:5`, `10–13:10`, Ace `1-2-3:5` vs `12-13-1:10` vs Ace-set `25`, Joker = represented tile) and validates initial batch (all tiles owned, each meld valid, ≥50, ≥1 run, no duplicate `TileId`, atomic).
- **`internal/rules/tile`** defines `TileInstance{ID,Colour,Rank,IsJoker}` with unique `TileInstanceId`.
- **`internal/setup`** creates 106-tile deck (104 numbered +2 jokers), seeded `Shuffle` (Fisher–Yates, injectable `Rand`), `Deal` (opener 15, others 14, remainder `Stock`), `NewRoundState` deterministic for tests, and `CheckTileConservation` invariant.
- **`internal/protocol`** defines `Version=1`, client `1..9` / server `100..199` opcodes, `Envelope{v,op,requestId,payload}`, `ParseEnvelope`/`ValidatePayload`/`ValidateEnvelope`, and `ErrorResponse` `OpServerError 102` with codes `bad_json`/`bad_version`/`unknown_opcode`/`bad_payload`/`not_member`/`not_your_turn`/`wrong_phase`/`invalid_tile`/`not_opened`/`already_opened`.

Hidden-information invariant: public snapshots expose only `RackCount`/`StockCount`/`DiscardRow`/`TableMelds`; never foreign `Rack` `TileInstanceId`s (verified by `internal/setup/redaction_test.go` exhaustive).

## Protocol Overview

Stable `Version=1`, `OpCode` never reused. Client `1..9`, server `100..199` (`internal/protocol/opcodes.go:1`).

| Dir | Op | Name | Payload | Valid Phase |
|-----|----|------|---------|-------------|
| C→S | 1 | `OpClientStart` | `{}` | `Waiting` (host Seat 0, ≥2 players) → `OpeningDiscard` |
| C→S | 2 | `OpClientDiscard` | `{tileId:"..."}` | `OpeningDiscard` (15→14 blocked) or `Playing/MeldOrDiscard` → advance anticlockwise `MustDraw` |
| C→S | 3 | `OpClientDrawStock` | `{}` | `Playing/MustDraw` → `MeldOrDiscard` (pop stock top, stock empty → `bad_request` per MVP `docs/rules-decisions.md:6.2`) |
| C→S | 4 | `OpClientDrawPreviousDiscard` | `{}` | `Playing/MustDraw` (opened, `HasOpened`, not opening `IsOpeningDiscard`, latest only) → `MeldOrDiscard` (discard→rack, `DiscardRow` pop) |
| C→S | 5 | `OpClientPickupDiscardForMeld` | `{discardIndex:int, tileIds:[2], jokerReps:{jokerId:{colour,rank}}, kind:"run"\|"set"}` | `Playing/MustDraw` (opened, `HasOpened`, non-opening `discardIndex` + exactly 2 `tileIds` owned → valid 3-tile meld with that discard, sweep all later discards to rack) → `MeldOrDiscard` (+ new `TableMeld`, `DiscardRow` truncated, `Racks` sweep) |
| C→S | 6 | `OpClientMeldInitial` | `{melds:[{id,kind:"run"/"set",tileIds:[3+],jokerReps:{jokerId:{colour,rank}}} ]}` | `Playing/MeldOrDiscard` (not yet opened, ≥50 with run, all tiles owned, each meld valid, no duplicate) → `HasOpened=true`, stays `MeldOrDiscard` or `RoundComplete` if `rack==0` |
| C→S | 7 | `OpClientMeldNew` | same shape | `Playing/MeldOrDiscard` (already opened, each meld valid, no score minimum, atomic, `meldId` not colliding) → stays `MeldOrDiscard` or `RoundComplete` if `rack==0` |
| C→S | 8 | `OpClientExtendMeld` | `{meldId:"...",tileIds:[1+],jokerReps:{jokerId:{colour,rank}}}` | `Playing/MeldOrDiscard` (opened, revalidates entire resulting meld, preserves `JokerReps` immutability, atomic, any owner) → stays `MeldOrDiscard` or `RoundComplete` if `rack==0` |
| C→S | 9 | `OpClientReplaceJoker` | `{targetMeldId:"...",tileId:"...",newMeldTiles:[2], jokerReps:{jokerId:{colour,rank}}, newMeldKind:"run"\|"set"}` | `Playing/MeldOrDiscard` (opened, exact tile for run or missing colour for set per `docs/rules-decisions.md:6.4` MVP, new meld `joker+2` tiles valid, atomic replacement + new meld) → stays `MeldOrDiscard` or `RoundComplete` if `rack==0` |
| S→C | 100 | `OpServerState` | `PrivateSnapshot` | match start / reconnect |
| S→C | 101 | `OpServerStatePublic` | `PublicSnapshot` | broadcast |
| S→C | 102 | `OpServerError` | `{code, message, details{}, requestId, op}` | on any `sendError` |
| S→C | 103 | `OpServerEvent` | `{"phase":..., "currentSeat":..., "op":...}` | start/discard/draw/meld |

Envelope: `{"v":1,"op":6,"requestId":"...","payload":{...}}` (`internal/protocol/envelope.go:7`). Errors echo `requestId`/`op`. See `docs/protocol.md` for schemas and `docs/state-machine.md` for phase matrix.

## How to Debug a Match

- **Logs:** `docker compose logs -f nakama | grep -E "Rummy|MatchLoop|Meld|Discard|DrawStock"` — look for `MatchJoin`, `MatchLoop op=`, `Sent error`, `Opening discard`, `MeldInitial/New`.
- **Console:** `http://127.0.0.1:7351` admin/password → Matches → `rummy` → inspect `State JSON` (public view). Private racks are never in public payload.
- **Loopback test:** use `go test ./internal/match -run TestMeldInitial -v` / `TestMeldNew` as canonical valid/invalid examples; they drive `RummyMatch.MatchLoop` via `mockDispatcher`/`protocol.MustEnvelope`.
- **Conservation:** every state-changing test calls `CheckTileConservation(state, allTiles106)` — duplicate/missing/not-in-deck fails.
- **DB:** `docker compose exec postgres psql -U postgres -d nakama -c "SELECT * FROM migrate"` or `make db-shell`.

## How to Add a New Command Safely

1. Stabilize opcode in `internal/protocol/opcodes.go` (never reuse).
2. Add payload schema to `ValidatePayload` in `internal/protocol/validator.go` with `bad_payload` details.
3. Add `AllowedOps` entry in `internal/match/phases.go` for correct `GamePhase`/`TurnPhase`.
4. Implement pure validation in `internal/rules/*` if meld/scoring logic needed (keep match thin).
5. Add `handleXxx` in `internal/match/` with phased checks (`ValidateActivePlayer`, `HasOpened`, ownership, `jokerReps`, `meld.New`, `scoring.Validate*`), atomic mutation (validate fully before touching `state`), and `sendError` with `requestId` correlation.
6. Wire handler in `MatchLoop` (`internal/match/rummy_match.go:268`) after envelope/phase validation.
7. Add view updates to `visibility.go` if public state changes.
8. Write tests: success, invalid atomic rollback, `not_opened`/`already_opened`/`not_your_turn`/`wrong_phase`, duplicate tile, joker immutability, conservation, redaction. Use `playingStateForMeldInitial`/`ForMeldNew` helpers as templates.
9. Run `make check` (`vet`+`fmt-check`+`test`) and `docker compose build`; update `docs/protocol.md` and `docs/state-machine.md` if phase/payload changed.

## Minimal Test Client (Day 24)

Minimal CLI for protocol validation without a full UI (`cmd/rummy-cli/main.go`, `go vet` clean). It runs a local simulation (in-process `RoundState` via `internal/match` + `internal/setup`) or can be pointed at a real Nakama with `--nakama` (device auth via `defaultkey`).

```bash
go run ./cmd/rummy-cli --help   # show commands
go run ./cmd/rummy-cli           # local 2-player simulation
make cli                         # alias
```

Inside the REPL:

- `state` — shows `PublicView` (`StockCount`, `DiscardRow`, `TableMelds`, `RackCount`) and `PrivateView.OwnRack` for the current seat (proves no leak: `json.Marshal(PublicView)` never contains `OwnRack` IDs, as in `visibility_test.go`).
- `switch <alice|bob>` — switch local view (simulates two browsers).
- `draw` — `DRAW_STOCK` (`MustDraw` → `MeldOrDiscard`).
- `prev` — `DRAW_PREVIOUS_DISCARD` (latest only, opened).
- `pickup <idx> <id1> <id2>` — `PICKUP_DISCARD_FOR_MELD` (2 tiles + discard at `idx` + sweep later).
- `discard <tileId>` — `DISCARD`.
- `meld <run|set> <id>...` — `MELD_INITIAL` if not opened else `MELD_NEW` (batch of 1 for demo).
- `extend <meldId> <id>...` — `EXTEND_MELD`.
- `replace <target> <tileId> <new1> <new2> [jokerId colour rank]` — `REPLACE_JOKER`.
- `winner` — show `Winner` if `RoundComplete`.

Two local users can play a manual flow: `alice: discard joker-2` → `bob: draw` → `bob: discard bF1` → `alice: draw` → `alice: meld run 5-6-7` etc. Server errors are printed as `OpServerError` JSON with `code`/`message`/`requestId`/`op`. No private data is in `Public JSON bytes` or `OpServerEvent` payloads.

For real Nakama, ensure `docker compose up --build -d` and use `--nakama` (currently falls back to local with a notice, but the envelope format is identical).

## How to Inspect

- **Nakama console:** `http://127.0.0.1:7351` user `admin` password `password` (local only).
- **Logs:** `docker compose logs -f nakama` — look for `InitModule` and `Found runtime modules`.
- **Healthcheck:** `docker inspect rummy_nakama --format '{{json .State.Health}}' | python3 -m json.tool` and `docker exec rummy_nakama /nakama/nakama healthcheck`.
- **DB:** `docker compose exec postgres pg_isready -U postgres -d nakama` or `psql` inside container.

## Troubleshooting

- **`plugin was built with a different version of ... pragma`** → protobuf mismatch; ensure `google.golang.org/protobuf v1.36.4` and `nakama-common v1.36.0` in `go.mod` and rebuild `docker compose build --no-cache`.
- **`Server key invalid` on device auth** → use `--user defaultkey:` not `defaulthttpkey:`.
- **`mkdir /nakama/data/modules: read-only file system`** (legacy TS bug) → fixed by mounting only `local.yml` not whole `data` dir (Go) — see `compose.yml:51`.
- **Port collision with tinybot** → rummy uses `5433` not `5432`; check `docker ps` and `.env.example`.
- **Nakama not loading new Go code** → you edited `main.go` but forgot `docker compose up --build -d`; host `go run` does not affect container.
- **Meld rejected `bad_payload`/`bad_request`** → `docker compose logs nakama | grep "Sent error"` shows `code`/`message`/`op`; valid meld needs `tileIds` owned in rack, `jokerReps` for every joker, `kind run/set`, `real>=2*joker`, for `MELD_INITIAL` also `≥50` and `≥1 run`.

## Docs & Decisions

- `docs/project-baseline.md` — Day 1 audit + language amendment §13.
- `docs/rules-decisions.md` — binding Romanian Tile Rummy decisions, MVP simplifications, deferred, and `TODO(product)` ambiguities (updated Day 19 + Day 15 `Kind` + Day 20 `HARDENING`).
- `docs/terminology.md` — shared vocabulary (`TileInstanceId`, `Seat`, `Rack`, `Stock`, `DiscardRow`, `HasOpened`, `GamePhase` etc.).
- `docs/protocol.md` — opcodes, envelope, payload schemas, snapshots, error codes, reconnection `OpServerState` per `rummy_match.go:79`.
- `docs/state-machine.md` — `GamePhase`/`TurnPhase` and `AllowedOps` matrix plus reconnection `MatchJoin`/`MatchLeave` and `RoundComplete`.
- `docs/testing.md` — deterministic test harness, helpers, conservation and redaction checks, now includes `TestDeterministicSimulation`.
- `docs/daily-log.md` — Handmade Hero daily slices (this README is a summary; log is authoritative).
- `AGENTS.md` — full product spec, tile set, meld rules, increment plan.
- `cmd/rummy-cli/main.go` — minimal CLI test client (Day 24, `go vet` clean, `make cli`).

## Next Steps (Roadmap)

After `MELD_INITIAL`/`MELD_NEW`/`EXTEND_MELD`/`DRAW_PREVIOUS`/`PICKUP`/`REPLACE_JOKER`/`ROUND_COMPLETE`/`SNAPSHOT HARDENING`/`DETERMINISTIC SIMULATION`/`MINIMAL CLIENT` (Days 13–21, 24, commits `de0d727`/`1b666c8`/`6b0d980`/`456c045`/`0da5f3a`/`8dd8ea9`/`2d278a5`/`3e1b92a`/`01f6a3c`/`this`):

- **Final backend regression** — `make check` (`vet`+`fmt-check`+`test` with `TestDeterministicSimulation` 7 subtests), `docker compose build`, `redaction` exhaustive, `win` invariants, `CheckTileConservation` after every action, no `OpServerEvent` leak of private rack.
- **Optional polish** — `docs/architecture.md` final, `Makefile` `cli` help, and tag `rummy-mvp-rc1` per `AGENTS.md:135` when `make smoke` passes on a fresh `docker compose up --build -d`.
- **Reconnection reference:** `internal/match/rummy_match.go:56`/`79` keeps `Players`/`Racks` on `MatchLeave` and re-sends `PrivateView` `OpServerState 100` to that `Seat` only; `docs/state-machine.md` and `docs/protocol.md` now document `PrivateSnapshot` versioned `1`.

---

*Stack: Nakama 3.26.0 (`go1.23.5`) on `postgres:15-alpine`, Go 1.23.5 plugin (`rummy_backend.so`), Docker Compose `rummy_backend:local`.*
