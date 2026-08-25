# Protocol — Client/Server JSON Messages

**Version:** `1` (`internal/protocol/opcodes.go:8` — bump only on breaking change).  
**Transport:** Nakama authoritative match `MatchLoop` via `dispatcher.BroadcastMessage`; client → server is `MatchData` with `Envelope` JSON, server → client is `OpServer*` broadcasts. No REST/RPC for gameplay (only `health`/`version` RPCs for smoke).

## Envelope

All client→server messages are a single JSON object `Envelope` (`internal/protocol/envelope.go:9`):

```json
{
  "v": 1,
  "op": 6,
  "requestId": "uuid-123",
  "payload": { ... }
}
```

- `v` must be `1`, else `bad_version`.
- `op` must be `1..9` `IsClientOp`, else `unknown_opcode`.
- `requestId` optional, echoed in `OpServerError` for correlation (`internal/protocol/errors.go:14`).
- `payload` is `json.RawMessage` validated per-opcode (`ValidatePayload`).

Helpers: `NewEnvelope(op,payload)` and `NewEnvelopeWithRequestId` for tests/clients; `MustEnvelope` panics on marshal error (tests only).

## Client → Server Opcodes (stable, never reuse)

| Op | Name | Payload schema | Notes |
|---|---|---|---|
| 1 | `OpClientStart` | `{}` or `{"extra":...}` allowed (empty) | `Waiting` only, host Seat 0, 2–4 players → `OpeningDiscard`. Also `MatchSignal "start"` for local dev. |
| 2 | `OpClientDiscard` | `{"tileId":"<TileInstanceId>"}` required non-empty | `OpeningDiscard` (15→14 blocked, `Index 0`) or `Playing/MeldOrDiscard` (append, `IsOpeningDiscard false`, advance `(current+1)%n` → `MustDraw`). |
| 3 | `OpClientDrawStock` | `{}` (empty) | `Playing/MustDraw` only; pops `Stock` top (last element) to `Racks[current]`; `TurnPhase→MeldOrDiscard`; `stock empty` → `bad_request` per `docs/rules-decisions.md:6.2`. |
| 4 | `OpClientDrawPreviousDiscard` | `{}` | `Playing/MustDraw` (opened, `HasOpened`, not opening `IsOpeningDiscard`, latest only) → `MeldOrDiscard` (discard tail → `Racks[seat]`, `DiscardRow` pop). |
| 5 | `OpClientPickupDiscardForMeld` | `{"discardIndex":int>=0, "tileIds":["...","..."], "jokerReps":{"jokerId":{"colour":"red","rank":7}}, "kind":"run"\|"set"}` exactly 2 | `Playing/MustDraw` (opened, `HasOpened`, non-opening `discardIndex` + exactly 2 `tileIds` owned → valid 3-tile `meld.New` with that discard; `jokerReps` for any joker among the three; `kind` optional or inferred `run`→`set`). Atomically `Racks = rack -2 + laterTiles` (`laterTiles = DiscardRow[discardIndex+1:]` swept in order), `DiscardRow = DiscardRow[:discardIndex]` reindexed, `TableMelds += newMeld` (`Kind` `run`/`set`), `TurnPhase→MeldOrDiscard` or `RoundComplete` if win. |
| 6 | `OpClientMeldInitial` | `{"melds":[{"id":"m1","kind":"run"\|"set","tileIds":["t1","t2","t3"], "jokerReps":{"jokerId":{"colour":"red","rank":7}}} , ...]}` ≥1 meld | `Playing/MeldOrDiscard`, `HasOpened==false`, each meld 3+ tiles owned, each `ValidateRun/Set` passes, total `>=50` with ≥1 run, no duplicate `TileId` across batch, `meldId` unique in batch. Atomically `Racks→TableMelds` and `HasOpened=true`; stays `MeldOrDiscard` or `RoundComplete` if `rack==0`. |
| 7 | `OpClientMeldNew` | same shape as 6 | `Playing/MeldOrDiscard`, `HasOpened==true`, each meld valid, ownership, no duplicate, `meldId` not colliding with existing `TableMelds`; no score minimum; atomic; supports multiple melds per batch; or `RoundComplete`. |
| 8 | `OpClientExtendMeld` | `{"meldId":"...","tileIds":["..."], "jokerReps":{"jokerId":{"colour":"red","rank":7}}}` ≥1 | `Playing/MeldOrDiscard` (opened, `HasOpened`, revalidates entire resulting meld `meld.New` + `ValidateRun/Set`, preserves `JokerReps` immutability, any owner (`OwnerSeat` kept), atomic `Racks→TableMeld`) → stays `MeldOrDiscard` or `RoundComplete`. |
| 9 | `OpClientReplaceJoker` | `{"targetMeldId":"...","tileId":"...","newMeldTiles":["...","..."], "jokerReps":{"jokerId":{"colour":"red","rank":7}}, "newMeldKind":"run"\|"set"}` exactly 2 | `Playing/MeldOrDiscard` (opened, `HasOpened`, exact tile for run or exact missing colour for set per `docs/rules-decisions.md:6.4` MVP, `tileId` owned real tile equals some joker’s `JokerReps` colour/rank, `newMeldTiles[2]` owned distinct real tiles; `jokerReps` must contain recovered joker’s new `colour`/`rank`; tries `run`→`set` or forced `newMeldKind`, validates both `updatedMeld` and `newMeld` via `meld.New`; atomically `Racks = rack -3` (`tileId`+2), `TableMelds[targetIdx]=updated` (joker→real, `JokerReps` without that joker) + `+=newMeld` (`Kind` stable), stays `MeldOrDiscard` or `RoundComplete`) |

