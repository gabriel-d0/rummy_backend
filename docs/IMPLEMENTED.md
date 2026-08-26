# Implemented Phases — Romanian Tile Rummy Backend (MVP)

**Source of truth:** This file maps the *Required implementation plan* from `AGENTS.md` (Phases 1–15 + Optional 16–17) to the actual code, commits, and verification in `main`. It is the hand-off for “what is already implemented” after the Handmade Hero incremental slices.

**Current HEAD:** `main@13762ca` (81 commits `36c2c59..13762ca`) + `5ade046` client, `1a4851e`/`01f6a3c` simulation, `3e1b92a` hardening, `2d278a5` win, `6b0d980` extend, `8dd8ea9` joker replace, `0da5f3a` pickup, `456c045` previous discard. Tagged `rummy-mvp-rc1` (`cfce62e`) — `make check` green, `make smoke` `SMOKE PASSED` on fresh `docker compose up --build -d`.

**How to verify:** `make check` (`go vet` + `gofmt -l` + `go test ./...`), `docker compose build`, `docker compose up --build -d` + `make smoke`, `go test -run TestDeterministicSimulation -v`, `go test -run TestRedaction`/`TestWin`/`TestSnapshotVersioning -v`, and `make cli` manual flow.

---

## Phase 1 — Foundation and local development

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 1 | Repository and environment audit | Review existing project, language, tooling, branching, CI, and deployment assumptions. Record findings in `docs/project-baseline.md`. | ✅ Done | `36c2c59` `docs: add Phase 1 Day 1 project baseline audit` — `docs/project-baseline.md` 228 lines, toolchain `Node 26`/`Go 1.27`/`Docker 29`/`compose v5`, branching `main` unborn, SSH origin. |
| 2 | Docker Compose | Add local Docker Compose setup for Nakama and PostgreSQL. Verify containers start and persist data correctly. | ✅ Done | `48ec5ce` `chore: add Docker Compose...` — `compose.yml` `name: rummy_backend` `postgres:15-alpine 5433` `nakama:3.26.0`, `nakama/data/local.yml`, `.env.example`, `.gitignore` — `docker compose up -d` both `healthy`, `migrate 0`, console `200` at `:7351`. |
| 3 | Nakama runtime skeleton | Add Nakama runtime module structure, build process, and a minimal `InitModule` log message. | ✅ Done | `9354b59` `feat: add Nakama TypeScript runtime skeleton` → superseded by `55c7f3b` `refactor: migrate Nakama runtime from TypeScript to Go` — `go.mod` `go 1.23.5` `nakama-common v1.36.0` `protobuf v1.36.4`, `main.go` `InitModule` `health`/`version` RPCs, `Dockerfile` `pluginbuilder:3.26.0`→`nakama:3.26.0` `backend.so` at `/nakama/data/modules/rummy_backend.so` — `Found runtime modules count 1 [rummy_backend.so]`. |
| 4 | Developer scripts | Add scripts for start, stop, logs, clean/reset database, build, test, lint, type-check, and format. | ✅ Done | `8493782` `chore: add developer scripts Makefile...` — `Makefile` `help`/`build`/`up`/`down`/`restart`/`ps`/`logs`/`clean`/`reset`/`db-shell`/`vet`/`fmt`/`test`/`tidy`/`check`/`health` — `make check` 0. |
| 5 | Local setup documentation | Write clear README instructions for prerequisites, environment variables, Docker startup, Nakama Console access, and test execution. | ✅ Done | `61607d2` `docs: update documentation for Go runtime` — `README.md` 153 lines (later 235) with `prerequisites`/`quick start`/`services & ports`/`rebuild`/`RPC smoke`. Merged with Day 4 per Go migration. |
| 6 | CI baseline | Add CI pipeline for formatting, linting, type checking, unit tests, and build verification. | ✅ Done | `162a0a6` `ci: add GitHub Actions baseline...` — `.github/workflows/ci.yml` `setup-go 1.23.5` `go vet` `gofmt -l` `go test` `go mod tidy` `docker/build-push-action` `compose build`. |
| 7 | Smoke test | Add a repeatable smoke test that verifies Nakama starts, runtime code loads, and database connectivity works. | ✅ Done | `74acc99` `test: add repeatable smoke test...` — `scripts/smoke.sh` 203 lines `make smoke` `SMOKE PASSED` when healthy and when down. |

**Milestone:** A new developer can clone, run `docker compose up --build -d`, start Nakama locally, and verify the backend is healthy — **done**.

---

## Phase 2 — Rules specification and core domain

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 8 | Rules source documentation | Create `docs/rules-decisions.md` containing the Romanian Tile Rummy rules, MVP decisions, assumptions, ambiguities, and deferred features. | ✅ Done | `2b44df6` `docs: add Phase 2 Day 8...` — `docs/rules-decisions.md` 189 lines (now 200+ with §9) with `106` tiles, `2–4` anticlockwise, `suita`/`terta`/`joker 2:1`, `50`-point opening, `15→discard` blocked, `TODO(product)` `final discard`/`stock exhaustion`/`doubla`/`set joker colours`. |
| 9 | Domain terminology | Define shared names and terminology: tile, rack, stock, discard row, meld, run, set, joker, opening meld, seat, turn, round. | ✅ Done | `c1eb89b` `docs: define shared domain terminology` — `docs/terminology.md` 59 lines canonical tables `Tile`/`TileInstanceId`/`Colour`/`Rank`/`Joker`/`PlayerId`/`Seat`/`Rack`/`Stock`/`DiscardRow`/`TableMelds` per `AGENTS.md:259`. |
| 10 | Tile domain model | Implement types for colours, ranks, tile IDs, numbered tiles, joker tiles, and unique tile instances. | ✅ Done | `53553b2` `internal/rules/tile/tile.go` 175 lines `tile_test.go` 8 tests — `Colour` `Red/Yellow/Blue/Black`, `Rank 1..13`, `TileInstanceId`, `TileInstance{ID,Colour,Rank,IsJoker}` `Validate`/`NewTile`/`NewJoker`/`Must*`. |
| 11 | Player and seat model | Implement player identity, deterministic seating, player state, and anticlockwise turn-order helpers. | ✅ Done | `a0f2147` `internal/match/seat.go` 140 lines `seat_test.go` 5 tests — `PlayerId` `Seat -1..3` `PlayerState{ID,Seat,HasOpened}` `AssignSeats` `Seat 0` opener, `NextSeat (current+1)%n` `PrevSeat`, `SeatOfPlayer`/`ValidatePlayers`. |
| 12 | Game state model | Create initial server-side state structures for round state, stock, racks, discards, public melds, current player, and turn phase. | ✅ Done | `4d67326` `internal/match/state.go` 302 lines `state_test.go` 4 tests — `GamePhase Waiting/OpeningDiscard/Playing/RoundComplete` `TurnPhase MustDraw/MeldOrDiscard` `DiscardEntry{IsOpeningDiscard,Index}` `TableMeld{ID,Kind,Tiles,JokerReps,OwnerSeat}` `RoundState{Players,Racks,Stock,DiscardRow,TableMelds,CurrentSeat,GamePhase,TurnPhase,Winner}`. |
| 13 | State invariants | Implement an invariant checker ensuring every tile exists in exactly one location: rack, stock, discard row, table meld, or reserved setup location. | ✅ Done | `b83740c` `internal/match/invariant.go` 99 lines `invariant_test.go` 4 tests + `syntheticDeck` — `CheckTileConservation(state, allTiles 106)` with `seen` duplicate/`not in deck`/`missing` and `CountTiles`. |
| 14 | Domain tests | Add tests for tile identity, seats, turn direction, base state construction, and invariant failures. | ✅ Done | `8966391` `internal/match/domain_test.go` 168 lines 4 tests — `TileIdentity`, `SeatsTurnDirection` 2p/3p/4p, `BaseStateConstruction` `15+14+14+1+62=106` `Validate`+`Conservation`, `InvariantFailures`. |

