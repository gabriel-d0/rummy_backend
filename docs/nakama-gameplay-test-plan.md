# Nakama Gameplay — Full Test Plan (Day-by-Day, Handmade Hero)

**Goal:** Prove the Romanian Tile Rummy authoritative backend is _actually_ multiplayer-correct: rules are enforced server-side, lobbies/rooms work, joins are safe, no hidden-info leaks, rounds can be won/lost deterministically, and the Nakama Docker stack is healthy end-to-end.

**Scope:** `internal/*` (Go, `go 1.23.5`, `nakama:3.26.0`, `postgres:15`), `client/` (`SvelteKit 5` + `nakama-js 2.8`) as thin proof-client, Docker Compose. Not payments, not ranking — just one **playable round** for 2–4 players, server-authoritative.

**Prereqs:** `docker compose up --build -d` → `rummy_nakama` healthy, `make smoke` `SMOKE PASSED`, `go test ./...` + `npm run test:e2e` green on `main`.

**Test pyramid:**

- **Unit (pure, fast):** `internal/rules/*`, `internal/setup` — no Nakama.
- **Integration (mockDispatcher):** `internal/match` — `protocol.MustEnvelope` → `MatchLoop` → `mockDispatcher` + `CheckTileConservation 106`.
- **E2E (real Nakama):** `docker` + `nakama-js` + `playwright` 2 browsers — real WS `createMatch`/`joinMatch`/`sendMatchState`, `OpServerState 100`/`101` + `OpServerError 102`.
- **Manual:** 2 laptops, `auth` → `create` → `join` → `start` → `draw/discard/meld/win` via UI + `docker compose logs nakama`.

**Invariants checked every day after every state-changing test:** `CheckTileConservation(state, allTiles106)` (106 unique `TileInstanceId`, exactly one location `Racks+Stock+DiscardRow+TableMelds`), `json.Marshal(PublicView)` contains no foreign `OwnRack` IDs (redaction), atomicity (invalid cmd leaves state unchanged), phase/turn ownership.

---

## Phase A — Rules core (Days 1–5, no Nakama)

### Day 1 — Baseline audit

**Goal:** Map what `docs/testing.md` already covers and where gaps are.

- Read `internal/setup/*`, `internal/rules/meld`, `scoring`, `internal/match/*`, `internal/protocol/*`, `client/tests/*`; run `go test ./... -v | grep -E "PASS|FAIL"`, `npm run test:unit -- --run`, `npm run test:e2e` (dry).
- Record in this file §A: table `module | file | coverage` (tile deck 106, shuffle seeded, deal 15/14, Run `1-2-3`/`12-13-1`/`13-1-2`, Set 3/4 colours, joker `real>=2*joker`, scoring 5/10/25, `MELD_INITIAL` 50+run).
- Output: `docs/nakama-gameplay-test-plan.md Day1` filled, no code.
- **Accept:** List of missing cases (e.g., ace-middle, joker at run end reinterpretation, empty stock).

### Day 2 — Tile set & deck

**Goal:** 106-tile correctness is provable.

- **Tasks:** `internal/setup/deck_test.go` asserts `104 numbered +2 jokers`, exactly 2× each `Colour 1..4 × Rank 1..13`, all `ID` unique, `Validate()` passes; `shuffle_test.go` Fisher–Yates deterministic for seed `42` and divergence for different seeds; `rand_test.go` injectable `Rand`.
- **Files:** `internal/setup/*`, `internal/rules/tile`.
- **Run:** `go test ./internal/setup -run TestDeck -v`, `go test ./internal/setup -run TestShuffle -v`.
- **Accept:** `106` never 105/107, no dup ID, no lost tiles after shuffle.

### Day 3 — Meld validation (pure)

**Goal:** Runs/sets are reject/accept exactly per `docs/rules-decisions.md:1.3` + `3`.

- **Tasks:** `internal/rules/meld/matrix_test.go` — valid `1-2-3`, `12-13-1`, `10-11-12-13-1`; invalid `13-1-2`, ace-middle; set `3 vs 4` distinct colours, duplicate colour rejected; joker filling gap with explicit `JokerReps`, `real>=2*joker` (3→1, 5→1, 6→2 max), run-end joker immutable; structured `ValidationError{code,field}` not bool.
- **Run:** `go test ./internal/rules/meld -v`.
- **Accept:** All ace edge cases + joker ratio + immutability covered.