`ValidatePayload` (`internal/protocol/validator.go:12`) enforces these schemas with `bad_payload` including field name. Unknown opcode or bad version never reaches `MatchLoop` payload handling — `ParseEnvelope` rejects before.

## Server → Client Opcodes

| Op | Name | Payload | When |
|---|---|---|---|
| 100 | `OpServerState` | `PrivateSnapshot` JSON (see below, `v:1` versioned, stable) | `MatchJoin` (initial join and **reconnect**: `rummy_match.go:79` sends `PrivateView(state, seat)` via `dispatcher.BroadcastMessage(OpServerState, json.Marshal(PrivateSnapshot), [presence], true)` to that presence only; others never receive foreign `OwnRack` — verified `snapshot_hardening_test.go:42` `TestReconnectionRestoresPrivateRack`) and `win.go:12` `checkWinAndComplete` final snapshot |
| 101 | `OpServerStatePublic` | `PublicSnapshot` | broadcast (if used, e.g., win `win.go:12` broadcasts `101` `RoundComplete` `Winner`) |
| 102 | `OpServerError` | `ErrorResponse{code,message,details,requestId,op}` JSON | on any rejection (`sendError` `rummy_match.go:494` to `[sender]` only, echoes `requestId`/`op`) |
| 103 | `OpServerEvent` | `{"phase":...,"currentSeat":...,"op":...,"melds":...,"winner":...}` JSON | `started`, `discard`, `drawStock`, `drawPreviousDiscard`, `pickupDiscardForMeld`, `meldInitial`, `meldNew`, `extendMeld`, `replaceJoker`, `roundComplete` (`win.go:12` `roundComplete`) |

Current implementation broadcasts `OpServerEvent` JSON strings (e.g., `{"op":"meldInitial","seat":0,"melds":3}`, `{"op":"drawPreviousDiscard","seat":0,"tileId":"..."}`, `{"op":"roundComplete","winner":0}`) plus `OpServerError` to sender only via `dispatcher.BroadcastMessage` with `presences=[sender]` (`internal/match/rummy_match.go:494`) and `OpServerState` snapshots on join/reconnect/win.

## Error Protocol

`ErrorResponse` (`internal/protocol/errors.go:12`):

