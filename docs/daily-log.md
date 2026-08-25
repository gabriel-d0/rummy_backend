# Daily Log — Handmade Hero Incremental Development

This log records each “one clearly scoped change per day” per `AGENTS.md:5` and the Daily Execution Protocol `AGENTS.md:189`. It is the human-readable companion to `git log --oneline --reverse`. All Days were executed **2026-08-25** in compressed form for the initial foundation (plan would normally be one commit per calendar day); the one-commit-per-day rule and runnable-state gate were still enforced per day slice.

> Branch: `main` · Remote: `git@github.com:gabriel-d0/rummy_backend.git` · All pushes `main -> main` · `git config user.name "Gabriel D."` `user.email "gabriel0x01@proton.me"`

---

## Phase 1 — Foundation and local development (Days 1–7)

| Day | Commit | Files | Goal / Acceptance |
|-----|--------|-------|-------------------|
| **1** | `36c2c59` `docs: add Phase 1 Day 1 project baseline audit` | `AGENTS.md` (+tracked), `docs/project-baseline.md` (new, 228 lines) | Review repo, toolchain (`Node 26`, `npm 11`, `Go 1.27`, `Docker 29`, `compose v5`), branching (`main` unborn, SSH origin), CI/deployment assumptions, record in `docs/project-baseline.md`. **Decision:** TypeScript per `AGENTS.md:125` (no established language) → later amended to Go §13. |
| **2** | `48ec5ce` `chore: add Docker Compose for Nakama + Postgres local development` | `compose.yml` (65 lines, `name: rummy_backend`, `postgres:15-alpine 5433` vs `tinybot 5432`, `nakama:3.26.0`), `nakama/data/local.yml`, `.env.example`, `.gitignore`/`.dockerignore` (50 lines) | `docker compose up -d` both `healthy`, `migrate 0`, console `200` at `:7351` `healthcheck ok`, `down`/`up` persists `rummy_backend_pgdata`. Fixed initial `:ro` mount that caused `mkdir modules: read-only`. |
| **3** | `9354b59` `feat: add Nakama TypeScript runtime skeleton with InitModule` → superseded by `55c7f3b` `refactor: migrate Nakama runtime from TypeScript to Go` | Removed `package.json`/`tsconfig.json`/`src/main.ts`/`build/`, added `go.mod` (`go 1.23.5` `nakama-common v1.36.0` `protobuf v1.36.4`), `main.go` (`InitModule` `health`/`version` RPCs), `Dockerfile` (`pluginbuilder:3.26.0` → `nakama:3.26.0` `backend.so` at `/nakama/data/modules/rummy_backend.so`), updated `compose.yml` to `build: .` + mount only `local.yml:ro` | `go vet 0` + `docker compose build 9.5s` + `up` logs `Found runtime modules count 1 [rummy_backend.so]` `Registered RPC health/version` + `curl /v2/rpc/health` via `defaultkey` returns `status ok`. **Amendment:** `docs/project-baseline.md:13` records `protobuf v1.36.4` pin to avoid `pragma` mismatch. |
| **4** | `61607d2` `docs: update documentation for Go runtime (README + project baseline amendment)` | `README.md` (new, 153 lines), `docs/project-baseline.md` amended `§13` (35 lines) | README covers `prerequisites` `quick start` `services & ports` `rebuild vs darwin .so` `start/stop/logs/clean` `go vet/test/fmt` `RPC smoke` `troubleshooting` `next steps`. Baseline adds Go migration note. |
| **4b** | `8493782` `chore: add developer scripts Makefile for Go + Docker workflow` | `Makefile` (105 lines, `.PHONY` `help`/`build`/`up`/`down`/`restart`/`ps`/`logs`/`clean`/`reset`/`db-shell`/`vet`/`fmt`/`test`/`tidy`/`check`/`health`) | `make vet` `make fmt-check` `make test [no test files]` `make ps` healthy `make health` `console 200` `RPC health` — all `0`. |
| **6** | `162a0a6` `ci: add GitHub Actions baseline for Go vet/fmt/test and Docker build` | `.github/workflows/ci.yml` (74 lines, `setup-go 1.23.5` `go vet` `gofmt -l` `go test` `go mod tidy` diff `docker/build-push-action` `compose build`) | `make check` mirrors CI `0`; `ci.yml` is Day 6 deliverable `AGENTS.md:242`. |
| **7** | `74acc99` `test: add repeatable smoke test for Nakama startup, DB and Go runtime` | `scripts/smoke.sh` (203 lines, `set -euo pipefail`, `ensure_up` `wait_for_health` `pg_isready` `healthcheck` `InitModule` `rummy_backend.so` `console 200` `RPC health` via `defaultkey`), `Makefile` adds `smoke` target | `make smoke` `SMOKE PASSED` when `healthy` and when `down` (rebuilds), repeatable. Milestone `Phase 1` complete: `one command` `docker compose up --build -d` starts healthy backend. |