### Day 4 — Scoring & initial batch

**Goal:** 50-point opening with at-least-one-run is pure and tested.

- **Tasks:** `scoring/scoring_test.go`, `run_test.go`, `set_test.go`: `2–9:5`, `10–13:10`, Ace low `1-2-3:5` vs high `12-13-1:10` vs set `A×3:25` each, `Joly = represented`; `ValidateInitialBatch` checks owned tiles, each meld valid, `total>=50` (49 rejected, 50 accepted), `≥1 run`, no duplicate `TileId` across batch, atomic.
- **Run:** `go test ./internal/rules/scoring -v`.
- **Accept:** Exactly 49 vs 50, duplicated tile rollback.

### Day 5 — Setup, seats, conservation

**Goal:** Deterministic round init.

- **Tasks:** `internal/setup/round_test.go`: `AssignSeats` join order `0..n-1`, anticlockwise `next=(cur+1)%n`; `NewRoundState(playerIds,seed)` gives opener 15, others 14, `stock = 106-(15+14*n)`, `GamePhase Waiting→OpeningDiscard`, `TurnPhase` `""`, `IsOpeningDiscard` flag prepared; `invariant_test.go`: `CheckTileConservation` for `n=2,3,4 × seeds 0,1,42,123,999,2026`.
- **Run:** `go test ./internal/setup -run TestNewRound -v && go test ./internal/match -run TestConservation -v`.
- **Accept:** `15/14` correct, `stockCount` correct, no dup/loss.

---

## Phase B — Match skeleton & protocol (Days 6–8)

### Day 6 — State machine & turn

**Goal:** Phases and allowed ops are explicit.

- **Tasks:** `internal/match/phases_test.go`: `GamePhase Waiting/OpeningDiscard/Playing/RoundComplete`, `TurnPhase MustDraw/MeldOrDiscard`, `AllowedOps` matrix (e.g., `RoundComplete` none, `Waiting` only `Start` by Seat 0, `MustDraw` only `DRAW_STOCK`/`DRAW_PREVIOUS`/`PICKUP`, `MeldOrDiscard` only `MELD_*`/`EXTEND`/`REPLACE`/`DISCARD`), `ValidateActivePlayer`/`ValidatePhaseOp` unit.
- **Run:** `go test ./internal/match -run TestAllowed -v`.
- **Accept:** Wrong-phase → `OpServerError wrong_phase` with `requestId` echo.

### Day 7 — Nakama lifecycle (mock)

**Goal:** Match can be created/joined/left without real Nakama.

- **Tasks:** `internal/match/rummy_match_test.go`: `MatchInit` → `Waiting`, `MatchJoinAttempt`/`MatchJoin` allocate `Seat` 0..3, `MatchLeave` keeps `Seat`+`Racks` (reconnect), `MatchTerminate` clears; `mockDispatcher`/`mockPresence` helpers drive `MatchLoop` via `protocol.MustEnvelope(v,op,requestId,payload)`.
- **Run:** `go test ./internal/match -run TestMatchJoin -v`.
- **Accept:** 5th player rejected, rejoin same `Seat` succeeds.

### Day 8 — Protocol & errors

**Goal:** `Version 1` and error format stable.

- **Tasks:** `internal/protocol/envelope_test.go`, `validator_test.go`, `errors_test.go`: `Version 1`, `OpClient 1..9`/`OpServer 100..103` never reused, `Envelope{v,op,requestId,payload}` parse rejects `bad_json`/`bad_version`/`unknown_opcode`/`bad_payload`, `ValidatePayload` per opcode, `ErrorResponse{code,message,details,requestId,op}` → `OpServerError 102` to sender only (`dispatcher.BroadcastMessage(102, EncodeError, [sender])`).
- **Run:** `go test ./internal/protocol -v && go vet ./...`.
- **Accept:** Unknown `op 999` → `bad_payload`, `requestId` echoed.

---

## Phase C — Lobby & rooms (Days 9–10)

### Day 9 — Lobby / waiting room