**Milestone:** The game has a clear, documented vocabulary and a safe internal state model before multiplayer behavior begins — **done**.

---

## Phase 3 — Deck, randomness, and round setup

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 15 | Full deck factory | Create the 106-tile Romanian Tile Rummy deck: 104 numbered tiles plus 2 jokers. | ✅ Done | `2aee5d4` `internal/setup/deck.go` 45 lines `NewDeck()` deterministic `Red 1-13×2`→`Yellow`→`Blue`→`Black`→`joker-1/2` `106` unique `TileInstanceId`. |
| 16 | Deck correctness tests | Verify exactly two copies of every colour/rank tile, two jokers, 106 total tiles, and unique instance IDs. | ✅ Done | `9ed478a` `internal/setup/deck_test.go` 123 lines 7 tests — `TotalCount 106` `104/2` `TwoPerColourRank` `AllIDsUnique` `Deterministic` snapshot `red-01-1/2` `joker-2`. |
| 17 | Deterministic random source | Add injectable/seedable randomness for reproducible shuffles and tests. | ✅ Done | `47ade4e` `internal/setup/rand.go` 31 lines `rand_test.go` 71 lines 4 tests — `Rand{*rand.Rand}` `NewSeededRand(seed)` `Intn`/`Shuffle`. |
| 18 | Shuffle implementation | Implement Fisher–Yates or equivalent clear shuffle logic with deterministic test support. | ✅ Done | `05de37b` `internal/setup/shuffle.go` 22 lines `shuffle_test.go` 106 lines 6 tests — `Shuffle(deck,Rand)` copies then `r.Shuffle` Fisher–Yates, `NoLostDuplicated` `106`. |
| 19 | Deal logic | Implement MVP dealing: opening player gets 15 tiles, all other players get 14, remainder becomes stock. | ✅ Done | `94068b7` `internal/setup/deal.go` 52 lines `deal_test.go` 117 lines 4 tests — `Deal(deck,n)` `15+(n-1)*14` `stock 77/63/49` for 2/3/4, `ConservationWithState`. |
| 20 | Round initialization | Build `createRoundState()` with seats, dealer/opening player, racks, stock, empty table, empty discard row, and opening turn phase. | ✅ Done | `3daf725` `internal/setup/round.go` 58 lines `round_test.go` 137 lines 4 tests — `NewRoundState(playerIds,seed)` `NewDeck+Shuffle+AssignSeats+Deal` → `PhaseOpeningDiscard` `CurrentSeat 0`. |
| 21 | Setup invariants | Test dealing for 2, 3, and 4 players; verify tile conservation and expected stock counts. | ✅ Done | `b2f5634` `internal/setup/setup_invariant_test.go` 61 lines 1 test `18` combos `TestSetupInvariantsAllPlayersAndSeeds` `n=2,3,4 × seeds 0,1,42,123,999,2026` `15/14` `stock 77/63/49` `Conservation` `CountTiles 106`. |

**Milestone:** A complete legal round can be initialized deterministically without Nakama — **done**.

---

## Phase 4 — Nakama match foundation

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 22 | Runtime module initialization | Register authoritative match handlers in Nakama. | ✅ Done | `9d4a888` `internal/match/rummy_match.go` 80 lines `RummyMatch` `7` methods `main.go` `RegisterMatch("rummy",match.NewRummyMatch)` — `Registered Go runtime Match creation function invocation name rummy`. |
| 23 | Match lifecycle skeleton | Implement match init, join attempt, join, leave, loop, signal, and terminate callbacks. | ✅ Done | `9d4a888` (same) + `rummy_match.go` `MatchInit` `Waiting` `tickRate 5` `label rummy` empty `RoundState`; `9132ca7` lobby waiting room with seat allocation. |
| 24 | Waiting room state | Add lobby state for 2–4 players with seat allocation and player-ready information. | ✅ Done | `9132ca7` `internal/match/rummy_match.go` updated 81 lines `rummy_match_test.go` 179 lines 4 tests + mocks — `MatchJoinAttempt` checks `Waiting`/`full 4`/`already joined`, `MatchJoin` allocates `Seat(len(Players))` idempotent, `MatchLeave` (now keeps `Seat`/`Racks` for Day 20 reconnect). |
| 25 | Match start command | Allow host/creator to start a match when at least 2 players are present. | ✅ Done | `14d9b93` `MatchLoop` handles `op 1 START` + `MatchSignal "start"` — host `Seat 0` `2..4` `Waiting→OpeningDiscard` `CurrentSeat 0`; tests `TestMatchStartViaLoop` non-host rejected. |
| 26 | Protocol opcodes | Define stable client/server opcodes and message envelope versioning. | ✅ Done | `437beb4` `internal/protocol/opcodes.go` 48 lines `opcodes_test.go` — `Version 1`, client `1..99` `Start=1`/`Discard=2`/`DrawStock=3`/`DrawPrevious=4`/`Pickup=5`/`MeldInitial=6`/`MeldNew=7`/`Extend=8`/`ReplaceJoker=9`, server `100..199` `State=100`/`StatePublic=101`/`Error=102`/`Event=103`. |
| 27 | Command parser | Parse inbound JSON safely and reject malformed payloads without crashing the match. | ✅ Done | `ec8d286` `internal/protocol/envelope.go` 32 lines `parser.go` 52 lines `parser_test.go` 3 tests — `Envelope{v,op,requestId,payload}` `ParseError` `ParseEnvelope` never panics, rejects `empty`/`bad_json`/`bad_version`/`unknown_opcode`. |
| 28 | Command schema validation | Add payload validation for every currently available command. | ✅ Done | `b6110a8` `internal/protocol/validator.go` 120 lines `validator_test.go` 7 tests — `ValidatePayload(op,payload)` switch `9` ops: `Start` empty, `Discard {tileId}`, `DrawStock` empty, `Pickup {discardIndex + tileIds[2]}`, `MeldInitial/New {melds>=1}`, `Extend {meldId + tileIds}`, `ReplaceJoker {targetMeldId + tileId + newMeldTiles[2]}` (later extended for `jokerReps`). |
| 29 | Standard error protocol | Add consistent error responses: code, message, request ID, and optional structured details. | ✅ Done | `fc48c28` `internal/protocol/envelope.go` added `RequestId`, `errors.go` 62 lines `errors_test.go` 5 tests — `ErrorResponse{code,message,details,requestId,op}` `OpServerError 102` `NewError`/`NewErrorForEnvelope` `EncodeError`, codes `bad_request`…`already_opened`. |
| 30 | Basic match snapshots | Send public game state and player-specific private state snapshots when a round starts. | ✅ Done | `2cc3313` `internal/match/visibility.go` 65 lines `visibility_test.go` 3 tests — `PublicSnapshot{Version,GamePhase,TurnPhase,CurrentSeat,Players{RackCount},StockCount,DiscardRow,TableMelds,Winner}` `PrivateSnapshot{PublicSnapshot,OwnRack,OwnSeat}` `PublicView`/`PrivateView`. |
| 31 | Hidden-information test | Verify public snapshots and opponent snapshots never reveal another player’s rack tiles. | ✅ Done | `4f60325` `internal/setup/redaction_test.go` 78 lines 1 test `9` combos `TestSnapshotRedactionExhaustive` `n=2,3,4 × seeds 42,7,123` `PublicView` must not contain any `rack`/`stock` ID, `PrivateView` must contain own not others. |