```json
{
  "code": "not_your_turn",
  "message": "not your turn: current seat-0, sender seat-1",
  "details": {"currentSeat":"0"},
  "requestId": "uuid-123",
  "op": 6
}
```

Stable codes (`errors.go:53`):

- `bad_request` — generic validation (score <50, invalid meld, stock empty)
- `bad_json` — empty/non-JSON
- `bad_version` — `v !=1`
- `unknown_opcode` — `op` not in 1..9
- `bad_payload` — missing/empty field per `ValidatePayload`
- `not_member` — `SeatOfPlayer` not found
- `not_your_turn` — `sender != CurrentSeat` in `OpeningDiscard`/`Playing`
- `wrong_phase` — `AllowedOps` forbids `op` in current `GamePhase`/`TurnPhase`
- `invalid_tile` — `tileId` not in rack
- `not_opened` — `MeldNew/Extend/Replace` requires `HasOpened`
- `already_opened` — `MeldInitial` when already opened

`NewError(code,message,requestId,details)` and `EncodeError` are used via `sendError` (`rummy_match.go:448`) which logs and sends `OpServerError` to sender only; `ValidateActivePlayer`/`ValidatePhaseOp` return `*ErrorResponse` with `RequestId`/`OpCode` echoed.

## Snapshots (visibility)

`PublicSnapshot` (`internal/match/visibility.go:9`):

```json
{
  "v": 1,
  "gamePhase": "Playing",
  "turnPhase": "MeldOrDiscard",
  "currentSeat": 0,
  "players": [{"id":"alice","seat":0,"hasOpened":true,"rackCount":6},{"id":"bob","seat":1,"hasOpened":false,"rackCount":14}],
  "stockCount": 76,
  "discardRow": [{"Tile":{"ID":"disc-open","Colour":1,"Rank":7,"IsJoker":false},"IsOpeningDiscard":true,"Index":0}],
  "tableMelds": [{"ID":"m1","Tiles":[...],"JokerReps":{"j1":{"ID":"rep-j1","Colour":1,"Rank":7,"IsJoker":false}},"OwnerSeat":0}],
  "winner": -1
}
```

- `Players[].RackCount` only, never tile IDs.
- `StockCount` only, never tile IDs.
- `DiscardRow` ordered, `Index` 0-based, `IsOpeningDiscard` only on first entry.
- `TableMelds` public with stable `ID` and explicit `JokerReps` (colour+rank, never joker).
- `Winner` is `SeatInvalid (-1)` until `RoundComplete`.

`PrivateSnapshot` (`visibility.go:30`) is `PublicSnapshot` plus `ownRack []TileInstance` and `ownSeat` (`v:1` versioned, stable `json.Marshal` deterministic `TestSnapshotVersioning`). Server sends `PrivateView(state, seat)` per player on `MatchJoin` and on reconnection (`rummy_match.go:79` keeps `Players`/`Racks` on `MatchLeave` and re-sends `PrivateSnapshot` to that `Seat` only); never sends foreign racks (verified by `internal/setup/redaction_test.go` exhaustive: `n=2,3,4 × seeds 42,7,123` and `internal/match/visibility_test.go` and `snapshot_hardening_test.go:42` `TestRedactionRoundComplete` which also verifies `RoundComplete` `Winner` public `winner` while `PublicView` still redacts `bob-secret` and `PrivateView` for winner empty `OwnRack`). `PublicSnapshot`/`PrivateSnapshot` both carry `v:1` and `winner` (`SeatInvalid -1` until `RoundComplete`).

## Adding a New Opcode

See `README.md#how-to-add-a-new-command-safely` and `opcodes.go:1` comment: pick next `1..99` not used, add to `IsClientOp` switch in `parser.go:42`, add schema to `validator.go`, add `AllowedOps` in `phases.go`, add handler and wire in `MatchLoop`.

---

*Code references: `internal/protocol/opcodes.go:8`, `envelope.go:9`, `parser.go:22`, `validator.go:12`, `errors.go:12`, `internal/match/phases.go:15`, `visibility.go:9`, `rummy_match.go:448`.*