**Goal:** 2–4 players can form a room and see each other.

- **Tasks:** Go: `match/lobby_test.go` — host creates, `listAvailableMatches` returns `Waiting` with `size` 1→4, second `joinMatch` same `matchId`, `ClearMatchId` on `leaveMatch`; Client: `client/src/lib/nakama/match.test.ts` + `tests/start.e2e.ts` — `TopBar START` visible only if `Waiting && ownSeat==0 && players>=2`, `OpClientStart 1 {}` via `sendMatchState`, `Waiting→Playing` after `Start`.
- **Run:** `go test ./internal/match -run TestLobby -v && npm run test:unit -- --run && npm run test:e2e -- tests/start.e2e.ts`.
- **Accept:** 1-player room not startable, guest sees no `START`, `PublicView.players` correct, `container` fix `5433`/`7350` not colliding with `tinybot`.

### Day 10 — Start + opening discard

**Goal:** First mandatory move.

- **Tasks:** `opening_discard_test.go`: opener discards 1 of 15 → `14`, `DiscardRow[0].IsOpeningDiscard=true` permanent `CanPickupPreviousDiscard` false for that entry, `CurrentSeat` advances `0→1` anticlockwise, phase `OpeningDiscard→Playing MustDraw`; reject non-opener, foreign `TileId`, empty rack, `already_opened` not relevant.
- **Run:** `go test ./internal/match -run TestOpeningDiscard -v`.
- **Accept:** Opening discard `15→14`, tile `106` still holds, opening entry never pickable.

---

## Phase D — Core turn loop (Days 11–16, the critical gameplay slice)

### Day 11 — Draw from stock + normal discard + turn rotation

**Goal:** Draw→meld→discard loop.

- **Tasks:** `draw_stock_test.go`: `MUST_DRAW` only current player, `stock top`→`rack`, `MUST_DRAW→MELD_OR_DISCARD`, empty stock `→ bad_request` or `stockExhausted` dead-round per `docs/rules-decisions.md:6.2`; `discard_test.go`: `MELD_OR_DISCARD` only via `OpClientDiscard 2 {tileId}` owned, appends `DiscardRow` ordered, `CurrentSeat=(cur+1)%n` → `MUST_DRAW`; verify `2,3,4` player rotation.
- **Run:** `go test ./internal/match -run TestDrawStock -v && go test ./internal/match -run TestDiscard -v`.
- **Accept:** `draw` before `meld`, `discard` after `draw`, row order preserved.

### Day 12 — MELD_INITIAL (50+ with run)

**Goal:** Open with server-trusted validation.

- **Tasks:** `meld_initial_test.go`: batch `melds:[{id,kind:run/set,tileIds:[3+],jokerReps:{jokerId:{colour,rank}}}]`, all owned, each meld `ValidateRun/Set` ok, `total>=50` with `≥1 run`, no duplicate `TileId` across batch, `HasOpened false→true`, stays `MELD_OR_DISCARD` (must discard), atomic on fail; `visibility_test.go` public sees `TableMelds` with `JokerReps` but not rack IDs.
- **Run:** `go test ./internal/match -run TestMeldInitial -v`.
- **Accept:** 49 rejected, 50 accepted, invalid batch leaves `Racks`+`TableMelds` unchanged.

### Day 13 — MELD_NEW & EXTEND_MELD

**Goal:** Opened player can build.

- **Tasks:** `meld_new_test.go`: `HasOpened==true` else `not_opened`, no score minimum, multiple melds per batch, `meldId` not colliding, atomic; `extend_meld_test.go`: `meldId` exists, `tileIds` owned, `combinedTiles = existing.Tiles+new` + `combinedReps` preserves `JokerReps` immutability, entire resulting meld revalidated via `meld.New`, allows extending others’ `OwnerSeat` unchanged, `OwnerSeat` stable.
- **Run:** `go test ./internal/match -run TestMeldNew -v && go test ./internal/match -run TestExtend -v`.
- **Accept:** Extending run at either legal end, set `3→4` colours, wrong colour/gap → no mutation.

### Day 14 — Discard pickup (previous + earlier)

**Goal:** Romanian special pickup.

