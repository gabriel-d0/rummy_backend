# Frontend ↔ Backend Audit — Day 1 (frontend-nakama-integration-test-plan)

**Date:** 2026-08-27 — after `client` reset to `docs-only` then rebuilt to Day 17 (`snapshot`).
**Branch:** `main@9c49ccc` — `client/src/lib/game/snapshot.ts` `Version 1`.
**Goal:** Map every Svelte touch point to its authoritative Go owner before wiring `store` → `actions` → `WS`.

## Scope
`client/` (`SvelteKit 2` `Svelte 5 runes` `Vite 6` `Tailwind 4` `TS 5` `@heroiclabs/nakama-js 2.8`) at `5173` vs `rummy_backend` (`go 1.23.5` `nakama:3.26.0` `postgres:15` `7350`/`7351` `defaultkey`) via `docs/protocol.md:1` `Envelope{v,op,requestId,payload}` `Version 1` `op 1..9`/`100..103`.

## Filesystem today
```
client/src/lib/nakama/client.ts            exists Day4 — getClient() singleton 127.0.0.1:7350 defaultkey
client/src/lib/config.ts                   exists Day6 — PUBLIC_NAKAMA_* via $env/static/public fallback
client/src/lib/ui/tokens.ts                exists Day3 — colors felt #0a2e1a etc.
client/src/lib/game/snapshot.ts            exists Day17 — Public/Private Snapshot, SnapshotVersion 1
client/src/lib/game/snapshot.test.ts       exists Day17 — 5 tests
client/src/components/Tile.svelte          exists Day11 — real Tile 64x90/52x72
client/src/components/TableBoard.svelte    exists Day13 — static melds 66/53/55 pct, no store yet
client/src/components/Rack.svelte          exists Day14 — 11 tiles static, no store yet
client/src/components/TopBar.svelte        exists Day15 — REMI ETALAT MASA, props only
client/src/routes/+page.svelte             exists Day16 — TopBar+TableBoard+Rack+Jurnal max-w-[1600px]
client/tests/smoke.e2e.ts                  exists Day7→16 — 200 + REMI ETALAT
MISSING — Day 18+ per svelte-vertical-slice: protocol.ts, auth.ts, socket.ts, match.ts, store.ts, errorStore.ts, actions.ts, reconnect.ts
docs/protocol.md, state-machine.md, visibility.go, opcodes.go — all exist backend
```

## Contract table — `frontend file | backend file | opcode | store | status`

| Frontend (Svelte) | Backend (Go) | Op | Store / Derived | Status |
|---|---|---|---|---|
| `client.ts:getClient()` `authenticate()` `createSocket()` | `compose.yml:1` `name: rummy_backend` `7350`/`7351` + Nakama `authenticateDevice` (external) + `main.go:12` `InitModule` | — (auth, not gameplay) | `localStorage rummy_device_id→rummy_token/rummy_userId` | **exists Day4** |
| `config.ts:NAKAMA_HOST` | `compose.yml:36` `nakama:3.26.0` `local.yml` `DEBUG` | — | `HOST/PORT/KEY/USE_SSL` | **exists Day6** |
| `ui/tokens.ts:colors.felt #0a2e1a` | — (UI only, no backend) | — | — | **exists Day3** |
| `components/Tile.svelte` `size rack 64x90 table 52x72` | `internal/rules/tile` `TileInstance{ID,Colour,Rank,IsJoker}` unique ID | — | — | **exists Day11** |
| `components/TableBoard.svelte` static `melds` | `internal/match/visibility.go:36` `PublicSnapshot.tableMelds[]` `TableMeld{ID,Kind,Tiles,JokerReps,OwnerSeat}` | `100 Private` `101 Public` (server→client) | `publicStore` (future Day23) | **placeholder Day13** — still `props melds`, not `derived($publicStore)` |
| `components/Rack.svelte` static 11 tiles | `internal/match/visibility.go:30` `PrivateSnapshot.ownRack/ownSeat` + `internal/rules/meld` + `scoring` | — | `privateStore` (future Day22) | **placeholder Day14** — still `props tiles`, not `derived($privateStore)` |
| `components/TopBar.svelte` `MASA n•x JUCĂTORI` props | `internal/match/state.go:12` `GamePhase Waiting` `Players[]` `CurrentSeat` | `100/101` `players` | `privateStore+publicStore` `players.length` (future Day43) | **placeholder Day15** — `players` prop, not `store` |
| `snapshot.ts:PublicSnapshot` `PrivateSnapshot` `SnapshotVersion 1` | `internal/match/visibility.go:9` `PublicSnapshot` `PrivateSnapshot` `Version 1` | `100 OpServerState Private` `101 OpServerStatePublic` | `isValidPublic/Private` + `checkNoLeak` | **exists Day17** — mirrors Go `v:1` |
| **MISSING** `protocol.ts:Envelope{v,op,requestId,payload}` `Version 1` `OpClient 1..9` `OpServer 100..103` `NewEnvelope` | `internal/protocol/opcodes.go:8` `IsClientOp 1..9` `IsServerOp 100..103` `Version 1` `envelope.go:9` `parser.go` `validator.go` `errors.go` | `1..9` `100..103` | `requestId` echo in `102` | **missing Day18** |
| **MISSING** `auth.ts:authStore→Session` `getOrCreateDeviceId` | `internal/match/state.go` `PlayerId` `Seat` | — | `authStore writable<Session\|null>` `localStorage rummy_device_id` | **missing Day19** |
| **MISSING** `socket.ts:socketStore` `setMatchDataHandler` `createSocket` `isTestEnv` mock | `internal/match/rummy_match.go:141` `MatchLoop` `dispatcher.BroadcastMessage` `OpServerState 100` via `visibility.go` | `100/101/102` | `socketStore` + `onmatchdata → handleMatchData` | **missing Day20** |
| **MISSING** `match.ts:matchStore` `createMatch rpc create_match → joinMatch` `listAvailableMatches` | `internal/match/rummy_match.go:56` `MatchJoinAttempt` `MatchJoin 79` `privateView` per `Seat` + `nakama/data/local.yml` | `100` on join | `rummy_matchId` `persistMatchId` | **missing Day21** |
| **MISSING** `store.ts:privateStore/publicStore` `onPrivateSnapshot` `privateBySeat` `handleMatchData` | `internal/match/visibility.go:36` `PublicView/PrivateView` + `internal/setup/redaction_test.go` | `100` `101` `102` | `lastPrivate` `isMyTurn = currentSeat===ownSeat` `myRack` | **missing Day22-23** |
| **MISSING** `errorStore.ts:onServerError` `errorStore 3s` | `internal/protocol/errors.go:12` `ErrorResponse{code,message,details,requestId,op}` `OpServerError 102` | `102` | `errorStore` `data-error-code` | **missing Day42/45** |
| **MISSING** `actions.ts:sendStart/Discard/DrawStock/...9 ops` `lastSent` `requestId` | `internal/match/rummy_match.go:141` `MatchLoop` `ValidateEnvelope→ValidateActivePlayer→ValidatePhaseOp→ValidatePayload→handler` + `win.go:12` | `1 Start, 2 Discard, 3 DrawStock, 4 DrawPrev, 5 Pickup, 6 MeldInitial, 7 MeldNew, 8 Extend, 9 Replace` → `102` on fail | `lastSent` | **missing Day29-38** |
| **MISSING** `reconnect.ts:reconnect()` `privateBySeat` | `internal/match/rummy_match.go:79` `MatchJoin` re-sends `PrivateView` per `Seat` only, `MatchLeave` keeps `Racks` | `100` on rejoin | `rummy_lastPrivate:${seat}` | **missing Day25-26** |
| **MISSING** `components/Toast.svelte` `WinnerOverlay.svelte` | `internal/protocol/errors.go:53` `ErrCode…` + `win.go:12` `RoundComplete` | `102` `100 winner` | `errorStore` `winner` | **missing Day39/45** |