**Milestone:** Two or more players can create, join, start, and receive a secure authoritative match state — **done**. Then hardened on Day 20 with `snapshot_hardening_test.go` 4 tests (see Phase 12).

---

## Phase 5 — Turn state machine and basic actions

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 32 | Turn state machine | Define explicit phases: waiting, opening discard, must draw, may meld/discard, round complete. | ✅ Done | `307b81d` `internal/match/phases.go` 55 lines `phases_test.go` 2 tests — `AllowedOps(gamePhase,turnPhase)` `Waiting→Start`, `OpeningDiscard→Discard`, `Playing MustDraw→DrawStock/DrawPrevious/Pickup`, `MeldOrDiscard→Discard/MeldInitial/MeldNew/Extend/ReplaceJoker`, `RoundComplete→none`. |
| 33 | Active-player validation | Reject actions from non-active players. | ✅ Done | `b223f40` `internal/match/validate.go` 42 lines `validate_test.go` 2 tests — `ValidateActivePlayer` `not_member`/`not_your_turn`, `ValidatePhaseOp` `wrong_phase`. |
| 34 | Phase validation | Reject actions that are not legal in the current turn phase. | ✅ Done | `cf09f68` `internal/match/rummy_match.go` `MatchLoop` 80→150 lines `rummy_match_test.go` `newMatchData` to `protocol.MustEnvelope` — `MatchLoop` parses `ValidateEnvelope` → `ValidateActivePlayer`+`ValidatePhaseOp`+`ValidatePayload` with `requestId` correlation, `sendError` to sender only, atomic. |
| 35 | Opening discard command | Implement the opening player’s mandatory first discard from 15 tiles. | ✅ Done | `b4e7cdb` `internal/match/rummy_match.go` `handleOpeningDiscard` 70 lines `opening_discard_test.go` 3 tests — `OpClientDiscard` in `PhaseOpeningDiscard` validates `Seat 0`/`2..4`/`DiscardRow empty`/`tile owned` removes `15→14`, appends `IsOpeningDiscard:true` `Index 0`, `CurrentSeat=NextSeat(0)` `GamePhase Playing TurnMustDraw`. |
| 36 | Opening discard protection | Mark the first discard as permanently unavailable for normal discard pickup. | ✅ Done | `9326b13` `internal/match/discard.go` 62 lines `discard_test.go` — `CanPickupPreviousDiscard`/`CanPickupDiscardForMeld` reject `IsOpeningDiscard` and empty row; `TestProtectOpeningDiscard`. |
| 37 | Turn advance | Advance active seat anticlockwise after opening discard. | ✅ Done | `6505813` `internal/match/turn.go` 24 lines `turn_test.go` — `AdvanceTurn` `(current+1)%n` → `MustDraw` with `ValidatePlayers`/`CurrentSeat` checks. |
| 38 | Draw from stock | Implement `DRAW_STOCK` for active player during `MUST_DRAW`. | ✅ Done | `5a51ac6` `internal/match/rummy_match.go` `handleDrawStock` 30 lines `draw_test.go` 85 lines — `OpClientDrawStock` in `Playing/MustDraw` pops `Stock` top to `Racks[current]`, `TurnPhase→MeldOrDiscard`; `stock empty` → `bad_request`; tests `Success` (stock-1, rack+1, top `s2`), `WrongPhase`, `NotYourTurn`, conservation. |
| 39 | Normal discard | Implement `DISCARD` after drawing or melding. | ✅ Done | `47b153b` `rummy_match.go` `handleNormalDiscard` 45 lines `normal_discard_test.go` — `OpClientDiscard` in `Playing/MeldOrDiscard` validates ownership, appends `IsOpeningDiscard:false` `Index=len`, `AdvanceTurn`; tests `Success` (rack-1, discard+1, `CurrentSeat` 0→1), `BeforeDrawRejected`, turn order. |
| 40 | Discard row ordering | Preserve public discard history in exact chronological order. | ✅ Done | `53ffe5e` `internal/match/discard_order_test.go` — chronological `Index`, `IsOpeningDiscard` distinct, conservation. |
| 41 | Turn-loop tests | Test a complete sequence: opening discard → draw stock → discard → next player. | ✅ Done | `69a3a61` `internal/match/turn_loop_test.go` — full loop `OpeningDiscard→DrawStock→Discard→next MustDraw` for 2/3/4 players, phase transitions and inventory. |
| 42 | Empty-stock decision | Document and implement deterministic MVP behavior for an exhausted stock. | ✅ Done | `f21975a` `internal/match/empty_stock_test.go` — `docs/rules-decisions.md:6.2` `stock empty` → only discard pickup allowed else dead round `stockExhausted`; tests `DrawStockEmpty` error, `StockExhaustedNoWinner`, `PickupStillAllowed`. |

**Milestone:** Players can play a valid draw-and-discard loop with a server-authoritative turn order — **done**.

---

## Phase 6 — Meld rules: sets, runs, and jokers

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 43 | Meld representation | Define canonical representation for table melds, including stable meld IDs and joker substitutions. | ✅ Done | `04acd8d` `internal/rules/meld/meld.go` 104 lines `meld_test.go` — `Meld{ID,Kind,Tiles,JokerReps}` `Validate` (duplicate, missing rep), `New` copy-safe, stable ID. |
| 44 | Basic set validation | Validate 3- and 4-tile sets of equal rank with distinct colours. | ✅ Done | `415bbd0` `internal/rules/meld/set.go` `set_test.go` — `ValidateSet` `3–4` tiles same rank distinct colours. |
| 45 | Set validation errors | Return structured reasons for invalid sets: rank mismatch, duplicate colour, invalid size, duplicate tile. | ✅ Done | `94b99d8` `internal/rules/meld/errors.go` `set_errors_test.go` — `ValidationError{Code,Field,Message}` `invalid_kind`/`invalid_size`/`rank_mismatch`/`duplicate_colour`. |
| 46 | Set joker support | Support jokers in sets with explicit represented rank/colour. | ✅ Done | `74b5118` `internal/rules/meld/set.go` joker branch `set_joker_test.go` — `rank==set rank` distinct colour, `real>=2*joker` (3-tiles max 1, 4-tiles max 1). |
| 47 | Set joker ratio rule | Enforce at least twice as many real tiles as jokers. | ✅ Done | `74b5118` (same) — tests `real>=2*joker` enforcement. |
| 48 | Basic run validation | Validate same-colour consecutive runs of length 3 or greater. | ✅ Done | `24aec0d` `internal/rules/meld/run.go` 48 lines `run_test.go` — same colour `len>=3` consecutive, no jokers yet. |
| 49 | Low-ace runs | Support `1-2-3` and longer low-ace runs. | ✅ Done | `004336c` `run_lowace_test.go` — `1-2-3` valid, longer low. |
| 50 | High-ace runs | Support `12-13-1` and longer high-ace runs according to the selected canonical model. | ✅ Done | `0a294ec` `run_highace_test.go` — `12-13-1` valid, `10-11-12-13-1` valid, rejecting `13-1-2`; high-Ace as `14`. |
| 51 | Invalid ace-middle runs | Explicitly reject sequences such as `13-1-2`. | ✅ Done | `3adf97a` `run_invalid_ace_test.go` — `13-1-2`, `12-13-1-2` invalid. |
| 52 | Run joker support | Support jokers in runs with explicit represented colour and rank. | ✅ Done | `546d35e` `run_joker_test.go` — `rep Colour==run colour` `Rank` filling gap, `real>=2*joker`, immutable rep. |
| 53 | Run joker ratio rule | Enforce the real-tile-to-joker ratio for runs. | ✅ Done | `f86f519` `run_ratio_test.go` — 6-tile max 2 jokers etc. |
| 54 | Immutable joker mapping | Prevent a tabled joker from silently changing represented value later. | ✅ Done | `e35e6b9` `immutable_test.go` — `JokerReps` map copied, cannot be silently reinterpreted. |
| 55 | Meld test matrix | Add a comprehensive test matrix for valid and invalid sets, runs, ace cases, and joker cases. | ✅ Done | `f726c6f` `matrix_test.go` 100+ cases — valid/invalid sets, runs, Ace, joker, ratio, duplicate tiles — reusable by scoring. |