- **Tasks:** `draw_previous_discard_test.go`: `MustDraw && HasOpened && last !IsOpeningDiscard` → latest `DiscardRow tail` → `Racks[seat]`, `DiscardRow pop`, `MUST_DRAW→MELD_OR_DISCARD`; `pickup_discard_for_meld_test.go`: `MustDraw && HasOpened` + `{discardIndex, tileIds:[2], jokerReps, kind:run|set}` where `discardTile+2 tiles` is valid 3-tile meld with that discard, sweep `DiscardRow[discardIndex+1:]` to rack in order, `DiscardRow=[:discardIndex]` reindexed, `+ new TableMeld`; opening index never selectable.
- **Run:** `go test ./internal/match -run TestDrawPrevious -v && go test ./internal/match -run TestPickup -v`.
- **Accept:** Latest only for `DRAW_PREVIOUS`, earlier needs immediate 2+discard valid meld, later discards swept, atomic.

### Day 15 — REPLACE_JOKER

**Goal:** Legal joker recovery only.

- **Tasks:** `replace_joker_test.go`: `MeldOrDiscard && HasOpened` + `{targetMeldId, tileId, newMeldTiles:[2], jokerReps, newMeldKind}` requires exact `tileId` equals joker’s represented `Colour`/`Rank` (run) or exact missing colour (set, 3-tile set with 1 joker per `6.4`), `updatedMeld` and `newMeld (joker+2)` both `ValidateRun/Set` ok, atomically `Racks[seat] -3 (tileId+2)`, `TableMelds[targetIdx]=updated` (keeps `OwnerSeat`), `+ new meld`, stays `MELD_OR_DISCARD`.
- **Run:** `go test ./internal/match -run TestReplaceJoker -v`.
- **Accept:** Wrong represented tile rejected, new mell needs `2` rack tiles, no silent `JokerReps` shift.

### Day 16 — Win / lose (RoundComplete)

**Goal:** Empty rack wins deterministically.

- **Tasks:** `win_test.go` + `win.go:checkWinAndComplete`: `len(Racks[seat])==0` after `DISCARD` or after any `MELD_*`/`EXTEND`/`REPLACE`/`PICKUP` per `docs/rules-decisions.md:6.1` (win without final discard allowed), `GamePhase→RoundComplete`, `Winner=seat`, `CurrentSeat=winner`, broadcasts `101 RoundComplete`/`103 roundComplete`, `AllowedOps` none, `CheckTileConservation 106` still holds, `PublicView` winner visible.
- **Run:** `go test ./internal/match -run TestWin -v`.
- **Accept:** `WinAfterMeldWithoutDiscard`, `WinAfterDiscardToZero`, `NoGameplayAfterRoundComplete`.

---

## Phase E — Security & determinism (Days 17–20)

### Day 17 — Visibility & redaction

**Goal:** No private leak, ever.

- **Tasks:** `visibility_test.go` + `setup/redaction_test.go` exhaustive `n=2,3,4 × seeds 42,7,123` with real `NewRoundState` 106: `json.Marshal(PublicView)` string search for every `rack`/`stock` `TileInstanceId` fails (only `stockCount`/`RackCount`/`DiscardRow`/`TableMelds` visible), each `PrivateView` contains own `OwnRack` and not others nor stock; client `src/lib/game/redaction.test.ts` same via `checkNoLeak`, `src/lib/game/store.test.ts` `PublicSnapshot 101` vs `PrivateSnapshot 100`.
- **Run:** `go test ./internal/setup -run TestRedaction -v && npm run test:unit -- --run`.
- **Accept:** `Public JSON` never contains foreign `OwnRack` IDs, `OpServerEvent` payloads redacted (verified).

### Day 18 — Reconnection

**Goal:** Returning player sees own rack again, others never.

- **Tasks:** `rummy_match.go:MatchLeave` keeps `Players`/`Racks`, `MatchJoin` re-sends `PrivateView 100` to that `Seat` only; `client/src/lib/nakama/reconnect.test.ts` + `game/store.test.ts` `privateBySeat` + `localStorage rummy_lastPrivate:${seat}` + `rummy_matchId`; `publicStore`/`privateStore` rehydrate from `OwnRack` not old, `Window` `reconnect()` test `new-1 not old-1`.
- **Run:** `go test ./internal/match -run TestReconnection -v && npm run test:unit -- --run`.
- **Accept:** Rejoin `Waiting` allowed (fix `6f2af5b`), `RoundComplete` winner still visible.