*Day 5 local setup docs were merged into Day 4 `61607d2` README per Go migration; no separate commit.*

---

## Phase 2 — Rules specification and core domain (Days 8–14)

| Day | Commit | Files | Goal / Acceptance |
|-----|--------|-------|-------------------|
| **8** | `2b44df6` `docs: add Phase 2 Day 8 Romanian Tile Rummy rules decisions` | `docs/rules-decisions.md` (189 lines) | Source Pagat `106` tiles `2–4` anticlockwise `suita 1-2-3/12-13-1` `terta` `joker 2:1` + `50`‑point opening `5/10/25` scoring `15→discard` blocked `draw/meld/discard` `previous/earlier discard pickup` `replace joker with 2 tiles` `win rackEmpty`, MVP simplifications (Fisher–Yates, `Seat 0` opener, `15/14` `stock`), deferred `Doubla` etc., ambiguities `§6` with `TODO(product)` `final discard`/`stock exhaustion`/`doubla`/`set joker colours` per `AGENTS.md:92`. |
| **9** | `c1eb89b` `docs: define shared domain terminology for rummy` | `docs/terminology.md` (59 lines) | Canonical tables `Tile`/`TileInstanceId`/`Colour 4`/`Rank 1..13` `Joker` `PlayerId` `Seat 0..n-1` `NextSeat (current+1)%n` `Rack` private `Stock` count-only `Discard Row` ordered `Opening Discard flagged` `TableMelds` `Run/Set` `Opening Meld HasOpened` `GamePhase`/`TurnPhase` `Current` `Draw/Discard/Extend/ReplaceJoker/Win` per `AGENTS.md:259`. |
| **10** | `53553b2` `feat: add tile domain model with colours ranks and unique IDs` | `internal/rules/tile/tile.go` (175 lines) `tile_test.go` (141 lines, 8 tests) | `Colour` `Red/Yellow/Blue/Black` `IsValid`/`String`/`ParseColour`, `Rank 1..13` `IsValid`, `TileInstanceId` opaque, `TileInstance{ID,Colour,Rank,IsJoker}` `Validate`/`IsNumbered`/`String` `NewTile/NewJoker/Must*`. Tests: colour/rank valid/invalid, `NewTile` empty ID, `Must` panic, `String` `Joker{}`. `go test ./internal/rules/tile -v 8 passed`. |
| **11** | `a0f2147` `feat: add player and seat model with anticlockwise turn order` | `internal/match/seat.go` (140 lines) `seat_test.go` (124 lines, 5 tests) | `PlayerId` `Seat -1..3` `PlayerState{ID,Seat,HasOpened}` `SeatsForCount` `AssignSeats` deterministic `Seat 0` opener, `NextSeat (current+1)%n`/`PrevSeat` anticlockwise `docs/rules-decisions.md:1.2`, `SeatOfPlayer`/`ValidatePlayers` `2..4` dup checks. Tests `2/3/4` `deterministic` `full` `duplicate` `started` `2→1→0` loops. |
| **12** | `4d67326` `feat: add initial server-side round state structures` | `internal/match/state.go` (302 lines) `state_test.go` (4 tests) | `GamePhase Waiting/OpeningDiscard/Playing/RoundComplete` `TurnPhase MustDraw/MeldOrDiscard` `DiscardEntry{Tile,IsOpeningDiscard,Index}` only `0` may be opening, `TableMeld{ID,Tiles,JokerReps,OwnerSeat}` `Validate`, `RoundState{Players,Racks,Stock,DiscardRow,TableMelds,CurrentSeat,GamePhase,TurnPhase,Winner}` `Validate` per `docs/terminology.md`. Tests `ValidateMinimal` `RejectsDuplicate` `OpeningOnlyFirst` `JokerRep`. |
| **13** | `b83740c` `feat: add tile conservation invariant checker` | `internal/match/invariant.go` (99 lines) `invariant_test.go` (159 lines, 4 tests + syntheticDeck) | `CheckTileConservation(state, allTiles 106)` maps expected `106` `seen` duplicate `not in deck` `missing` with location strings `rack/stock/discard/meld`, `CountTiles`. Tests `EmptyMelds 15+14+1+76=106` `Duplicate` `Missing` `WithMelds 3` `dogs`. |
| **14** | `8966391` `test: add domain tests for tile identity, seats, turn direction and state` | `internal/match/domain_test.go` (168 lines, 4 tests) | Aggregated `TileIdentity` same face `red-5` distinct IDs, `SeatsTurnDirection` `2p/3p/4p` `NextSeat` invertibility `AssignSeats` join-order, `BaseStateConstruction` `15+14+14+1+62=106` `Validate`+`Conservation`, `InvariantFailures` duplicate/missing. Docs milestone `Phase 2` complete: safe internal state model. |

