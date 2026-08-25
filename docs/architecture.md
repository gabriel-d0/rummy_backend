# Architecture — Rummy Backend

**Stack:** Go 1.23.5 plugin on Nakama 3.26.0 (`postgres:15-alpine`), Docker Compose `rummy_backend:local`. Small, explicit, server-authoritative.

## Layout

```
main.go                    # InitModule: health/version RPCs + RegisterMatch("rummy")
internal/
  match/                   # authoritative state machine + handlers
    state.go               # RoundState, GamePhase, TurnPhase, DiscardEntry, TableMeld
    seat.go                # PlayerId, Seat, AssignSeats, NextSeat/PrevSeat
    invariant.go           # CheckTileConservation (106)
    visibility.go          # PublicView/PrivateView projection
    phases.go              # AllowedOps matrix + ValidatePhase
    validate.go            # ValidateActivePlayer / ValidatePhaseOp
    turn.go                # AdvanceTurn
    discard.go             # CanPickupPreviousDiscard / CanPickupDiscardForMeld
    rummy_match.go         # RummyMatch (MatchInit/Join/Loop/Signal/Terminate) + dispatch + sendError
    meld_initial.go        # handleMeldInitial (50+ with run)
    meld_new.go            # handleMeldNew (opened, no minimum)
    # next: extend_meld.go, joker_replace.go, win.go
  rules/
    tile/tile.go           # Colour, Rank, TileInstanceId, TileInstance, Joker
    meld/
      meld.go              # Meld{ID,Kind,Tiles,JokerReps}, Validate, New
      set.go / run.go      # ValidateSet / ValidateRun (pure, ratio, Ace, jokers)
      errors.go            # ValidationError{Code,Field,Message}
    scoring/
      scoring.go           # ScoreTile with Ace/joker context
      run.go / set.go      # ScoreRun/Set
      batch.go / validate.go # Batch, TotalScore, ValidateBatchOwnership/Score/HasRun/InitialBatch
  setup/
    deck.go / shuffle.go / rand.go / deal.go / round.go  # 106-tile factory, seeded Rand, Deal 15/14/stock, NewRoundState
  protocol/
    opcodes.go / envelope.go / parser.go / validator.go / errors.go  # Version 1, opcodes 1..9/100..199, Envelope, Parse/Validate, ErrorResponse
docs/
  rules-decisions.md       # binding product decisions + TODO(product)
  terminology.md           # canonical names
  protocol.md              # wire format + schemas
  state-machine.md         # GamePhase/TurnPhase + AllowedOps + transitions
  testing.md               # deterministic harness
  daily-log.md             # Handmade Hero slices
```

## Data Flow

```
Client JSON --Envelope{v,op,requestId,payload}--> MatchLoop
  → ValidateEnvelope (ParseEnvelope + ValidatePayload)
  → ValidateActivePlayer (not_member / not_your_turn)
  → ValidatePhaseOp (wrong_phase via AllowedOps)
  → ValidatePayload (bad_payload)
  → handleXxx (pure rules validation)
      → meld.ValidateRun/Set or scoring.ValidateInitialBatch / ValidateBatchOwnership
      → on success: mutate RoundState atomically (Racks→TableMelds/Stock→Racks/DiscardRow, AdvanceTurn, HasOpened)
      → broadcast OpServerEvent / OpServerError (requestId echo)
      → PublicView / PrivateView for snapshots
  → CheckTileConservation (test-only, but invariant holds in prod)
```

Rules modules are **pure** (`no I/O, no Nakama imports, injectable Rand`), match handlers are **thin orchestrators** (no meld math). All randomness is seedable for tests.

## State

`RoundState` is the single authoritative source (authoritative match `state interface{}`), kept in memory per match, not yet persisted beyond the process. `Racks` are private per `Seat`, `Stock` ordered (top = last element), `DiscardRow` ordered with `Index`, `TableMelds` public with stable `ID` and explicit immutable `JokerReps`, `CurrentSeat`/`GamePhase`/`TurnPhase`/`Winner`.

## Visibility and Security

- `visibility.go` centralizes `PublicView` (counts only) and `PrivateView` (ownRack copy). No other player’s `TileInstanceId` appears in public JSON (verified by `visibility_test.go` and `redaction_test.go` exhaustive).
- Server never trusts client meld validity; `tileIds` are resolved via `Racks[seat]` → `TileInstance` (colour/rank from server, not client). Joker `represented` tile is explicit in `payload.jokerReps` and revalidated.
- `sendError` sends `OpServerError 102` to sender only (`presences=[sender]`), with `requestId` correlation, so clients can retry idempotently; future days will add `RequestId` deduplication.

## Nakama Lifecycle

`NewRummyMatch` factory → `MatchInit` (Waiting, tickRate 5, label `rummy`) → `MatchJoinAttempt` (check Waiting, 4 max, duplicate) → `MatchJoin` (deterministic Seat 0..n-1, `Racks` init empty) → `MatchLoop` (envelope dispatch) → `MatchSignal` (`"start"` / `"start:alice"` for local dev) → `MatchLeave` (delete Racks, keep seat numbers stable) → `MatchTerminate`. `MatchLabelUpdate` keeps lobby listing `rummy:2` etc.

## Adding a New Command

(See `README.md#how-to-add-a-new-command-safely` for checklist.) In short: stabilize opcode, add `ValidatePayload` schema, add `AllowedOps`, write pure validator if needed, write `handleXxx` with full pre-validation and atomic mutation, wire in `MatchLoop`, update `visibility` if needed, write success/atomic/phase/duplicate/joker/conservation/redaction tests, `make check` + `docker compose build`, update `docs/protocol.md`/`state-machine.md`.

## Deferred

No wall geometry, Doubla, exposed-tile bonus, multi-round scoring/dealer rotation, spectators, chat, friends, bots, spectators, production infra (Terraform, TLS, WAF, backups) — see `docs/rules-decisions.md:7` and roadmap Phase 15+.

---

*Code: `main.go:12`, `internal/match/state.go:12`, `protocol/opcodes.go:8`, `rules/meld/run.go:16`, `rules/scoring/scoring.go:16`, `setup/deck.go:10`.*
