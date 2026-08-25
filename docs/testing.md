# Testing — Deterministic Harness

**Runner:** `go test ./...` / `go vet ./...` / `go fmt ./...` — pinned `go 1.23.5` per `go.mod:3` (tested on `1.27` locally, `1.23` in CI and `nakama-pluginbuilder:3.26.0`). No Jest/Vitest; Go is authoritative after `55c7f3b`.

## Commands

```bash
go vet ./...          # static checks, required before commit
go test ./...         # all packages, deterministic
go test ./internal/match -run TestMeldInitial -v   # single suite
go test ./internal/match -run TestMeldNew -v
go test ./internal/rules/meld -run TestValidate -v
go test ./internal/setup -v
make vet / make test / make check   # vet + fmt-check + test (CI baseline)
make fmt-check        # fails if go fmt would change
gofmt -l .            # list unformatted
```

CI mirrors `make check` in `.github/workflows/ci.yml` (`setup-go 1.23.5`, `go vet`, `gofmt -l`, `go test`, `go mod tidy` diff, `docker/build-push-action` for `compose build`).

## Determinism

- **Shuffled deck:** `internal/setup/rand.go:9` `NewSeededRand(seed)` wraps `math/rand.Rand`; `Shuffle(deck,Rand)` is Fisher–Yates via `r.Shuffle` (pure, original not mutated). Tests use fixed seeds (`42`, `123`, `12345`, `2026`) and assert `DeterministicFixedSeed` equal, `DifferentSeedsDiverge`.
- **Dealing:** `Deal(deck,n)` and `NewRoundState(playerIds,seed)` deterministic from seed and `AssignSeats` join order (Seat 0 = opener). Tests cover `n=2,3,4 × seeds 0,1,42,123,999,2026` and conservation.
- **Match flow:** helpers build synthetic `RoundState` with known `TileInstanceId`s (e.g., `t1`, `j1`, `alice-fill-09`) and `Stock` top = last element; `MatchLoop` driven via `protocol.MustEnvelope(op,payload)` and `mockDispatcher`/`mockPresence`. No network, no Nakama DB.

## Helpers

- `playingStateForMeldInitial(currentRack []TileInstance) (*RoundState, []TileInstance)` (`meld_initial_test.go:19`) — builds `Playing/MeldOrDiscard` with `CurrentSeat 0`, `HasOpened false`, `alice 15` (newTiles + Black filler), `bob 14` Yellow, `stock 76` Blue (so `15+14+76+1=106`), `discard IsOpening true`. Returns `allTiles106` for `CheckTileConservation`.
- `playingStateForMeldNew(opened bool, existingMelds []TableMeld, newTiles []TileInstance, seat Seat)` (`meld_new_test.go:10`) — as above but `HasOpened` flag, `existingMelds` on table, `stock =106 - (15+14+1+meldTiles)`, `CurrentSeat` param. Same conservation guarantee.
- `openingStateWithDeal()` (`opening_discard_test.go:15`) — synthetic 106 deck `t-000…` + `joker-1/2`, deal `15/14/77` for opening discard tests.
- `playingStateWithStock()` / `playingStateForDiscard()` — small deterministic states for `DRAW_STOCK` / `DISCARD`.

All helpers fill `allTiles106` as the union of `Racks`+`Stock`+`DiscardRow`+`TableMelds` with distinct IDs (no reuse across colours/ranks except intentional `red-01-1`/`red-01-2` etc. for deck tests).

## Conservation and Invariants

Every state-changing test calls `CheckTileConservation(state, allTiles106)` (`internal/match/invariant.go:18`):

- `len(allTiles)==106` and each `Validate()` passes, no duplicate `ID` in `allTiles`.
- Every `ID` in `state` appears exactly once across `Racks`/`Stock`/`DiscardRow`/`TableMelds` with location string (`rack seat 0 index 2`, `stock index 5`, `discard index 0`, `meld "m1"`).
- Duplicate, not-in-deck, missing all error with pinpointed location.
- Also `state.Validate()` checks per-location duplicate, `IsOpeningDiscard` only on index 0, `TableMeld.Validate()` checks `JokerReps` map (missing rep, rep for non-joker, invalid colour/rank).

After rejected commands, tests assert `len(Racks[Seat])` and `len(TableMelds)` unchanged and `CheckTileConservation` still passes (atomicity).

## Redaction

- `internal/match/visibility_test.go:9` — `PublicView` must not contain any `TileInstanceId` from racks/stock (via `json.Marshal` string search `HidesRacks`/`OpponentNotLeaked`), `PrivateView` must contain own rack and not others.
- `internal/setup/redaction_test.go:9` — exhaustive `n=2,3,4 × seeds 42,7,123` with real `NewRoundState` 106 deck: `PublicView` must not contain any `rack`/`stock` `ID` (only `stockCount`), each `PrivateView` must contain own and not others nor stock. Placed in `setup` to avoid `setup→match` import cycle.

## Meld Matrix

`internal/rules/meld/matrix_test.go` — comprehensive valid/invalid sets, runs, Ace (`1-2-3` vs `12-13-1` vs `13-1-2`), joker filling gap at end, `real>=2*joker` ratio, immutable `JokerReps`. Scoring tests (`scoring/scoring_test.go`, `run_test.go`, `set_test.go`, `joker_test.go`) cover `5/10/25` and high/low Ace.

## Running Focused Suites

```bash
go test ./internal/rules/meld -run TestValidate -v
go test ./internal/rules/scoring -run TestTotalScore -v
go test ./internal/match -run TestMeldInitialSuccess -v
go test ./internal/match -run TestMeldNewWithJoker -v
go test ./internal/setup -run TestSetupInvariantsAllPlayersAndSeeds -v
```

Keep tests small, explicit, and deterministic; prefer named tiles (`t1` Red 5) and builders over raw IDs. After every `MatchLoop` action in a simulation, call `CheckTileConservation` and assert `TurnPhase`/`CurrentSeat`/`HasOpened`/`RackCount`.

---

*Code: `internal/setup/rand.go:9`, `shuffle.go:9`, `deal.go:9`, `round.go:9`, `invariant.go:18`, `visibility.go:9`, `meld_initial_test.go:19`, `meld_new_test.go:10`.*