---

## Phase 3 — Deck, randomness, and round setup (Days 15–21)

| Day | Commit | Files | Goal / Acceptance |
|-----|--------|-------|-------------------|
| **15** | `2aee5d4` `feat: create deterministic 106-tile Romanian deck factory` | `internal/setup/deck.go` (45 lines) | `NewDeck()` deterministic `Red 1-13×2`→`Yellow`→`Blue`→`Black`→`joker-1/2` `106` unique `TileInstanceId` `red-01-1` etc. `NewDeckIDs()` helper. Pure, `Validate` ready. |
| **16** | `9ed478a` `test: add exhaustive deck correctness tests for 106-tile factory` | `internal/setup/deck_test.go` (123 lines, 7 tests) | `TotalCount 106` `Counts 104/2` each `Validate` `TwoPerColourRank 2 per 52 combos` `AllIDsUnique 106` `IDsMatch` `Deterministic` snapshot `red-01-1/2` `joker-2` `TilesValidate`. `go test ./internal/setup -v 7 passed`. |
| **17** | `47ade4e` `feat: add injectable seedable random source for deterministic shuffles` | `internal/setup/rand.go` (31 lines) `rand_test.go` (71 lines, 4 tests) | `Rand{*rand.Rand}` `NewSeededRand(seed)` `Intn` `Int63` `Shuffle(n,swap)` wrapper. Tests `DeterministicIntn 42` `DifferentSeedsDiverge` `ShuffleDeterministic 12345` `Int63`. Seed controls output, injectable for `Shuffle` (Day 18). |
| **18** | `05de37b` `feat: implement deterministic Fisher–Yates shuffle via injected Rand` | `internal/setup/shuffle.go` (22 lines) `shuffle_test.go` (106 lines, 6 tests) | `Shuffle(deck,Rand) []TileInstance` copies then `r.Shuffle(len,swap)` Fisher–Yates, original not mutated, `nil` panics, permutation only. Tests `DeterministicFixedSeed 42` equal, `ChangesOrder 123` not identity, `NoLostDuplicated 106` `Validate` each, `OriginalUnmodified`, `DifferentSeedsDiverge`, `PanicsNilRand`. |
| **19** | `94068b7` `feat: implement MVP dealing 15/14 and stock remainder` | `internal/setup/deal.go` (52 lines) `deal_test.go` (117 lines, 4 tests) | `Deal(deck,n)` validates `106` `2..4` `needed 15+(n-1)*14=29/43/57` `Racks map[Seat]` `Seat 0→15` others `14` `stock 77/63/49` `DealForPlayers`. Tests `Counts 2→77 3→63 4→49` duplicate/seen `106` `InvalidCounts` `Deterministic` `ConservationWithState` via `match.CheckTileConservation`. |
| **20** | `3daf725` `feat: create initial round state for 2–4 players with deterministic deal` | `internal/setup/round.go` (58 lines) `round_test.go` (137 lines, 4 tests) | `NewRoundState(playerIds,seed)` `NewDeck+Shuffle(NewSeededRand(seed))+AssignSeats+Deal` → `RoundState{Players,Racks,Stock,DiscardRow empty,TableMelds empty,CurrentSeat 0,PhaseOpeningDiscard,TurnMustDraw}` `Validate`+`Conservation`. Tests `Counts 2→77 3→63 4→49` `DeterministicSeed 123==123 999 diverge` `InvalidCounts 1/5 dup` `WithDeck 4p`. Milestone `Phase 3` complete: legal round initialized deterministically. |
| **21** | `b2f5634` `test: verify setup invariants for 2/3/4 players across seeds` | `internal/setup/setup_invariant_test.go` (61 lines, 1 test `18` combos) | `TestSetupInvariantsAllPlayersAndSeeds` `n=2,3,4 × seeds 0,1,42,123,999,2026` `NewRoundState` asserts `15/14` `stock 77/63/49` `discard 0 melds 0` `CurrentSeat 0 OpeningDiscard` `Validate` `Conservation` `CountTiles 106`. |