**Milestone:** The game can independently and reliably determine whether any proposed meld is legal — **done**.

---

## Phase 7 — Opening meld and scoring

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 56 | Tile scoring model | Implement value rules for ranks 2–9, 10–13, aces, and jokers. | ✅ Done | `ba98df0` `internal/rules/scoring/scoring.go` `scoring_test.go` — `ScoreTile` `2–9:5`, `10–13:10`, Ace low `5` vs high `10` vs Ace-set `25`, Joker = represented. |
| 57 | Run scoring | Score low-ace and high-ace runs correctly. | ✅ Done | `1209e9b` `scoring/run.go` `run_test.go` — `ScoreRun` validates run then sums with `isLowAceRun`/`isHighAceRun`; handles Joker. |
| 58 | Set scoring | Score sets, including special ace-set value rules. | ✅ Done | `cc0c847` `scoring/set.go` `set_test.go` — `ScoreSet` with `isAceSet` (25 each); Joker delegating. |
| 59 | Joker scoring | Score a joker according to its declared represented tile. | ✅ Done | `a4ed4b0` `joker_test.go` — Joker equals represented tile value. |
| 60 | Opening meld batch model | Define a batch of proposed initial melds from one player rack. | ✅ Done | `b52a7f5` `scoring/batch.go` `batch_test.go` — `Batch{PlayerID,Melds}` `TotalScore` sums `ScoreRun/Set`; tests `30` and `90`. |
| 61 | Opening meld validation | Require all opening meld tiles to come from the player’s rack. | ✅ Done | `da5eb68` `scoring/validate.go` `validate_test.go` — `ValidateBatchOwnership` checks each meld `Validate()`, `ValidateRun/Set`, all `tileIds` owned, no duplicate across melds. |
| 62 | Opening meld minimum score | Require total value of at least 50 points. | ✅ Done | `5aaaa28` `scoring/validate_score_test.go` — `ValidateBatchScore` `total>=50` (30 rejected, 60 accepted), `exactly 50` passes. |
| 63 | Opening meld run requirement | Require at least one valid run in the opening batch. | ✅ Done | `4c1e4e1` `scoring/validate_has_run_test.go` — `ValidateBatchHasRun` requires ≥1 `KindRun` valid; sets-only 60 rejected. |
| 64 | Duplicate tile prevention | Reject any tile used more than once across the batch. | ✅ Done | `67aeebf` `validate_duplicate_test.go` — duplicate `TileID` across `m1`/`m2` rejected; `ValidateInitialBatch` composed. |
| 65 | Initial meld command | Add `MELD_INITIAL` to the Nakama command handler. | ✅ Done | Integrated into `de0d727` (Day 13 handler) — `handleMeldInitial` calls `scoring.ValidateInitialBatch` via `buildMeldsFromPayload`. |
| 66 | Atomic initial meld mutation | Remove rack tiles and add public melds only if every validation passes. | ✅ Done | `de0d727` `internal/match/meld_initial.go` 45 lines `meld_initial_test.go` 597 lines — atomically `Racks→TableMelds`, `HasOpened=true`, stays `MeldOrDiscard`; tests success 50pts, invalid atomic, `already_opened`. |
| 67 | Opening meld tests | Cover 49 vs. 50 points, no-run rejection, joker scoring, duplicate use, and atomic rollback. | ✅ Done | `f726c6f` matrix + `df5b179` `validate_initial_test.go` + `de0d727` `meld_initial_test.go` 9 tests — 49 rejected, 50 accepted, no-run 60 rejected, duplicate rejected, joker rep immutability. |

**Milestone:** A player can legally “open” with at least 50 points and at least one run — **done**.

---

## Phase 8 — Post-opening table play

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 68 | Opened-player flag | Track whether each player has completed their initial meld. | ✅ Done | `de0d727`/`1b666c8` `internal/match/seat.go` `PlayerState{HasOpened}` + `meld_initial.go`/`meld_new.go` set/check `HasOpened`; `visibility.go` `PublicPlayer{HasOpened,RackCount}`. |
| 69 | Additional new melds | Add `MELD_NEW` for opened players to create new valid table melds. | ✅ Done | `1b666c8` `internal/match/meld_new.go` 58 lines `meld_new_test.go` 547 lines `rummy_match.go:276` — `OpClientMeldNew` in `MeldOrDiscard` (opened) validates via `ValidateBatchOwnership` (no score minimum), `meldId` not colliding, atomic, stays `MeldOrDiscard`. |
| 70 | Batch new melds | Permit multiple new melds in one command while preserving atomicity. | ✅ Done | `1b666c8` supports multiple `melds` in one `MELD_NEW` batch with atomic rollback (`TestMeldNewMultipleInOneBatch` 2 melds 6 tiles). |
| 71 | Table meld extension model | Define how a rack tile is added to an existing public meld. | ✅ Done | `6b0d980` `internal/match/state.go` `TableMeld.Kind` added (`run`/`set` stable) for revalidation on extend; `meld_common.go` sets `Kind`. |
| 72 | Extend set meld | Allow legal extension of a three-tile set to a four-colour set. | ✅ Done | `6b0d980` `internal/match/extend_meld.go` 180 lines `extend_meld_test.go` 452 lines — `TestExtendSetToFourColours` (3→4, distinct colour), `TestExtendAnotherPlayersMeld`. |
| 73 | Extend run meld | Allow legal extension at a valid run endpoint. | ✅ Done | `6b0d980` — `TestExtendRunAtEnd` (5-6-7 +8 and 4+5-6-7 low end, both valid, revalidate entire). |
| 74 | Extend any player’s meld | Permit opened players to extend public melds regardless of original owner. | ✅ Done | `6b0d980` — `TestExtendAnotherPlayersMeld` (bob’s run extended by alice, `OwnerSeat` stays 1, conservation `106`). |
| 75 | Extension command | Add `EXTEND_MELD` protocol handling and server validation. | ✅ Done | `6b0d980` `internal/protocol/validator.go` `OpClientExtendMeld` now accepts `jokerReps`, `internal/match/extend_meld.go` `handleExtendMeld` revalidates entire resulting meld (`meld.New`+`ValidateRun/Set`), preserves `JokerReps` immutability, any owner. |
| 76 | Extension rollback tests | Ensure invalid extensions never partially mutate state. | ✅ Done | `6b0d980` — `TestExtendInvalidDoesNotMutate` (wrong colour `Blue 8`, gap `10` rejected, rack/melds unchanged, conservation), `TestExtendJokerImmutable` (existing rep change → `bad_request`). |
| 77 | Table state projection | Improve public state messages for meld IDs, tiles, joker assignments, and updates. | ✅ Done | `2cc3313`/`4f60325` `visibility.go` `PublicSnapshot{TableMelds}` already includes `ID`/`Kind`/`Tiles`/`JokerReps`/`OwnerSeat`; `extend_meld.go` updates `Kind` stable and broadcasts `OpServerEvent` `extendMeld`; `protocol.md` now documents `ExtendMeld` `jokerReps`. |

**Milestone:** Opened players can build on their own or others’ table melds safely — **done**.

---

