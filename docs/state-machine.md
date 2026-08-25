# State Machine — GamePhase, TurnPhase, Allowed Ops

**Source:** `internal/match/state.go:12`, `phases.go:15`, `rummy_match.go:141`, `docs/rules-decisions.md:5`, `docs/terminology.md`.

## Phases

`GamePhase` (`state.go:12`):

- `Waiting (0)` — lobby, 2–4 players not yet started. `CurrentSeat = SeatInvalid`, `TurnPhase = MustDraw` placeholder, `DiscardRow` empty, `TableMelds` empty.
- `OpeningDiscard (1)` — opening player (Seat 0, deterministic host/first joiner) must discard one tile from 15 to reach 14. No other ops. After success: `GamePhase→Playing`, `CurrentSeat` advances anticlockwise to next seat, `TurnPhase=MustDraw`, first discard flagged `IsOpeningDiscard=true` index 0.
- `Playing (2)` — normal loop, governed by `TurnPhase`.
- `RoundComplete (3)` — winner or dead round, no gameplay ops. `Winner = Seat` or `SeatInvalid` for stock-exhausted dead round (future `docs/rules-decisions.md:6.2`).

`TurnPhase` (`state.go:42`) only meaningful when `GamePhase==Playing`:

- `MustDraw (0)` — at turn start, must draw: `DRAW_STOCK` or (if opened) `DRAW_PREVIOUS_DISCARD` / `PICKUP_DISCARD_FOR_MELD`.
- `MeldOrDiscard (1)` — after draw, may `MELD_INITIAL` (once, if not opened), `MELD_NEW`, `EXTEND_MELD`, `REPLACE_JOKER`, must end with `DISCARD` unless win.

## Seats and Turn Advance

- Seats deterministic join order `0..n-1` (`internal/match/seat.go:71` `AssignSeats`). Opening player is Seat 0.
- Anticlockwise = `NextSeat(current,n) = (current+1)%n` (`seat.go:93`). `PrevSeat` is clockwise for tests.
- `AdvanceTurn(state)` (`turn.go:9`) moves `CurrentSeat→NextSeat` and resets `TurnPhase=MustDraw`; validates `2..4` players and `CurrentSeat` valid.

## Allowed Ops Matrix (`phases.go:15` `AllowedOps`)

| GamePhase | TurnPhase | Allowed `op` |
|---|---|---|
| `Waiting` | — (`MustDraw` placeholder) | `1 OpClientStart` (host Seat 0, ≥2 players) |
| `OpeningDiscard` | — | `2 OpClientDiscard` (only Seat 0, 15→14, blocked flag) |
| `Playing` | `MustDraw` | `3 DrawStock`, `4 DrawPreviousDiscard` (opened, not opening), `5 PickupDiscardForMeld` (opened) |
| `Playing` | `MeldOrDiscard` | `2 Discard` (must end turn), `6 MeldInitial` (once), `7 MeldNew` (opened), `8 ExtendMeld` (opened), `9 ReplaceJoker` (opened) |
| `RoundComplete` | — | none |

`ValidatePhaseOp(state,op)` (`validate.go:28`) returns `wrong_phase` with `gamePhase`/`turnPhase` details if not in matrix. `ValidateActivePlayer(state,senderId)` (`validate.go:13`) returns `not_member` or `not_your_turn` for `OpeningDiscard`/`Playing` (not enforced in `Waiting`/`RoundComplete`).

## Transitions

```
Waiting --[host Start 1, 2..4 players]--> OpeningDiscard (CurrentSeat=0)
OpeningDiscard --[Discard 2 by Seat0, tile owned, 15→14]--> Playing/MustDraw (CurrentSeat=NextSeat(0))
Playing/MustDraw --[DrawStock 3, stock>0]--> Playing/MeldOrDiscard (same CurrentSeat, TurnPhase→MeldOrDiscard, Stock-1, Racks+1)
Playing/MustDraw --[DrawPreviousDiscard 4, opened, not opening]--> Playing/MeldOrDiscard
Playing/MustDraw --[PickupDiscardForMeld 5, opened, meld valid with 2 rack tiles, sweep later discards]--> Playing/MeldOrDiscard
Playing/MeldOrDiscard --[MeldInitial 6, not opened, batch 50+ with run, each meld valid, owned, no dup]--> Playing/MeldOrDiscard (same seat, HasOpened=true, Racks→TableMelds)
Playing/MeldOrDiscard --[MeldNew 7, opened, each meld valid, ownership, no dup, meldId not colliding]--> Playing/MeldOrDiscard (same seat, Racks→TableMelds)
Playing/MeldOrDiscard --[ExtendMeld 8, opened, whole resulting meld revalidated]--> Playing/MeldOrDiscard
Playing/MeldOrDiscard --[ReplaceJoker 9, opened, exact represented tile, then new meld with 2 rack tiles]--> Playing/MeldOrDiscard
Playing/MeldOrDiscard --[Discard 2, owned, append DiscardRow false, AdvanceTurn]--> Playing/MustDraw (CurrentSeat=NextSeat, TurnPhase→MustDraw)
Playing/MeldOrDiscard --[any melding that empties rack]--> RoundComplete (future win.go, Winner set) — also Discard that empties rack per docs/rules-decisions.md:6.1
Playing/MustDraw --[stock empty and no discard pickup legal]--> RoundComplete (dead round, no winner, reason stockExhausted, docs/rules-decisions.md:6.2) — pending
```