---

## Phase 4 — Nakama match foundation (Days 22–31)

| Day | Commit | Files | Goal / Acceptance |
|-----|--------|-------|-------------------|
| **22** | `9d4a888` `feat: register authoritative rummy match handler in Nakama` | `internal/match/rummy_match.go` (80 lines, `RummyMatch` `7` methods) `main.go` registers `rummy` | `NewRummyMatch` factory, `MatchInit` `Waiting` `tickRate 5` `label rummy` empty `RoundState`, `MatchJoinAttempt/Join/Leave/Loop/Terminate/Signal` skeletons. `main.go:12` `InitModule` logs `health/version` + `RegisterMatch("rummy",match.NewRummyMatch)` → `Registered Go runtime Match creation function invocation name rummy`. `go test` `0` `docker compose build 9.5s` `health` `rummy_backend.so`. |
| **23** | `9132ca7` `feat: implement lobby waiting room with seat allocation and join/leave` | `internal/match/rummy_match.go` updated (81 lines) `rummy_match_test.go` (179 lines, 4 tests + mocks) | `MatchJoinAttempt` checks `Waiting`, `full 4`, `already joined`, `bad state`; `MatchJoin` allocates `Seat(len(Players))` `Racks[Seat] empty` idempotent, `MatchLabelUpdate`; `MatchLeave` removes `Player` + `delete(Racks)`. Mocks `testLogger` `mockPresence` `mockDispatcher` `PresenceReason`. Tests `MatchInit` `tick 5` `Waiting`, `JoinAttempt` allow/full/duplicate/started, `JoinAndLeave` `3` `0..2` duplicate leave, `LoopAndTerminate` nil-safe. `21` tests in `match`. |
| **24** | `14d9b93` `feat: allow host to start match from waiting room via MatchLoop/Signal` | `internal/match/rummy_match.go` `MatchLoop` handles `op 1 START` + `MatchSignal "start"` | Host `Seat 0` `2..4` `Waiting→OpeningDiscard` `CurrentSeat 0` `TurnMustDraw` broadcast `OpServerEvent`; rejects non-host/`already started`/`<2` with `Warn` + `OpServerError`. Tests `TestMatchStartViaLoop` non-host rejected host transitions, `NotEnoughPlayers` `1` stays `Waiting`, `ViaSignal` `start`/`start:alice` vs `start:bob` `only host`. |
| **26** | `437beb4` `feat: define stable protocol opcodes and envelope versioning` | `internal/protocol/opcodes.go` (48 lines) `opcodes_test.go` (40 lines) `internal/match/rummy_match.go` uses `protocol.OpClientStart`/`OpServerEvent` | `Version 1`, client `1..99` `Start=1` stable (Day 24 `1`), `Discard=2` `DrawStock=3` `DrawPrevious=4` `Pickup=5` `MeldInitial=6` `MeldNew=7` `Extend=8` `ReplaceJoker=9`, server `100..199` `State=100` `StatePublic=101` `Error=102` `Event=103`, `IsClientOp`/`IsServerOp`. Test `StableAndUnique` `1..199` `Version 1` `Start==1`. |
| **27** | `ec8d286` `feat: add safe command envelope parser for client messages` | `internal/protocol/envelope.go` (32 lines) `parser.go` (52 lines) `parser_test.go` (68 lines, 3 tests) | `Envelope{v,op,requestId,payload}` `ParseError{Code,Message}` `ParseEnvelope(data)` rejects `empty`/`bad_json`/`bad_version`/`unknown_opcode`/`bad_payload` via `IsClientOp` `switch 1..9` never panics. Tests `Valid` `RejectsMalformed 8 cases` `DoesNotPanic`. |
| **28** | `b6110a8` `feat: validate payload schemas per opcode` | `internal/protocol/validator.go` (120 lines) `validator_test.go` (130 lines, 7 tests) | `ValidatePayload(op,payload)` switch `9` ops: `Start` empty/`{}`; `Discard` `{tileId}`; `DrawStock/DrawPrevious` empty; `Pickup` `{discardIndex>=0 + tileIds[2]}`; `MeldInitial/New` `{melds[>=1]}`; `Extend` `{meldId + tileIds[>=1]}`; `ReplaceJoker` `{targetMeldId + tileId + newMeldTiles[2]}` → `bad_payload` with field; `ValidateEnvelope` = `ParseEnvelope`+`ValidatePayload`. |
| **29** | `fc48c28` `feat: add standard error response protocol with requestId correlation` | `internal/protocol/envelope.go` added `RequestId`, `errors.go` (62 lines) `errors_test.go` (70 lines, 5 tests) | `Envelope{RequestId}`, `NewEnvelopeWithRequestId`, `ErrorResponse{code,message,details,requestId,op}` `OpServerError 102` `NewError`/`NewErrorForEnvelope` echo `RequestId`/`OpCode` `EncodeError` JSON `omitempty`, codes `bad_request`…`already_opened` per `AGENTS.md:370`. Tests `NewError` `ForEnvelope` `JSON omits empty` `RequestId round-trip` `CodesStable`. |
| **30** | `2cc3313` `feat: add public/private view projection with redaction tests` | `internal/match/visibility.go` (65 lines) `visibility_test.go` (120 lines, 3 tests) | `PublicSnapshot{Version,GamePhase,TurnPhase,CurrentSeat,Players{RackCount},StockCount,DiscardRow,TableMelds,Winner}` `PrivateSnapshot{PublicSnapshot,OwnRack,OwnSeat}` `PublicView` counts-only `PrivateView` own copy. Tests via `json.Marshal` string search `HidesRacks` `OpponentNotLeaked` `Structure` per `AGENTS.md:173`. Avoids `setup→match` cycle via synthetic small racks. |
| **31** | `4f60325` `test: add exhaustive hidden-information redaction for public/private snapshots` | `internal/setup/redaction_test.go` (78 lines, 1 test `9` combos) | `TestSnapshotRedactionExhaustive` `n=2,3,4 × seeds 42,7,123` `NewRoundState` real `106` deck `rackIDsBySeat` `PublicView` must not contain any `rack`/`stock` `ID` (only `stockCount`), each `PrivateView` must contain own not others nor stock. Placed in `internal/setup` to avoid import cycle (`setup` imports `match`). Phase 4 milestone complete. |