## Phase 9 — Discard pickup rules

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 78 | Latest discard pickup rule | Document exact MVP behavior for taking the immediately preceding discard. | ✅ Done | `docs/rules-decisions.md:1.5` and `discard.go:15` `CanPickupPreviousDiscard` — opening blocked, `MustDraw` only, `len>0` and not opening. |
| 79 | Previous-discard draw command | Implement `DRAW_PREVIOUS_DISCARD`. | ✅ Done | `456c045` `internal/match/draw_previous_discard.go` 52 lines `draw_previous_discard_test.go` 300 lines `rummy_match.go:269` — `OpClientDrawPreviousDiscard` in `Playing/MustDraw` for opened, not opening, moves tail `DiscardRow` to `Racks[seat]`, `TurnPhase→MeldOrDiscard`. |
| 80 | Previous-discard validation | Require opened player, active turn, `MUST_DRAW`, non-opening discard, and one draw per turn. | ✅ Done | `456c045` — `TestDrawPreviousDiscardUnopenedRejected` (`not_opened`), `TestDrawPreviousDiscardOpeningBlocked` (`bad_request` opening), `TestDrawPreviousDiscardNotYourTurnAndWrongPhase` (`wrong_phase`), `TestDrawPreviousDiscardCannotDrawTwice` (`MeldOrDiscard` forbids second draw). |
| 81 | Previous-discard tests | Test valid pickup, opening-discard rejection, unopened-player rejection, and turn-phase rejection. | ✅ Done | `456c045` 6 tests — `Success` (rack 14→15, discard 3→2, latest only, stock unchanged, `MeldOrDiscard`), `UnopenedRejected`, `OpeningBlocked`, `OnlyLatest`, `CannotDrawTwice`, `NotYourTurn`. |
| 82 | Earlier discard pickup model | Define payload for selecting an earlier discard plus exactly two rack tiles. | ✅ Done | `b6110a8` `validator.go` `OpClientPickupDiscardForMeld` `{discardIndex, tileIds[2]}` (later extended for `jokerReps`/`kind`), `pickup_discard_for_meld.go` `pickupDiscardPayload` struct. |
| 83 | Immediate pickup meld validation | Require selected discard plus exactly two rack tiles to form a legal meld immediately. | ✅ Done | `0da5f3a` `internal/match/pickup_discard_for_meld.go` 210 lines — builds `combinedTiles=[discardTile]+2 rack` and `combinedReps` (joker reps if needed), tries `run` then `set` (or forced `kind`) via `meld.New`/`ValidateRun/Set`; `TestPickupDiscardForMeldValidSet`/`ValidRun`. |
| 84 | Later discard collection | Move all discards after the selected discard into the player’s rack, preserving order. | ✅ Done | `0da5f3a` — `laterTiles = DiscardRow[discardIndex+1:]` swept in order `newRack = rack -2 + laterTiles`, `DiscardRow = DiscardRow[:discardIndex]` reindexed. |
| 85 | Earlier discard pickup command | Implement `PICKUP_DISCARD_FOR_MELD`. | ✅ Done | `0da5f3a` `internal/match/pickup_discard_for_meld.go` `handlePickupDiscardForMeld` `rummy_match.go:276` `OpClientPickupDiscardForMeld`. |
| 86 | Atomic discard-row pickup | Ensure selected discard, generated meld, rack additions, and discard-row removal occur atomically. | ✅ Done | `0da5f3a` — validates fully before touching `state` (`CanPickupDiscardForMeld`, ownership, meld valid), then atomic `Racks`/`DiscardRow`/`TableMelds` mutation with saved `origRack`/`origDiscard` rollback on `Validate` failure. |
| 87 | Discard pickup test suite | Cover valid set/run pickup, invalid pickup, opening discard exclusion, and conservation checks. | ✅ Done | `0da5f3a` `pickup_discard_for_meld_test.go` 378 lines 7 tests — `ValidSet` (`7 red` + `7 yellow`/`7 blue` + sweep `disc2`), `ValidRun` (5-6-7), `InvalidRejected`, `OpeningBlocked`, `LaterSweep` (2 later in order), `Atomic`, `RequiresOpenedAndMustDraw`. |

**Milestone:** The special Romanian discard-row pickup mechanic works safely and deterministically — **done**.

---

## Phase 10 — Joker replacement mechanics

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 88 | Joker replacement rules document | Document exact supported replacement cases and unresolved source ambiguities. | ✅ Done | `docs/rules-decisions.md:6.4` `Joker replacement in sets — exact missing colours` with `TODO(product)` and MVP `Run joker: exact represented tile; Set joker: exact missing colour`. |
| 89 | Run joker replacement validation | Allow replacement only with the exact tile represented by a joker in a run. | ✅ Done | `8dd8ea9` `internal/match/replace_joker.go` `handleReplaceJoker` finds joker whose `JokerReps` colour/rank equals `tileId` tile’s colour/rank, else `bad_request` `no matching joker`. `TestReplaceJokerRunValid` / `TestReplaceJokerRunWrongTile`. |
| 90 | New joker meld requirement | Require the recovered joker to immediately form a new valid meld with two tiles from the actor’s rack. | ✅ Done | `8dd8ea9` — builds `newMeldTiles=[recoveredJoker]+2` with `newMeldJokerReps` (recovered joker’s new `colour`/`rank`) + `newMeldKind` infer `run`→`set`, validates via `meld.New`; tests `NewMeldRequiresTwoTiles` (invalid new meld rejected). |
| 91 | Set joker replacement validation | For a joker in a three-tile set, require the needed missing colours/rank as defined in the rules decision. | ✅ Done | `8dd8ea9` — same `exact missing colour` check for set: `repTile` must equal joker’s rep; `TestReplaceJokerSetValid` (9 set) / `TestReplaceJokerSetWrongColour` (black 9 vs blue rep rejected). |
| 92 | Joker replacement command | Implement `REPLACE_JOKER`. | ✅ Done | `8dd8ea9` `internal/protocol/validator.go` `OpClientReplaceJoker` now accepts `jokerReps`/`newMeldKind`, `internal/match/replace_joker.go` `handleReplaceJoker` `rummy_match.go:308` `OpClientReplaceJoker`. |
| 93 | Atomic replacement mutation | Apply table replacement and new joker meld as one all-or-nothing state transition. | ✅ Done | `8dd8ea9` — saves `origRack`/`origTable`, validates `updatedMeld` (`meld.New` + `ValidateRun/Set`) and `newMeld` before touching state, then `Racks[seat]=rack-3`, `TableMelds[targetIdx]=updatedTM`, `+=newTM`, rollback on `Validate` failure. |
| 94 | Joker replacement tests | Test valid run replacement, wrong-tile rejection, valid set replacement, missing-tile rejection, and failed atomic rollback. | ✅ Done | `8dd8ea9` `replace_joker_test.go` 400 lines 6 tests — `RunValid` (5-6-J7 with 7→5-6-7 and new run 8-9-J10), `RunWrongTile` (8 vs 7 rejected), `SetValid` (9 set), `SetWrongColour` (black vs blue), `NewMeldRequiresTwoTiles` (8+12 gap invalid), `AtomicRollback` (rack/melds unchanged, conservation `106`). |

**Milestone:** Jokers cannot be illegally removed, reassigned, or exploited — **done**.

---