### Day 19 — Deterministic simulation harness

**Goal:** One readable end-to-end Go simulation proves the whole flow.

- **Tasks:** `internal/match/deterministic_simulation_test.go` — builders for named tiles/racks/stock/discard, `NewRoundState` with fixed `seed 2026`, `RunDeterministicSimulation` steps: `Start` → `OpeningDiscard` → `DrawStock→Discard` loop → `MeldInitial 50+` → `Extend` → `DrawPrevious` → `PickupDiscardForMeld` → `ReplaceJoker` → `RoundComplete`, after every action `CheckTileConservation` + `TurnPhase`/`CurrentSeat`/`HasOpened` assertions, failure pinpoints transition/rule.
- **Run:** `go test ./internal/match -run TestDeterministicSimulation -v`.
- **Accept:** 7 subtests pass, `MeldOrDiscard` stays until `DISCARD` or win.

### Day 20 — Invariant stress / fuzz

**Goal:** Random input cannot corrupt state.

- **Tasks:** `go test -run TestInvariantStress` — `1000` random legal/illegal `MustEnvelope` sequences with `seed` table `t1`/`j1`, assert no panic, `CheckTileConservation`, `Validate()` holds; `protocol/fuzz_test.go` — `go test -fuzz FuzzParseEnvelope -fuzztime 30s` for `bad_json`/`unknown_opcode`; `internal/rules/meld` fuzz for `real>=2*joker`.
- **Run:** `go test ./... -count=100 -run TestInvariantStress && go test ./internal/protocol -fuzz=.`.
- **Accept:** No `TileId` dup/loss, all malformed → `OpServerError` not crash.

---

## Phase F — Real Nakama (Days 21–24, Docker)

### Day 21 — Real Nakama smoke

**Goal:** Backend boots and answers.

- **Tasks:** `make smoke` (`scripts/smoke.sh`): `pg_isready`, `nakama healthcheck`, `console 7351 200`, `Found runtime modules [rummy_backend.so]`, `InitModule Rummy backend`, `RPC health`/`version` via `defaultkey` `device auth` (UUID→token), `rpc/list_matches` returns.
- **Run:** `docker compose up --build -d && docker compose logs nakama --tail=100 | grep -E "Rummy|modules|Registered" && make smoke`.
- **Accept:** `SMOKE PASSED` with `health/version` `rummy` match handler, no `plugin was built with` error (`protobuf v1.36.4` pinned).

### Day 22 — Two real clients join & start

**Goal:** 2 browsers can create/join via `nakama-js`.

- **Tasks:** `client/src/lib/nakama/match.test.ts` mock → `client/tests/e2e-actions.e2e.ts` 2 `BrowserContext` `alice`/`bob` `createMatch()`→`joinMatch(matchId)` same `matchId`, `PrivateView 100` per seat `OwnRack` distinct, `PublicView 101` shared `TableMelds`; manual: `docker exec` `nakama` + `client` at `5173` `Camere disponibile` live.
- **Run:** `npm run test:e2e -- tests/e2e-actions.e2e.ts` + manual lobby.
- **Accept:** `alice` vs `bob` different `Rack`, `STOCK` same count, no leak via Network payload.

### Day 23 — Real gameplay to win/lose via WS

**Goal:** Full round over WS ends with winner.

- **Tasks:** `client/tests/visual-actions.e2e.ts` + `client/tests/winner.e2e.ts` → `2 browsers` `START` → `OpeningDiscard 15→14` → `DrawStock 3` → `MeldInitial 6` (50+ run) → `MeldNew 7` → `Extend 8` → `DrawPrevious 4` / `Pickup 5` → `Replace 9` → `Discard 2` loop → `RoundComplete Winner` overlay `RESTART`; `internal/match/win_test.go` `NoGameplayAfterRoundComplete` → further `DISCARD` → `OpServerError`.
- **Run:** `npm run test:e2e -- tests/visual-actions.e2e.ts tests/winner.e2e.ts` and `go test ./internal/match -run TestWin -v` with real `MatchLoop`.
- **Accept:** Winner recorded, `Winner 0` visible, replay shows same with fixed `seed`.