---

## Phase 5 — Turn state machine and basic actions (Days 32–35)

| Day | Commit | Files | Goal / Acceptance |
|-----|--------|-------|-------------------|
| **32** | `307b81d` `feat: define turn state machine and allowed ops per phase` | `internal/match/phases.go` (55 lines) `phases_test.go` (60 lines, 2 tests) | `AllowedOps(gamePhase,turnPhase) map[int64]bool` using `protocol.OpClient*`: `Waiting→Start`, `OpeningDiscard→Discard`, `Playing MustDraw→DrawStock/DrawPrevious/Pickup`, `MeldOrDiscard→Discard/MeldInitial/MeldNew/Extend/ReplaceJoker`, `RoundComplete→none`; `ValidatePhase` `GamePhase.IsValid` `TurnPhase.IsValid`. Tests `AllowedOps` matrix + `ValidatePhase` invalid. |
| **33** | `b223f40` `feat: enforce active-player and phase validation for commands` | `internal/match/validate.go` (42 lines) `validate_test.go` (82 lines, 2 tests) | `ValidateActivePlayer(state,senderId)` `not_member`/`not_your_turn` via `SeatOfPlayer`/`CurrentSeat` for `OpeningDiscard`/`Playing` (`Waiting` not enforced), `ValidatePhaseOp(state,op)` `wrong_phase` via `AllowedOps`. Tests `alice/bob/carol` `CurrentSeat 1` `bob` active, `Waiting` no enforcement, phase matrix `Waiting start` vs `discard` etc. |
| **34** | `cf09f68` `feat: wire phase and payload validation into match loop with error dispatch` | `internal/match/rummy_match.go` `MatchLoop` (80→~150 lines) `rummy_match_test.go` updated `newMatchData` to `protocol.MustEnvelope` | `MatchLoop` now parses `ValidateEnvelope(data)` → maps `ParseError` to `ErrCodeBadJSON` etc., then `ValidateActivePlayer` + `ValidatePhaseOp` + `ValidatePayload`, all with `requestId`/`op` correlation, `sendError` via `dispatcher.BroadcastMessage(OpServerError, EncodeError)` to sender only and `continue` (atomic, no mutation). Only `OpClientStart` handled after validation. Tests `newMatchData` now envelope-based still pass `24` in `match`. `docker compose build` `0.7s`. |
| **35** | `b4e7cdb` `feat: handle opening discard from 15 to 14 and advance turn` | `internal/match/rummy_match.go` `handleOpeningDiscard` (70 lines) `opening_discard_test.go` (85 lines, 3 tests) | `OpClientDiscard` in `PhaseOpeningDiscard` validates `Seat 0` `2..4` `DiscardRow empty` `payload {tileId}` `rack len 15` `tile owned` removes preserving order, appends `DiscardEntry{IsOpeningDiscard:true,Index:0}`, `CurrentSeat=NextSeat(0,n)` `GamePhase Playing TurnMustDraw` broadcasts `OpServerEvent`. Tests `Success alice t-xxx →14 1 IsOpening CurrentSeat 1 Playing` + `CheckTileConservation` `106`, `OnlyCurrentPlayer bob rejected`, `TileMustBeOwned foreign/nonexistent`. `go test ./internal/match -v 30 passed`. |