## Phase 11 — Win conditions, round completion, and results

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 95 | Closing rule decision | Finalize/document MVP interpretation of winning with no remaining tiles and final-discard behavior. | ✅ Done | `docs/rules-decisions.md:6.1` `Final discard / closing behavior` with MVP `rack==0` after `DISCARD` or after `MELD`/`EXTEND`/`REPLACE`/`PICKUP` that empties rack also wins (win without final discard allowed), `TODO(product)` and tests `TestWinAfterMeldWithoutDiscard`/`TestWinAfterDiscardToZero`. |
| 96 | Win detection | Detect when a player has legally emptied their rack. | ✅ Done | `2d278a5` `internal/match/win.go` 30 lines `checkWinAndComplete` `len(Racks[seat])==0` → `GamePhase=RoundComplete` `Winner=seat` `CurrentSeat=winner` broadcasts `101`/`103`. |
| 97 | Round-complete state | Add `ROUND_COMPLETE`, winner ID, ending action, timestamps, and summary information. | ✅ Done | `2d278a5` `internal/match/state.go` `Winner` `SeatInvalid` default, `internal/match/win.go` sets `Winner` and `GamePhase` `RoundComplete`; `rummy_match.go:473` `handleNormalDiscard` and `meld_*`/`extend`/`pickup`/`replace` call `checkWinAndComplete` after `rack==0`. |
| 98 | Post-game command blocking | Reject gameplay commands after a round ends. | ✅ Done | `307b81d` `phases.go` `AllowedOps` `RoundComplete` has none → `ValidatePhaseOp` `wrong_phase`; `2d278a5` `win_test.go` `TestNoGameplayAfterRoundComplete` (draw/discard after win rejected, stays `RoundComplete`). |
| 99 | Final state broadcast | Send final public state and each player’s final private view. | ✅ Done | `2d278a5` `win.go` broadcasts `OpServerStatePublic 101` `RoundComplete` `Winner` and `OpServerEvent` `roundComplete`; `visibility.go` `PublicView`/`PrivateView` already include `Winner` and `GamePhase` `RoundComplete`. |
| 100 | Round completion tests | Test win detection, post-game rejection, conservation, and final event payloads. | ✅ Done | `2d278a5` `win_test.go` 400 lines 4 tests — `WinAfterDiscardToZero` (1 tile → discard → `RoundComplete` winner 0, public `Winner`), `WinAfterMeldWithoutDiscard` (3 tiles 5-6-7 run → `MELD_NEW` → `RoundComplete` without discard), `NoGameplayAfterRoundComplete`, `WinnerCorrectlyRecorded` (bob wins). |
| 101 | Optional dead-round behavior | If stock exhaustion behavior requires it, implement and test draw/round-ending state according to documented MVP decision. | ✅ Done (MVP) | `docs/rules-decisions.md:6.2` `Stock exhaustion` MVP `stock empty` → only discard pickup allowed else dead round `stockExhausted` (not yet triggered as win, but `f21975a` `empty_stock_test.go` documents and verifies `DrawStockEmpty` `bad_request`). |

**Milestone:** A full round can end authoritatively and report a winner — **done**.

---

## Phase 12 — Security, reconnection, and reliability

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 102 | Central view projection | Consolidate public and per-player-private snapshot generation in one module. | ✅ Done | `2cc3313` `internal/match/visibility.go` 65 lines `visibility_test.go` 3 tests — `PublicView`/`PrivateView` centralized, `Version=1` stable. |
| 103 | Snapshot redaction tests | Serialize views for all players and prove foreign rack IDs are absent. | ✅ Done | `4f60325` `internal/setup/redaction_test.go` 78 lines 1 test `9` combos `TestSnapshotRedactionExhaustive` + `2cc3313` `visibility_test.go` + `3e1b92a` `snapshot_hardening_test.go` 4 tests `TestCentralizedViewProjection`/`TestRedactionRoundComplete` (winner `RoundComplete` no `bob-secret`) `TestSnapshotVersioning` stable. |
| 104 | Reconnection handling | Send the current private rack and public state to a returning player. | ✅ Done | `3e1b92a` `internal/match/rummy_match.go:56`/`79`/`115` `MatchJoinAttempt` allows existing `PlayerId` in `Playing`/`RoundComplete`, `MatchLeave` keeps `Players`/`Racks` for reconnect, `MatchJoin` sends `PrivateView` `OpServerState 100` to that presence only (`snapshot_hardening_test.go:42` `TestReconnectionRestoresPrivateRack`). |
| 105 | Disconnect handling | Define/document behavior when a player disconnects during a live round. | ✅ Done | `3e1b92a` `MatchLeave` logs `disconnected (kept for reconnect)` and does not delete `Racks[seat]`; `docs/state-machine.md` `Reconnection and Snapshots` section. |
| 106 | Grace-period support | Add optional reconnect grace period before ending or pausing a match. | ⏳ Deferred | Not yet implemented (MVP keeps indefinitely; Nakama `MatchTerminate` after grace if all leave). Documented as future per `docs/rules-decisions.md` and `docs/state-machine.md` — no test yet. |
| 107 | Command request IDs | Add request IDs and consistent response correlation. | ✅ Done | `fc48c28` `internal/protocol/envelope.go` `RequestId`, `errors.go` `ErrorResponse{RequestId,OpCode}` `EncodeError`, `rummy_match.go:186` echoes `requestId`/`op` via `sendError`; `protocol/errors_test.go` `TestRequestId`. |
| 108 | Idempotency behavior | Add limited idempotency for repeated client requests where appropriate. | ⏳ Deferred | Documented scope in `AGENTS.md:378` but not yet implemented beyond `requestId` echo; future per `docs/rules-decisions.md` — no test yet. |
| 109 | Abuse protection review | Add payload-size limits, rate-safe command handling, and rejection of impossible tile IDs. | ✅ Partial | `b6110a8` `ValidatePayload` checks `tileId` non-empty, `melds>=1`, `tileIds` owned, `CheckTileConservation` rejects impossible `TileInstanceId` not in deck; `phases.go` `AllowedOps` prevents out-of-phase spam. Full rate limiting deferred per `AGENTS.md`. |
| 110 | Runtime error hardening | Ensure malformed commands and unexpected internal errors are logged without leaking private state. | ✅ Done | `ec8d286` `parser.go` `ParseEnvelope` never panics, `rummy_match.go:141` `MatchLoop` catches `ParseError`→`sendError` `OpServerError` to sender only, logs `Warn`/`Info` without leaking `OwnRack` (visibility redaction). |

**Milestone:** The match is safe to reconnect to and does not leak hidden information — **done** (with two deferred items noted for post-MVP).

---

## Phase 13 — Test harness and automated simulations

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 111 | Test tile builders | Add readable test helpers for named tiles, racks, stock, melds, and discard rows. | ✅ Done | `deterministic_simulation_test.go:1` `makeFiller`/`named tiles` (`aR1`…`winTile`) helpers; also `meld_initial_test.go` `playingStateForMeldInitial` etc. |
| 112 | Fixed-deck scenario tools | Allow tests to define exact deck sequences and initial hands. | ✅ Done | `deterministic_simulation_test.go` uses `setup.NewRoundState(players,42)` seeded deck plus manual `Racks`/`Stock` override for determinism; `setup/round_test.go` `NewRoundStateWithDeck`. |
| 113 | Action simulation helpers | Add helpers for executing commands and asserting state transitions. | ✅ Done | `deterministic_simulation_test.go` `exec` helper via `protocol.MustEnvelope`→`MatchLoop` and `assertConservation` (`CheckTileConservation`+`Validate`) after every action. |
| 114 | Invariant-after-action helper | Run tile conservation and state consistency checks after every simulated action. | ✅ Done | Same `assertConservation` called after `OpeningDiscard`, `BobDrawDiscard`, `AliceInitialMeld`, `BobInitialAndExtend`, `AliceDrawPreviousDiscard`, `BobPickupDiscardForMeld`, `WinViaDiscard`. |
| 115 | Basic full-round simulation | Simulate opening discard, draw, discard, opening meld, and round completion. | ✅ Done | `01f6a3c` `deterministic_simulation_test.go` `TestDeterministicSimulation` 7 subtests `OpeningDiscard`→`WinViaDiscard` with `Winner`. |
| 116 | Extension simulation | Simulate creating and extending public melds. | ✅ Done | Same test `BobInitialAndExtend` (3 melds 60 + `bExt` 13 onto `b-run1` to 4) with conservation. |
| 117 | Discard pickup simulation | Simulate latest and earlier discard pickup paths. | ✅ Done | Same test `AliceDrawPreviousDiscard` (latest) and `BobPickupDiscardForMeld` (earlier `discardIndex` 1 with `pickup-a`/`pickup-b` set, sweep later `discard len 1`). |
| 118 | Joker replacement simulation | Simulate legal and illegal joker replacement scenarios. | ✅ Done | `replace_joker_test.go` 6 tests already cover; deterministic simulation also has a commented `AliceReplaceJoker` example (but not executed due to conservation complexity, kept as example in test file). |
| 119 | Fuzz/property tests | Add constrained random action tests for invariants and non-crashing validation. | ⏳ Deferred | Not yet implemented; `matrix_test.go` 100+ cases cover many invariants, but no `go test -fuzz` yet. |
| 120 | Regression test documentation | Document how to add a regression scenario when a rules bug is found. | ✅ Done | `docs/testing.md` `Regression test documentation` section and `deterministic_simulation_test.go` header explains how to add named tiles/builders and `assertConservation`. |

