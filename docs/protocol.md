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
| 4 | `OpClientDrawPreviousDiscard` | `{}` | `Playing/MustDraw` (opened, not opening) — *implemented `CanPickupPreviousDiscard` but handler pending*. |
| 5 | `OpClientPickupDiscardForMeld` | `{"discardIndex":int>=0, "tileIds":["...","..."]}` exactly 2 | `Playing/MustDraw` (opened) — *handler pending*. `discardIndex` must not be opening, `tileIds` owned. |
| 6 | `OpClientMeldInitial` | `{"melds":[{"id":"m1","kind":"run"\|"set","tileIds":["t1","t2","t3"], "jokerReps":{"jokerId":{"colour":"red","rank":7}}} , ...]}` ≥1 meld | `Playing/MeldOrDiscard`, `HasOpened==false`, each meld 3+ tiles owned, each `ValidateRun/Set` passes, total `>=50` with ≥1 run, no duplicate `TileId` across batch, `meldId` unique in batch. Atomically `Racks→TableMelds` and `HasOpened=true`; stays `MeldOrDiscard`. |
| 7 | `OpClientMeldNew` | same shape as 6 | `Playing/MeldOrDiscard`, `HasOpened==true`, each meld valid, ownership, no duplicate, `meldId` not colliding with existing `TableMelds`; no score minimum; atomic; supports multiple melds per batch. |
| 8 | `OpClientExtendMeld` | `{"meldId":"...","tileIds":["..."]}` ≥1 | `Playing/MeldOrDiscard` (opened) — next: revalidate entire resulting meld, any owner’s meld. |
| 9 | `OpClientReplaceJoker` | `{"targetMeldId":"...","tileId":"...","newMeldTiles":["...","..."]}` exactly 2 | `Playing/MeldOrDiscard` (opened) — deferred; run exact tile or set missing colour, then new meld with 2 rack tiles. |

`ValidatePayload` (`internal/protocol/validator.go:12`) enforces these schemas with `bad_payload` including field name. Unknown opcode or bad version never reaches `MatchLoop` payload handling — `ParseEnvelope` rejects before.

## Server → Client Opcodes

| Op | Name | Payload | When |
|---|---|---|---|
| 100 | `OpServerState` | `PrivateSnapshot` (see below) | match start / reconnect (per-player) |
| 101 | `OpServerStatePublic` | `PublicSnapshot` | broadcast (if used) |
| 102 | `OpServerError` | `ErrorResponse{code,message,details,requestId,op}` JSON | on any rejection |
| 103 | `OpServerEvent` | `{"phase":...,"currentSeat":...,"op":...,"melds":...}` JSON | `started`, `discard`, `drawStock`, `meldInitial`, `meldNew` |

Current implementation broadcasts `OpServerEvent` JSON strings (e.g., `{"op":"meldInitial","seat":0,"melds":3}`) plus `OpServerError` to sender only via `dispatcher.BroadcastMessage` with `presences=[sender]` (`internal/match/rummy_match.go:448`).

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

`PrivateSnapshot` (`visibility.go:30`) is `PublicSnapshot` plus `ownRack []TileInstance` and `ownSeat`. Server sends `PrivateView(state, seat)` per player; never sends foreign racks (verified by `internal/setup/redaction_test.go` exhaustive: `n=2,3,4 × seeds 42,7,123` and `internal/match/visibility_test.go`).

## Adding a New Opcode

See `README.md#how-to-add-a-new-command-safely` and `opcodes.go:1` comment: pick next `1..99` not used, add to `IsClientOp` switch in `parser.go:42`, add schema to `validator.go`, add `AllowedOps` in `phases.go`, add handler and wire in `MatchLoop`.

---

*Code references: `internal/protocol/opcodes.go:8`, `envelope.go:9`, `parser.go:22`, `validator.go:12`, `errors.go:12`, `internal/match/phases.go:15`, `visibility.go:9`, `rummy_match.go:448`.*