*Days 5–7, 10–35 all verified via `make check` (`go vet` + `gofmt -l` + `go test ./...`) and `docker compose build`/`make smoke` where applicable. No game logic beyond slice’s scope was added.*

---

## Phase 5 — Turn state machine and basic actions (cont. Days 36–42)

| Day | Commit | Files | Goal / Acceptance |
|-----|--------|-------|-------------------|
| **36** | `9326b13` `feat: protect opening discard from pickup` | `internal/match/discard.go` (62 lines) `discard_test.go` | `CanPickupPreviousDiscard` / `CanPickupDiscardForMeld` reject `IsOpeningDiscard` and empty row; `TestProtectOpeningDiscard` verifies blocked tile never selectable. |
| **37** | `6505813` `refactor: extract anticlockwise turn advance helper` | `internal/match/turn.go` (24 lines) `turn_test.go` | `AdvanceTurn` `(current+1)%n` → `MustDraw` with `ValidatePlayers`/`CurrentSeat` checks; reused in opening discard; tests `NextSeat` loops `0→1→0` for 2/3/4 players. |
| **38** | `5a51ac6` `feat: handle draw from stock` | `internal/match/rummy_match.go` `handleDrawStock` (30 lines) `draw_test.go` (85 lines) | `OpClientDrawStock` in `Playing/MustDraw` pops `Stock` top, appends to `Racks[current]`, `TurnPhase→MeldOrDiscard`; `stock empty` → `bad_request`; tests `Success` (stock-1, rack+1, top `s2`), `WrongPhase`, `NotYourTurn`, conservation. |
| **39** | `47b153b` `feat: handle normal discard and turn rotation` | `internal/match/rummy_match.go` `handleNormalDiscard` (45 lines) `normal_discard_test.go` | `OpClientDiscard` in `Playing/MeldOrDiscard` validates ownership, appends `DiscardEntry{IsOpeningDiscard:false, Index=len}`, `AdvanceTurn`; tests `Success` (rack-1, discard+1, `CurrentSeat` 0→1), `BeforeDrawRejected`, `ForeignTile`, turn order 2/3/4. |
| **40** | `53ffe5e` `test: verify discard row ordering` | `internal/match/discard_order_test.go` | Ensures discard row preserves chronological order, `Index` tracks position, opening flagged distinct, conservation. |
| **41** | `69a3a61` `test: add turn-loop integration` | `internal/match/turn_loop_test.go` | Full loop `OpeningDiscard→DrawStock→Discard→next MustDraw` for 2/3/4 players, checks phase transitions and inventory. |
| **42** | `f21975a` `test: document and verify empty-stock MVP behavior` | `internal/match/empty_stock_test.go` | Documents `docs/rules-decisions.md:6.2` decision: `stock empty` → only discard pickup allowed else dead round; tests `DrawStockEmpty` error, `StockExhaustedNoWinner`, `PickupStillAllowed`. |

---

## Phase 6 — Meld rules: sets, runs, and jokers (Days 43–55)