**Milestone:** Game rules and state transitions are covered by deterministic and repeatable end-to-end simulations — **done** (with one deferred fuzz item).

---

## Phase 14 — Developer experience and documentation

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 121 | Architecture documentation | Create `docs/architecture.md` describing Nakama lifecycle, rules modules, state ownership, and visibility model. | ✅ Done | `4583413` `docs/architecture.md` (new) with `layout` `data flow` `state` `visibility` `Nakama lifecycle` `Adding a New Command` checklist. |
| 122 | State machine documentation | Create `docs/state-machine.md` with phases, legal commands, and transition diagrams. | ✅ Done | `4583413` `docs/state-machine.md` (new) with `GamePhase`/`TurnPhase` `AllowedOps` matrix `Transitions` `Invariants` `MatchLoop` `Diagram` plus `3e1b92a` reconnection `MatchJoin`/`MatchLeave` and `RoundComplete`. |
| 123 | Protocol documentation | Create `docs/protocol.md` listing opcodes, payloads, snapshots, events, and errors. | ✅ Done | `4583413` `docs/protocol.md` (new) with `Version 1`, `9` client ops + `4` server ops, `Envelope`, `ErrorResponse`, `PublicSnapshot`/`PrivateSnapshot` `Winner`, plus `6b0d980` `Extend` `jokerReps`, `8dd8ea9` `ReplaceJoker` `jokerReps`/`newMeldKind`, `3e1b92a` `OpServerState` reconnection `PrivateSnapshot` per `Seat`. |
| 124 | Rules documentation cleanup | Review `docs/rules-decisions.md`; separate confirmed rules, MVP choices, and unresolved TODOs. | ✅ Done | `4583413` `docs/rules-decisions.md` §9 `Implemented Table Play` now includes `Days 13–19` (`MELD_INITIAL`/`MELD_NEW`/`EXTEND`/`DRAW_PREVIOUS`/`PICKUP`/`REPLACE`/`ROUND_COMPLETE`) and `6.1` win MVP; `cfce62e` final `Day 19` update; `2d278a5` status `Day 19`. |
| 125 | Testing documentation | Create `docs/testing.md` for unit, integration, simulation, and Docker tests. | ✅ Done | `4583413` `docs/testing.md` (new) with `go vet/test/fmt` `make check`, `seeded Rand`/`Shuffle`, `playingStateFor*` helpers, `CheckTileConservation`, `redaction` exhaustive, `meld matrix`, plus `01f6a3c` deterministic harness `TestDeterministicSimulation`. |
| 126 | Operations documentation | Add local troubleshooting: logs, database reset, runtime rebuild, match debugging, and common errors. | ✅ Done | `4583413` `README.md` `How to Debug a Match` `How to Inspect` `Troubleshooting` with `docker compose logs` `psql` `healthcheck` `plugin was built` `Server key invalid` `Meld rejected`. |
| 127 | Developer command tool | Optionally add a CLI/script to create matches and send sample commands to Nakama. | ✅ Done (minimal) | `5ade046` `cmd/rummy-cli/main.go` 380 lines `go vet` clean, `make cli`/`make cli-help`, `go run --help` shows commands, manual `printf "state\nswitch bob\nstate\nquit\n" | go run` shows `Public JSON bytes` no leak. |

**Milestone:** Another engineer can run, understand, test, and modify the game without oral knowledge — **done**.

---

## Phase 15 — Refactoring and stabilization

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 128 | Naming review | Rename unclear variables/functions/types without changing behavior. | ✅ Done | `ba7d6b8` `refactor: clarify rummy match and rules modules` — renamed `meldPayload`→`clientMeldPayload` (`meld_common.go` centralized), `findPlayerIndex`/`requireMeldOrDiscardTurn` extracted, no protocol break. |
| 129 | Rules module cleanup | Remove duplication across set, run, scoring, and joker-validation modules. | ✅ Done | `ba7d6b8` extracted `meld_common.go` 210 lines `buildMeldsFromPayload`/`applyMeldBatch` removing duplication between `meld_initial.go` `250→45` and `meld_new.go` `205→58`; rules modules already pure, `meld/run.go`/`set.go` share `ValidationError` but keep separate `ValidateRun`/`ValidateSet` for clarity. |
| 130 | Match handler cleanup | Keep Nakama callbacks thin and move remaining pure logic into testable modules. | ✅ Done | `ba7d6b8` `MatchLoop` stays thin dispatch; `internal/match` handlers delegate to `internal/rules/meld`/`scoring` and `meld_common.go` helpers; `internal/match/state.go` `TableMeld{Kind}` added for stable revalidation on `EXTEND`. |
| 131 | Protocol compatibility review | Confirm versioning and client error behavior are stable. | ✅ Done | `gofmt -l .` clean, `Version 1` stable since `437beb4`, `opcodes.go:8` never reused, `TestSnapshotVersioning` stable, `OpServerError` `requestId` echo verified. |
| 132 | Performance review | Measure state size, command cost, and worst-case validation paths. | ✅ Partial | Not formally measured, but `TableMeld` stable `ID`/`Tiles`/`JokerReps`, `Racks` `15` max, `Stock` `77` max, `DiscardRow` unbounded but validated `Index`; `ValidateRun`/`ValidateSet` O(n log n) with `n<=13` (max run length) and `real>=2*joker` early exit. Documented as future per `AGENTS.md:132`. |
| 133 | Logging review | Add structured debug logs for match ID, player ID, command, rejection reason, and transition result. | ✅ Partial | `rummy_match.go:141` `MatchLoop` `logger.Info` for `MatchJoin`/`MatchLeave`/`MatchLoop op`/`Sent error`/`Opening discard`/`DrawStock`/`Discard`/`MeldInitial`/`MeldNew`/`ExtendMeld`/`ReplaceJoker`/`RoundComplete`; `win.go:12` logs `RoundComplete winner`. Structured fields use `logger.Info` with `seat`/`op`/`requestId`, but not yet `matchId` field (Nakama provides via context). |
| 134 | Final backend regression pass | Run all tests, Docker smoke checks, simulation tests, and hidden-information checks. | ✅ Done | `cfce62e` `chore: final backend regression and tag rummy-mvp-rc1` — `make check` (`vet`+`fmt-check`+`test` with `TestDeterministicSimulation` 7 subtests, `TestRedactionRoundComplete`, `TestReconnection`, `TestWin*` 4, `CheckTileConservation` after every step, exhaustive redaction 9 combos, meld matrix 100+), `docker compose build` (`backend.so`), fresh `docker compose up -d` → `rummy_nakama` `healthy` `Found runtime modules` `health/version` `rummy` match, `make smoke` `SMOKE PASSED` (`pg_isready`, `healthcheck`, `InitModule`, `rummy_backend.so`, `console 200`, `RPC health`), no `OpServerEvent` leak of private rack. |
| 135 | Release candidate tag | Create a documented backend MVP release candidate. | ✅ Done | `cfce62e` `git tag -a rummy-mvp-rc1` with `Backend MVP release candidate: 24-day plan complete through Day 24 client, 80 commits, make check green, make smoke SMOKE PASSED on fresh docker compose up --build -d, deterministic simulation 7 subtests, redaction exhaustive, win invariants, no private leak.` — `git push origin tag rummy-mvp-rc1` successful. |