## Invariants (must hold after every command, `invariant.go:18`)

- `CheckTileConservation(state, allTiles106)` — every of 106 `TileInstanceId` appears exactly once across `Racks` + `Stock` + `DiscardRow` + `TableMelds` (with location strings). Duplicate, not-in-deck, missing all error.
- No duplicate `TileId` in any single `Racks[seat]` or `TableMeld.Tiles` (`state.go:97` `Validate`).
- Only discard index 0 may have `IsOpeningDiscard=true`.
- `GamePhase`/`TurnPhase` valid per `IsValid`.
- `JokerReps` map key is joker `TileId` present in `Tiles`, value is `Colour.IsValid && Rank.IsValid` and `IsJoker==false`, immutable after meld (unless `REPLACE_JOKER`).
- Atomicity: invalid commands do not partially mutate state (validate fully before touching `state`).

## MatchLoop Dispatch (`rummy_match.go:141`)

`MatchLoop` loops `messages []MatchData`: `ValidateEnvelope` → `ValidateActivePlayer` (not enforced in `Waiting`/`RoundComplete`, allows reconnection) → `ValidatePhaseOp` → `ValidatePayload` (all with `requestId` correlation, `sendError` to sender only) → handle `Start`, `OpeningDiscard`, `NormalDiscard`, `DrawStock`, `DrawPreviousDiscard` (`4`), `PickupDiscardForMeld` (`5`), `MeldInitial` (`6`), `MeldNew` (`7`), `ExtendMeld` (`8`), `ReplaceJoker` (`9`), win check via `win.go:12` `checkWinAndComplete` after any `rack==0` (discard or meld), else `not implemented`.

## Reconnection and Snapshots (Day 20)

- `MatchJoinAttempt` (`rummy_match.go:56`) allows existing `PlayerId` to rejoin in any phase (`Waiting` returns `already joined`, `Playing`/`RoundComplete` returns success for same `Seat`).
- `MatchJoin` (`rummy_match.go:79`) keeps `Players`/`Racks` on `MatchLeave` (Day 20) for reconnect; on rejoin (already in `Players`) it does not allocate new `Seat` and sends `PrivateView(state, seat)` via `dispatcher.BroadcastMessage(OpServerState 100, json.Marshal(PrivateSnapshot), [presence], true)` to that presence only — other players never receive foreign `OwnRack` (verified by `snapshot_hardening_test.go:42` `TestReconnectionRestoresPrivateRack`).
- `MatchLeave` (`rummy_match.go:115`) logs `disconnected (kept for reconnect)` and does not delete `Racks[seat]`; Nakama terminates after grace if all leave.
- Snapshots are versioned `PublicSnapshot.Version=1` (`visibility.go:48`) and stable (`json.Marshal` deterministic, tested `TestSnapshotVersioning`).
- Client sync expectation: on join or rejoin, client receives `OpServerState 100` `PrivateSnapshot` (public + `OwnRack` + `OwnSeat` + `Winner`); public state also broadcast via `OpServerStatePublic 101` if needed. Client must not trust local rack after disconnect — rehydrate from snapshot. See `docs/protocol.md` and `docs/testing.md`.

## Diagram (text)

```
Waiting ─1─► OpeningDiscard ─2─► Playing/MustDraw ─3/4/5─► Playing/MeldOrDiscard ─2─► Playing/MustDraw ─► ...
                               │                          ▲ 6/7/8/9 (stay)
                               │                          └─2 or 6/7/8/9/5(empty rack → RoundComplete via win.go:12)
                               │                          └─reconnect: MatchJoin sends PrivateView to that Seat only
```

See `docs/rules-decisions.md:5` for turn flow narrative and `docs/protocol.md` for payload schemas per transition. See `docs/architecture.md` for visibility centralization.

---

*Code: `internal/match/state.go:12`, `phases.go:15`, `turn.go:9`, `validate.go:13`, `rummy_match.go:56`/`79`/`115`/`141`, `visibility.go:36`, `win.go:12`, `discard.go:9`.*