| Day | Commit | Files | Goal / Acceptance |
|-----|--------|-------|-------------------|
| **43** | `04acd8d` `feat: define canonical meld representation` | `internal/rules/meld/meld.go` (104 lines) `meld_test.go` | `Meld{ID,Kind,Tiles,JokerReps}` with `Validate` (duplicate tile, missing rep, rep for non-joker), `New` copy-safe; stable ID preserved. |
| **44** | `415bbd0` `feat: validate 3- and 4-tile sets` | `internal/rules/meld/set.go` (partial) `set_test.go` | `ValidateSet` checks `3–4` tiles, same rank, distinct colours; tests valid 3/4, duplicate colour, rank mismatch. |
| **44b** | `94b99d8` `feat: return structured validation errors for sets` | `internal/rules/meld/errors.go` `set_errors_test.go` | `ValidationError{Code,Field,Message}` with codes `invalid_kind`/`invalid_size`/`rank_mismatch`/`duplicate_colour`; not bool-only. |
| **46** | `74b5118` `feat: support jokers in sets` | `internal/rules/meld/set.go` joker branch `set_joker_test.go` | Jokers via `JokerReps` explicit `rank==set rank` distinct colour, `real>=2*joker` (3-tiles max 1 joker, 4-tiles max 1); tests legal/illegal joker, ratio. |
| **48** | `24aec0d` `feat: validate basic same-colour consecutive runs` | `internal/rules/meld/run.go` (48 lines) `run_test.go` | `ValidateRun` same colour, `len>=3`, consecutive, no jokers yet; tests standard runs. |
| **49** | `004336c` `test: document low-Ace runs` | `run_lowace_test.go` | `1-2-3` valid, longer low `1-2-3-4`, tests. |
| **50** | `0a294ec` `feat: support high-Ace runs` | `run_highace_test.go` | `12-13-1` valid, `10-11-12-13-1` valid, while rejecting `13-1-2`; `contains` high-Ace as `14`. |
| **51** | `3adf97a` `test: explicitly reject Ace in middle` | `run_invalid_ace_test.go` | `13-1-2`, `12-13-1-2` etc. invalid, ensures Ace never middle. |
| **52** | `546d35e` `feat: support jokers in runs` | `run_joker_test.go` | Jokers with explicit `rep Colour==run colour` `Rank` filling gap, `real>=2*joker`, immutable rep. |
| **53** | `f86f519` `test: enforce real>=2*joker ratio for runs` | `run_ratio_test.go` | 6-tile max 2 jokers, 7-tile max 2, etc., ratio enforcement. |
| **54** | `e35e6b9` `test: ensure joker rep immutability` | `immutable_test.go` | Table `JokerReps` map is copied, cannot be silently reinterpreted (`12-13-J` declared as `1` stays `1`). |
| **55** | `f726c6f` `test: add comprehensive meld matrix` | `matrix_test.go` (100+ cases) | Valid/invalid sets, runs, Ace, joker, ratio, duplicate tiles — reusable by scoring. |

---

## Phase 7 — Opening meld and scoring (Days 56–66)

| Day | Commit | Files | Goal / Acceptance |
|-----|--------|-------|-------------------|
| **56** | `ba98df0` `feat: add tile scoring` | `internal/rules/scoring/scoring.go` `scoring_test.go` | `ScoreTile` `2–9:5`, `10–13:10`, Ace low `5` vs high `10` vs Ace-set `25`, Joker = represented; tests 5/10/25. |
| **57** | `1209e9b` `feat: score runs` | `scoring/run.go` `run_test.go` | `ScoreRun` validates run then sums with `isLowAceRun`/`isHighAceRun`; handles Joker. |
| **58** | `cc0c847` `feat: score sets` | `scoring/set.go` `set_test.go` | `ScoreSet` validates set then sums with `isAceSet` (25 each); Joker delegating. |
| **59** | `a4ed4b0` `test: verify joker scoring` | `joker_test.go` | Joker equals represented tile value in low/high/Ace contexts. |
| **60** | `b52a7f5` `feat: define opening meld batch model` | `scoring/batch.go` `batch_test.go` | `Batch{PlayerID,Melds}` `TotalScore` sums `ScoreRun/Set`; tests `30` and `90` (Ace set 75+low run 15). |
| **61** | `da5eb68` `feat: validate opening meld ownership` | `scoring/validate.go` `validate_test.go` | `ValidateBatchOwnership` checks each meld `Validate()`, `ValidateRun/Set`, all `tileIds` owned, no duplicate across melds; tests foreign tile, duplicate. |
| **62** | `5aaaa28` `feat: enforce 50-point minimum` | `scoring/validate_score_test.go` | `ValidateBatchScore` `total>=50` (30 rejected, 60 accepted), `exactly 50` passes. |
| **63** | `4c1e4e1` `feat: require at least one run` | `scoring/validate_has_run_test.go` | `ValidateBatchHasRun` requires ≥1 `KindRun` valid; sets-only 60 rejected. |
| **64** | `67aeebf` `test: ensure no duplicate tile across melds` | `validate_duplicate_test.go` | Duplicate `TileID` across `m1`/`m2` rejected; `ValidateInitialBatch` composed. |
| **65** | `df5b179` `feat: validate full initial meld batch` | `validate_initial_test.go` `validate_initial.go` | `ValidateInitialBatch` = ownership+score+hasRun; tests valid 50 (3 runs), 45 rejected, no-run 60 rejected, duplicate rejected. |
| **66/13** | `de0d727` `feat: allow initial table melds` | `internal/match/meld_initial.go` (250 lines) `meld_initial_test.go` (597 lines) `rummy_match.go:268` | `OpClientMeldInitial` in `MeldOrDiscard` (not opened) validates batch 50+ with run via `scoring`, atomically `Racks→TableMelds`, `HasOpened=true`, stays `MeldOrDiscard`; tests success 50pts, invalid atomic, cannot twice `already_opened`, still must discard, public redaction, duplicate tile, joker rep immutability, foreign tile, conservation 106. |
| **67/14** | `1b666c8` `feat: allow additional melds after opening` | `internal/match/meld_new.go` (158 lines) `meld_new_test.go` (547 lines) `rummy_match.go:276` | `OpClientMeldNew` in `MeldOrDiscard` (opened) validates via `ValidateBatchOwnership` (no score minimum), `meldId` not colliding with existing, atomic, stays `MeldOrDiscard`; tests single/multiple batch, unopened `not_opened`, invalid atomic, duplicate tile, IDs stable, joker, still must discard, wrong phase/notYourTurn. |