### Day 24 — Negative & error cases over real Nakama

**Goal:** Invalid client input never crashes or mutates.

- **Tasks:** `tests/error.e2e.ts` `DISCARD tileId not in rack` → `OpServerError 102 {code:bad_payload, requestId, op:2}`, `MELD_INITIAL 49` → `OpServerError`, `not_your_turn`, `wrong_phase`, `not_opened`/`already_opened`, duplicate `TileId`, `jokerReps` missing; `client/src/components/Toast.svelte` shows `bg #dc2626` `3s` with `data-error-code`; `go test ./internal/match -run Test*Rejected -v` assert atomic rollback.
- **Run:** `npm run test:e2e -- tests/error.e2e.ts && go test ./internal/match -run TestInvalid -v`.
- **Accept:** All errors `102` to sender only, state `CheckTileConservation` holds, toast `LEAKED` not present.

---

## Phase G — Confidence (Days 25–26)

### Day 25 — Load & parallel

**Goal:** 4 players × 10 concurrent matches don’t leak or dead-lock.

- **Tasks:** `go test -run TestParallelMatches -count=10` — `10` goroutines each `NewRoundState` with `n=4`, parallel `MatchLoop` `Start→Win`, `CheckTileConservation`; `client` `playwright` 4 browsers joins same `Waiting` → `4/4` `Masa 1 • 4 JUCĂTORI`.
- **Run:** `go test ./internal/match -run TestParallel -v -race`.
- **Accept:** No `TileId` dup under `-race`, no blocked `TickRate 5`.

### Day 26 — Manual QA checklist + CI

**Goal:** Anyone can run it without oral knowledge.

- **Tasks:** Create `docs/manual-qa.md` 1-page checklist (auth → create → 2nd join → START → opening 15→14 → draw → meld 50+ → extend others’ meld → IA ULTIMA → RIDICĂ → ÎNLOCUIEȘTE JOLY → ARUNCĂ → WIN → RESTART → reconnect → error toast); wire `make check` (`vet+fmt-check+test`) + `make build` (`compose build`) + `make smoke` + `npm run check && npm run test:unit -- --run && npm run build && npm run test:e2e` in `.github/workflows/ci.yml` for both Go and client.
- **Run:** `make check && docker compose build && npm run check && npm run build`.
- **Accept:** `README` `How to test` lists exactly `make check`/`npm run test:e2e`/`make smoke` with expected logs.

---

## Appendix — Commands cheat sheet

```bash
# Go authoritative (fast, deterministic)
go vet ./... && go test ./... -v
go test ./internal/rules/meld -run TestValidate -v
go test ./internal/match -run TestDeterministicSimulation -v
go test ./internal/setup -run TestRedaction -v

# Docker real Nakama
docker compose up --build -d
docker compose logs nakama --tail=100 | grep -E "Rummy|modules|Registered|health"
make smoke  # pg_isready + healthcheck + InitModule + console 200 + RPC health

# Client stores
npm run check          # svelte-check
npm run lint           # prettier --check + eslint
npm run test:unit -- --run   # vitest (snapshot, protocol, store, redaction)
npm run test:e2e       # playwright (all demos + 2-browser alice/bob + win)
npm run test:e2e -- tests/winner.e2e.ts  # single
npm run build

# Full gate (CI)
make check && docker compose build && npm run check && npm run test:unit -- --run && npm run build && npm run test:e2e
```

**Tag:** when Days 1–26 green, `git tag nakama-gameplay-verified` (like `rummy-mvp-rc1`). Next is client polish (`client/docs/roadmap.md:53`) — no gameplay logic changes without new day.

---

_Refs: `docs/rules-decisions.md` §1–10 incl. 15.9 `Kind` + `win.go` + `HARDENING`, `docs/protocol.md` opcodes `1..9`/`100..103` `Version 1`, `docs/state-machine.md` `AllowedOps`, `docs/architecture.md` visibility, `AGENTS.md` 796 lines timeline `Day 24` CLI, `internal/match/*_test.go` `CheckTileConservation 106`._