## Gaps before Day 18 (must close before any real WS test)

1. **No `protocol.ts`** — `Rack:canDraw` etc. cannot build `NewEnvelope(3,{})` with `Version 1`; `store:handleMatchData` cannot `JSON.parse` `Envelope`. Backend rejects `bad_version`/`unknown_opcode` but frontend cannot prove `requestId` echo.
2. **No `socket.ts` `isTestEnv` mock** — `vitest` would hit real `7350` if we `createSocket` now; `match.ts` already handled this with `isTestEnv` branch at Day21 but `socket.ts` guard is prerequisite.
3. **No `store.ts`** — `TableBoard`/`Rack`/`TopBar` still `props`, not `derived($privateStore)`; `publicStore` leak test cannot run; `initGameStore` auto-wire missing.
4. **No `auth.ts` `rummy_device_id`** — `client.ts:getOrCreateDeviceId` is inside `client.ts` now, but `authStore` separation (Day19) needed for `createMatch rpc` to get `session` via `get(authStore)`.
5. **No `errorStore`** — `OpServerError 102` would be swallowed (current `handleMatchData` missing case `102`).

## How to verify “frontend ↔ backend connected” today (manual)

```bash
docker compose up --build -d && make smoke # backend up? SMOKE PASSED
cat docs/protocol.md | grep -E "^\| (1|2|3|100|101)" # opcodes stable
grep -R "SnapshotVersion" client/src/lib/game/snapshot.ts # 1
grep -R "isValidPublic" client/src/lib/game/snapshot.test.ts # 5 tests
npm run test:unit -- --run # 4 files 12 tests (tokens/snapshot/client)
npm run build # vite build green (no store yet, so no leak)
# Day1 has no WS yet — next Day 18 will add protocol + handleMatchData and prove 100/101 parse
```

## Acceptance for Day 1

- [x] Table `frontend file | backend file | opcode | store` recorded (this file)
- [x] `grep -R OpClient client/src` today = 0 (expected — Day 18 will add 1..9)
- [x] `grep -R "publicStore" client/src` today = 0 (expected — Day 22-23)
- [x] `SnapshotVersion 1` matches `internal/match/visibility.go:36` `Version 1` (proven by `snapshot.test.ts`)
- [x] Missing list 1..9 + `isTestEnv` + `errorStore` explicitly listed so Day 18+ has clear dependency

## Next
**Day 2 of frontend-nakama-integration-test-plan:** `Env & Docker smoke` — prove `VITE_NAKAMA_*` → `7350` `defaultkey` actually hits Nakama (`docker compose logs nakama | grep rummy_backend.so`, `make smoke` `RPC health`). Depends on this audit table so Day 2 knows which env vars to check.