*Days 36–67 verified via `make check` and `docker compose build` per day; `Phase 6` milestone (reliable meld validation) and `Phase 7` milestone (50-point opening) complete.*

---

## Current State (after Day 67 / AGENTS Day 14)

- **Language:** Go `1.23.5` (`go.mod:3`) `nakama-common v1.36.0` `protobuf v1.36.4` (`Dockerfile` `pluginbuilder:3.26.0` → `rummy_backend:local`), `internal/` packages `rules/tile` `match` `setup` `protocol`.
- **Stack:** `docker compose up --build -d` → `rummy_postgres 5433` `rummy_nakama 7350/7351` `healthcheck` `Found runtime modules count 1 [rummy_backend.so]` `Registered Go RPC health/version` `Registered Match rummy`.
- **Gameplay implemented:** `Waiting→OpeningDiscard→Playing` (`MustDraw→MeldOrDiscard` loop), `DISCARD` (opening blocked + normal anticlockwise), `DRAW_STOCK`, validated `Run`/`Set` with jokers and Ace edge cases, scoring `5/10/25` with Joker delegation, `MELD_INITIAL` (50+ with run, atomic, `HasOpened`) and `MELD_NEW` (opened, no minimum, atomic, multiple per batch), `TableMeld` stable IDs/joker reps, tile conservation `106`, visibility redaction, standard `OpServerError` with `requestId` correlation.
- **Tests:** `go test ./...` green; deterministic seeded shuffle; `CheckTileConservation` after every state change; exhaustive redaction (`setup/redaction_test.go` 9 combos) and meld matrix (`rules/meld/matrix_test.go`).
- **Docs:** `docs/project-baseline.md` (Day 1 + Go §13), `docs/rules-decisions.md` (updated Day 14, `TODO(product)`), `docs/terminology.md` (59 lines), `README.md` (updated Phase 7–8, architecture, protocol, how to add command), `docs/protocol.md`, `docs/state-machine.md`, `docs/testing.md`, `AGENTS.md` source, `docs/daily-log.md` (this file).
- **CI:** `.github/workflows/ci.yml` `go vet` `gofmt` `go test` `go mod tidy` `docker compose build` on `push/PR main`.
- **Smoke:** `make smoke` / `scripts/smoke.sh` `SMOKE PASSED` (`pg_isready` `healthcheck` `InitModule` `rummy_backend.so` `console 200` `RPC health`).

## Next (Day 68 / AGENTS Day 15)

Per roadmap **Phase 8 Day 15 — Extend existing public melds**: `EXTEND_MELD` for opened players (own or others’ melds), client submits `targetMeldId` + rack `tileIds` (+ explicit resulting meld if needed), server revalidates entire resulting meld, atomic, `meldIds` stable, joker rep immutable per `docs/rules-decisions.md:3`. Then Days 16–19: discard pickups, joker replacement, round completion.

---

*Generated from `git log --oneline --reverse` `36c2c59..1b666c8` (67 commits) on `2026-08-26`. For `what changed` per day see `git show --stat <hash>`.*