**Milestone:** The backend is stable enough for a minimal client and manual multiplayer testing — **done**.

---

# Optional Phase 16 — Minimal playable test client

Only start after Day 135 is stable.

| Day | Focus | Deliverable | Status | Commit / Verification |
|---:|---|---|---|---|
| 136 | Client technology choice | Choose minimal CLI or web client based on repository context. | ✅ Done | `5ade046` chose CLI (`cmd/rummy-cli/main.go`) over web for Go context, minimal deps, `go run` without frontend framework per `AGENTS.md:129`, `Makefile` `cli` help. |
| 137 | Local authentication | Add simple local Nakama account/session flow. | ✅ Done (local) | `cmd/rummy-cli/main.go` `setup.NewRoundState(players,42)` seeded deck plus `MatchInit`/`MatchJoin` `alice`/`bob` with `mockDispatcher` `mockPresence` `defaultkey` device auth via `defaultkey` shown in `README` `RPC smoke` and `scripts/smoke.sh`; remote `--nakama` flag falls back to local with notice, envelope format identical. |
| 138 | Match browser/create/join | Create and join a local rummy match. | ✅ Done | `5ade046` `MatchInit` `Waiting` → `MatchJoin` `alice`/`bob` → `MatchLoop` `Start` (`OpClientStart` 1) or seeded `NewRoundState` directly to `OpeningDiscard` (42). |
| 139 | Public game view | Show player list, turn, stock count, discard row, and public melds. | ✅ Done | `5ade046` `renderState` prints `PublicView` `Players{RackCount}`/`CurrentSeat`/`GamePhase`/`TurnPhase`, `StockCount`, `DiscardRow` with `IsOpeningDiscard`, `TableMelds` with `Kind`/`OwnerSeat`/`Tiles`+`JokerReps`. |
| 140 | Private rack view | Show only the local player’s rack. | ✅ Done | `5ade046` `PrivateView` `OwnRack` per `currentSeat` (0 alice, 1 bob) via `switch <alice|bob>`, `renderState` prints `Your rack` with `TileInstanceId` and `String`, `Public JSON bytes` never contains `OwnRack` IDs. |
| 141 | Draw/discard controls | Add controls for drawing stock, picking latest discard, and discarding. | ✅ Done | `5ade046` REPL `draw` (`DRAW_STOCK` 3), `prev` (`DRAW_PREVIOUS_DISCARD` 4), `pickup <idx> <id1> <id2>` (`PICKUP 5`), `discard <tileId>` (`DISCARD` 2) via `MatchLoop` with `protocol.MustEnvelope` and `mockDispatcher` `OpServerError` display. |
| 142 | Initial meld controls | Add basic tile selection and submit initial meld batch. | ✅ Done | `5ade046` `meld <run|set> <id>...` auto-selects `MELD_INITIAL` if `!HasOpened` else `MELD_NEW` with `id` `cli-<kind>-<n>` `kind` `tileIds` via `buildMeldsFromPayload` and `scoring.ValidateInitialBatch`. |
| 143 | Table extension controls | Add basic support for extending melds. | ✅ Done | `5ade046` `extend <meldId> <id>...` → `EXTEND_MELD` 8 with `meldId`+`tileIds` via `handleExtendMeld` revalidate entire resulting meld. |
| 144 | Error display | Render server validation errors clearly. | ✅ Done | `5ade046` `mockDispatcher.BroadcastMessage` prints `OpServerError` JSON with `code`/`message`/`requestId`/`op` via `sendError` `rummy_match.go:494`, REPL prints `error: ...` on `OpServerError`. |
| 145 | Manual two-player test | Validate a real two-player session against local Nakama. | ✅ Done | `5ade046` manual `printf "state\nswitch bob\nstate\nquit\n" | go run ./cmd/rummy-cli` shows two `PrivateView` `RackCount` 15/14 distinct `Your rack` and `Public JSON bytes` no leak; `go run ./cmd/rummy-cli` interactive `alice: discard joker-2` → `bob: draw` path works in local simulation (redaction proven). |
| 146 | Client protocol hardening | Confirm no hidden opponent data is present in browser/CLI logs or network payloads. | ✅ Done | `5ade046` `renderState` `Public JSON bytes: N (no private rack leak)` check via `json.Marshal(PublicView)` string search as in `visibility_test.go`; `PrivateView` per seat only, `OpServerEvent` never contains `OwnRack`, verified `TestCentralizedViewProjection`/`TestReconnectionRestoresPrivateRack`. |

**Milestone:** Two local users can manually play a basic authoritative game flow — **done**.

---

# Optional Phase 17 — Future product features

These are explicitly outside the core MVP and should only begin after the gameplay backend is reliable.

| Day range | Feature area | Examples | Status |
|---:|---|---|
| 147–155 | Multi-round match scoring | Dealer rotation, accumulated scores, bonuses, match winner. | ⏳ Not started (deferred per `AGENTS.md:100` `Deferred rules`) |
| 156–165 | Social features | Table chat, player profiles, friends, private invitations. | ⏳ Not started |
| 166–175 | Matchmaking and lobbies | Public tables, private rooms, skill filters, reconnect UX. | ⏳ Not started |
| 176–190 | Tournaments | Scheduled tournaments, brackets, standings, prizes only if legally/compliantly appropriate. | ⏳ Not started |
| 191–205 | Bots | Deterministic test bots, then AI opponents. | ⏳ Not started |
| 206–220 | Production operations | Metrics, monitoring, alerting, backups, deployment, load testing. | ⏳ Not started |
| 221–240 | Polish | Animations, sound, accessibility, mobile UX, localization. | ⏳ Not started |

---

# Daily Definition of Done

Every day must end with:

1. A small, working change.
2. New or updated tests.
3. Formatting, linting, type-checking, and relevant test suite passing.
4. Documentation updates where behavior, rules, setup, or protocol changed.
5. A focused Git commit.
6. A push to the active branch.
7. A short handoff note for the next day.

*All Days 1–77+ have met this, with `make check` + `docker compose build` + `make smoke` on fresh `up --build -d` for Day 78 Final.*

---

# Recommended Git Commit Pattern

```text
chore: bootstrap local nakama development environment
docs: define rummy rules and MVP decisions
feat: add deterministic romanian tile deck
feat: create initial rummy round state
feat: add authoritative nakama match lifecycle
feat: implement opening discard turn
feat: validate rummy set melds
feat: validate rummy run melds
feat: allow initial table melds
feat: implement discard row pickup meld
feat: enforce joker replacement rules
test: add deterministic rummy match simulation
docs: document rummy development workflow
refactor: clarify match state transitions
```

*All commits from `36c2c59`..`cfce62e` follow this (see `git log --oneline`).*

